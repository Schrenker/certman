package output

import (
	"fmt"
	"net/smtp"

	"github.com/schrenker/certman/internal/jsonparse"
)

//Sendmail fetches byte array from prepareBodyBytes, authenticates using credentials from seetings.json and sends email
//
//First argument takes a pointer to settings struct, used for SMTP authentication
//
//Second arguments is an array of vhost arrays, already segregated according to days array
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
					messages[i][j].Certificate.NotAfter.Format("02-01-2006 15:04")))...)
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

//PrepareBody is used when email is not send by certman. Takes care of formatting whole output as string.
//
//First argument is uint16 array that contain days "breakpoints" used as range for composing the final output.
//
//Second arguments is an array of vhost arrays, already segregated according to days array
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
					messages[i][j].Certificate.NotAfter.Format("02-01-2006 15:04"))
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
