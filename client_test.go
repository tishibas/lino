package lino_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/tishibas/lino"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(f RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(f),
	}
}

func TestClient_Notify(t *testing.T) {

	pStr := func(s string) *string { return &s }

	type fields struct {
		accessToken string
		httpClient  *http.Client
	}
	type args struct {
		r *RequestNotify
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success_only_message",
			fields: fields{
				accessToken: "VALID_ACCESS_TOKEN",
				httpClient: NewTestClient(
					func(req *http.Request) *http.Response {
						if req.URL.String() != "https://notify-api.line.me/api/notify" {
							panic("url is different")
						}
						if req.Method != "POST" {
							panic("method is not POST")
						}
						header := make(http.Header)
						header["Content-Type"] = []string{"application/json"}
						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(bytes.NewBufferString(`{"status":200,"message":"ok","target":"foo"}`)),
							Header:     header,
						}
					},
				),
			},
			args: args{
				r: &RequestNotify{
					Message: "This is test.",
				},
			},
		},
		{
			name: "success_images",
			fields: fields{
				accessToken: "VALID_ACCESS_TOKEN",
				httpClient: NewTestClient(
					func(req *http.Request) *http.Response {
						header := make(http.Header)
						header["Content-Type"] = []string{"application/json"}
						return &http.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(bytes.NewBufferString(`{"status":200,"message":"ok","target":"foo"}`)),
							Header:     header,
						}
					},
				),
			},
			args: args{
				r: &RequestNotify{
					Message:              "",
					ImageThumbnail:       pStr("https://example.com/foo.jpg"),
					ImageFullsize:        pStr("https://example.com/bar.jpg"),
					NotificationDisabled: true,
				},
			},
		},
		{
			name: "invalid_access_token",
			fields: fields{
				accessToken: "INVALID_ACCESS_TOKEN",
				httpClient: NewTestClient(
					func(req *http.Request) *http.Response {
						header := make(http.Header)
						header["Content-Type"] = []string{"application/json"}
						return &http.Response{
							StatusCode: 401,
							Body:       ioutil.NopCloser(bytes.NewBufferString(`{"status":401,"message":"Invalid access token"}`)),
							Header:     header,
						}
					},
				),
			},
			args: args{
				r: &RequestNotify{
					Message: "This is test.",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := New(
				&Config{
					AccessToken: tt.fields.accessToken,
					HttpClient:  tt.fields.httpClient,
				},
			)
			if err := c.Notify(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Client.Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
