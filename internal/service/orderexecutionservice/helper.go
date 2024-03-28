package orderexecutionservice

import (
	"github.com/shopspring/decimal"
	"math"
	"strconv"
)

func multiplyByDecimals(amountString string, decimalString string) string {
	amount, _ := decimal.NewFromString(amountString)
	decimals, _ := strconv.Atoi(decimalString)
	decimalPlaces := decimal.NewFromFloat(math.Pow(10, float64(decimals)))
	result := amount.Mul(decimalPlaces).Round(18).String()
	return result

}
