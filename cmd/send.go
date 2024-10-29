package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wneessen/go-mail"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send a custom email as html file",
	Run:   sendRun,
}

func sendRun(Cmd *cobra.Command, args []string) {
	validateConfig(conf)
	sendMail()
}

func sendMail() {
	htmlData, readFileErr := os.ReadFile(conf.HTMLFile)
	if readFileErr != nil {
		fmt.Println("Error Reading HTML File:", readFileErr)
		return
	}
	m := mail.NewMsg()

	if err := m.From(conf.From); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}

	m.Subject(conf.Subject)
	m.SetBodyString(mail.TypeTextHTML, string(htmlData))

	c, err := mail.NewClient(conf.Host, mail.WithPort(conf.Port), mail.WithTLSPolicy(mail.NoTLS))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}

	for _, str := range conf.To {
		if err := m.To(str); err != nil {
			log.Fatalf("failed to set To address: %s", err)
		}

		if err := c.DialAndSend(m); err != nil {
			log.Fatalf("failed to send mail: %s", err)
		}
	}
	fmt.Println("Mails successfully send!")
}
