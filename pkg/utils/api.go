package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type RequestLog struct {
	Method       string `json:"method"`
	Host         string `json:"host"`
	Path         string `json:"path"`
	RequestBody  string `json:"request_body"`
	ResponseBody string `json:"response_body"`
	TimeLength   int    `json:"time_length"`
}

type Header struct {
	Key   string
	Value string
}

type Response struct {
	Body       []byte
	Duration   time.Duration
	StatusCode int
}

var defHTTPTimeout = 5 * time.Second
var httpClient = &http.Client{Timeout: defHTTPTimeout}

func DoRequest(method, path string, headers []Header, body interface{}, v interface{}) (*Response, error) {
	response := Response{}
	var byteJson []byte

	if body != nil {
		byteJson, _ = json.Marshal(body)
	}

	reqBody := bytes.NewBuffer(byteJson)

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		log.Println("Request creation failed: ", err)
		return nil, err
	}

	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}

	start := time.Now()
	log.Println("Request ", req.Method, ": ", req.URL.Host, req.URL.Path)

	RequestLog := RequestLog{
		Method:      req.Method,
		Host:        req.URL.Host,
		Path:        req.URL.Path,
		RequestBody: Stringify(body),
	}
	log.Println(Stringify(RequestLog))

	res, err := httpClient.Do(req)
	executionTime := time.Since(start)
	response.Duration = executionTime

	if err != nil {
		log.Println("Cannot send request : ", err)
		return &response, err
	}

	log.Println("Completed in ", executionTime)
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	response.Body = resBody
	response.StatusCode = res.StatusCode

	if err != nil {
		log.Println("Cannot read response body: ", err)
		return &response, err
	}

	if res.StatusCode != http.StatusOK {
		errMessage := fmt.Sprintf("Cannot send request and get status code %d from %s body %v", res.StatusCode, req.URL, string(resBody))
		return &response, errors.New(errMessage)
	}

	if v != nil {
		err := json.Unmarshal(resBody, v)
		RequestLog.TimeLength = int(time.Since(start))
		log.Println(Stringify(RequestLog))

		return &response, err
	}

	return &response, nil
}

func Stringify(m interface{}) string {
	json, _ := json.Marshal(m)
	return string(json)
}
