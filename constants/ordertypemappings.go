package constants

const (
	NFT_TRANSFER = "NFT_TRANSFER"
	OG_MINT      = "NFT_OG_MINT"
	SS_MINT      = "NFT_W3A_MINT"
	NFT_MINT     = "NFT_MINT"
)

var ORDER_TYPES = map[string]bool{NFT_TRANSFER: true, OG_MINT: true, SS_MINT: true, NFT_MINT: true}
