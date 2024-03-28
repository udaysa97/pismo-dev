package types

type CrossMintMintResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	ID      string      `json:"id"`
	OnChain OnChainData `json:"onChain"`
}

type OnChainData struct {
	Status          string `json:"status"`
	Chain           string `json:"chain"`
	ContractAddress string `json:"contractAddress"`
}

type CrossMintError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

type CrossMintBYOCMintRequest struct {
	Recipient         string            `json:"recipient"`
	ContractArguments ContractArguments `json:"contractArguments"`
}

type CrossMintMintRequest struct {
	Recipient string `json:"recipient"`
	MetaData  string `json:"metadata"`
}

type ContractArguments struct {
	URI string `json:"uri"`
}

type CrossMintStatusResponse struct {
	ID       string      `json:"id"`
	Metadata interface{} `json:"metadata"`
	OnChain  OnChainInfo `json:"onChain"`
	Error    bool        `json:"error"`
	Message  string      `json:"message"`
}

type OnChainInfo struct {
	Status          string `json:"status"`
	TokenID         string `json:"tokenId"`
	Owner           string `json:"owner"`
	TxID            string `json:"txId"`
	ContractAddress string `json:"contractAddress"`
	Chain           string `json:"chain"`
}
