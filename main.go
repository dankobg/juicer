// package main

// import (
// 	"github.com/dankobg/juicer/cmd"
// 	"log"
// )

//	func main() {
//		if err := cmd.Run(); err != nil {
//			log.Fatalf("failed to run juicer chess server")
//		}
//	}
package main

import (
	"bytes"
	"context"
	"log"
	"mime"
	"net/mail"

	"github.com/dankobg/juicer/mailer"
)

func main() {
	c := mailer.SmtpClient{
		Host:       "localhost",
		Port:       465,
		Username:   "danko",
		Password:   "danko1995bgd",
		TLS:        true,
		AuthMethod: "PLAIN",
	}

	b := &bytes.Buffer{}
	b.WriteRune('😀')
	b.WriteString("hello dudes")

	m := mailer.Message{
		From:    mail.Address{Name: "Danko Petrovic", Address: "danbkop@gmail.com"},
		To:      []mail.Address{{Name: "Danko Petrovic", Address: "dankop@gmail.com"}},
		Subject: "Naslov",
		ReplyTo: "kurac@kurac.com",
		HTML:    "<h1>big html</h1>",
		Text:    "some plain text hyaa",
		Bcc:     []mail.Address{{Name: "Kurva mac", Address: "kurva@gmail.com"}, {Name: "Wtf shit", Address: "wtf@gmail.com"}},
		Cc:      []mail.Address{{Name: "Rofel bobic", Address: "rofel@gmail.com"}, {Name: "Le maoo", Address: "lemao@gmail.com"}},
		Headers: map[string]string{
			"X-Mailer":   "Custom OG lang",
			"User-Agent": "Agent Smith",
		},
		Attachments: []mailer.Attachment{
			{
				FileName:    "rofl.png",
				Description: "some file",
				Content:     b,
				Inline:      false,
				ContentType: mime.TypeByExtension(".png"),
			},
		},
	}

	if err := c.Send(context.Background(), &m); err != nil {
		log.Fatalln(err)
	}
}
