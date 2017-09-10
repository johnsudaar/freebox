package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type freeboxClient struct {
	BaseUrl         string
	AppID           string
	AppVersion      string
	ApiBasePath     string
	ApiMajorVersion string
}

func NewFreeboxClient(baseUrl, appID, appVersion string) (Client, error) {
	client := &freeboxClient{
		BaseUrl:    baseUrl,
		AppID:      appID,
		AppVersion: appVersion,
	}

	versions, err := client.Version()
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize client")
	}

	// "X.Y" => ["X", "Y"] => "X"
	// TODO: Check that this is an version string
	client.ApiMajorVersion = strings.Split(versions.ApiVersion, ".")[0]
	client.ApiBasePath = versions.ApiBaseUrl

	return client, nil
}

func (c *freeboxClient) fullUrl(path string, opts RequestOpts) string {
	var fullUrl string

	if opts.ExpandPath {
		fullUrl = fmt.Sprintf("http://%s%sv%s/%s", c.BaseUrl, c.ApiBasePath, c.ApiMajorVersion, path)
	} else {
		fullUrl = fmt.Sprintf("http://%s/%s", c.BaseUrl, path)
	}

	fmt.Println(fullUrl)

	return fullUrl
}

func (c *freeboxClient) Get(path string, opts RequestOpts) ([]byte, error) {
	resp, err := http.Get(c.fullUrl(path, opts))
	if err != nil {
		return nil, errors.Wrapf(err, "fail to make freebox get request to %s", path)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to fetch response data")
	}

	if opts.CheckResult {
		res, err = checkResult(res)
		if err != nil {
			return nil, errors.Wrap(err, "errors while checking response")
		}
	}
	return res, nil
}

func (c *freeboxClient) Post(path string, data interface{}, opts RequestOpts) ([]byte, error) {
	params, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "fail to marshal json")
	}

	resp, err := http.Post(c.fullUrl(path, opts), "application/json", bytes.NewBuffer(params))
	if err != nil {
		return nil, errors.Wrapf(err, "fail to make freebox get request to %s", path)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to fetch response data")
	}

	if opts.CheckResult {
		res, err = checkResult(res)
		if err != nil {
			return nil, errors.Wrap(err, "errors while checking response")
		}
	}

	return res, nil
}

func checkResult(data []byte) ([]byte, error) {
	var response APIResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, errors.Wrap(err, "invalid response from server")
	}

	if !response.Success {
		return nil, errors.New(fmt.Sprintf("request failed, error_code:%s message:%s", response.ErrorCode, response.Message))
	}

	return response.Result.MarshalJSON()
}
