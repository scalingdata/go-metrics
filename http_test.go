package metrics

import (
  "io/ioutil"
  "encoding/json"
  "net/http"
  "testing"
)

func TestHttpServer(t *testing.T) {
  addr := "localhost:43563"
  go HttpJson(DefaultRegistry, addr)

  /* Add some metrics */
  c := NewCounter()
  Register("foo.bar", c)
  c.Inc(47)
  c =NewCounter()
  Register("foo.baz", c)
  c.Inc(74)

  timer := NewTimer()
  Register("baz.time", timer)
  timer.Update(15)

  c = NewCounter()
  Register("somecounter", c)
  c.Inc(22)

  resp, err := http.Get("http://" + addr)
  if nil != err {
    t.Fatal(err)
  }
  if nil == resp {
    t.Fatal("No response!")
  }
  if "application/json" != resp.Header.Get("content-type") {
    t.Fatal("Unexpected content-type: ", resp.Header.Get("content-type"))
  }
  
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if nil != err {
    t.Fatal(err)
  }
  if nil == body {
    t.Fatal("No body returned")
  }
  decodedData := make(map[string]interface{})
  err = json.Unmarshal(body, &decodedData)
  if nil != err {
    t.Fatal(err)
  }

  if _, ok := decodedData["foo.bar"]; !ok {
    t.Fatal("Missing expected key foo.bar")
  }
  if _, ok := decodedData["baz.time"]; !ok {
    t.Fatal("Missing expected key baz.time")
  }
  if _, ok := decodedData["somecounter"]; !ok {
    t.Fatal("Missing expected key somecounter")
  }
}
