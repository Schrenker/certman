package output

import (
	"fmt"
	"net/smtp"

	"github.com/schrenker/certman/internal/jsonparse"
)

// Sendmail takes care of authentication and sends request for message to be formatted, then sends it
func Sendmail(settings *jsonparse.Settings, messages [][]*jsonparse.Vhost) error {

	message := prepareBodyBytes(settings.Days, messages)

	auth := smtp.PlainAuth("", settings.EmailAddr, settings.EmailPass, settings.EmailServer)
	err := smtp.SendMail(
		settings.EmailServer+":"+settings.EmailPort,
		auth,
		settings.EmailAddr,
		settings.EmailDest,
		message,
	)

	if err != nil {
		return err
	}
	return nil
}

func prepareBodyBytes(days []uint16, messages [][]*jsonparse.Vhost) []byte {
	buffer := make([]byte, 0)

	buffer = append(buffer, []byte("Subject: Certificate expiration check")...)

	if len(days) > 0 {
		for i := range days {
			buffer = append(buffer, []byte(fmt.Sprintf("Certificates expiring in %v days\n\n", days[i]))...)
			for j := range messages[i] {
				buffer = append(buffer, []byte(fmt.Sprintf("%v %v:%v expires on %v\n",
					messages[i][j].Hostname,
					messages[i][j].Domain,
					messages[i][j].Port,
					messages[i][j].Certificate.NotAfter))...)
			}
			buffer = append(buffer, []byte("\n")...)
		}
	}

	buffer = append(buffer, "Errors:\n\n"...)
	lastEl := len(messages) - 1
	for i := range messages[lastEl] {
		buffer = append(buffer, []byte(fmt.Sprintf("%v %v:%v - %v\n",
			messages[lastEl][i].Hostname,
			messages[lastEl][i].Domain,
			messages[lastEl][i].Port,
			messages[lastEl][i].Error))...)
	}

	return buffer
}

//PrepareBody ...
func PrepareBody(days []uint16, messages [][]*jsonparse.Vhost) string {
	buffer := ""

	if len(days) > 0 {
		for i := range days {
			buffer = buffer + fmt.Sprintf("Certificates expiring in %v days\n\n", days[i])
			for j := range messages[i] {
				buffer = buffer + fmt.Sprintf("%v %v:%v expires on %v\n",
					messages[i][j].Hostname,
					messages[i][j].Domain,
					messages[i][j].Port,
					messages[i][j].Certificate.NotAfter)
			}
			buffer = buffer + "\n"
		}
	}

	buffer = buffer + "Errors\n\n"
	lastEl := len(messages) - 1
	for i := range messages[lastEl] {
		buffer = buffer + fmt.Sprintf("%v %v:%v - %v\n",
			messages[lastEl][i].Hostname,
			messages[lastEl][i].Domain,
			messages[lastEl][i].Port,
			messages[lastEl][i].Error)
	}

	return buffer
}
