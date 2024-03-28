package types

type TransferNotificationKafkaDataInterface struct {
	JobId      string `json:"jobId"`
	UserId     string `json:"userId"`
	EntityId   string `json:"entityId,omitempty"`
	EntityType string `json:"entityType,omitempty"`
	OrderType  string `json:"orderType"`
	VendorId   string `json:"vendorId,omitempty"`
	NetworkId  string `json:"networkId"`
}

type AmplitudeTransferEventInterface struct {
	AppName          string `json:"app_name"`
	Source           string `json:"source"`
	Network          string `json:"network"`
	Token            string `json:"token"`
	TokenId          string `json:"token_id"`
	NftType          string `json:"nft_type"`
	ErcType          string `json:"erc_type"`
	Type             string `json:"type"`
	OrderId          string `json:"order_id"`
	Status           string `json:"status"`
	DeviceId         string `json:"device_id"`
	UserId           string `json:"user_id"`
	DeviceType       string `json:"device_type"`
	ReceiversAddress string `json:"receivers_wallet_address"`
	TokenCount       int    `json:"no_of_tokens"`
	CollectionId     string `json:"collection_id"`
	CollectionName   string `json:"collection_name"`
	Chain            string `json:"chain"`
	VendorId         string `json:"vendor_id"`
}

type AmplitudeMintEventInterface struct {
	AppName          string `json:"app_name"`
	Network          string `json:"network"`
	TokenId          string `json:"token_id"`
	NftType          string `json:"nft_type"`
	Product          string `json:"product"`
	Type             string `json:"type"`
	Status           string `json:"status"`
	ReceiversAddress string `json:"receivers_wallet_address"`
	TokenCount       int    `json:"no_of_tokens"`
	CollectionId     string `json:"collection_id"`
	CollectionName   string `json:"collection_name"`
	Chain            string `json:"chain"`
	VendorId         string `json:"vendor_id"`
}

type AmplitudeEventInterface struct {
	UserId          string      `json:"user_id"`
	Eventtype       string      `json:"event_type"`
	EventProperties interface{} `json:"event_properties"`
}

type UserDetailsInterface struct {
	Id                string                `json:"id"`
	Source            string                `json:"source"`
	ReloginPin        string                `json:"reloginPin,omitempty"`
	UserOTP           string                `json:"userOTP,omitempty"`
	AuthToken         string                `json:"authToken,omitempty"`
	UserWalletAddress string                `json:"userWalletAddress,omitempty"`
	DeviceId          string                `json:"deviceId"`
	DeviceDetails     DeviceDetailInterface `json:"devieDetails"`
}

type NFTTransferDetailsInterface struct {
	NftId                    string `json:"tokenId"`
	NetworkId                string `json:"networkId"`
	ErcType                  string `json:"ercType"`
	DestinationWalletAddress string `json:"destinationWalletAddress"`
	Amount                   string `json:"amount"`
	TransferType             string `json:"transferType"`
	IsGsnRequired            bool   `json:"isGsnRequired"`
	GsnIncludeToken          string `json:"gsnIncludeToken"`
	GsnIncludeNetworkId      string `json:"gsnIncludeNetworkId"`
	GsnIncludeMaxAmount      string `json:"gsnIncludeMaxAmount"`
	OrderId                  string `json:"orderId,omitempty"`
}

type OpenAPINFTTransferDetailsInterface struct {
	UserId                 string `json:"userId"`
	VendorId               string `json:"vendorId"`
	NetworkId              string `json:"networkId"`
	NFTId                  string `json:"nftId"`
	NFTTokenId             string `json:"nftTokenId"`
	CollectionAddress      string `json:"collectionAddress"`
	CollectionId           string `json:"collectionId"`
	RecipientWalletAddress string `json:"recipientWalletAddress"`
	ErcType                string `json:"ercType"`
	Quantity               string `json:"quantity"` // in gwei
	OperationType          string `json:"operationType"`
	IsGsnRequired          bool   `json:"isGsnRequired"`
	GsnIncludeToken        string `json:"gsnIncludeToken"` // address
	GsnIncludeNetworkId    string `json:"gsnIncludeNetworkId"`
	GsnIncludeMaxAmount    string `json:"gsnIncludeMaxAmount"`
}

type SQSJobInterface struct {
	JobId string `json:"jobId"`
}

type UserAgentDetails struct {
	Browser    string
	Version    string
	Os         string
	OsVersion  string
	DeviceType string
	Platform   string
}

type DeviceDetailInterface struct {
	IPAddress  string `json:"ipAddress"`
	Device     string `json:"device"`
	DeviceType string `json:"deviceType"`
	Source     string `json:"source"`
	UserAgent  string `json:"userAgent"`
}
