package setup

import (
	"find-doctor/service/operations/send_user_message"

	"github.com/victormf2/gox"
	"github.com/victormf2/gox/httpx"
)

func Setup() gox.IApplication {
	app := gox.New()
	app.AddEndpoint(httpx.Endpoint("POST /messages", send_user_message.NewOperation))

	return app.Build()
}
