package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		return err
	}

	res, err := http.Post(ServicesURL, "application/json", buf) // 发送一个post请求进行注册
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to registry service. Registry service "+
			"responded with code %v", res.StatusCode)
	}
	return nil
}

func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesURL,
		bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregister service. Registry "+
			"service responded with code %v", res.StatusCode)
	}
	return nil
}
