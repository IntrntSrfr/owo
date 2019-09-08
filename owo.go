package owo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type Client struct {
	token  string
	client *http.Client
}

func NewClient(tkn string) *Client {
	return &Client{
		token:  tkn,
		client: &http.Client{},
	}
}

func (o *Client) Upload(text string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="files[]"; filename="text.txt"`)
	h.Set("Content-Type", "text/plain;charset=utf-8")

	part, err := writer.CreatePart(h)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, bytes.NewReader([]byte(text)))
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.awau.moe/upload/pomf", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", o.token)

	res, err := o.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	result := Result{}
	err = json.Unmarshal(resbody, &result)
	if err != nil {
		return "", err
	}

	if !result.Success {
		return "", errors.New(result.Description)
	}

	if len(result.Files) > 0 {
		return "https://chito.ge/" + result.Files[0].URL, nil
	}
	return "", nil
}
