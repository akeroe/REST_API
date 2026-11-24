package storage

import "errors" 
// тут просто ошибки для нашего storage 
var (
	ErrURLNotFound = errors.New("url not found")
	ErrExists      = errors.New("url exists")
)
