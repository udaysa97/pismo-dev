package regexUtil

import "regexp"

func IsValidBlockchainAddress(address string) bool {
	pattern := `^(0x)?[0-9a-fA-F]{40}$`

	regex := regexp.MustCompile(pattern)

	return regex.MatchString(address)
}
func IsValidAptosBlockchainAddress(address string) bool {
	pattern := `^(0x)?[0-9a-fA-F]{64}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(address)
}
