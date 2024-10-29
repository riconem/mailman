package cmd

import (
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var validate *validator.Validate
var conf MailmanConfig

var rootCmd = &cobra.Command{
	Use:   "mailman",
	Short: "Mail Automation on Console",
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfig()
		fmt.Printf("htmlOutputFile: %s\n", config.HTMLOutputFile)
		fmt.Printf("markdownFile: %s\n", config.MarkdownFile)
		fmt.Printf("htmlTemplateFile: %s\n", config.HTMLTemplateFile)
		fmt.Printf("htmlFile: %s\n", config.HTMLFile)
		fmt.Printf("Host: %s\n", config.Host)
		fmt.Printf("Port: %d\n", config.Port)
		fmt.Printf("Subject: %s\n", config.Subject)
		fmt.Printf("From: %s\n", config.From)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./mailman.yaml)")
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(runCmd)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("mailman")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found.")
		os.Exit(1)
	}
	conf = getConfig()
}

func getConfig() MailmanConfig {
	var config = MailmanConfig{}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
		os.Exit(1)
	}
	return config
}

func validateConfig(config MailmanConfig) {
	validate = validator.New()
	if err := validate.Struct(config); err != nil {
		fmt.Printf("%s\n", err.(validator.ValidationErrors))
		os.Exit(1)
	}
	fmt.Println("Validation successful")
}
