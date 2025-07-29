package email

import (
	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
	"zavrsni/yo-yo-car/runtimebag"
)

const (
	ReservationMadeSubject      string = "Reservation made"
	ReservationConfirmedSubject string = "Reservation confirmed"
)

func SendEmail(to string, subject string, body string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	message := gomail.NewMessage()

	message.SetHeader("From", runtimebag.GetEnvString("APPLICATION_EMAIL", ""))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	message.SetBody("text/plain", body)

	dialer := gomail.NewDialer(runtimebag.GetEnvString("MAILTRAP_HOST", ""), int(runtimebag.GetEnvInt("MAILTRAP_PORT", 587)), runtimebag.GetEnvString("MAILTRAP_USERNAME", "api"), runtimebag.GetEnvString("MAILTRAP_PASSWORD", ""))

	if err = dialer.DialAndSend(message); err != nil {
		return err
	} else {
		return nil
	}
}
