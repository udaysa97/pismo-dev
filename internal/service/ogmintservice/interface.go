package ogmintservice

import (
	"context"
	"pismo-dev/internal/types"
	"sync"
)

type OGMintServiceInterface interface {
	PushOrderForProcessing(ctx context.Context, userId, collectionId, networkId, orderType, customUriString string) (string, error)
	SetRequiredServices(services RequiredServices)
	SetRequiredRepos(repos RequiredRepos)
	InitiateSqsConsumer(key string)
	DqlCollectionInteraction(wg *sync.WaitGroup, channel chan types.CollectionDqlChannel, ctx context.Context, identifier string)
	NotifyAmplitudeEvent(orderId, userId string, event types.AmplitudeMintEventInterface) error
}
