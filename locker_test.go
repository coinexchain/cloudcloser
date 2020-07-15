package cloudlocker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

var testPath = "./tmp"
var testUrl = ":33300"
var testUrl4client = "http://127.0.0.1" + testUrl

func TestHandlers(t *testing.T) {
	s, err := NewLockerServer(testPath, testUrl)
	if err != nil {
		t.Fatal(err)
	}
	go s.Start()
	e := entry{
		K: string([]byte{0x01}),
		V: string([]byte{0x02}),
	}
	b, _ := json.Marshal(e)
	_, err = http.Post(testUrl4client+"/set", "application/json", strings.NewReader(string(b)))
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(testUrl4client+"/get", "application/json", strings.NewReader(string(e.K)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	out, _ := ioutil.ReadAll(resp.Body)
	if len(out) == 0 || out[0] != e.V[0] {
		panic("value not exactly set")
	}
	s.Stop()
	_, err = http.Post(testUrl4client+"/set", "application/json", strings.NewReader(string(b)))
	if err == nil {
		panic("server should has been stop")
	}
	_ = os.RemoveAll(testPath)
}