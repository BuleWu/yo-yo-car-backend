package email

import (
	"fmt"
	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
	"log"
	"zavrsni/yo-yo-car/runtimebag"
)

func SendEmail(from string, to string, subject string, body string) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file...")
	}

	message := gomail.NewMessage()

	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	message.SetBody("text/plain", body)

	fmt.Println("Mailtrap host: ", runtimebag.GetEnvString("MAILTRAP_HOST", ""))
	fmt.Println("Mailtrap port: ", int(runtimebag.GetEnvInt("MAILTRAP_PORT", 587)))
	fmt.Println("Mailtrap username: ", runtimebag.GetEnvString("MAILTRAP_USERNAME", ""))
	fmt.Println("Mailtrap api token: ", runtimebag.GetEnvString("MAILTRAP_PASSWORD", ""))

	dialer := gomail.NewDialer(runtimebag.GetEnvString("MAILTRAP_HOST", ""), int(runtimebag.GetEnvInt("MAILTRAP_PORT", 587)), runtimebag.GetEnvString("MAILTRAP_USERNAME", "api"), runtimebag.GetEnvString("MAILTRAP_PASSWORD", ""))

	if err = dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
