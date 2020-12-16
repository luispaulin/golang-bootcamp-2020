package repository

import "fmt"

type errorHTTP struct {
	status string
	code   int
}

// Error string returned
func (he *errorHTTP) Error() string {
	return fmt.Sprintf("Http error: %v, %v", he.status, he.code)
}
