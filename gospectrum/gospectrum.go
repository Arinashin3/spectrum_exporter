package gospectrum

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"spectrum_exporter/gospectrum/api"
	"time"
)

type SpectrumClient struct {
	endpoint   string
	username   string
	password   string
	token      string
	tryLoginAt time.Time
	client     *http.Client
	success    bool
}

func NewClient(endpoint string, username string, password string, insecure bool) *SpectrumClient {
	return &SpectrumClient{
		endpoint: endpoint,
		username: username,
		password: password,
		token:    "",
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
		},
	}
}

func (_c *SpectrumClient) send(req *http.Request) ([]byte, error) {
	if _c.token == "" {
		return nil, errors.New("no token provided")
	}
	req.Header.Add("X-Auth-Token", _c.token)
	resp, err := _c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = _c.checkHttpCode(resp.StatusCode)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, nil

}

func (_c *SpectrumClient) checkHttpCode(code int) error {
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
	}
	return nil
}

func (_c *SpectrumClient) HealthCheck() bool {
	return _c.success
}

func (_c *SpectrumClient) newRequest(path api.SpectrumAPIPath, data []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", _c.endpoint+string(path), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.URL, err = req.URL.Parse(_c.endpoint + string(path))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if path == api.SpectrumAPIAuth {
		req.Header.Add("X-Auth-Username", _c.username)
		req.Header.Add("X-Auth-Password", _c.password)
	} else {
		req.Header.Add("X-Auth-Token", _c.token)
	}
	return req, nil
}
