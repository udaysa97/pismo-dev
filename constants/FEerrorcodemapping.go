package constants

var FE_ERROR_CODES_MAPPING = map[string]FEErrorType{
	"FE0001": {
		ErrorCode: "FE0001",
		Message:   "Insufficient funds to pay for the gas fee",
	},
	"FE0002": {
		ErrorCode: "FE0002",
		Message:   "Insufficient funds in your wallet",
	},
	"FE0003": {
		ErrorCode: "FE0003",
		Message:   "Insufficient fund to support GSN",
	},
	"FE0004": { //Deadline exceeded
		ErrorCode: "FE0004",
		Message:   "Deadline exceeded",
		SubErrorName: map[string]string{
			"SIGN_CONSUMER":    "Deadline exceeded - user did not complete signing",
			"SEND_FOR_SIGNING": "Deadline exceeded - system did not complete signing",
			"TxLifeCycle":      "Deadline exceeded - transaction failed on chain after signing",
		},
	},
	"FE0005": { //Downstream failure
		ErrorCode: "FE0005",
		Message:   "Downstream failure",
	},
	"FE0006": {
		ErrorCode: "FE0006",
		Message:   "Path not found to perform Cross chain",
	},
	"FE0007": {
		Message:   "Transaction failed on chain after successful submission",
		ErrorCode: "TTS0007",
	},
	"FE0008": {
		Message:   "Order rejected as the different transaction was encountered during processing",
		ErrorCode: "FE0008",
	},
}

type FEErrorType struct {
	ErrorCode    string
	Message      string
	SubErrorName map[string]string
}

type SubErrorType struct {
	SIGN_CONSUMER    string
	SEND_FOR_SIGNING string
	TxLifeCycle      string
}
