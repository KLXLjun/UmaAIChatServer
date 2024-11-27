package vitsfast

import (
	utils "UmaAIChatServer/Utils"
	"UmaAIChatServer/Utils/logx"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

const tag = "API.VITS-fast"

const APIEndPoint = "http://127.0.0.1:7860"

var Characters = make([]string, 0)
var Languages = make([]string, 0)
var sessionHash = utils.GenerateID(11)
var IsOk = false

func Init() {
	chara, Lang := GetCharacterAndLanguage()
	if chara == nil || Lang == nil {
		return
	}
	if len(Languages) != 3 {
		return
	}
	IsOk = true
	utils.LangList = [3]string(Lang)
	Characters = chara
	Languages = Lang
	logx.Info(tag, fmt.Sprintf("Characters %d", len(Characters)))
}

func GetCharacterAndLanguage() ([]string, []string) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "text/html").
		Get(APIEndPoint)
	if err != nil {
		logx.Warn(tag, "Failed to request VITS-fast, please check if it has been started. 请求 VITS-fast 失败，请检查是否已启动", err.Error())
		return nil, nil
	}

	if resp.StatusCode() != 200 {
		return nil, nil
	}

	t := string(resp.Body())
	start := strings.Index(t, "<script>window.gradio_config = ")
	if start != -1 {
		start += 31
	} else {
		logx.Warn(tag, "An exception was encountered when requesting VITS-fast; please check if it is running properly. 请求VITS-fast时发现异常，请检查是否正常运行")
		return nil, nil
	}
	t = t[start:]
	end := strings.Index(t, ";</script>")
	t = t[:end]

	result := make([]string, 0)
	result2 := make([]string, 0)

	components := gjson.Get(t, "components").Array()
	for _, v := range components {
		label := v.Get("props.label")
		if label.Str == "character" {
			choices := v.Get("props.choices").Array()
			for _, v := range choices {
				result = append(result, v.Str)
			}
		}
		if label.Str == "language" {
			choices := v.Get("props.choices").Array()
			for _, v := range choices {
				result2 = append(result2, v.Str)
			}
		}
	}
	return result, result2
}

func GenerateAudio(index int, text string, language int, speed float32) []byte {
	client := resty.New()
	if index >= len(Characters) {
		return nil
	}

	charaName := Characters[index]
	ps := VitsPostData{
		Data: []interface{}{
			text,
			charaName,
			Languages[language],
			speed,
		},
		Index:       index,
		SessionHash: sessionHash,
	}

	resp, err := client.R().
		SetBody(ps).
		SetHeader("Accept", "application/json").
		Post(APIEndPoint + "/run/predict/")
	if err != nil {
		logx.Warn(tag, "Request Failed", err.Error())
	}

	if resp.StatusCode() == 200 {
		respContent := string(resp.Body())
		rsl := gjson.Get(respContent, "data").Array()
		if rsl[0].Str == "Success" {
			filePath := rsl[1].Get("name").Str
			ok, fileBytes := utils.EasyFileRead(filePath)
			if ok {
				return fileBytes
			}
			return nil
		}

		return nil
	}
	logx.Warn(tag, "Generate Failed", resp.StatusCode(), text, "Body", string(resp.Body()))
	return nil
}
