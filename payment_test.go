package validationbenchmark

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_createTag(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
	}{
		{
			name: "fail empty json",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "fail invalid json",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": 10.50"}
				`)),
			},
			code: http.StatusInternalServerError,
		},
		{
			name: "success string",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": "10.50"}
				`)),
			},
			code: http.StatusOK,
		},
		{
			name: "success float",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": 10.5}
				`)),
			},
			code: http.StatusOK,
		},
		{
			name: "success int",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": 10}
				`)),
			},
			code: http.StatusOK,
		},
		{
			name: "fail negative float",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": -10.5}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "fail negative int",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": -10}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "fail negative string",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": "-10.5"}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "fail with 0",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": 0}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createByTag(tt.args.w, tt.args.r)
			rr := tt.args.w.(*httptest.ResponseRecorder)
			if rr.Code != tt.code {
				t.Errorf("got: %v, want: %v", rr.Code, tt.code)
			}
		})
	}
}

func Test_createSafe(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
	}{
		{
			name: "fail empty json",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "fail invalid json",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": 10.50"}
				`)),
			},
			code: http.StatusInternalServerError,
		},
		{
			name: "success string",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": "10.50"}
				`)),
			},
			code: http.StatusOK,
		},
		{
			name: "success float",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": 10.5}
				`)),
			},
			code: http.StatusOK,
		},
		{
			name: "success int",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": 10}
				`)),
			},
			code: http.StatusOK,
		},
		{
			name: "fail negative float",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": -10.5}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "fail negative int",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": -10}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "fail negative string",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": "-10.5"}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			name: "fail with 0",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
					{"amount": 0}
				`)),
			},
			code: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createByTypeSafe(tt.args.w, tt.args.r)
			rr := tt.args.w.(*httptest.ResponseRecorder)
			if rr.Code != tt.code {
				t.Errorf("got: %v, want: %v", rr.Code, tt.code)
			}
		})
	}
}

func BenchmarkCreateTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createByTag(
			httptest.NewRecorder(),
			httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"amount": 50}`),
			),
		)
	}
}

func BenchmarkCreateTypeSafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createByTypeSafe(
			httptest.NewRecorder(),
			httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"amount": 50}`),
			),
		)
	}
}
