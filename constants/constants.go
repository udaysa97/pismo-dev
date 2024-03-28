package constants

import (
	"net/http"
	"time"
)

var (
	DEVELOPMENT_ENV               = "DEVELOPMENT"
	AMPLITUDE_APP_NAME            = "okto"
	AMPLITUDE_TRANSFER_EVENT_TYPE = "nft_send"
	AMPLITUDE_MINT_EVENT_TYPE     = "buy"
	AMPLITUDE_PRODUCT_TYPE        = "nft"
	AMPLITUDE_EVENT_NAME          = "order_placed"
	USER_ID                       = "user_id"
	HEADER_DEVICE_ID              = "X-Adjust-Device-Id"
	HEADER_DEVICE_TYPE            = "X-Source"
	AUTH_TOKEN                    = "Authorization"
	DefaultErrorMap               = map[int]bool{
		500: true,
		501: true,
		502: true,
		503: true,
	}
	CacheErrorConfigKey                       = "errorConfig"
	API_IDENTIFIER                            = "apiIdentifier"
	DEFAULT_PAGE_SIZE                         = 20
	DEFAULT_SORT_BY                           = "created_at"
	DEFAULT_ORDER_BY                          = "DESC"
	DEFAULT_SORTING_ORDER                     = "DESC"
	REDIS_PREFIX                              = "OKTO::NFTMS"
	SOMETHING_WENT_WRONG                      = "Something went wrong. Please try again"
	BAD_REQUEST_ERROR                         = "BAD_REQUEST"
	DATA_NOT_FOUND_ERROR                      = "DATA_NOT_FOUND"
	DB_ERROR                                  = "DB_ERROR"
	INVALID_DATA_ERROR                        = "INVALID_DATA"
	PROCESS_ERROR                             = "PROCESS"
	RETRY_ERROR                               = "RETRY_ERROR"
	UNAUTHORISED_ERROR                        = "UNAUTHORIZED"
	TIMEOUT_ERROR                             = "TIMEOUT"
	VALIDATION_ERROR                          = "VALIDATION_ERROR"
	METHOD_NOT_ALLOWED_ERROR                  = "METHOD_NOT_ALLOWED"
	NOT_ABLE_TO_CALL_DOWNSTREAM_SERVICE_ERROR = "NOT_ABLE_TO_CALL_DOWNSTREAM_SERVICE"
	LOCKED_BY_OTHER_PROCESS_ERROR             = "LOCKED_BY_OTHER_PROCESS"
	UNPROCESSABLE_ENTITY_ERROR                = "UNPROCESSABLE_ENTITY"
	QUANTITY_ERROR                            = "QUANTITY_ERROR"
	INVALID_OTP_ERROR                         = "INVALID_OTP_ERROR"
	INVALID_PIN_ERROR                         = "INVALID_PIN_ERROR"
	SELF_TRANSFER_ERROR                       = "SELF_TRANSFER_ERROR"
	INVALID_ADDRESS_ERROR                     = "INVALID_ADDRESS_ERROR"
	PENDING_ORDER_ERROR                       = "PENDING_ORDER_ERROR"
	COLLECTION_LIMIT_ERROR                    = "COLLECTION_LIMIT_ERROR"
	USER_OG_LIMIT_ERROR                       = "USER_OG_LIMIT_ERROR"
	WALLET_NOT_BACKED_UP_ERROR                = "WALLET_NOT_BACKED_UP_ERROR"

	DEFAULT_RETRY_ATTEMPT        = 3
	DEFAULT_RETRY_FIXED_INTERVAL = time.Second * 2
	LOGOUT_COUNT                 = 10
	COOL_OFF_COUNT               = 5
	NFT_TRANSFER_PURPOSE         = "NFT_TRANSFER"
	NFT_MINT_PURPOSE             = "NFT_MINT"
	NETWORK_CACHE_PURPOSE        = "NETWORK_ID_NAME_MAP"

	COMMUNICATION_SERVICE_SOURCE         = "okto"
	COMMUNICATION_SERVICE_PURPOSE_PREFIX = "okto_"
	POLYGON_DQL_ID                       = "ae506585-0ba7-32f3-8b92-120ddf940198"

	WALLET_NOT_BACKEDUP_SS_ERROR = "wallet found in incorrect state"
	NETWORK_TYPE_APTOS           = "APT"
	NFT_MINT_APTOS_HASH_TYPE     = "0x4::collection::MintEvent"
)

var MAGIC_OTP = map[string]string{
	"testing": "999147",
}

