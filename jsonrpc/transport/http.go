package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/umbracle/ethgo/jsonrpc/codec"
)

// HTTP is an http transport
type HTTP struct {
	addr    string
	client *http.Client
	headers map[string]string
}

func newHTTP(addr string, headers map[string]string) *HTTP {
	return &HTTP{
		addr:    addr,
		client: &http.Client{},
		headers: headers,
	}
}

// Close implements the transport interface
func (h *HTTP) Close() error {
	return nil
}

// Call implements the transport interface
func (h *HTTP) Call(method string, out interface{}, params ...interface{}) error {
	// Encode json-rpc request
	request := codec.Request{
		JsonRPC: "2.0",
		Method:  method,
	}
	if len(params) > 0 {
		data, err := json.Marshal(params)
		if err != nil {
			return err
		}
		request.Params = data
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", h.addr, bytes.NewReader(raw))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}

	res, err := h.client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Decode json-rpc response
	var response codec.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error
	}

	if err := json.Unmarshal(response.Result, out); err != nil {
		return err
	}
	return nil
}

// SetMaxConnsPerHost sets the maximum number of connections that can be established with a host
func (h *HTTP) SetMaxConnsPerHost(count int) {
}
