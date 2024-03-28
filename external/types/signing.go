package types

type SigningSvcResponse struct {
	Address string `json:"address"`
	Error   Error  `json:"error"`
}

type Error struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
