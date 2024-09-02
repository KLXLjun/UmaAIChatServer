package config

type Config struct {
	Port            int          `toml:"port"`
	ChatConfig      ClientConfig `toml:"chat"`
	EmotionConfig   ClientConfig `toml:"emotion"`
	TranslateConfig ClientConfig `toml:"translate"`
}

type ClientConfig struct {
	APIUrl string `toml:"api_url"`
	APIKey string `toml:"api_key"`
	Model  string `toml:"model"`
	Proxy  string `toml:"proxy"`
}
