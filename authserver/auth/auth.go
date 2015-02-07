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
	Host   string
	Port   string
	Client *http.Client
}

func NewAuthClient(host string, port string, timeoutMS int) (ac *AuthClient) {
	t := time.Millisecond * time.Duration(timeoutMS)
	c := &http.Client{Timeout: t}
	ac = &AuthClient{Host: host, Port: port, Client: c}
	return
}

func (ac *AuthClient) Get(uuid string) (name string, err error) {
	log.Info("auth: Get called.")
    params := map[string]string{"cookie": uuid}
    name , err = ac.request("get", params)
    log.Debug("auth: Get complete.")
    return
}

func (ac *AuthClient) Set(uuid string, name string) (err error) {
	log.Info("auth: Set called.")
    params := map[string]string{"cookie": uuid, "name": name}
    _ , err = ac.request("set", params)
    log.Debug("auth: Set complete.")
    return
}

func (ac *AuthClient) request(path string, params map[string]string) (contents string, err error){
	log.Debug("auth: request called.")

    uri := url.URL{Scheme: AUTH_SCHEME, Host: ac.Host + ac.Port, Path: path}
    values := url.Values{}
    for k, v := range params {
        values.Add(k, v)
    }
    uri.RawQuery = values.Encode()

    log.Debug("auth: URI string - " + uri.String())
    resp, getErr := ac.Client.Get(uri.String())
    if getErr != nil {
    	log.Error(err)
        err = getErr
        return
    }

    // Call to close response body will cause
    // panic unless error on call to client.Get
    // is non-nil. Calling here, after error checking
    // ensures response is valid.
    defer resp.Body.Close()
    body, readErr := ioutil.ReadAll(resp.Body)
    if readErr != nil {
    	log.Error(err)
        err = readErr
        return
    }
    contents = string(body)
    log.Debug("auth: request complete.")
    return
}
