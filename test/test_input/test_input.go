package test_input

type TestReadCloser struct {
	ReadError  error
	CloseError error
}

func (r *TestReadCloser) Read(p []byte) (n int, err error) {
	return 0, r.ReadError
}

func (r *TestReadCloser) Close() error { return r.CloseError }

type TestReader struct {
}

func (r *TestReader) Read(p []byte) (n int, err error) {
	return 0, nil
}
