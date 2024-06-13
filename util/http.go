package util

import (
	"bytes"
	"encoding/json"
	"github.com/JoeyZeYi/source/log"
	"github.com/JoeyZeYi/source/log/zap"
	"io"
	"net/http"
)

func HttpGet[T any](url string) (*T, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error("HttpGet", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("HttpGet", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	result := new(T)
	if err = json.Unmarshal(respBytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func HttpPost[T any](url, contentType string, req any) (*T, error) {
	reqBytes, _ := json.Marshal(req)

	resp, err := http.Post(url, contentType, bytes.NewReader(reqBytes))
	if err != nil {
		log.Error("HttpPost", zap.String("url", url), zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("HttpPost", zap.String("url", url), zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	result := new(T)
	if err = json.Unmarshal(respBytes, result); err != nil {
		return nil, err
	}
	return result, nil

}
