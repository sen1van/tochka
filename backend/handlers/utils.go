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
	ErrBadJson      = errors.New("bad JSON")
	ErrInvalidField = errors.New("invalid field")
	ErrTooLarge     = fmt.Errorf("request body is too large (max %d bytes)", MaxBodySize)
	ErrValidation   = errors.New("validation error")
	ErrDeveloper    = errors.New("developer error")
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
	if status == http.StatusNoContent {
		w.WriteHeader(status)

		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		slog.Error("failed to encode JSON", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte(`{"error": "internal server error"}`))
		if err != nil {
			slog.Error("failed to write error response", "error", err)
		}

		return
	}

	w.WriteHeader(status)

	_, err = w.Write(bytes)
	if err != nil {
		slog.Error("failed to write JSON response", "error", err)
	}
}

func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, map[string]string{"message": message})
}

func validateFields(dst any) error {
	val := reflect.ValueOf(dst).Elem()

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("%w: expected a struct, got %s", ErrDeveloper, val.Kind())
	}

	typ := val.Type()

	for i := range val.NumField() {
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

				return fmt.Errorf("%w: field '%s' is required", ErrValidation, jsonName)
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
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("%w: destination is not a pointer or is nil", ErrDeveloper)
	}

	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError

		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("%w, syntax error", ErrBadJson)
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("%w, incorrect type for field %s", ErrInvalidField, unmarshalTypeError.Field)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("%w, request body is empty", ErrBadJson)
		default:
			if fieldName, ok := IsUnknownFieldError(err); ok {
				return fmt.Errorf("%w: unknown field %s", ErrBadJson, fieldName)
			}

			return fmt.Errorf("%w, some unknown error", ErrBadJson)
		}
	}

	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w, request body must contain only one JSON object and nothing else", ErrBadJson)
	}

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	err := decodeJSON(w, r, dst)
	if err == nil {
		err = validateFields(dst)
	}

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
