package sections

type Validation struct {
	Email        string `yaml:"email"`
	Alpha        string `yaml:"alpha"`
	AlphaNum     string `yaml:"alphanum"`
	AlphaUnicode string `yaml:"alphaunicode"`
	Required     string `yaml:"required"`
	Oneof        string `yaml:"oneof"`
}