var ERROR_TYPES = map[string]ErrorType{
	BAD_REQUEST_ERROR: {
		ErrorCode:  "ER-TECH-0001",
		HttpStatus: http.StatusBadRequest,
	},
	DATA_NOT_FOUND_ERROR: {
		ErrorCode:  "ER-TECH-0002",
		HttpStatus: http.StatusNotFound,
	},
	DB_ERROR: {
		ErrorCode:  "ER-TECH-0003",
		HttpStatus: http.StatusInternalServerError,
	},
	INVALID_DATA_ERROR: {
		ErrorCode:  "ER-TECH-0004",
		HttpStatus: http.StatusBadRequest,
	},
	PROCESS_ERROR: {
		ErrorCode:  "ER-TECH-0005",
		HttpStatus: http.StatusInternalServerError,
	},
	RETRY_ERROR: {
		ErrorCode:  "ER-TECH-0006",
		HttpStatus: http.StatusInternalServerError,
	},
	UNAUTHORISED_ERROR: {
		ErrorCode:  "ER-TECH-0007",
		HttpStatus: http.StatusUnauthorized,
	},
	TIMEOUT_ERROR: {
		ErrorCode:  "ER-TECH-0008",
		HttpStatus: http.StatusRequestTimeout,
	},
	VALIDATION_ERROR: {
		ErrorCode:  "ER-TECH-0009",
		HttpStatus: http.StatusBadRequest,
	},
	METHOD_NOT_ALLOWED_ERROR: {
		ErrorCode:  "ER-TECH-0010",
		HttpStatus: http.StatusMethodNotAllowed,
	},
	NOT_ABLE_TO_CALL_DOWNSTREAM_SERVICE_ERROR: {
		ErrorCode:  "ER-TECH-0011",
		HttpStatus: http.StatusServiceUnavailable,
	},
	LOCKED_BY_OTHER_PROCESS_ERROR: {
		ErrorCode:  "ER-TECH-0012",
		HttpStatus: http.StatusLocked,
	},
	UNPROCESSABLE_ENTITY_ERROR: {
		ErrorCode:  "ER-TECH-0013",
		HttpStatus: http.StatusUnprocessableEntity,
	},
	QUANTITY_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0001",
		HttpStatus: http.StatusBadRequest,
	},
	INVALID_PIN_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0002",
		HttpStatus: http.StatusBadRequest,
	},
	INVALID_OTP_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0003",
		HttpStatus: http.StatusBadRequest,
	},
	SELF_TRANSFER_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0004",
		HttpStatus: http.StatusBadRequest,
	},
	INVALID_ADDRESS_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0005",
		HttpStatus: http.StatusBadRequest,
	},
	PENDING_ORDER_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0006",
		HttpStatus: http.StatusBadRequest,
	},
	COLLECTION_LIMIT_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0007",
		HttpStatus: http.StatusBadRequest,
	},
	USER_OG_LIMIT_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0008",
		HttpStatus: http.StatusBadRequest,
	},
	WALLET_NOT_BACKED_UP_ERROR: {
		ErrorCode:  "NFTMS-ER-TECH-0009",
		HttpStatus: http.StatusBadRequest,
	},
}

var ERROR_CODE_TO_STATUS = map[string]int{
	"ER-TECH-0003": http.StatusInternalServerError,
	"ER-TECH-0002": http.StatusNotFound,
	"ER-TECH-0013": http.StatusUnprocessableEntity,
	"ER-TECH-0012": http.StatusLocked,
	"ER-TECH-0011": http.StatusServiceUnavailable,
	"ER-TECH-0010": http.StatusMethodNotAllowed,
	"ER-TECH-0008": http.StatusRequestTimeout,
	"ER-TECH-0007": http.StatusUnauthorized,
}

var ERROR_CODES_TO_TYPES_MAPPING = map[int]string{
	400: "ER-TECH-0001",
	404: "ER-TECH-0002",
	500: "ER-TECH-0003",
	401: "ER-TECH-0007",
	408: "ER-TECH-0008",
	405: "ER-TECH-0010",
	503: "ER-TECH-0011",
	423: "ER-TECH-0013",
}

type ErrorType struct {
	ErrorCode  string
	HttpStatus int
}

var (
	SQSTraceDriverName               = "nft-ms-sqs"
	SQSDelaySecondsDefault           = "120"
	SQSMessageRetentionPeriodDefault = "86400"
	SQSVisibilityTimeoutDefault      = "180"
	DelaySeconds                     = "DelaySeconds"
	MessageRetentionPeriod           = "MessageRetentionPeriod"
	VisibilityTimeout                = "VisibilityTimeout"
)

var ErcTypes = map[string]bool{"ERC1155": true, "ERC721": true, "NFT": true}

var OPERATION_TYPE = FM_FLOW_TYPE
