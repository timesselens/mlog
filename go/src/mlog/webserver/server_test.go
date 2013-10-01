package webserver

import (
    "bufio"
    "bytes"
    "testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"mlog/core"
)

func Test_processorHandler_1(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(processorHandler));
    ts2 := httptest.NewServer(http.HandlerFunc(pipeHandler));
    ts3 := httptest.NewServer(http.HandlerFunc(processorDeleteHandler));
    defer ts.Close()
    defer ts2.Close()
    defer ts3.Close()

    fmt.Println(ts.URL)
    res, err := http.Get(ts.URL)
    if err != nil {
        t.Error(err)
    }
    res.Body.Close()

    http.Post(ts.URL,"application/json",bufio.NewReader(bytes.NewBufferString(`{"type":"stdin"}`)))
    http.Post(ts.URL,"application/json",bufio.NewReader(bytes.NewBufferString(`{"type":"grep", "re": "foo"}`)))
    http.Post(ts.URL,"application/json",bufio.NewReader(bytes.NewBufferString(`{"type":"debug", "prefix": "test-"}`)))
    http.Post(ts2.URL,"application/json",bufio.NewReader(bytes.NewBufferString(`{"i": 0, "j": 1}`)))
    http.Post(ts2.URL,"application/json",bufio.NewReader(bytes.NewBufferString(`{"i": 1, "j": 2}`)))
    http.Post(ts3.URL + "?i=0","application/x-www-form-urlencoded",nil)
    http.Post(ts3.URL + "?i=0","application/x-www-form-urlencoded",nil)
    http.Post(ts3.URL + "?i=0","application/x-www-form-urlencoded",nil)



    if len(core.Processors.List) > 0 {
        fmt.Printf("++++ %#v\n",core.Processors.List)
        t.Error("still procs left")
    }

    // <-time.After(1*time.Minute)
}
