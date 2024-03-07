package common

import "fmt"

type ErrorJson struct {
	Message     string `json:"message"`
	ErrorSystem string `json:"error"`
}

func (e *ErrorJson) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.ErrorSystem)
}

func NewErrorJson(msg string, err error) *ErrorJson {
	ej := &ErrorJson{
		Message:     msg,
		ErrorSystem: "",
	}

	if err != nil {
		ej.ErrorSystem = err.Error()
	}

	return ej
}
