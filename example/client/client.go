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
