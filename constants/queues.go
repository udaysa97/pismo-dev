package constants

var (
	RECONQUEUE     = "nftTransferReconQueue"
	NFTPORTQUEUE   = "nftPortQueue_"
	CROSSMINTQUEUE = "crossMintQueue_"
)

type OGOPTYPE string

const (
	MINT        OGOPTYPE = "MINT"
	STATUSCHECK OGOPTYPE = "STATUS"
)
