# shu

shu is a project for http client request for restful api

# usage

```bash
go get github.com/damonchen/shu/cmd/shugen
```

```go
package client

import "github.com/damonchen/shu"

//go:generate shugen generate.go -o client.go

var (
	// using bundle.SetOption(client.WithAuth(auth), client.WithConfig(config)) to set
	// the auth and config outside
	bundle = shu.NewBundle()
)

type UserLogin struct {
	Name     string `param:"name;pos:query"`
	Password string `param:"password"`
}

type UserLoginResp struct {
	Status string `json:"status"`
}

func init() {
	bundle.Client(shu.Client{
		Name: "auth",
		APIs: []*shu.API{
			{
				Name:     "login",
				Path:     "/api/v1/login",
				Method:   shu.POST,
				Params:   UserLogin{},
				Response: UserLoginResp{},
			},
		},
	},
	)
	bundle.Init()
}

func GetBundle() *shu.Bundle {
	return bundle
}
```

when run go generate, it will generate the following code:

```go
package client

type AuthClient struct {
}

func (c AuthClient) Login(userLogin UserLogin) (*UserLoginResp, error) {
	resp, err := bundle.Get("auth").Call("login", userLogin)
	if err != nil {
		return nil, err
	}
	return resp.(*UserLoginResp), nil
}

func GetAuthClient() AuthClient {
	return AuthClient{}
}

```


if we use the client as like below, using: 
```go

import (
    "xxxx/client"
)

config := shu.Config{
    VerifyCert:  false,
    Server:      []string{"http://localhost:5000"},
    HTTP:        shu.HTTP{},
    Trace:       false,
    TraceFile:   "",
    Marshaler:   nil,
    Unmarshaler: nil,
}
client.GetBundle().SetOption(
    shu.WithConfig(&config))

authClient := client.GetAuthClient()
userLogin := client.UserLogin{
    Name:     "damon",
    Password: "chen",
}
resp, err := authClient.Login(userLogin)
if err != nil {
    log.Fatalf("login response error %s", err)
    return
}

userLogin := client.UserLogin{}
resp, err := client.GetAuthClient().Login(userLogin)
if err != nil {
    ... 
}
...

```

## TODO

- [ ] base running 
- [ ] generate tool
- [ ] rest api parser and geneartor
- [ ] http client pool for performance
- [ ] refect performance

