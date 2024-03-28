package otpservice

import (
	"pismo-dev/external/auth"
	"pismo-dev/external/email"
)

type RequiredServices struct {
	AuthSvc  auth.AuthInterface
	EmailSvc email.EmailInterface
}
