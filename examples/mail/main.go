package main

import (
	"fmt"
	"time"

	"github.com/cryptopay-dev/yaga/mail"
)

var messages = []string{
	"foobar",
	"foobar", // will not be send
	"aaaaa",
	"bbbbb",
	"foobar", // will not be send
}

var defferedMessage = "foobar"

func main() {
	service, err := mail.New(mail.Options{
		APIKey:            "xxx",
		Recipients:        []string{"foo@bar.com", "no.reply@example.com"},
		FromEmail:         "person@place.com",
		FromName:          "Boss Man",
		SendUniqTimeout:   time.Second * 5,
		RetryErrorTimeout: time.Second,
	})
	if err != nil {
		fmt.Println("Error instantiating client")
		return
	}

	events := service.Events("Welcome aboard!", nil)

	for _, msg := range messages {
		events.Send(msg)
	}

	time.Sleep(time.Second * 6)
	events.Send(defferedMessage) // will be send
}
