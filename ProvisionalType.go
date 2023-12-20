package main

type DisplayError struct {
	Message string
}

func (err DisplayError) Error() string {
	return err.Message
}
