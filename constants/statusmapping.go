package constants

var AMPLITUDE_STATUS_MAPPING = map[string]string{
	CREATED:               AMPLITUDE_EVENT_STATUS["RECEIVED"],
	REJECTED:              AMPLITUDE_EVENT_STATUS["FAILED"],
	SUCCESS:               AMPLITUDE_EVENT_STATUS["SUCCESS"],
	FAILED:                AMPLITUDE_EVENT_STATUS["FAILED"],
	RUNNING:               AMPLITUDE_EVENT_STATUS["SIGNING"],
	WAITING_FOR_SIGNATURE: AMPLITUDE_EVENT_STATUS["SIGNING"],
	SUBMITTED:             AMPLITUDE_EVENT_STATUS["RECEIVED"],
}

var AMPLITUDE_EVENT_STATUS = map[string]string{
	"RECEIVED": "txn_received",
	"FAILED":   "txn_signing_failed",
	"SUCCESS":  "txn_signing_successful",
	"SIGNING":  "txn_signing",
}

var (
	CREATED               = "CREATED"
	CONFIRMED             = "CONFIRMED"
	REJECTED              = "REJECTED"
	SUCCESS               = "SUCCESS"
	FAILED                = "FAILED"
	RUNNING               = "RUNNING"
	DISPUTED              = "DISPUTED"
	WAITING_FOR_SIGNATURE = "WAITING_FOR_SIGNATURE"
	SUBMITTED             = "SUBMITTED"
)

var STATUS_LEVELS = map[string]int{
	CREATED:               0,
	RUNNING:               1,
	WAITING_FOR_SIGNATURE: 2,
	CONFIRMED:             3,
	REJECTED:              4,
	SUCCESS:               4,
	FAILED:                4,
}

var CHECK_ALLOW_OG_MINT = "check_allow_og"

var CHECK_ALLOW_OG_MINT_STATUS = []string{CREATED, SUCCESS, RUNNING, SUBMITTED}
