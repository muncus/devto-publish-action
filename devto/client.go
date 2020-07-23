package devto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var defaultBaseURL = "https://dev.to/api"

// Client handles calling the dev.to api and parsing responses into the appropriate structs.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// NewClient creates a new client, given a dev.to api key.
func NewClient(apikey string) *Client {
	c := &Client{
		apiKey:  apikey,
		baseURL: defaultBaseURL,
	}
	c.http = http.DefaultClient
	return c
}

// do: wraps http.Client.Do() to include api key.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("api-key", c.apiKey)
	req.Header.Set("content-type", "application/json")
	return c.http.Do(req)
}

// UpsertArticle updates an existing article. Specified article struct must contain an Id.
func (c *Client) UpsertArticle(art *Article, body io.Reader) (*Article, error) {
	// If a body is not specified as io.Reader, read it from the struct.
	if body == nil {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(art)
		if err != nil {
			return nil, fmt.Errorf("Error encoding request object: %s", err)
		}
		body = buf
	}
	var req *http.Request
	var err error
	if art.ID > 0 {
		req, err = http.NewRequest("PUT", fmt.Sprintf("%s/articles/%d", c.baseURL, art.ID), body)
	} else {
		req, err = http.NewRequest("POST", fmt.Sprintf("%s/articles", c.baseURL), body)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to create Request object: %s", err)
	}

	// dump, err := httputil.DumpRequest(req, true)
	// fmt.Println(string(dump))

	r, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %s", err)
	}
	if r.StatusCode < 200 || r.StatusCode > 300 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("Error parsing api response: %#v", err)
		}
		return nil, fmt.Errorf("dev.to api replied with %s: %s", r.Status, body)
	}
	defer r.Body.Close()
	// dump, err = httputil.DumpResponse(r, true)
	// fmt.Println(string(dump))

	newArticle := &Article{}
	err = json.NewDecoder(r.Body).Decode(newArticle)
	if err != nil {
		return nil, fmt.Errorf("Error parsing response json: %s", err)
	}
	return newArticle, nil
}
