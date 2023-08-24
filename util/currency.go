package util

const (
	USD = "USD"
	INR = "INR"
	EUR = "EUR"
)

func IsSuppertedCurrency(c string) bool {
	switch c {
	case USD, INR, EUR:
		return true
	}
	return false
}
