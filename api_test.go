package aixifan

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tidwall/gjson"
)

func Test_requestDouga(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	body, err := requestDouga(client, "41")
	if err != nil {
		t.Fatal(err)
		return
	}
	json, err := html2json(body)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(json)
	os.WriteFile("test/outcome/Douga.json", []byte(json), 0644)

	if gjson.Get(json, "title").String() != "Nihilum公会全1-球FD伊利丹视频" {
		t.Fatal("title not match")
	}
}

func Test_getDougaMultiP(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	json, err := GetDouga(client, "32749_1")
	if err != nil {
		t.Fatal(err)
		return
	}
	os.WriteFile("test/outcome/DougaMultiP_1.json", []byte(json), 0644)

	if !strings.Contains(gjson.Get(json, "title").String(), "东方伪装天") {
		t.Fatal("title not match")
	}
	if gjson.Get(json, "currentVideoId").Int() != 659501 {
		t.Fatal("currentVideoId p1 not match")
	}

	json, err = GetDouga(client, "32749_2")
	if err != nil {
		t.Fatal(err)
		return
	}
	os.WriteFile("test/outcome/DougaMultiP_2.json", []byte(json), 0644)
	if gjson.Get(json, "currentVideoId").Int() != 659502 {
		t.Fatal("currentVideoId p2 not match")
		return
	}
}
