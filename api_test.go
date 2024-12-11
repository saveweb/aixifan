package aixifan

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/tidwall/gjson"
)

func Test_requestVideoInfo(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	body, err := requestVideoInfo(client, "46638469")
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
	os.WriteFile("test/outcome/requestVideoInfo.json", []byte(json), 0644)

	if gjson.Get(json, "title").String() != "以军杀向叙利亚首都。别慌，只是想谈笔买卖【岩论469期】" {
		t.Fatal("title not match")
	}
}
