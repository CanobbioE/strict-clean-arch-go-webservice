package presenter_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/interfaces/presenter"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/test/testlog"
)

func TestErrorPresenter_Present(t *testing.T) {
	logger := testlog.NewTestLogger()
	tests := []struct {
		name   string
		err    error
		code   int
		expect func(*httptest.ResponseRecorder)
	}{
		{
			name: "skips nil error",
			err:  nil,
			expect: func(res *httptest.ResponseRecorder) {
				wantCode := http.StatusOK
				if res.Code != wantCode {
					t.Errorf("expected response code to be %d, got %d", wantCode, res.Code)
				}
			},
		},
		{
			name: "overwrites status code for domain error ErrBookNotFound",
			err:  domain.ErrBookNotFound,
			code: http.StatusInternalServerError,
			expect: func(res *httptest.ResponseRecorder) {
				wantCode := http.StatusNotFound
				if res.Code != wantCode {
					t.Errorf("expected response code to be %d, got %d", wantCode, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"book not found","status":"Not Found"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "overwrites status code for domain error ErrInvalidBookID",
			err:  domain.ErrInvalidBookID,
			code: http.StatusInternalServerError,
			expect: func(res *httptest.ResponseRecorder) {
				wantCode := http.StatusBadRequest
				if res.Code != wantCode {
					t.Errorf("expected response code to be %d, got %d", wantCode, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"invalid book id","status":"Bad Request"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "redacts internal server errors",
			err:  errors.New("sensitive implementation data"),
			code: http.StatusInternalServerError,
			expect: func(res *httptest.ResponseRecorder) {
				wantCode := http.StatusInternalServerError
				if res.Code != wantCode {
					t.Errorf("expected response code to be %d, got %d", wantCode, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"internal server error","status":"Internal Server Error"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "handles status codes",
			err:  errors.New("return error to caller"),
			code: http.StatusConflict,
			expect: func(res *httptest.ResponseRecorder) {
				wantCode := http.StatusConflict
				if res.Code != wantCode {
					t.Errorf("expected response code to be %d, got %d", wantCode, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"return error to caller","status":"Conflict"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "defaults to internal server error with no http status code",
			err:  errors.New("some error"),
			expect: func(res *httptest.ResponseRecorder) {
				wantCode := http.StatusInternalServerError
				if res.Code != wantCode {
					t.Errorf("expected response code to be %d, got %d", wantCode, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"internal server error","status":"Internal Server Error"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			w := httptest.NewRecorder()
			p := presenter.NewErrorPresenter(logger)
			p.Present(w, tt.err, tt.code)
			if tt.expect != nil {
				tt.expect(w)
			}
		})
	}
}
