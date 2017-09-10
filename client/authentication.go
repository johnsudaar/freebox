package client

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	AuthResponsePending = "pending"
	AuthResponseGranted = "granted"
)

type AppTokenParams struct {
	AppID      string `json:"app_id"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
}

type AuthorizationResponse struct {
	AppToken string `json:"app_token"`
	TrackID  int    `json:"track_id"`
}

type AuthorizationStatusResponse struct {
	Status    string `json:"status"`
	Challenge string `json:"challenge"`
}

func (c *freeboxClient) RequestAppToken(appName string) (string, chan error, error) { // APPToken, RegistrationDone, error
	deviceName, err := os.Hostname()
	if err != nil {
		deviceName = "freebox-go-client"
	}
	params := AppTokenParams{
		AppID:      c.AppID,
		AppName:    appName,
		AppVersion: c.AppVersion,
		DeviceName: deviceName,
	}

	resp, err := c.Post("login/authorize/", params, RequestOpts{ExpandPath: true, CheckResult: true})
	if err != nil {
		return "", nil, errors.Wrap(err, "fail to make authentication request")
	}

	var authResp AuthorizationResponse
	err = json.Unmarshal(resp, &authResp)
	if err != nil {
		return "", nil, errors.Wrap(err, "invalid response from server")
	}

	responseChan := make(chan error)

	go c.checkAuthorizationState(authResp.TrackID, responseChan)

	return authResp.AppToken, responseChan, nil
}

func (c *freeboxClient) checkAuthorizationState(trackID int, responseChan chan error) {
	path := fmt.Sprintf("login/authorize/%d", trackID)
	for {
		resp, err := c.Get(path, RequestOpts{ExpandPath: true, CheckResult: true})
		if err != nil {
			responseChan <- errors.Wrap(err, "fail to make check request")
			return
		}

		var status AuthorizationStatusResponse
		err = json.Unmarshal(resp, &status)
		if err != nil {
			responseChan <- errors.Wrap(err, "invalid response from server")
			return
		}

		if status.Status == AuthResponseGranted {
			responseChan <- nil
			return
		}

		if status.Status != AuthResponsePending {
			responseChan <- errors.Wrapf(err, "invalid response state:%s", status.Status)
			return
		}
		time.Sleep(2 * time.Second)
	}
}
