package extractor

import (
	"os"
	"testing"

	"github.com/tidwall/gjson"
)

func Test_html2json(t *testing.T) {
	raw, err := os.ReadFile("test_html2json.html")
	if err != nil {
		t.Fatal(err)
		return
	}
	json, err := Html2json(raw)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(json)

	err = os.WriteFile("test_html2json.json", []byte(json), 0644)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func Test_empty_json(t *testing.T) {
	// read from out.json file
	if gjson.Valid("") {
		t.Fatal("empty json is not valid")
	}
}

func Test_GetKsPlayJson(t *testing.T) {
	wantKsPlayJSON, err := os.ReadFile("test_ksPlayJson.json")
	raw, err := os.ReadFile("test_html2json.html")
	if err != nil {
		t.Fatal(err)
		return
	}
	json, err := Html2json(raw)
	if err != nil {
		t.Fatal(err)
		return
	}
	ksPlayJson, err := GetKsPlayJson(json)
	if err != nil {
		t.Fatal(err)
		return
	}
	if string(wantKsPlayJSON) != ksPlayJson {
		t.Fatal("ksPlayJson not match")
	}
}

func Test_GetM3u8s(t *testing.T) {
	raw, err := os.ReadFile("test_html2json.html")
	if err != nil {
		t.Fatal(err)
		return
	}
	json, err := Html2json(raw)
	if err != nil {
		t.Fatal(err)
		return
	}
	ksPlayJson, err := GetKsPlayJson(json)
	if err != nil {
		t.Fatal(err)
		return
	}
	m3u8s, err := GetM3U8s(ksPlayJson)
	if err != nil {
		t.Fatal(err)
		return
	}
	if "720P" != m3u8s[0].QualityLabel {
		t.Fatal("m3u8s not match")
	}
}
