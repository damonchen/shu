package main

import (
	"fmt"
	"github.com/damonchen/shu"
	"github.com/damonchen/shu/example/client"
	"log"
)

func main() {
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

	fmt.Printf("%#v", resp)
}
