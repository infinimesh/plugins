package api

import "fmt"

func NewFailedHTTPRequestError(method, endpoint string, err error) error {
	return fmt.Errorf("failed to make http request for endpoint %s %s: %w", method, endpoint, err)
}

func NewUnexpectedHTTPStatusCodeError(method, endpoint string, statusCode int) error {
	return fmt.Errorf("unexpected http status code: received status code %d from endpoint %s %s", statusCode, method, endpoint)
}
