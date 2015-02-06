package auth

import (
	log "github.com/cihub/seelog"
	"net/http"
)

const (
	AUTH_SCHEME = "http"
    HOST_PORT_SEPARATOR = ":"
)

type AuthClient struct {
	Host   string
	Port   string
	Client *http.Client
}

func NewAuthClient(host string, port string, timeoutMS int) (ac *AuthClient, err error) {
	t := time.Millisecond * time.Duration(timeoutMS)
	c := &http.Client{Timeout: t}
	ac = &AuthClient{Host: host, Port: port, Client: c}
	return
}

func (ac *AuthClient) Get(uuid string) (name string, err error) {
	log.Debug("Auth.Get called.")

    params := map[string]string{"cookie": uuid}
    name , err = ac.request("get", params)
    return
}

func (ac *AuthClient) Set(uuid string, name string) (err error) {
	log.Debug("Auth.Set called.")

    params := map[string]string{"cookie": uuid, "name": name}
    _ , err = ac.request("set", params)
    return
}

func (ac *AuthClient) request(path string, params map[string]string) (contents string, err error){
	log.Debug("Auth.request called.")

    uri := url.URL{Scheme: AUTH_SCHEME, Host: ac.Host + HOST_PORT_SEPARATOR + ac.Port, Path: path}
    values := url.Values{}
    for k, v := range params {
        values.Add(k, v)
    }
    uri.RawQuery = values.Encode()

    resp, getErr := ac.Client.Get(uri.String())
    defer resp.Body.Close()
    if getErr != nil {
        err = getErr
        return
    }

    body, readErr := ioutil.ReadAll(resp.Body)
    if readErr != nil {
        err = readErr
        return
    }

    contents = string(body)
    return
}
