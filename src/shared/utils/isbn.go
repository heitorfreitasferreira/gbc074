package utils

type ISBN string

func (isbn ISBN) Validate() bool {
	return true
}
