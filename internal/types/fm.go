package types

// TODO: as this type is being used in both internal and external, it can be moved to a common type later
type FMEstimateRequest struct {
	UserId    string          `json:"userId"`
	FlowType  string          `json:"flowType"`
	Operation string          `json:"operation"`
	Payload   EstimatePayload `json:"payload"`
}

type EstimatePayload struct {
	NftContractAddress string `json:"nftContractAddress"`
	NftId              string `json:"nftId"`
	NftType            string `json:"nftType"`
	Amount             string `json:"amount"`
	NetworkId          string `json:"networkId"`
	SenderAddress      string `json:"senderAddress"`
	RecepientAddress   string `json:"recepientAddress"`
}

type FMEstimateResponse struct {
	Status string           `json:"status"`
	Result FMEstimateResult `json:"result"`
}
type FMEstimateResult struct {
	JobId  string                 `json:"jobId"`
	Status string                 `json:"status"`
	Output FMEstimateResultOutput `json:"output"`
	Error  Error                  `json:"error,omitempty"`
}

type Error struct {
	ErrorCode         string      `json:"errorCode"`
	Name              string      `json:"name"`
	StatusCode        interface{} `json:"statusCode"`
	StandardErrorCode string      `json:"standardErrorCode"`
	Details           interface{} `json:"details,omitempty"`
	Message           string      `json:"message"`
}

type FMEstimateResultOutput struct {
	TransactionFee    []FMEstimateTxFee  `json:"transactionFee"`
	IsGsnRequired     bool               `json:"isGsnRequired"`
	IsGsnPossible     bool               `json:"isGsnPossible"`
	GsnWithdrawTokens []GsnWithdrawToken `json:"gsnWithdrawTokens"`
	Success           bool               `json:"success"`
}

type FMEstimateTxFee struct {
	NetworkId   string `json:"networkId"`
	NetworkType string `json:"networkType"`
	ChainId     string `json:"chainId"`
	UserAddress string `json:"userAddress"`
	GasAmount   string `json:"gasAmount"`
	TokenId     string `json:"tokenId"`
}

type GsnWithdrawToken struct {
	TokenAddress               string `json:"tokenAddress"`
	UserAddress                string `json:"userAddress"`
	TokenAmount                string `json:"tokenAmount"`
	NetworkId                  string `json:"networkId"`
	TokenAmountWithoutSlippage string `json:"tokenAmountWithoutSlippage"`
	UsdValueWithoutSlippage    string `json:"usdValueWithoutSlippage"`
	UsdValue                   string `json:"usdValue"`
	WithdrawGas                string `json:"withdrawGas"`
	TokenId                    string `json:"tokenId"`
}

type FMExecuteRequest struct {
	UserId    string         `json:"userId"`
	FlowType  string         `json:"flowType"`
	Operation string         `json:"operation"`
	Payload   ExecutePayload `json:"payload"`
	JobId     string         `json:"jobId"`
}

type ExecutePayload struct {
	NftContractAddress  string `json:"nftContractAddress"`
	NftId               string `json:"nftId"`
	NftType             string `json:"nftType"`
	Amount              string `json:"amount"`
	NetworkId           string `json:"networkId"`
	SenderAddress       string `json:"senderAddress"`
	RecepientAddress    string `json:"recepientAddress"`
	IsGsnRequired       bool   `json:"isGsnRequired"`
	Deadline            string `json:"deadline"`
	GsnIncludeToken     string `json:"gsnIncludeToken"`
	GsnIncludenetworkId string `json:"gsnIncludenetworkId"`
	GsnIncludeMaxAmount string `json:"gsnIncludeMaxAmount"`
}

type OpenAPINFTOrder struct {
	ExecuteData
	NetworkID string `json:"networkId,omitempty"`
	UserId    string `json:"userId"`
	JobId     string `json:"jobId"`
	Sponsor   bool   `json:"sponsor"`
}

type ExecuteData struct {
	Payload  any   `json:"payload,omitempty"` // required for aptos
	Deadline int64 `json:"deadline"`          // required for aptos
}

type OpenAPIExecutePayload struct {
	NftContractAddress string `json:"nftContractAddress"`
	NftId              string `json:"nftId"`
	NftType            string `json:"nftType"`
	Quantity           string `json:"quantity"`
	NetworkId          string `json:"networkId"`
	SenderAddress      string `json:"senderAddress"`
	RecepientAddress   string `json:"recepientAddress"`
}

type OpenAPIMintPayload struct {
	NftContractAddress string      `json:"nftContractAddress"`
	NftType            string      `json:"nftType"`
	Quantity           string      `json:"quantity"`
	NetworkId          string      `json:"networkId"`
	SenderAddress      string      `json:"senderAddress"`
	RecepientAddress   string      `json:"recepientAddress"` // Will be same as SenderAddress for now
	MetaData           NftMetaData `json:"metaData"`
}

type NftMetaData struct {
	CollectionName string             `json:"collectionName"`
	NftName        string             `json:"nftName"`
	Uri            string             `json:"uri"`
	Description    string             `json:"description"`
	Properties     *[]ExtraProperties `json:"properties,omitempty"`
}

type ExtraProperties struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Value string `json:"value"`
}

type FMExecuteResponse struct {
	Status string          `json:"status"`
	Result FMExecuteResult `json:"result"`
	Error  Error           `json:"error,omitempty"`
}
type FMExecuteResult struct {
	JobId  string `json:"jobId"`
	Status string `json:"status"`
	Error  Error  `json:"error,omitempty"`
}

type FMGetStatusResult struct {
	JobId    string           `json:"jobId"`
	FlowId   string           `json:"flowId"`
	UserId   string           `json:"userId"`
	VendorId string           `json:"vendorId,omitempty"`
	Status   string           `json:"status"`
	Output   *GetStatusOutput `json:"output,omitempty"`
	Error    Error            `json:"error,omitempty"`
	Metadata OrderMetadata    `json:"metadata"`
}

type OrderMetadata struct {
	Transactions []TransactionsData `json:"transactions,omitempty"` //strongly type for the required ones
}

type TransactionsData struct {
	JobId            string                 `json:"jobId"`
	TransactionHash  string                 `json:"transactionHash"`
	TransactionState string                 `json:"transactionState"`
	OrderTxType      string                 `json:"orderTxType"`
	PayloadType      string                 `json:"payloadType"`
	GasUsed          int                    `json:"gasUsed"`
	GasPrice         any                    `json:"gasPrice"`
	TokenTransfers   []map[string]any       `json:"tokenTransfers"`
	NetworkId        string                 `json:"networkId"`
	NetworkType      string                 `json:"networkType"`
	WalletAddress    string                 `json:"walletAddress"`
	Receipt          map[string]interface{} `json:"receipt"`
}

type FMGetStatusResponse struct {
	Status string            `json:"status"`
	Result FMGetStatusResult `json:"result"`
}

type GetStatusOutput struct {
	Error *GetStatusOutputError `json:"error,omitempty"`
}

type GetStatusOutputError struct {
	Detail    map[string]any `json:"detail"`
	Message   string         `json:"message"`
	Success   string         `json:"success"`
	ErrorCode string         `json:"errorCode"`
}
