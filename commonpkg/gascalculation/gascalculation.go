package gascalculation

import (
	"math"

	"github.com/shopspring/decimal"
)

func CalculateGasFees(gasAmount string, decimals int) string {
	number1, _ := decimal.NewFromString(gasAmount)
	dividend := decimal.NewFromFloat(math.Pow(10, float64(decimals)))
	result := number1.Div(dividend).Round(18).String()
	return result

}
