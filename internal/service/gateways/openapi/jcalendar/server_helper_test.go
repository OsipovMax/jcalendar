package jcalendar

import (
	"bytes"
	"io"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func prepareContext(method, path string, headers map[string]string, payload []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()

	var body io.Reader
	if payload != nil {
		body = bytes.NewBuffer(payload)
	}

	req := httptest.NewRequest(method, path, body)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()

	return e.NewContext(req, rec), rec
}
