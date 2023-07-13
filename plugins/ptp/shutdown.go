package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"
	"sync/atomic"

	"github.com/openrelayxyz/plugeth-utils/core"
)

var globalId int64

var client = &http.Client{Transport: &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConnsPerHost:   16,
	MaxIdleConns:          16,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}}

type Call struct {
	Version string            `json:"jsonrpc"`
	ID      json.RawMessage   `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

func toRawMessages(items ...interface{}) ([]json.RawMessage, error) {
	result := make([]json.RawMessage, len(items))
	for i, item := range items {
	  d, err := json.Marshal(item)
	  if err != nil { return nil, err }
	  result[i] = (json.RawMessage)(d)
	}
	return result, nil
}

func PreTrieCommit(node core.Hash) {

	id, err := toRawMessages(atomic.AddInt64(&globalId, 1))
  	if err != nil {
		log.Error("json marshalling error, id", "err", err)
	}

	call := &Call{
		Version: "2.0",
		ID : id[0],
		Method: "plugeth_capturePreTrieCommit",
		Params: []json.RawMessage{},
	  }

	backendURL := "http://127.0.0.1:9546"

	callBytes, _ := json.Marshal(call)

	request, _ := http.NewRequestWithContext(context.Background(), "POST", backendURL, bytes.NewReader(callBytes))
	request.Header.Add("Content-Type", "application/json")

	_, err = client.Do(request)

	if err != nil {
		log.Error("Error calling passive node from PreTrieCommit", "err", err)
	}

}

func PostTrieCommit(node core.Hash) {

	id, err := toRawMessages(atomic.AddInt64(&globalId, 1))
  	if err != nil {
		log.Error("json marshalling error, id", "err", err)
	}

	call := &Call{
		Version: "2.0",
		ID : id[0],
		Method: "plugeth_capturePostTrieCommit",
		Params: []json.RawMessage{},
	  }

	backendURL := "http://127.0.0.1:9546"

	callBytes, _ := json.Marshal(call)

	request, _ := http.NewRequestWithContext(context.Background(), "POST", backendURL, bytes.NewReader(callBytes))
	request.Header.Add("Content-Type", "application/json")

	_, err = client.Do(request)

	if err != nil {
		log.Error("Error calling passive node from PostTrieCommit", "err", err)
	}

}

func OnShutdown() {

	id, err := toRawMessages(atomic.AddInt64(&globalId, 1))
  	if err != nil {
		log.Error("json marshalling error, id", "err", err)
	}

	call := &Call{
		Version: "2.0",
		ID : id[0],
		Method: "plugeth_captureShutdown",
		Params: []json.RawMessage{},
	  }

	backendURL := "http://127.0.0.1:9546"

	callBytes, _ := json.Marshal(call)

	request, _ := http.NewRequestWithContext(context.Background(), "POST", backendURL, bytes.NewReader(callBytes))
	request.Header.Add("Content-Type", "application/json")

	_, err = client.Do(request)

	if err != nil {
		log.Error("Error calling passive node from OnShutdown", "err", err)
	}

}