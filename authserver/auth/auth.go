package auth

import (
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	AUTH_SCHEME = "http"
)

type AuthClient struct {
	host   string
	port   string
	client *http.Client
}

func NewAuthClient(host string, port string, timeoutMS time.Duration) (ac *AuthClient) {
	t := timeoutMS
	c := &http.Client{Timeout: t}
	ac = &AuthClient{host: host, port: port, client: c}
	return
}

func (ac *AuthClient) Get(uuid string) (name string, err error) {
	log.Trace("auth: Get called.")
    params := map[string]string{"cookie": uuid}
    name , err = ac.request("get", params)
    log.Trace("auth: Get complete.")
    return
}

func (ac *AuthClient) Set(uuid string, name string) (err error) {
	log.Trace("auth: Set called.")
    params := map[string]string{"cookie": uuid, "name": name}
    _ , err = ac.request("set", params)
    log.Trace("auth: Set complete.")
    return
}

func (ac *AuthClient) request(path string, params map[string]string) (contents string, err error){
	log.Trace("auth: Request called.")

    var resp *http.Response
    var body []byte

    uri := url.URL{Scheme: AUTH_SCHEME, Host: ac.host + ac.port, Path: path}
    values := url.Values{}
    for k, v := range params {
        values.Add(k, v)
    }
    uri.RawQuery = values.Encode()

    log.Debug("auth: Requesting URI - " + uri.String())
    if resp, err = ac.client.Get(uri.String()) ; err != nil {
        return
    }

    // Call to close response body will cause
    // panic unless error on call to client.Get
    // is non-nil. Calling here, after error checking
    // ensures response is valid.
    defer resp.Body.Close()
    if body, err = ioutil.ReadAll(resp.Body) ; err != nil {
        return
    }
    contents = string(body)
    log.Trace("auth: Request complete.")
    return
}
