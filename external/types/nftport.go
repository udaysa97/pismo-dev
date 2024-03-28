package types

type NftPortMintResponse struct {
	Response               string        `json:"response"`
	Error                  *NftPortError `json:"error,omitempty"`
	Chain                  string        `json:"chain"`
	ContractAddress        string        `json:"contract_address"`
	TransactionHash        string        `json:"transaction_hash"`
	TransactionExternalUrl string        `json:"transaction_external_url"`
	MetadataUri            string        `json:"metadata_uri"`
	MintToAddress          string        `json:"mint_to_address"`
}

type NftPortError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

type NftPortMintRequest struct {
	Chain           string `json:"chain"`
	ContractAddress string `json:"contract_address"`
	MetadataURI     string `json:"metadata_uri"`
	MintToAddress   string `json:"mint_to_address"`
}

type NftPortStatusResponse struct {
	Response        string        `json:"response"`
	Error           *NftPortError `json:"error,omitempty"`
	Chain           string        `json:"chain"`
	ContractAddress string        `json:"contract_address"`
	TokenID         string        `json:"token_id"`
}
