package valorant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	Version          = "v0.1.0"
	defaultBaseURL   = "https://na.api.riotgames.com/val/"
	defaultUserAgent = "go-valorant" + "/" + Version
)

type Region string

const (
	RegionAP    Region = "ap"
	RegionBR    Region = "br"
	RegionEU    Region = "eu"
	RegionKR    Region = "kr"
	RegionLATAM Region = "latam"
	RegionNA    Region = "na"
)

var errNonNilContext = errors.New("context must be non-nil")

type Client struct {
	client *http.Client

	// Base URL for API requests. Defaults to the NA region.
	BaseURL *url.URL

	// User agent used when communicating with the Valorant API.
	UserAgent string

	common service

	// Services used for talking to different parts of the Valorant API.
	ContentService *ContentService
	MatchService   *MatchService
	RankedService  *RankedService
	StatusService  *StatusService
}

func (c *Client) Client() *http.Client {
	clientCopy := *c.client
	return &clientCopy
}

type service struct {
	client *Client
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient2 := *httpClient
	c := &Client{client: &httpClient2}
	c.initialize()
	return c
}

func (c *Client) WithAuthToken(token string) *Client {
	c.client.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.Header.Set("X-Riot-Token", token)
		return http.DefaultTransport.RoundTrip(req)
	})
	return c
}

func (c *Client) WithRegion(region Region) *Client {
	c.BaseURL, _ = url.Parse("https://" + string(region) + ".api.riotgames.com/val/")
	return c
}

func (c *Client) initialize() {
	if c.client == nil {
		c.client = &http.Client{}
	}
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(defaultBaseURL)
	}
	if c.UserAgent == "" {
		c.UserAgent = defaultUserAgent
	}
	c.common.client = c
	c.ContentService = (*ContentService)(&c.common)
	c.MatchService = (*MatchService)(&c.common)
	c.RankedService = (*RankedService)(&c.common)
	c.StatusService = (*StatusService)(&c.common)
}

// RequestOption represents an option that can modify a http.Request.
type RequestOption func(req *http.Request)

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	decErr := json.NewDecoder(resp.Body).Decode(v)
	if decErr == io.EOF {
		decErr = nil // ignore EOF errors caused by empty response body
	}
	if decErr != nil {
		err = decErr
	}

	return resp, err
}

func (c *Client) BareDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	err = CheckResponse(resp)

	return resp, err
}

type ErrorResponse struct {
	Response *http.Response `json:"-"`
	Status   struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	} `json:"status"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Status.Message)
}

func CheckResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		_ = json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}

// roundTripperFunc creates a RoundTripper (transport)
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
