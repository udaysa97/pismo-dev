package types

type DQLEntityResponseInterface struct {
	Id      string     `json:"id"`
	Details DQLDetails `json:"details"`
	Error   Error      `json:"error,omitempty"`
}

type DQLDetails struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	NetworkId string `json:"network_id"`
	Address   string `json:"address"`
	//Type          string `json:"type"` //NFT
	ErcType          string `json:"type,omitempty"`
	TokenId          string `json:"nft_token_id"`
	NativeTokenId    string `json:"native_token_id,omitempty"`
	NativeCurrencyId string `json:"native_currency_id,omitempty"`
	Decimals         string `json:"decimals,omitempty"`
	CollectionName   string `json:"collection_name,omitempty"`
	CollectionId     string `json:"collection_id,omitempty"`
}

type DQLByAddressResponseInterface struct {
	Entities []DQLEntityResponseInterface `json:"entities"`
	Error    Error                        `json:"error,omitempty"`
}

type DQLNftCollectionResponseInterface struct {
	Id         string               `json:"id"`
	EntityType string               `json:"entityType"`
	Error      Error                `json:"error,omitempty"`
	Details    DQLCollectionDetails `json:"details"`
}

type DQLCollectionDetails struct {
	IsActive         bool                  `json:"is_active"`
	IsPublished      bool                  `json:"is_published"`
	CreatedAt        string                `json:"created_at"`
	ContractAddress  string                `json:"contract_address"`
	FundedBy         string                `json:"funded_by"`
	NetworkID        string                `json:"network_id"`
	CollectionSymbol string                `json:"collection_symbol"`
	UpdatedAt        string                `json:"updated_at"`
	ContractMetadata DQLCollectionMetaData `json:"contract_metadata"`
	MintType         string                `json:"mint_type"`
	IsMintEnabled    bool                  `json:"is_mint_enabled"`
	ID               string                `json:"id"`
	IsTransferable   bool                  `json:"is_transferable"`
	NftType          string                `json:"nft_type"`
	CollectionName   string                `json:"collection_name"`
}

type DQLCollectionMetaData struct {
	TxHash                string `json:"tx_hash"`
	NftMetadataURI        string `json:"nft_metadata_uri"`
	MintLimitPerWallet    string `json:"mint_limit_per_wallet"`
	NftMintLimit          string `json:"nft_mint_limit"`
	CrossmintCollectionId string `json:"crossmint_collection_id"`
	VendorName            string `json:"vendor_name"`
	// MintingStartTime      string `json:"minting_start_time"`
	// MintingEndTime        string `json:"minting_end_time"`
	//Benefits              string `json:"benefits"`
	//CollectionDescription string `json:"collection_description"`
	//	Image                 string `json:"image"`
}
