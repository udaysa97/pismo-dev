package appconfig

import (
	commontypes "pismo-dev/commonpkg/types"
	"pismo-dev/constants"
)

func GetMintConfigs() map[string]commontypes.MintConfig {
	var mintconfig = map[string]commontypes.MintConfig{
		constants.POLYGON_NFTPORT_UID: {
			MintApiCacheLockKey:   NFTPORT_CACHE_LOCK_KEY_MINT,
			MintApiCacheLockTTL:   NFTPORT_CACHE_LOCK_TTL_MINT,
			StatusApiCacheLockKey: NFTPORT_CACHE_LOCK_KEY_STATUS,
			StatusApiCacheLockTTL: NFTPORT_CACHE_LOCK_TTL_STATUS,
			MaxRetries:            4,
		},
		constants.POLYGON_CROSSMINT_UID: {
			MintApiCacheLockKey:   CROSSMINT_CACHE_LOCK_KEY_MINT,
			MintApiCacheLockTTL:   CROSSMINT_CACHE_LOCK_TTL_MINT,
			StatusApiCacheLockKey: CROSSMINT_CACHE_LOCK_KEY_STATUS,
			StatusApiCacheLockTTL: CROSSMINT_CACHE_LOCK_TTL_STATUS,
			MaxRetries:            4,
		},
	}
	return mintconfig
}
