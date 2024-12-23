package extractor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
)

func extractJsonFromHtml(_html string) (string, error) {
	// json in html
	parser := html.NewTokenizer(strings.NewReader(_html))
	var tag_found = false
	for {
		tokenType := parser.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := parser.Token()
		if token.Type == html.StartTagToken {
			if token.Data == "script" {
				for _, attr := range token.Attr {
					if attr.Key == "class" && attr.Val == "videoInfo" {
						tag_found = true
						break
					}
				}
			}
		}
		if tag_found && token.Type == html.TextToken {
			return token.Data, nil
		}
		if tag_found && token.Type == html.EndTagToken {
			tag_found = false
		}
	}

	return "", errors.New("json not found")
}

func Html2json(body []byte) (string, error) {
	// trim the json body
	// {...}/*<!-- fetch-stream -->*/
	var body_trimed = body
	if !bytes.HasSuffix(body, []byte("}")) {
		body_trimed = body[:bytes.LastIndex(body, []byte("}"))+1]
	}

	// html in json
	html := gjson.GetBytes(body_trimed, "html").String()

	// json in html
	json_block, err := extractJsonFromHtml(html)
	if err != nil {
		return "", err
	}

	// trim the json
	json := json_block[strings.Index(json_block, "{") : strings.LastIndex(json_block, "}")+1]

	if !gjson.Valid(json) {
		return "", fmt.Errorf("invalid json: %s", json)
	}

	return json, nil
}

func GetTitles(part string) (dougaTitle string, partTitle string, err error) {
	// douga title
	result := gjson.Get(part, "title")
	if !result.Exists() {
		return "", "", fmt.Errorf("title not found")
	}
	if result.Type != gjson.String {
		return "", "", fmt.Errorf("title is not a string")
	}
	dougaTitle = result.String()

	// part title
	result = gjson.Get(part, "currentVideoInfo.title")
	if !result.Exists() {
		return "", "", fmt.Errorf("title not found")
	}
	if result.Type != gjson.String {
		return "", "", fmt.Errorf("title is not a string")
	}
	partTitle = result.String()

	return dougaTitle, partTitle, nil
}

type Tag struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type CurrentVideoInfo struct {
	Id        string `json:"id"` // video id (类似 B 站的 cid)
	PartTitle string `json:"title"`
}

type PartInfo struct {
	CoverUrl         string           `json:"coverUrl"` // 封面
	DougaId          string           `json:"dougaId"`
	DougaTitle       string           `json:"title"`
	CurrentVideoInfo CurrentVideoInfo `json:"currentVideoInfo"`
	Description      string           `json:"description"`
	User             Tag              `json:"user"` // 共用 Tag，反正字段一致
	TagList          []Tag            `json:"tagList"`
	CreateTime       string           `json:"createTime"`
}

func GetPartInfo(part string) (PartInfo, error) {
	var info PartInfo
	err := json.Unmarshal([]byte(part), &info)
	if err != nil {
		return info, err
	}
	return info, nil
}

// part: video info json
func GetKsPlayJson(part string) (string, error) {
	result := gjson.Get(part, "currentVideoInfo.ksPlayJson")

	if !result.Exists() {
		return "", fmt.Errorf("ksPlayJson not found")
	}
	if result.Type != gjson.String {
		return "", fmt.Errorf("ksPlayJson is not a string")
	}

	ksPlayJson := result.String()
	if !gjson.Valid(ksPlayJson) {
		return "", fmt.Errorf("ksPlayJson is not a valid json")
	}

	return ksPlayJson, nil
}

type M3u8 struct {
	Url          string
	QualityLabel string
}

func GetM3U8s(ksPlayJson string) (m3u8s []M3u8, err error) {
	urls, err := GetM3u8Urls(ksPlayJson)
	if err != nil {
		return
	}
	qulityLabels, err := GetM3u8QualityLabels(ksPlayJson)
	if err != nil {
		return
	}
	for i, url := range urls {
		m3u8s = append(m3u8s, M3u8{Url: url, QualityLabel: qulityLabels[i]})
	}
	return
}

func GetM3u8Urls(ksPlayJson string) ([]string, error) {
	result := gjson.Get(ksPlayJson, "adaptationSet.0.representation.#.url")
	if !result.Exists() {
		return nil, fmt.Errorf("url not found")
	}
	if result.Type != gjson.JSON {
		return nil, fmt.Errorf("url is not a json")
	}

	urls := make([]string, 0)
	for _, url := range result.Array() {
		urls = append(urls, url.String())
	}

	return urls, nil
}

func GetM3u8QualityLabels(ksPlayJson string) ([]string, error) {
	result := gjson.Get(ksPlayJson, "adaptationSet.0.representation.#.qualityLabel")
	if !result.Exists() {
		return nil, fmt.Errorf("url not found")
	}
	if result.Type != gjson.JSON {
		return nil, fmt.Errorf("url is not a json")
	}

	qualityLabels := make([]string, 0)
	for _, url := range result.Array() {
		qualityLabels = append(qualityLabels, url.String())
	}

	return qualityLabels, nil
}
