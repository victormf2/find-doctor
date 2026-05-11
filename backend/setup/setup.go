package setup

import (
	"find-doctor/service/operations/send_user_message"

	"github.com/victormf2/framework"
	"github.com/victormf2/framework/httpx"
)

func Setup() framework.IApplication {
	app := framework.New()
	app.AddEndpoint(httpx.Endpoint("POST /messages", send_user_message.NewOperation))

	return app.Build()
}
