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

type Logout struct{}

type UserRegister struct {
	Name     string `param:"name;pos:query"`
	Password string `param:"password"`
}
type UserRegisterResp struct {
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
			{
				Name:     "logout",
				Path:     "/api/v1/logout",
				Method:   shu.POST,
				Params:   Logout{},
				Response: UserLoginResp{},
			},
		},
	},
		shu.Client{
			Name: "user",
			APIs: []*shu.API{
				{
					Name:     "register",
					Path:     "/api/v1/register",
					Method:   shu.POST,
					Params:   UserRegister{},
					Response: UserRegisterResp{},
				},
			},
		},
	)
	bundle.Init()
}

func GetBundle() *shu.Bundle {
	return bundle
}
