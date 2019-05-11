package iexcloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type IEX struct {
	client *http.Client
	apiKey string
}

func NewIEX(client *http.Client, apiKey string) *IEX {
	return &IEX{
		client: client,
		apiKey: apiKey,
	}
}

func (iex *IEX) get(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, iex.fullURL(endpoint), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := iex.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("received unexpected status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}
