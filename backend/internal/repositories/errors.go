package repositories

import "errors"

var ErrNotFound = errors.New("resource not found")
var ErrConflict = errors.New("resource already exists")
var ErrEmptyCart = errors.New("cart is empty")
