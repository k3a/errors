package errors

import (
	"fmt"
	"net/http"
)

// HTTPError holds public-facing error message, with optional
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`

	Internal error `json:"-"` // internal error, not to be presented to the user
}

// Error implements `error` interface but doesn't reveal internal error.
// This error can be presented to the user.
func (e *HTTPError) Error() string {
	// format nicely but do not reveal internal error!
	str := fmt.Sprintf("code=%d", e.Code)
	if e.Message != "" {
		str += ", message=" + e.Message
	}
	return str
}

// Unwrap satisfies the Go 1.13 error wrapper interface to access internal error
func (e *HTTPError) Unwrap() error {
	return e.Internal
}

// UnwrapHTTPError is like Unwrap except it unwraps only HTTPError type errors.
// If there is multiple *HTTPError instances in the error chain, this function
// returns the first one (the deepest one) in an attempt to reveal the root cause
// of the error presentable to the user.
// It may return the same error e or another *HTTPError instance.
func (e *HTTPError) UnwrapHTTPError() *HTTPError {
	err := e.Internal
	if err != nil {
		if err, ok := err.(*HTTPError); ok {
			return err.UnwrapHTTPError()
		}
	}
	return e
}

// HTTP creates a new HTTP error with http status code and message,
// wrapping an optional intenral error as well.
// If msg is empty, http.StatusText for the code is used.
func HTTP(err error, code int, msg string) error {
	if msg == "" {
		msg = http.StatusText(code)
	}

	return &HTTPError{
		Code:     code,
		Message:  msg,
		Internal: err,
	}
}

// NewHTTPError is a deprecated alias of HTTP
var NewHTTPError = HTTP

// HTTP creates a new HTTP error with http status code and formatted message,
// wrapping an optional intenral error as well.
func HTTPf(err error, code int, template string, args ...interface{}) error {
	return &HTTPError{
		Code:     code,
		Message:  fmt.Sprintf(template, args...),
		Internal: err,
	}
}

// NewHTTPErrorf is a deprecated alias of HTTPf
var NewHTTPErrorf = HTTPf

func HTTPUnsupportedMediaType(err error, msg string) error {
	return HTTP(err, http.StatusUnsupportedMediaType, msg)

}

func HTTPUnsupportedMediaTypef(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusUnsupportedMediaType, template, args...)
}

func HTTPNotFound(err error, msg string) error {
	return HTTP(err, http.StatusNotFound, msg)

}

func HTTPNotFoundf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusNotFound, template, args...)
}

func HTTPUnauthorized(err error, msg string) error {
	return HTTP(err, http.StatusUnauthorized, msg)

}

func HTTPUnauthorizedf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusUnauthorized, template, args...)
}

func HTTPForbidden(err error, msg string) error {
	return HTTP(err, http.StatusForbidden, msg)

}

func HTTPForbiddenf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusForbidden, template, args...)
}

func HTTPMethodNotAllowed(err error, msg string) error {
	return HTTP(err, http.StatusMethodNotAllowed, msg)

}

func HTTPMethodNotAllowedf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusMethodNotAllowed, template, args...)
}

func HTTPRequestEntityTooLarge(err error, msg string) error {
	return HTTP(err, http.StatusRequestEntityTooLarge, msg)

}

func HTTPRequestEntityTooLargef(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusRequestEntityTooLarge, template, args...)
}

func HTTPTooManyRequests(err error, msg string) error {
	return HTTP(err, http.StatusTooManyRequests, msg)

}

func HTTPTooManyRequestsf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusTooManyRequests, template, args...)
}

func HTTPBadRequest(err error, msg string) error {
	return HTTP(err, http.StatusBadRequest, msg)

}

func HTTPBadRequestf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusBadRequest, template, args...)
}

func HTTPBadGateway(err error, msg string) error {
	return HTTP(err, http.StatusBadGateway, msg)

}

func HTTPBadGatewayf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusBadGateway, template, args...)
}

func HTTPInternalServerError(err error, msg string) error {
	return HTTP(err, http.StatusInternalServerError, msg)

}

func HTTPInternalServerErrorf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusInternalServerError, template, args...)
}

func HTTPRequestTimeout(err error, msg string) error {
	return HTTP(err, http.StatusRequestTimeout, msg)

}

func HTTPRequestTimeoutf(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusRequestTimeout, template, args...)
}

func HTTPServiceUnavailable(err error, msg string) error {
	return HTTP(err, http.StatusServiceUnavailable, msg)

}

func HTTPServiceUnavailablef(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusServiceUnavailable, template, args...)
}

func HTTPGone(err error, msg string) error {
	return HTTP(err, http.StatusGone, msg)

}

func HTTPGonef(err error, template string, args ...interface{}) error {
	return HTTPf(err, http.StatusGone, template, args...)
}
