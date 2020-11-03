package sendmail

import (
	"fmt"
	"net/smtp"

	"github.com/schrenker/certman/internal/jsonparse"
)

//Sendmail takes care of authentication and sends request for message to be formatted, then sends it
func Sendmail(settings *jsonparse.Settings, messages [][]*jsonparse.Vhost) {

	message := prepareMail(settings.Days, messages)

	auth := smtp.PlainAuth("", settings.EmailAddr, settings.EmailPass, settings.EmailServer)
	err := smtp.SendMail(
		settings.EmailServer+":"+settings.EmailPort,
		auth,
		settings.EmailAddr,
		settings.EmailDest,
		message,
	)

	if err != nil {
		fmt.Println(err)
	}
}

func prepareMail(days []uint16, messages [][]*jsonparse.Vhost) []byte {
	buffer := make([]byte, 0)

	subject := "Subject: Weekly certificate check!\n"

	buffer = append(buffer, []byte(subject)...)

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
