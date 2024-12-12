package extractor

import (
	"bytes"
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
