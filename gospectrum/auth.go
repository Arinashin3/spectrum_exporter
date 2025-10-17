package gospectrum

import (
	"encoding/json"
	"net/http"
	"spectrum_exporter/gospectrum/api"
	"time"
)

type AuthToken struct {
	Token string
}

func (_c *Client) Login() (bool, error) {
	// Create Request
	var err error
	var req *http.Request
	var path = api.SpectrumCommandAuth.String("")
	var updateToken bool

	// Create Token
	if _c.token == "" {
		_c.tryLoginAt = time.Now()
		if req, err = _c.newRequest(path, nil); err != nil {
			return updateToken, err
		}

		// Send & Parsing
		var data *AuthToken
		var resp *Response
		if resp, err = _c.send(req); err == nil {
			err = json.Unmarshal(resp.Body, &data)
		}
		if err != nil {
			return updateToken, err
		}

		// Update Token
		_c.token = data.Token
		updateToken = true
		_c.isLogin = true

	}
	return updateToken, err
}
