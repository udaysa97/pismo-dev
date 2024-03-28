package types

import extTypes "pismo-dev/external/types"

type DQLChannel struct {
	Result DQLEntityResponseInterface
	Error  error
}

type NetworkChannel struct {
	Result string
	Error  error
}

type SigningServiceChannel struct {
	Result extTypes.SigningSvcResponse
	Error  error
}

type CollectionDqlChannel struct {
	Result DQLNftCollectionResponseInterface
	Error  error
}
