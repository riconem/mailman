package cmd

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a custom html file from a markdown File and a html template",
	Run:   generateRun,
}

func generateRun(Cmd *cobra.Command, args []string) {
	validateConfig(conf)
	generateMail(conf.HTMLTemplateFile, conf.MarkdownFile, conf.Subject, conf.HTMLOutputFile)
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func generateMail(htmlTemplateFile string, markdownFile string, subject string, htmlOutputFile string) {
	markDownContent, readFileErr := os.ReadFile(markdownFile)
	if readFileErr != nil {
		fmt.Println("Error Reading File:", readFileErr)
		return
	}
	messageData := []byte(markDownContent)
	messageDataHTML := mdToHTML(messageData)

	templateData, readFileErr := os.ReadFile(htmlTemplateFile)
	if readFileErr != nil {
		fmt.Println("Error Reading Template File:", readFileErr)
		return
	}
	tmpl, _ := template.New("htmlMail").Parse(string(templateData))
	htmlMessage := HTMLMessage{Subject: subject, Message: string(messageDataHTML)}
	var tpl bytes.Buffer

	var err = tmpl.Execute(&tpl, htmlMessage)
	if err != nil {
		panic(err)
	}

	writeErr := os.WriteFile(htmlOutputFile, tpl.Bytes(), 0644)
	if writeErr != nil {
		fmt.Println("Error Writing HTML File:", writeErr)
		return
	}
	fmt.Println("HTML file has been written successfully!")

}
