package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Mock struct {
	Url string
	HttpMethod string
	Response *http.Response
	Err error
}

var (
	enableMock = false
	mocks = make(map[string]*Mock)
)

func getMockId(httpMethod, url string) string {
	return fmt.Sprintf("%s_%s", httpMethod, url)
}

func AddMockBehavior(mock Mock) {
	mocks[getMockId(mock.HttpMethod, mock.Url)] = &mock
}

func StartMockup() {
	enableMock = true
}

func RestoreMockup() {
	mocks = make(map[string]*Mock)
}

func StopMockup() {
	enableMock = false
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {

	if enableMock {
		mockFound := mocks[getMockId(http.MethodPost, url)]
		if mockFound == nil {
			return nil, errors.New("could not find a valid mock")
		}

		return mockFound.Response, mockFound.Err
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	timeout := 2 * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	request.Header = headers
	return client.Do(request)
}