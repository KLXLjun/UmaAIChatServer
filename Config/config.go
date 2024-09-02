package config

import (
	utils "UmaAIChatServer/Utils"
	"bytes"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var Conf = Config{
	Port: 32679,
	ChatConfig: ClientConfig{
		APIUrl: "https://api.openai.com/v1",
		APIKey: "sk-test",
		Model:  "gpt-4o-mini",
		Proxy:  "",
	},
	EmotionConfig: ClientConfig{
		APIUrl: "https://api.openai.com/v1",
		APIKey: "sk-test",
		Model:  "gpt-4o-mini",
		Proxy:  "",
	},
	TranslateConfig: ClientConfig{
		APIUrl: "https://api.openai.com/v1",
		APIKey: "sk-test",
		Model:  "gpt-4o-mini",
		Proxy:  "",
	},
}

func LoadConfig(run string) (bool, error) {
	var b bytes.Buffer
	e := toml.NewEncoder(&b)
	e.Encode(Conf)

	_, err := utils.ReadOrCreateFile(path.Join(run, "config.toml"), b.Bytes())
	if err != nil {
		logrus.Warn("The configuration file `config.toml` does not exist; it has now been created. 配置文件 config.toml 不存在,现在已创建")
		return false, nil
	}
	if _, err := toml.DecodeFile(path.Join(run, "config.toml"), &Conf); err != nil {
		return false, err
	}
	return true, nil
}
