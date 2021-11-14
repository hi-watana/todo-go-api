package main

type IllegalIdError struct {
}

func (e *IllegalIdError) Error() string {
	return "Illegal ID"
}

type InternalError struct {
}

func (e *InternalError) Error() string {
	return "Internal error"
}