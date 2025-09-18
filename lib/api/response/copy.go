package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Responsee struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOKK    = "OK"
	StatusErrorr = "Error"
)

func OKK() *Response {
	return &Response{Status: StatusOKK}
}

func Errorr(msg string) *Response {
	return &Response{Status: StatusErrorr, Error: msg}
}

func ValidationErrorss(errs validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %v is required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %v is not a valid url", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %v is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusErrorr,
		Error:  strings.Join(errMsgs, ", "),
	}
}
