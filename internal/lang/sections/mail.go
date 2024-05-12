package sections

type Mail struct {
	Registration Registration `yaml:"registration"`
	Login        Login        `yaml:"login"`
	Activation   Activation   `yaml:"activation"`
}

type Registration struct {
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}

type Login struct {
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}

type Activation struct {
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}
