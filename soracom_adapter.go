package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type SoracomAdapter struct {
	credential *SoracomCredential
	token      *SoracomToken
	target     string
}

func (soracom *SoracomAdapter) GetSoracomToken() error {
	buf, err := soracom.RequestHttp("POST", "https://api.soracom.io/v1/auth", soracom.credential)
	if err != nil {
		return err
	}
	var token = &SoracomToken{}
	err = json.Unmarshal(buf, token)
	if err != nil {
		return errors.New("fail to parse token")
	}
	soracom.token = token
	return nil
}

func (soracom *SoracomAdapter) StartNapter(request *CreatePortMappingRequest) (*PortMapping, error) {
	buf, err := soracom.RequestHttp("POST", "https://api.soracom.io/v1/port_mappings", request)
	if err != nil {
		return nil, err
	}
	response := &PortMapping{}
	err = json.Unmarshal(buf, response)
	if err != nil {
		return nil, errors.New("fail to parse create port mapping")
	}
	return response, nil
}

func (soracom *SoracomAdapter) StopNapter(request *PortMapping) error {
	_, err := soracom.RequestHttp("DELETE", "https://api.soracom.io/v1/port_mappings/"+request.IpAddress+"/"+strconv.Itoa(request.Port), nil)
	if err != nil {
		return err
	}
	return nil
}

func (soracom *SoracomAdapter) RequestHttp(method, url string, data interface{}) ([]byte, error) {
	var req *http.Request
	var err error
	if data != nil {
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			return nil, errors.New("fail to serialize data")
		}
		body := bytes.NewReader(payloadBytes)
		req, err = http.NewRequest(method, url, body)
		if err != nil {
			return nil, errors.New("fail to create http request")
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, errors.New("fail to create http request")
		}
	}
	req.Header.Set("Accept", "application/json")
	if soracom.token != nil {
		req.Header.Set("X-Soracom-Api-Key", soracom.token.ApiKey)
		req.Header.Set("X-Soracom-Token", soracom.token.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("fail to access via http")
	}
	defer resp.Body.Close()

	buf := make([]byte, 65536)
	readLen, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return nil, errors.New("fail to read http body")
	}
	if resp.StatusCode < 300 {
		return buf[:readLen], nil
	} else {
		return nil, errors.New("fail to request API\n" + string(buf[:readLen]))
	}
}
