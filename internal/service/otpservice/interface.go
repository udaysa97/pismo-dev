package otpservice

import (
	"context"
	"pismo-dev/internal/types"
)

type OTPServiceInterface interface {
	GenerateOTP(ctx context.Context, userDetails types.UserDetailsInterface, nftTokenDetails types.NFTTransferDetailsInterface) (bool, error)
	SetRequiredServices(services RequiredServices)
	CheckEligibility(ctx context.Context, user types.UserDetailsInterface) bool
	MatchReloginPin(ctx context.Context, user types.UserDetailsInterface) (bool, error)
	MatchTransferPin(ctx context.Context, user types.UserDetailsInterface) bool
}
