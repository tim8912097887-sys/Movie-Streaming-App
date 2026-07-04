package validation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)


func BindAndValidate[T any](r *http.Request) (T, error) {
	var input T

	// Check if the request body is missing entirely
	if r.Body == nil || r.Header.Get("Content-Length") == "" {
        return input, errors.New("request body is empty")
    }

	// Bind the request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return input, err
	}

	if err := validate.Struct(input); err != nil {
        return input, err
    }

    return input, nil
}

var validate = validator.New()

func Validate[T any](input T) (T, error) {
	if err := validate.Struct(input); err != nil {
		return input, err
	}
	return input, nil
}