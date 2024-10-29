package cmd

type MailmanConfig struct {
	HTMLOutputFile   string   `mapstructure:"htmlOutputFile" validate:"required"`
	MarkdownFile     string   `mapstructure:"markdownFile" validate:"required"`
	HTMLTemplateFile string   `mapstructure:"htmlTemplateFile" validate:"required"`
	HTMLFile         string   `mapstructure:"htmlFile" validate:"required"`
	Host             string   `mapstructure:"host" validate:"required"`
	Port             int      `mapstructure:"port" validate:"required"`
	Subject          string   `mapstructure:"subject" validate:"required"`
	From             string   `mapstructure:"from" validate:"required"`
	To               []string `mapstructure:"to" validate:"required"`
}

type HTMLMessage struct {
	Subject string
	Message string
}
