package constants

var (
	FM_FLOW_TYPE          = "NFT_TRANSFER"
	FM_ESTIMATE_OPERATION = "estimate"
	FM_EXECUTE_OPERATION  = "execute"
)

type DAPPOrderType string

const (
	NFT_TRANSFER_ORDER DAPPOrderType = "NFT_TRANSFER"
	NFT_MINT_ORDER     DAPPOrderType = "NFT_MINT"
)
