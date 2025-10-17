package gospectrum

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"spectrum_exporter/gospectrum/api"
	"strings"
	"time"
)

type Client struct {
	endpoint   string
	username   string
	password   string
	token      string
	tryLoginAt time.Time
	client     *http.Client
	isLogin    bool
}

func NewTransport(insecure bool) *http.Transport {
	return &http.Transport{
		MaxConnsPerHost: 1,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return tls.Dial(network, addr, &tls.Config{
				InsecureSkipVerify: insecure,
			})
		},
	}
}

func NewClient(endpoint string, username string, password string, tr *http.Transport) (*Client, error) {
	var required []string
	if endpoint == "" {
		required = append(required, "endpoint")
	}
	if username == "" {
		required = append(required, "username")
	}
	if password == "" {
		required = append(required, "password")
	}
	if tr == nil {
		required = append(required, "transport")
	}
	if len(required) > 0 {
		return nil, errors.New("parameters are not exists: [" + strings.Join(required, ",") + "]")
	}

	return &Client{
		endpoint: endpoint,
		username: username,
		password: password,
		token:    "",
		client:   &http.Client{Transport: tr},
	}, nil
}

type Response struct {
	StatusCode int
	Body       []byte
}

func (_c *Client) send(req *http.Request) (*Response, error) {
	var tmp *Response
	// Request Data
	resp, err := _c.client.Do(req)

	// Output Data
	var body []byte
	if resp != nil {
		body, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		// Check Status Code
		err = _c.checkHttpCode(resp.StatusCode)
		tmp = &Response{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	return tmp, err
}

func (_c *Client) checkHttpCode(code int) error {
	switch code {
	case http.StatusForbidden:
		_c.token = ""
		return errors.New("where the command did not send a valid authentication token to interact with the specified URL")
	case http.StatusNotFound:
		return errors.New("where the command tried to issue a request to a URL that does not exist")
	case http.StatusMethodNotAllowed:
		return errors.New("where the command tried to use an HTTP method that is invalid for the specified URL")
	case http.StatusInternalServerError:
		return errors.New("where a Spectrum Virtualize command error is forwarded from the RESTful API")
	default:
		return nil
	}
}

func (_c *Client) newRequest(path string, data []byte) (*http.Request, error) {
	var req *http.Request
	var err error

	if req, err = http.NewRequest("POST", _c.endpoint+string(path), bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	// Set Headers...
	req.Header.Add("Content-Type", "application/json")
	if path == api.SpectrumCommandAuth.String("") {
		req.Header.Add("X-Auth-Username", _c.username)
		req.Header.Add("X-Auth-Password", _c.password)
		_c.client.Jar, err = cookiejar.New(nil)
	} else if _c.token != "" {
		req.Header.Add("X-Auth-Token", _c.token)
	} else {
		return nil, errors.New("no authentication provided")
	}

	return req, nil
}

func (_c *Client) Endpoint() string {
	return _c.endpoint
}

func (_c *Client) IsLogin() bool {
	return _c.isLogin
}
