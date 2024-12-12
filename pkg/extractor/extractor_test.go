package extractor

import (
	"os"
	"testing"

	"github.com/tidwall/gjson"
)

func Test_html2json(t *testing.T) {
	raw, err := os.ReadFile("test_html2json.html")
	if err != nil {
		// print pwd
		t.Fatal(os.Getwd())
		t.Fatal(err)
		return
	}
	_json, err := Html2json(raw)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(_json)

	err = os.WriteFile("test_html2json.json", []byte(_json), 0644)
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
