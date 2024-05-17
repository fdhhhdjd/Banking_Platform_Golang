package helpers

import (
	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/constants"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case constants.USD, constants.EUR, constants.VND:
		return true
	}
	return false
}
