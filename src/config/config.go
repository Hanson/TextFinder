package config

type Config struct {
	Path       string
	Action     string
	Find       string
	Replace    string
	Level      int64
	Except     string
	ExceptList []string
	Only       string
	OnlyList   []string
	Debug      bool
}

var Cfg = &Config{}
