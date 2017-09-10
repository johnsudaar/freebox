package client

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Version struct {
	UID            string `json:"uid"`
	DeviceName     string `json:"device_name"`
	ApiVersion     string `json:"api_version"`
	ApiBaseUrl     string `json:"api_base_url"`
	DeviceType     string `json:"device_type"`
	ApiDomain      string `json:"api_domain"`
	HttpsAvailable bool   `json:"https_available"`
	HttpsPort      int    `json:"https_port"`
}

func (c *freeboxClient) Version() (Version, error) {
	var res Version
	body, err := c.Get("api_version", RequestOpts{})
	if err != nil {
		return res, errors.Wrap(err, "fail to fetch version")
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, errors.Wrap(err, "invalid response from server")
	}

	return res, nil
}
