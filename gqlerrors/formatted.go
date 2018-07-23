package gqlerrors

import (
	"encoding/json"
	"errors"

	"github.com/equinux/graphql/language/location"
)

type ExtendedError interface {
	error
	Extensions() map[string]interface{}
}

type FormattedError struct {
	Message    string                    `json:"message"`
	Locations  []location.SourceLocation `json:"locations"`
	Path       []interface{}             `json:"path,omitempty"`
	Extensions map[string]interface{}    `json:"extensions,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for the `FormattedError` type
// in order to place the `ErrorExtensions` at the top level.
func (g FormattedError) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	if g.Extensions != nil {
		for k, v := range g.Extensions {
			m[k] = v
		}
		m["extensions"] = g.Extensions
	}
	m["message"] = g.Message
	m["locations"] = g.Locations
	return json.Marshal(m)
}

func (g FormattedError) Error() string {
	return g.Message
}

func NewFormattedError(message string) FormattedError {
	err := errors.New(message)
	return FormatError(err)
}

func FormatError(err error) FormattedError {
	switch err := err.(type) {
	case FormattedError:
		return err
	case *Error:
		ret := FormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
			Path:      err.Path,
		}
		if err := err.OriginalError; err != nil {
			if extended, ok := err.(ExtendedError); ok {
				ret.Extensions = extended.Extensions()
			}
		}
		return ret
	case Error:
		return FormatError(&err)
	default:
		return FormattedError{
			Message:   err.Error(),
			Locations: []location.SourceLocation{},
		}
	}
}

func FormatErrors(errs ...error) []FormattedError {
	formattedErrors := []FormattedError{}
	for _, err := range errs {
		formattedErrors = append(formattedErrors, FormatError(err))
	}
	return formattedErrors
}
