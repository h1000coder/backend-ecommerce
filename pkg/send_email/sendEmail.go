package sendemail

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mailersend/mailersend-go"
)


func SendEmail(email string) error {
	emailFrom := os.Getenv("EMAIL")

	ms := mailersend.NewMailersend(os.Getenv("EMAIL_TOKEN"))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)

	defer cancel()

	subject := "Compra SoulStreet"
	text := "Confirmação de compra"
	html := "<h1>Você receberá seu codigo de rastreio em breve!</h1>"

	from := mailersend.From {
		Name: "SoulStreet",
		Email: emailFrom,
	}

	recipients := []mailersend.Recipient {
		{
			Name: "Cliente SoulStreet",
			Email: email,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)
	
	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		return errors.New("Erro ao enviar o email de confirmação de pagamento")
	}

	fmt.Println(res.StatusCode)

	return nil

}