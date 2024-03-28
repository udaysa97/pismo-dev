package types

type JobExecutionKafka struct {
	UserId    string        `json:"userId"`
	JobId     string        `json:"jobId"`
	Status    string        `json:"status"`
	SubStatus string        `json:"subStatus"`
	Error     *FMKafkaError `json:"error"`
	Metadata  *Transactions `json:"metadata,omitempty"`
}

type FMKafkaError struct {
	ErrorCode string       `json:"errorCode"`
	Details   *interface{} `json:"details"`
	Name      string       `json:"name"`
}

type Transactions struct {
	Transactions *[]TransactionKafkaInterface `json:"transactions,omitempty"`
}

type TransactionKafkaInterface struct {
	Id                string                   `json:"id,omitempty"`
	UserId            string                   `json:"userId"`
	JobId             string                   `json:"jobId"`
	NetworkId         string                   `json:"networkId"`
	NetworkType       string                   `json:"networkType"` //TODO: USE ENUM const (enum someconstant = "value")
	ContractAddress   string                   `json:"contractAddress"`
	PayloadType       string                   `json:"payloadType"`
	OrderTxType       string                   `json:"orderTxType"`
	Signer            any                      `json:"signer"`
	TransactionState  string                   `json:"transactionState"` //TODO: USE ENUM const (enum someconstant = "value")
	WalletAddress     string                   `json:"walletAddress"`
	TransactionHash   string                   `json:"transactionHash"`
	TransactionStatus bool                     `json:"transactionStatus"`
	Nonce             int                      `json:"nonce"`
	GasPrice          any                      `json:"gasPrice"`
	GasUsed           int                      `json:"gasUsed"`
	Block             int                      `json:"block"`
	Timestamp         int64                    `json:"timestamp"`
	TokenTransfers    []map[string]interface{} `json:"tokenTransfers,omitempty"` //will not be present for non EVM events
	Topics            []string                 `json:"topics"`
	Receipt           map[string]interface{}   `json:"receipt,omitempty"` //will not be present for EVM events
}

// type ITxReceiptLogParsed struct {
// 	Name      string             `json:"name"`
// 	Signature string             `json:"signature"`
// 	Address   string             `json:"address"`
// 	Topic     string             `json:"topic"`
// 	Args      ITokenTransferArgs `json:"args,omitempty"`
// }

// type ITokenTransferArgs struct {
// 	From  string `json:"from"`
// 	To    string `json:"to"`
// 	Value string `json:"value"`
// }

type IAptosTransactionReceipt struct {
	Hash      string               `json:"hash"`
	Type      string               `json:"type"`
	Events    []IAptosReceiptEvent `json:"events"`
	Sender    string               `json:"sender"`
	Success   bool                 `json:"success"`
	Timestamp int64                `json:"timestamp"`
}

type IAptosReceiptEvent struct {
	Data           map[string]interface{} `json:"data"`
	Guid           map[string]interface{} `json:"guid"`
	Type           string                 `json:"type"`
	SequenceNumber string                 `json:"sequence_number"`
}
