package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
)

const (
	B   = 1
	KiB = B << 10
	MiB = B << 20

	MaxBodySize = 1 * MiB
)

var (
	ErrBadJson      = errors.New("Bad JSON")
	ErrInvalidField = errors.New("Invalid field")
	ErrTooLarge     = fmt.Errorf("Request body is too large (max %d bytes)", MaxBodySize)
	ErrValidation   = errors.New("Validation error")
	ErrDeveloper    = errors.New("Developer error")
)

func IsUnknownFieldError(err error) (string, bool) {
	if err == nil {
		return "", false
	}

	errStr := err.Error()
	prefix := "json: unknown field "

	fieldName, ok := strings.CutPrefix(errStr, prefix)
	if ok {
		return strings.Trim(fieldName, `"`), true
	}

	return "", false
}

func SendJSON(w http.ResponseWriter, status int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("failed to encode JSON", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	w.WriteHeader(status)
	w.Write(bytes)
}

func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, map[string]string{"error": message})
}

func validateFields(dst any) error {
	val := reflect.ValueOf(dst).Elem()

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("%w: expected a struct, got %s", ErrDeveloper, val.Kind())
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Tag.Get("required") == "true" {
			if fieldType.Type.Kind() != reflect.Pointer {
				return fmt.Errorf("%w: field %s must be a pointer to use required tag", ErrDeveloper, fieldType.Name)
			}

			if fieldVal.IsNil() {
				jsonTag := fieldType.Tag.Get("json")
				jsonName, _, _ := strings.Cut(jsonTag, ",")
				if jsonName == "" || jsonName == "-" {
					jsonName = fieldType.Name
				}
				return fmt.Errorf("field '%s' is required", jsonName)
			}
		}
	}
	return nil
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("%w: destination is not a pointer", ErrDeveloper)
	}
	if v.IsNil() {
		return fmt.Errorf("%w: destination is a nil pointer", ErrDeveloper)
	}

	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("%w, syntax error", ErrBadJson)
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("%w, incorrect type for field %s", ErrInvalidField, unmarshalTypeError.Field)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("%w, request body is empty", ErrBadJson)
		case errors.As(err, &maxBytesError):
			return ErrTooLarge
		default:
			if fieldName, ok := IsUnknownFieldError(err); ok {
				return fmt.Errorf("%w: unknown field %s", ErrBadJson, fieldName)
			}
			return fmt.Errorf("%w, some unknown error", ErrBadJson)
		}
	}

	err = validateFields(dst)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w, request body must contain only one JSON object and nothing else", ErrBadJson)
	}

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	err := decodeJSON(w, r, dst)
	if err != nil {
		if errors.Is(err, ErrDeveloper) {
			SendError(w, http.StatusInternalServerError, "Some developer error")
			return false
		} else {
			SendError(w, http.StatusBadRequest, err.Error())
			return false
		}
	}
	return true
}
