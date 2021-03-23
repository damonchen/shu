package shu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type APIFieldType int

const (
	Bool   = 1
	Int    = 2
	Number = 4
	String = 8
	Bytes  = 16
	Array  = 32
	Map    = 64
)

type FieldPosition int

const (
	Body  = 0
	Path  = 1
	Query = 2
)

type APIParam struct {
	FieldName string
	FieldType APIFieldType
	Name      string
	Position  FieldPosition
	// when array and map exits
	Children []APIParam
}

type APIParamValue struct {
	APIParam
	// APIFieldType will interpret
	value interface{}
}

// API api
type API struct {
	// Name name of api
	Name string
	// Path api path
	Path string
	// Method api method
	Method Method
	// Params params description
	Params interface{}
	// Response response description
	Response interface{}
	// Config when config is nil, use client config
	Config *Config
	// Auth for api, when auth is nil, use client auth
	Auth Auth

	apiParams []APIParam
}

func getApiFieldType(t reflect.Type) APIFieldType {
	switch t.Kind() {
	case reflect.Bool:
		return Bool
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Uint:
		return Int
	case reflect.Float32, reflect.Float64:
		return Number
	case reflect.String:
		return String
	case reflect.Array:
		return Array
	case reflect.Map:
		return Map
	case reflect.Struct:
		return Map
	default:
		return 0
	}
}

func nameFormat(name string) string {
	titleName := strings.ToTitle(name)
	return strings.ToLower(titleName[:1]) + titleName[1:]
}

func getFieldPosition(position string) FieldPosition {
	var pos FieldPosition
	switch position {
	case "body":
		pos = Body
	case "path":
		pos = Path
	case "query":
		pos = Query
	default:
		// TODO: which line position
		panic(fmt.Sprintf("not support position %s", position))
	}
	return pos
}

func getNameAndPosition(param string) (name string, pos FieldPosition, err error) {
	pos = Body

	paramDefs := strings.Split(param, ";")
	for _, p := range paramDefs {
		p = strings.TrimSpace(p)
		kv := strings.Split(p, ":")
		kvLen := len(kv)

		if kvLen == 0 {
			name = ""
		} else if kvLen == 1 {
			name = strings.TrimSpace(kv[0])
		} else if kvLen == 2 {
			if strings.TrimSpace(kv[0]) == "name" {
				name = strings.TrimSpace(kv[1])
			} else if strings.TrimSpace(kv[0]) == "pos" {
				pos = getFieldPosition(strings.TrimSpace(kv[1]))
			}
		} else {
			return "", 0, errors.New("not support definition")
		}

	}
	return
}

func generateAPIParams(t reflect.Type) []APIParam {
	fields := t.NumField()
	apiParams := make([]APIParam, fields)

	for i := 0; i < fields; i++ {
		field := t.Field(i)

		var pos FieldPosition
		var name string
		var err error

		tag := field.Tag
		param, ok := tag.Lookup("param")
		if !ok {
			name = nameFormat(field.Name)
			pos = Body
		} else {
			name, pos, err = getNameAndPosition(param)
			if err != nil {
				return nil
			}

			if name == "" {
				name = nameFormat(field.Name)
			}
		}

		var children []APIParam
		if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Array {
			children = generateAPIParams(field.Type)
		}

		apiParam := APIParam{
			FieldName: field.Name,
			FieldType: getApiFieldType(field.Type),
			Name:      name,
			Position:  pos,
			Children:  children,
		}
		apiParams = append(apiParams, apiParam)
	}
	return apiParams
}

// api init
func (api *API) Init() {
	t := reflect.TypeOf(api.Params)
	if t.Kind() != reflect.Struct {
		panic("must define param using struct")
	}

	api.apiParams = generateAPIParams(t)
}

// APIOption api option
//type APIOption struct {
//	// Name name
//	Name string
//	// Params params
//	Params interface{}
//}

type restAPI struct {
	// API api
	API    *API
	config *Config
}

var (
	re = regexp.MustCompile(`\$\{\w+\}`)
)

func (ra restAPI) getServer() string {
	return ra.config.Server[0]
}

func paramValueToString(param APIParamValue) string {
	if param.FieldType == String {
		return param.value.(string)
	}

	if param.FieldType == Int {
		return strconv.Itoa(param.value.(int))
	}

	if param.FieldType == Bool {
		v := param.value.(bool)
		if v {
			return "true"
		} else {
			return "false"
		}
	}

	return ""
}

func (ra restAPI) getPath(paramValues []APIParamValue) (string, error) {
	path := ra.API.Path

	pathUrlMap := map[string]APIParamValue{}
	for _, param := range paramValues {
		if param.Position == Path {
			pathUrlMap[param.Name] = param
		}
	}

	//var err error
	result := re.ReplaceAllStringFunc(path, func(s string) string {
		param := pathUrlMap[s]
		return paramValueToString(param)
		//err = errors.New("not support")
		//return s
	})
	return result, nil
}

