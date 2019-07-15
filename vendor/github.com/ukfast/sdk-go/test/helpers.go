package test

import "net/http"

type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

type TestReadCloser struct {
	ReadError  error
	CloseError error
}

func (r *TestReadCloser) Read(p []byte) (n int, err error) {
	return 0, r.ReadError
}

func (r *TestReadCloser) Close() error { return r.CloseError }
