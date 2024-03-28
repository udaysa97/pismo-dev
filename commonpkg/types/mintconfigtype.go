package commontypes

type MintConfig struct {
	MintApiCacheLockKey   string
	MintApiCacheLockTTL   int
	StatusApiCacheLockKey string
	StatusApiCacheLockTTL int
	MaxRetries            int
}
