package constants

var ENDPOINTS = map[string]map[string]string{
	// Add endpoints like below:
	// "ServiceName": {
	// 	"EndpointType": "Endpoint",
	// }
	"FM": {
		"Execute":   "/api/v1/execute",
		"Estimate":  "/api/v1/execute",
		"GetStatus": "/api/v1/execute/job/%s",
	},
	"DAPP": {
		"Execute": "/vpc/api/v2/order/execute/%s",
	},
	"PORTFOLIO": {
		"GetUserBalance": "/api/vpc/v1/defi/users/%s/nft-balance",
	},
	"DQL": {
		"GetEntityById": "/v1/query/entities/%s",
		"GetEntity":     "/v1/query/entities",
	},
	"AUTH": {
		"VERIFY_RELOGIN_PIN": "/api/vpc/v1/verify_relogin_pin",
		"FORCE_LOGOUT":       "/api/vpc/v1/force_logout",
	},
	"MAIL": {
		"SEND_MAIL": "/api/vpc/v1/emails/send",
	},
	"SIGNING": {
		"GET_WALLET_ADDRESS": "/api/v2/internal/address",
	},
	"NFTPORT": {
		"MINT":   "/v0/mints/customizable",
		"STATUS": "/v0/mints/%s?chain=%s",
	},
	"CROSSMINT": {
		"MINT":   "/api/2022-06-09/collections/%s/nfts",
		"STATUS": "/api/2022-06-09/collections/%s/nfts/%s",
	},
}
