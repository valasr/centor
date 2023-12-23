package grpc_server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mrtdeh/centor/proto"
)

type APIResponse struct {
	Body       []byte
	Error      string
	StatusCode int
}

func (s *agent) CallAPI(ctx context.Context, req *proto.APIRequest) (*proto.APIResponse, error) {
	res, err := makeInternalRequest(req.Method, req.Addr, req.Body)
	if err != nil {
		return &proto.APIResponse{
			Error: err.Error(),
		}, err
	}
	return &proto.APIResponse{
		Body:   string(res.Body),
		Error:  res.Error,
		Status: int32(res.StatusCode),
	}, nil
}

func makeInternalRequest(method, addr, bodyStr string) (*APIResponse, error) {
	var res map[string]interface{}

	// make [POST] request to local service
	var err error
	var body []byte
	var stat int

	switch strings.ToLower(method) {
	case "post":

		body, stat, err = makePostRequest(addr, []byte(bodyStr))
		if err != nil {
			log.Printf("error in make post request in client : %s\n", err.Error())
		}
	case "get":
		body, stat, err = makeGetRequest(addr)
		if err != nil {
			log.Printf("error in make get request in client : %s\n", err.Error())
		}
	}

	fmt.Printf("Make request to localhost route=%v addr=%v method=%v\n", res["route"], res["addr"], res["method"])
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return &APIResponse{
		Body:       body,
		Error:      errStr,
		StatusCode: stat,
	}, nil
}

func makePostRequest(url string, data []byte) ([]byte, int, error) {

	client := &http.Client{}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, 0, fmt.Errorf("error create request : %s", err.Error())
	}
	r.Header.Add("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		return nil, 0, fmt.Errorf("error in request : %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusCreated {
		return nil, res.StatusCode, fmt.Errorf("error with status : %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error in read body : %s", err.Error())
	}

	return body, res.StatusCode, nil
}

func makeGetRequest(url string) ([]byte, int, error) {

	client := &http.Client{}
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("error create request : %s", err.Error())
	}

	res, err := client.Do(r)
	if err != nil {
		return nil, 0, fmt.Errorf("error in request : %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusCreated {
		return nil, res.StatusCode, fmt.Errorf("error with status : %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error in read body : %s", err.Error())
	}

	// get := &map[string]interface{}{}
	// derr := json.NewDecoder(res.Body).Decode(get)
	// if derr != nil {
	// 	return nil, 0, fmt.Errorf("error in decode : %s", derr.Error())

	// }

	return body, res.StatusCode, nil
}
