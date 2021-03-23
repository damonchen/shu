package shu

import (
	"fmt"
)

// Client client api definition
type Client struct {
	// Name of client
	Name string
	// APIs api array
	APIs     []*API
	apiIndex map[string]*API
	// client config, when config is nil, will use global config
	Config *Config
	// Auth when is nil, use global Auth
	Auth Auth
}

// Init init
func (c *Client) Init() {
	if c.apiIndex == nil {
		c.apiIndex = map[string]*API{}
	}
	for _, api := range c.APIs {
		api.Init()
		c.apiIndex[api.Name] = api
	}
}

func (c Client) getAPI(name string) *API {
	return c.apiIndex[name]
}

// Bundle client bundle
type Bundle struct {
	// Clients client definitions
	Clients []*Client
	// client index for clients
	clientIndex map[string]*Client
	// client config, must set it
	Config *Config
	// Auth when is nil, no auth
	Auth Auth
}

// NewBundle new bundle
func NewBundle(opts ...BundleOption) *Bundle {
	bundle := &Bundle{
		clientIndex: map[string]*Client{},
	}
	bundle.SetOption(opts...)
	return bundle
}

// SetOption set option for client bundle
func (b *Bundle) SetOption(opts ...BundleOption) *Bundle {
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// BundleOption bundle option
type BundleOption func(client *Bundle)

// WithAuth with auth
func WithAuth(auth Auth) BundleOption {
	return func(client *Bundle) {
		client.Auth = auth
	}
}

// WithConfig with config
func WithConfig(config *Config) BundleOption {
	return func(client *Bundle) {
		client.Config = config
	}
}

// Init init
func (b *Bundle) Init() {
	for _, client := range b.Clients {
		client.Init()
		b.clientIndex[client.Name] = client
	}
}

// Client client
func (b *Bundle) Client(clients ...Client) *Bundle {
	for _, client := range clients {
		if _, ok := b.clientIndex[client.Name]; ok {
			panic(fmt.Sprintf("client %s has already exists, duplicated", client.Name))
		}
		b.clientIndex[client.Name] = &client

		b.Clients = append(b.Clients, &client)
	}

	return b
}

// Get get
func (b *Bundle) Get(name string) *RestClient {
	if client, ok := b.clientIndex[name]; ok {
		return &RestClient{
			client: client,
			config: b.Config,
		}
	}
	return nil
}

// RestClient rest client impl
type RestClient struct {
	client *Client
	config *Config
}

func (rc RestClient) getRestAPI(name string) restAPI {
	api := rc.client.getAPI(name)

	var config *Config
	if rc.client.Config != nil {
		config = rc.client.Config
	} else {
		config = rc.config
	}

	return restAPI{
		API:    api,
		config: config,
	}
}

// Call call
func (rc RestClient) Call(name string, params interface{}) (interface{}, error) {
	api := rc.getRestAPI(name)
	return api.Call(params)
}
