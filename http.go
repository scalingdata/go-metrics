package metrics

import (
  "net/http"
)

// Configuration for the HTTP Server
type HttpConfig struct {
  Registry Registry    // Registry to be exported
  Addr string          // Network address to connect to
  Encoder HttpEncoder  // Function to convert Registry to wire format
  ContentType string   // Value for the 'content-type' header specific to your encoding
}

type HttpEncoder func (r Registry, req *http.Request) ([]byte, error)

// Blocking function that starts an HTTP server that responds with
// a JSON encoded copy of metrics in the registry.
func HttpJson(r Registry, addr string) error {
  return HttpFromConfig(HttpConfig{
    Registry: r,
    Addr: addr,
    Encoder: jsonHttpEncoder,
    ContentType: "application/json",
  })
}

// Blocking function that starts an HTTP server that responds with an encoded
// version of all metrics in the registry.
func Http(r Registry, addr string, encoder HttpEncoder, contentType string) error {
  return HttpFromConfig(HttpConfig{
    Registry: r,
    Addr: addr,
    Encoder: encoder,
    ContentType: contentType,
  })
}

// Same as Http() but accepts a HttpConfig instead of individual arguments
func HttpFromConfig(cfg HttpConfig) error {
  http.HandleFunc("/", makeHttpHandler(cfg.Registry, cfg.Encoder, cfg.ContentType))
  return http.ListenAndServe(cfg.Addr, nil)
}

func makeHttpHandler(r Registry, encode HttpEncoder, contentType string) func(http.ResponseWriter, *http.Request) {
  return func (w http.ResponseWriter, req *http.Request) {
    header := w.Header()
    if "" != contentType {
      header.Add("content-type", contentType)
    }
    data, _ := encode(r, req)
    w.Write(data)
  }
}

func jsonHttpEncoder (r Registry, req *http.Request) ([]byte, error) { 
  return MarshalJSON(r) 
}
