package sections

type Mail struct {
	Registration Registration `yaml:"registration"`
	Login        Login        `yaml:"login"`
}

type Registration struct {
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}

type Login struct {
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}