func (ra restAPI) getQuery(paramValues []APIParamValue) string {
	query := url.Values{}
	for _, param := range paramValues {
		if param.FieldType == Query {
			value := paramValueToString(param)
			query.Set(param.Name, value)
		}
	}
	return query.Encode()
}

func (ra restAPI) generateBody(params interface{}, apiParams []APIParam) (map[string]interface{}, error) {
	t := reflect.TypeOf(params)
	fields := t.NumField()

	// TODO: place it in init
	nameParams := make(map[string]APIParam)
	for _, param := range apiParams {
		nameParams[param.FieldName] = param
	}

	v := reflect.ValueOf(params)
	r := map[string]interface{}{}
	for i := 0; i < fields; i++ {
		field := t.Field(i)
		fieldName := field.Name
		param := nameParams[fieldName]

		valueField := v.Field(i)
		valueKind := valueField.Kind()

		var value interface{}
		var err error
		// TODO: slice, pointer

		switch valueKind {
		case reflect.String:
			value = valueField.String()
		case reflect.Int:
			value = valueField.Int()
		case reflect.Struct:
			fallthrough
		case reflect.Array:
			vv := valueField.Interface()
			value, err = ra.generateBody(vv, param.Children)
			if err != nil {
				return nil, err
			}
		}

		r[param.Name] = value
	}
	return r, nil
}

// ParamMerge param merge
func (ra restAPI) paramMerge(params interface{}, apiParams []APIParam) []APIParamValue {
	t := reflect.TypeOf(params)
	fields := t.NumField()

	// TODO: place it in init
	nameParams := make(map[string]APIParam)
	for _, param := range apiParams {
		nameParams[param.FieldName] = param
	}

	v := reflect.ValueOf(params)
	var apiParamValues []APIParamValue
	for i := 0; i < fields; i++ {
		field := t.Field(i)
		fieldName := field.Name
		param := nameParams[fieldName]

		valueField := v.Field(i)
		valueKind := valueField.Kind()

		var value interface{}
		// TODO: slice, pointer

		switch valueKind {
		case reflect.String:
			value = valueField.String()
		case reflect.Int:
			value = valueField.Int()
		case reflect.Struct:
			vv := valueField.Interface()
			value = ra.paramMerge(vv, param.Children)
		case reflect.Array:
			vv := valueField.Interface()
			value = ra.paramMerge(vv, param.Children)
		}

		apiParamValue := APIParamValue{
			APIParam: param,
			value:    value,
		}

		apiParamValues = append(apiParamValues, apiParamValue)
	}
	return apiParamValues
}

func (ra restAPI) getUrlAddr(path string) string {
	server := ra.getServer()
	if strings.HasSuffix(server, "/") {
		server = strings.TrimRight(server, "/")
	}

	return server + path
}

func (ra restAPI) getResponseObject() interface{} {
	t := reflect.TypeOf(ra.API.Response)
	v := reflect.New(t)
	return v.Interface()
}

// Call call
func (ra restAPI) Call(params interface{}) (interface{}, error) {
	// merge the params call
	paramValues := ra.paramMerge(params, ra.API.apiParams)
	path, err := ra.getPath(paramValues)
	if err != nil {
		return nil, err
	}

	query := ra.getQuery(paramValues)
	if len(query) > 0 {
		path = path + "?" + query
	}

	urlAddr := ra.getUrlAddr(path)

	var transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(ra.config.HTTP.ConnectTime) * time.Second,
			KeepAlive: time.Duration(ra.config.HTTP.KeepAlive) * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: time.Second * time.Duration(ra.config.HTTP.TLSHandshakeTimeout),
	}

	client := http.Client{
		Transport: transport,
		Timeout:   time.Duration(ra.config.HTTP.Timeout) * time.Second,
	}

	var req *http.Request
	var resp *http.Response
	switch ra.API.Method {
	case GET:
		req, err = http.NewRequest("GET", urlAddr, nil)
	case POST:
		body, err := ra.generateBody(params, ra.API.apiParams)
		if err != nil {
			return nil, err
		}

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", urlAddr, bytes.NewBuffer(bodyBytes))
	case PUT:
		body, err := ra.generateBody(params, ra.API.apiParams)
		if err != nil {
			return nil, err
		}

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("PUT", urlAddr, bytes.NewBuffer(bodyBytes))
	case DELETE:
		req, err = http.NewRequest("DELETE", urlAddr, nil)
	case OPTIONS:
		req, err = http.NewRequest("OPTIONS", urlAddr, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// set the req auth information
	if ra.API.Auth != nil {
		err = ra.API.Auth.Auth(req)
		if err != nil {
			return nil, err
		}
	}

	userAgent := ra.config.HTTP.UserAgent
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respObj := ra.getResponseObject()
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return nil, err
	}

	return respObj, nil
}
