package main

import (
        "fmt"
        "log"
        "net/http"
        "net/http/httputil"
        "net/url"
)

type DebugTransport struct{}

func (DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
        b, err := httputil.DumpRequestOut(r, false)
        if err != nil {
                return nil, err
        }
        fmt.Println(string(b))

        resp, err := http.DefaultTransport.RoundTrip(r)
      	if err != nil {
		      log.Fatal(err)
	      }

        if (resp.StatusCode == 302) {
          location, _ := resp.Location()
          fmt.Printf("New location: %s", location)
          resp, err = http.Get(location.String())
          if err != nil {
            fmt.Printf("Get failed: %s", err.Error())
            return resp, err
          }
          fmt.Printf("Resp: %+v", resp)
        }

        dump, err := httputil.DumpResponse(resp, true)
      	if err != nil {
		      log.Fatal(err)
	      }
        fmt.Println(string(dump))

        return resp, nil
}

func followRedirect(resp *http.Response) (err error) {
  if (resp.StatusCode == 302) {
    location, _ := resp.Location()
    fmt.Printf("New location: %s", location)
    resp, err = http.Get(location.String())
    if err != nil {
      fmt.Printf("Get failed: %s", err.Error())
      return err
    }
    fmt.Printf("Resp: %+v", resp)
  }
  return nil
}

func main() {
        target, _ := url.Parse("http://flibusta.is")
        log.Printf("forwarding to -> %s\n", target)

        proxy := httputil.NewSingleHostReverseProxy(target)
        proxy.Transport = DebugTransport{}

        http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
                req.Host = req.URL.Host

                proxy.ServeHTTP(w, req)
        })

        log.Fatal(http.ListenAndServe(":8080", nil))
}