package cmd

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a local email server to test emails",
	Run:   runRun,
}

func runRun(Cmd *cobra.Command, args []string) {
	// Set up a WaitGroup to ensure all goroutines complete
	var wg sync.WaitGroup
	wg.Add(1)

	go runEmailServer(&wg)

	// Wait for a signal to stop the server
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)

	<-sigterm // Wait for the interrupt signal

	// Server graceful shutdown
	log.Println("Shutting down email server...")
	wg.Done()
	wg.Wait()
	log.Println("Email server stopped.")
}

func runEmailServer(wg *sync.WaitGroup) {
	// Set up the local SMTP server
	conn, err := net.Listen("tcp", ":2525") // Change the port as needed
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Local email server started on port 2525.")

	// Handle incoming connections
	for {
		c, err := conn.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			break
		}

		// Handle each connection concurrently
		go handleConnection(c)
	}

	log.Println("Email server shutting down...")
	wg.Done()
}

func handleConnection(c net.Conn) {
	defer c.Close()

	// Set a read deadline to avoid hanging indefinitely
	err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Printf("Failed to set read deadline: %v", err)
		return
	}

	// Read the data from the connection
	buf := make([]byte, 4096)
	n, err := c.Read(buf)
	if err != nil {
		log.Printf("Failed to read from connection: %v", err)
		return
	}

	// Extract the email message from the received data
	receivedEmail := string(buf[:n])

	// Process the received email as needed (e.g., display, log, etc.)
	log.Println("Received email:")
	log.Println(receivedEmail)

	// Extract the email sender, recipients, subject, and HTML content
	recipients := getEmailRecipients(receivedEmail)
	sender := getEmailSender(receivedEmail)
	subject := getEmailSubject(receivedEmail)
	htmlContent := extractHTMLContent(receivedEmail)

	// Respond with a standard SMTP response
	smtpResponse := "250 OK\r\n"

	if err = c.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Printf("Failed to set write deadline: %v", err)
		return
	}

	// Send the SMTP response to the client
	_, err = c.Write([]byte(smtpResponse))
	if err != nil {
		log.Printf("Failed to send response: %v", err)
		return
	}

	// Process the received email as needed (e.g., forward, save, etc.)
	processEmail(sender, recipients, subject, htmlContent)
}

func extractHTMLContent(emailData string) string {
	contentStart := strings.Index(emailData, "<html>")
	if contentStart == -1 {
		return ""
	}

	contentEnd := strings.Index(emailData, "</html>")
	if contentEnd == -1 {
		return ""
	}

	return emailData[contentStart : contentEnd+len("</html>")]
}

func getEmailRecipients(emailData string) []string {
	recipients := make([]string, 0)

	lines := strings.Split(emailData, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "RCPT TO:") {
			email := strings.TrimPrefix(line, "RCPT TO:")
			email = strings.TrimSpace(email)
			recipients = append(recipients, email)
		}
	}

	return recipients
}

func getEmailSubject(emailData string) string {
	lines := strings.Split(emailData, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Subject:") {
			subject := strings.TrimPrefix(line, "Subject:")
			subject = strings.TrimSpace(subject)
			return subject
		}
	}

	return ""
}

func getEmailSender(emailData string) string {
	lines := strings.Split(emailData, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MAIL FROM:") {
			email := strings.TrimPrefix(line, "MAIL FROM:")
			email = strings.TrimSpace(email)
			return email
		}
	}

	return ""
}

func processEmail(sender string, recipients []string, emailData string, htmlContent string) {
	// Implement your desired email processing logic here
	// (e.g., forward the email, save it to a database, etc.)
	log.Printf("Processing email from: %s", sender)
	log.Printf("Recipients: %v", recipients)
	log.Println("Email Content:")
	log.Println(emailData)
	log.Println("htmlContent:")
	log.Println(htmlContent)
}
