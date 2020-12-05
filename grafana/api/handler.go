package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Handler interface {
	GetNamespaces() (*NamespacesRes, error)
	GetObjects(namespaceID string) (*ObjectsRes, error)
	GetDevicesStateStream(deviceID string) (<-chan *DevicesStateRes, error)
}

type handlerImpl struct {
	client       *http.Client
	tokenHandler TokenHandler
}

func NewHandler(tokenHandler TokenHandler) Handler {
	return &handlerImpl{
		client:       &http.Client{},
		tokenHandler: tokenHandler,
	}
}

func (h *handlerImpl) authRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+h.tokenHandler.GetToken())
	return h.client.Do(req)
}

type NamespacesRes struct {
	Namespaces []Namespace `json:"namespaces"`
}

type Namespace struct {
	DeleteInitiationTime string `json:"deleteinitiationtime"`
	ID                   string `json:"id"`
	MarkForDeletion      bool   `json:"markfordeletion"`
	Name                 string `json:"name"`
}

func (h *handlerImpl) GetNamespaces() (*NamespacesRes, error) {
	res, err := h.authRequest(context.Background(), http.MethodGet, endpointNamespaces, nil)
	if err != nil {
		return nil, NewFailedHTTPRequestError(http.MethodGet, endpointNamespaces, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, NewUnexpectedHTTPStatusCodeError(http.MethodGet, endpointNamespaces, res.StatusCode)
	}
	b, _ := ioutil.ReadAll(res.Body)
	tmp := &NamespacesRes{}
	_ = json.Unmarshal(b, tmp)
	return tmp, nil
}

type ObjectsRes struct {
	Objects []Object `json:"objects"`
}

type Object struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func (h *handlerImpl) GetObjects(namespaceID string) (*ObjectsRes, error) {
	url := fmt.Sprintf("%s?namespace=%s&recurse=1", endpointObjects, namespaceID)
	res, err := h.authRequest(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, NewFailedHTTPRequestError(http.MethodGet, endpointObjects, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, NewUnexpectedHTTPStatusCodeError(http.MethodGet, url, res.StatusCode)
	}
	b, _ := ioutil.ReadAll(res.Body)
	tmp := &ObjectsRes{}
	_ = json.Unmarshal(b, tmp)
	return tmp, nil
}

type DevicesStateRes struct {
	Result DevicesStateResult `json:"result"`
}

type DevicesStateResult struct {
	ReportedState DeviceState `json:"reportedState"`
}

type DeviceState struct {
	Version   string                 `json:"version"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func (h *handlerImpl) GetDevicesStateStream(deviceID string) (<-chan *DevicesStateRes, error) {
	url := fmt.Sprintf(endpointDevicesStateStream, deviceID) + "?only_delta=false"
	res, err := h.authRequest(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, NewFailedHTTPRequestError(http.MethodGet, url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, NewUnexpectedHTTPStatusCodeError(http.MethodGet, url, res.StatusCode)
	}

	rdr := bufio.NewReader(res.Body)
	buf := make([]byte, 4*1024)
	ret := make(chan *DevicesStateRes)
	go func() {
		for {
			l, err := rdr.Read(buf)
			if l > 0 {
				b := buf[:l]
				tmp := &DevicesStateRes{}
				_ = json.Unmarshal(b, tmp)
				ret <- tmp
			}
			if err != nil {
				close(ret)
				res.Body.Close()
				break
			}
		}
	}()
	return ret, nil
}
