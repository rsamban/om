package fakes

import "net/http"

type Client struct {
	DoCall struct {
		Receives struct {
			Request *http.Request
		}
		Returns struct {
			Response *http.Response
			Error    error
		}
	}
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	c.DoCall.Receives.Request = request

	return c.DoCall.Returns.Response, c.DoCall.Returns.Error
}