package api

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/saveweb/aixifan/pkg/extractor"
)

// dougaId: int-string
func requestDouga(client *http.Client, dougaId string) ([]byte, error) {
	if strings.Contains(dougaId, "ac") {
		return nil, errors.New("dougaId should not contain 'ac'")
	}
	url := "https://www.acfun.cn/v/ac" + dougaId + "?quickViewId=videoInfo_new&ajaxpipe=1"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "aixifanfan/0.0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetDouga(client *http.Client, dougaId string) (string, error) {
	body, err := requestDouga(client, dougaId)
	if err != nil {
		return "", err
	}
	return extractor.Html2json(body)
}
