package types

type TokenBalanceInterface struct {
	Id                string `json:"nftId"`
	Quantity          string `json:"quantity"`
	UserId            string `json:"userId"`
	UserWalletAddress string `json:"networkWalletId"`
	NetworkId         string `json:"networkId"`
	EntityId          string `json:"entityId"`
}

type TokenBalanceDataInterface struct {
	Rows []TokenBalanceInterface `json:"rows"`
}

type TokenBalanceResponse struct {
	Status bool                      `json:"success"`
	Result TokenBalanceDataInterface `json:"data,omitempty"`
	Error  *PortfolioErrorInterface  `json:"error,omitempty"`
}

type PortfolioErrorInterface struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
	Status    string `json:"status"`
}
