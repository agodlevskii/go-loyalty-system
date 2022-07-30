package order

import "github.com/vivekmurali/luhn"

func validateOrderNumber(n string) bool {
	return luhn.Validate(n)
}
