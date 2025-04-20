package http

import (
	"context"
	"net/http"
	"time"

	"github.com/mercor/payment-service/constants"
	"github.com/go-resty/resty/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Client struct {
	clientSvc   string
	restyClient *resty.Client
}

// Option represents the client options
type Option func(*Client) error

// NewHTTPClient returns a new client with given options
func NewHTTPClient(clientService string, baseUrl string, transport *http.Transport, options ...Option) (*Client, error) {
	restyClient := resty.New()
	restyClient.SetTransport(newrelic.NewRoundTripper(transport)).SetBaseURL(baseUrl).SetTimeout(constants.DefaultTimeout)

	// Client with default Config
	client := &Client{
		clientSvc:   clientService,
		restyClient: restyClient,
	}

	for _, option := range options {
		err := option(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

// WithTimeout option to set request Timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.restyClient.SetTimeout(timeout)
		return nil
	}
}

// Get executes a HTTP GET request.
func (c *Client) Get(request *Request, result interface{}) (*resty.Response, error) {
	return c.Execute(context.Background(), APIGet, request, result)
}

// Post executes a HTTP POST request.
func (c *Client) Post(request *Request, result interface{}) (*resty.Response, error) {
	return c.Execute(context.Background(), APIPost, request, result)
}

// Put executes a HTTP PUT request.
func (c *Client) Put(request *Request, result interface{}) (*resty.Response, error) {
	return c.Execute(context.Background(), APIPut, request, result)
}

// Patch executes a HTTP PATCH request.
func (c *Client) Patch(request *Request, result interface{}) (*resty.Response, error) {
	return c.Execute(context.Background(), APIPatch, request, result)
}

// Delete executes a HTTP DELETE request.
func (c *Client) Delete(request *Request, result interface{}) (*resty.Response, error) {
	return c.Execute(context.Background(), APIDelete, request, result)
}

func (c *Client) Execute(ctx context.Context, method APIMethod, request *Request, result interface{}) (*resty.Response, error) {
	req := c.constructRequest(request, result)

	if request.Timeout != 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, request.Timeout)
		defer cancel()
		req.SetContext(timeoutCtx)
	} else {
		req.SetContext(ctx)
	}

	response, apiErr := req.Execute(method.String(), request.Url)
	if apiErr != nil {
		return response, apiErr
	}

	return response, nil
}

// ConstructRequest creates a new request.
func (c *Client) constructRequest(request *Request, result interface{}) *resty.Request {

	// Setting default headers
	headers := request.Headers
	if headers == nil {
		headers = make(map[string][]string, 0)
	}
	headers["x-client-service"] = []string{c.clientSvc}

	req := c.restyClient.R().
		SetBody(request.Body).
		SetPathParams(request.PathParams).
		SetQueryParamsFromValues(request.QueryParams).
		SetHeaderMultiValues(headers)

	if result != nil {
		req.SetResult(result)
		req.SetError(result)
	}

	return req
}
