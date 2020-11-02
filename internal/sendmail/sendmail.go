package sendmail

import (
	"fmt"
	"net/smtp"
)

//Sendmail takes care of authentication and sends request for message to be formatted, then sends it
func Sendmail(sender, password, server, port string, dest []string) {
	message := []byte("Subject: Let's see if this works now!\n" +
		"test mail, let's hope it is not all in subject! :D")

	auth := smtp.PlainAuth("", sender, password, server)
	err := smtp.SendMail(server+":"+port, auth, sender, dest, message)
	if err != nil {
		fmt.Println(err)
	}
}
