package shu

import "encoding/json"

type HTTP struct {
	// Timeout The timeout includes connection time, any
	// redirects, and reading the response body
	Timeout uint `yaml:"timeout"`
	// ConnectTimeout connect timeout
	ConnectTime uint `yaml:"connectTime"`
	// KeepAlive keep alive
	KeepAlive uint `yaml:"keepAlive"`
	// TLSHandshakeTimeout, unit second
	TLSHandshakeTimeout uint `yaml:"omitempty"`
	// UserAgent user agent
	UserAgent string `yaml:"userAgent"`
}

// Config config
type Config struct {
	// VerifyCert verify certificate when use https
	VerifyCert bool `yaml:"verifyCert,omitempty"`
	// Server servers
	Server []string `yaml:"server"`
	// ConnectTime connect time
	HTTP HTTP `yaml:"http,omitempty"`
	// Trace the http client request
	Trace bool `yaml:"trace,omitempty"`
	// Trace filename,if not exist, use stdout
	TraceFile string `yaml:"traceFile,omitempty"`
	// Marshaler marshaler
	Marshaler json.Marshaler `yaml:"omitempty"`
	// Unmarshaler unmarshaler
	Unmarshaler json.Unmarshaler `yaml:"omitempty"`
}
