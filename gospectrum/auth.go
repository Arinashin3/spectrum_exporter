package gospectrum

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"spectrum_exporter/gospectrum/api"
	"time"
)

func (_c *SpectrumClient) Login() error {
	// 기존에 토큰이 존재할 경우, 연결 시도
	if _c.token != "" {
		req, err := _c.newRequest(api.SpectrumAPILsSystem, nil)
		if err != nil {
			return err
		}
		resp, err := _c.client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 0 {
			// 접근 실패
			_c.success = false
			return nil
		} else if resp.StatusCode != 403 {
			// 로그인 성공
			return nil
		}
		_c.token = ""
		_c.success = false
	}

	// 토큰 갱신 요청
	_c.tryLoginAt = time.Now()
	req, err := _c.newRequest(api.SpectrumAPIAuth, nil)
	if err != nil {
		return err
	}
	resp, err := _c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return errors.New("invalid username or password")
	}

	var respToken struct {
		Token string `json:"token"`
	}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &respToken)
	if err != nil {
		return err
	}
	_c.token = respToken.Token
	_c.success = true

	return nil
}
