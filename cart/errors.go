package cart

import (
	"errors"
)

var (
	ErrCartNotFound          = errors.New("cart not found")
	ErrProductInCartNotFound = errors.New("product in cart not found")
)
