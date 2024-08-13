package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Most of the code is tested in [../response] package

func Test_ValidationError_match_ErrValidation(t *testing.T) {
	err := BadRequest(errors.New("some-error"), "super message")

	require.ErrorIs(t, err, ErrBadRequest)
	require.EqualError(t, err, "bad request: some-error")
}

func TestErrorMsgFormat(t *testing.T) {
	tests := []struct {
		Name          string
		Err           error
		UserJSON      string
		InternalError string
	}{
		{
			Name:          "BadRequest with the default message",
			Err:           BadRequest(errors.New("some-error")),
			UserJSON:      `{"message": "bad request"}`,
			InternalError: "bad request: some-error",
		},
		{
			Name:          "BadRequest with a custom message",
			Err:           BadRequest(errors.New("some-error"), "some details: %d", 42),
			UserJSON:      `{"message": "some details: 42"}`,
			InternalError: "bad request: some-error",
		},
		{
			Name:          "Unauthorized with the default message",
			Err:           Unauthorized(errors.New("some-error")),
			UserJSON:      `{"message": "unauthorized"}`,
			InternalError: "unauthorized: some-error",
		},
		{
			Name:          "Unauthorized with a custom message",
			Err:           Unauthorized(errors.New("some-error"), "some details: %d", 42),
			UserJSON:      `{"message": "some details: 42"}`,
			InternalError: "unauthorized: some-error",
		},
		{
			Name:          "NotFound with the default message",
			Err:           NotFound(errors.New("some-error")),
			UserJSON:      `{"message": "not found"}`,
			InternalError: "not found: some-error",
		},
		{
			Name:          "NotFound with a custom message",
			Err:           NotFound(errors.New("some-error"), "some details: %d", 42),
			UserJSON:      `{"message": "some details: 42"}`,
			InternalError: "not found: some-error",
		},
		{
			Name:          "Unhandled with the default message",
			Err:           Unhandled(errors.New("some-error")),
			UserJSON:      `{"message": "internal error"}`,
			InternalError: "unhandled: some-error",
		},
		{
			Name:          "Validation with the default message",
			Err:           Validation(errors.New("some-error")),
			UserJSON:      `{"message": "some-error"}`,
			InternalError: "validation: some-error",
		},
		{
			Name:          "Internal with the default message",
			Err:           Internal(errors.New("some-error")),
			UserJSON:      `{"message": "internal error"}`,
			InternalError: "internal: some-error",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var innerErr *Error
			ok := errors.As(test.Err, &innerErr)
			require.True(t, ok, "invalid error struct")

			raw, err := innerErr.MarshalJSON()
			require.NoError(t, err)

			assert.JSONEq(t, test.UserJSON, string(raw))
			require.EqualError(t, innerErr, test.InternalError)
		})
	}
}
