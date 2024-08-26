package openai

import (
	config "UmaAIChatServer/Config"
	"UmaAIChatServer/Utils/logx"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/otiai10/openaigo"
)

const tag = "API.OpenAI"

var messageTaskQueue = goconcurrentqueue.NewFIFO()
var Quit = make(chan os.Signal)
var tokenStack *Queue = NewQueue()

var systemPrompt = "You are a helpful assistant. You can help me by answering my questions. You can also ask me questions."

var maxToken = 8192

var client *openaigo.Client
var emotionClient *openaigo.Client

var tokenLock = sync.Mutex{}

func SetSystemPrompt(input string) {
	systemPrompt = input
}

func InitClient(proxyUrl string, Url string, key string) {
	client = openaigo.NewClient(key)
	client.BaseURL = Url
	var transport *http.Transport
	proxy, proxyError := url.Parse(proxyUrl)
	if proxyError == nil && proxyUrl != "" {
		transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	} else {
		if proxyUrl != "" {
			logx.Warn(tag, "The proxy address is incorrect. 代理地址不正确")
		}
		transport = &http.Transport{}
	}
	client.HTTPClient = &http.Client{Transport: transport}
}

func InitEmotionClient(proxyUrl string, Url string, key string) {
	emotionClient = openaigo.NewClient(key)
	emotionClient.BaseURL = Url
	var transport *http.Transport
	proxy, proxyError := url.Parse(proxyUrl)
	if proxyError == nil && proxyUrl != "" {
		transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	} else {
		if proxyUrl != "" {
			logx.Warn(tag, "The proxy address is incorrect. 代理地址不正确")
		}
		transport = &http.Transport{}
	}
	emotionClient.HTTPClient = &http.Client{Transport: transport}
}

func AddTask(push QueuePrompt) {
	err := messageTaskQueue.Enqueue(push)
	if err != nil {
		logx.Warn(tag, "The queue is locked. 队列被锁定")
		return
	}
}

func ChatTask() {
	for {
		value, _ := messageTaskQueue.DequeueOrWaitForNextElement()
		currentTask := value.(QueuePrompt)
		tokenLock.Lock()
		ScanOutOfToken(tokenStack)
		var resq = GenPromptNew(currentTask.PromptGroup, tokenStack, true)
		var requestOK, Message, useToken = RequestChat(resq)
		if requestOK {
			result := ChatResult{
				Emotion: "默认",
			}
			var requestEmotionOk, EmotionSelect, _ = RequestEmotion(openaigo.Message{
				Role:    ChatMessageRoleUser,
				Content: fmt.Sprintf("问：%s\n答：%s\n请分析上面对话发生时回答者的状态。只需要告诉我状态，不要回复标点符号或者其他内容。可供选择的状态：%s。", currentTask.PromptGroup.Content, Message.Content, currentTask.Emotion),
			})

			if requestEmotionOk {
				eSelect := strings.Split(currentTask.Emotion, ",")
				for _, v := range eSelect {
					if strings.Contains(EmotionSelect.Content, v) {
						result.Emotion = v
						break
					}
				}
			}

			//储存输入
			tokenStack.Enqueue(SavePrompt{
				TokenUse:    useToken.PromptTokens,
				PromptGroup: currentTask.PromptGroup,
			})

			//存储返回
			tokenStack.Enqueue(SavePrompt{
				TokenUse:    useToken.CompletionTokens,
				PromptGroup: Message,
			})

			result.Message = Message.Content
			tokenLock.Unlock()
			currentTask.CallBack <- result
		} else {
			tokenLock.Unlock()
		}
		time.After(time.Millisecond * 50)
	}
}

func ScanOutOfToken(tokenStack *Queue) {
	//最大Token限制
	countToken := 0

	if len(tokenStack.Items) == 0 {
		return
	}

	for {
		countToken = 0
		for _, i2 := range tokenStack.Items {
			countToken += i2.TokenUse
		}
		if countToken < maxToken {
			break
		} else {
			tokenStack.Dequeue()
			tokenStack.Dequeue()
		}
	}
}

func ClearToken() {
	tokenStack.Clear()
}

func GenPromptNew(input openaigo.Message, tokenStack *Queue, mutil bool) []openaigo.Message {
	rsl := make([]openaigo.Message, 0)
	rsl = append(rsl, openaigo.Message{
		Role:    "system",
		Content: systemPrompt,
	})

	for _, item := range tokenStack.Items {
		rsl = append(rsl, item.PromptGroup)
	}

	rsl = append(rsl, input)
	return rsl
}

func RequestEmotion(input openaigo.Message) (bool, openaigo.Message, openaigo.Usage) {
	if emotionClient == nil {
		return false, openaigo.Message{}, openaigo.Usage{}
	}
	request := openaigo.ChatRequest{
		Model: config.Conf.EmotionConfig.Model,
		Messages: []openaigo.Message{
			{Role: ChatMessageRoleSystem, Content: "You are a helpful assistant. You can help me by answering my questions. You can also ask me questions."},
			input,
		},
	}
	ctx := context.Background()
	chat, err := emotionClient.Chat(ctx, request)
	if err != nil {
		return false, openaigo.Message{}, openaigo.Usage{}
	}
	return true, chat.Choices[0].Message, chat.Usage
}

func RequestChat(input []openaigo.Message) (bool, openaigo.Message, openaigo.Usage) {
	if client == nil {
		return false, openaigo.Message{}, openaigo.Usage{}
	}
	request := openaigo.ChatRequest{
		Model:    config.Conf.ChatConfig.Model,
		Messages: input,
	}
	ctx := context.Background()
	chat, err := client.Chat(ctx, request)
	if err != nil {
		return false, openaigo.Message{}, openaigo.Usage{}
	}
	return true, chat.Choices[0].Message, chat.Usage
}
