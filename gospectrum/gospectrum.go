package gospectrum

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
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

func (_c *SpectrumClient) login() error {
	defer func() {
		_c.tryLoginAt = time.Now()
	}()
	// 토큰 존재 여부 확인
	if _c.token != "" {
		return nil
	}
	// 실패했을 시, 최소 1분 후 로그인 재시도
	if time.Since(_c.tryLoginAt) < 1*time.Minute {
		return nil
	}
	// 요청 생성
	req, err := api.SpectrumAPIAuth.NewRequest(_c.endpoint, nil)
	if err != nil {
		return err
	}

	// 헤더 추가
	req.Header.Add("X-Auth-Username", _c.username)
	req.Header.Add("X-Auth-Password", _c.password)
	_c.client.Jar, _ = cookiejar.New(nil)

	// 전송
	body, err := _c.send(req)
	if err != nil {
		return err
	}

	// 토큰 갱신
	var loginResp struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return err
	}
	_c.token = loginResp.Token

	return nil
}
func (_c *SpectrumClient) send(req *http.Request) ([]byte, error) {
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
	if body == nil {
		return nil, errors.New(resp.Status)
	}
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
