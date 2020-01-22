package lino

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type Config struct {
	AccessToken string
	HttpClient  *http.Client
}

func New(c *Config) *Client {
	client := &Client{
		accessToken: c.AccessToken,
		httpClient:  c.HttpClient,
	}
	if c.HttpClient == nil {
		client.httpClient = http.DefaultClient
	}
	return client
}

type Client struct {
	accessToken string
	httpClient  *http.Client
}

type RequestNotify struct {
	Message              string
	ImageThumbnail       *string
	ImageFullsize        *string
	StickerPackageID     *int
	StickerID            *int
	NotificationDisabled bool
}

type responseNotify struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	TargetType string `json:"targetType"`
	Target     string `json:"target"`
}

func (c *Client) Notify(r *RequestNotify) error {
	form := url.Values{}
	form.Add("message", r.Message)

	if r.ImageThumbnail != nil {
		form.Add("imageThumbnail", *r.ImageThumbnail)
	}
	if r.ImageFullsize != nil {
		form.Add("imageFullsize", *r.ImageFullsize)
	}
	if r.StickerPackageID != nil {
		form.Add("stickerPackageId", strconv.Itoa(*r.StickerPackageID))
	}
	if r.StickerID != nil {
		form.Add("stickerId", strconv.Itoa(*r.StickerID))
	}
	if r.NotificationDisabled {
		form.Add("notificationDisabled", "true")
	}

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res responseNotify
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	if res.Status != 200 {
		return xerrors.New(res.Message)
	}
	return nil
}
