package auth

import (
	"context"
	"pismo-dev/internal/types"
)

type AuthInterface interface {
	ForceLogout(ctx context.Context, user types.UserDetailsInterface)
	VerifyReloginPin(ctx context.Context, user types.UserDetailsInterface) (bool, error)
}
