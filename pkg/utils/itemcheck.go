package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CheckIAItemExist(identifier string) (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://archive.org/services/check_identifier.php", nil)
	req.Header.Set("User-Agent", GetUA())
	params := req.URL.Query()
	params.Add("identifier", identifier)
	params.Add("output", "json")
	req.URL.RawQuery = params.Encode()
	respRaw, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer respRaw.Body.Close()

	type Response struct {
		Type string `json:"type"`
		Code string `json:"code"`
	}
	resp := Response{}
	err = json.NewDecoder(respRaw.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	if resp.Type == "success" {
		if resp.Code == "available" {
			return false, nil
		} else if resp.Code == "not_available" {
			return true, nil
		}
	}
	return false, fmt.Errorf("unexpected response: %+v", resp)

}
