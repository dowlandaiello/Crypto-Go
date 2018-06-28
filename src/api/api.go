package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/mitsukomegumi/Crypto-Go/src/common"
)

// AvailableRequestTypes - defined set of available re
var AvailableRequestTypes = []string{"GET", "POST", "DELETE"}

// RequestElement - struct defining characteristics of specific requested element
type RequestElement struct {
	ElementName         string `json:"elementname"`         // Name of element (e.g. 'BTC-USD')
	BaseElementLocation string `json:"BaseElementLocation"` // Link to request (e.g. '/trade')
	ElementRequestType  string `json:"requesttype"`         // Type of request (e.g. 'post')
	ElementContents     string `json:"requestdata"`         // Contents of request
}

// Handle - attempt to serve specified data
func (request RequestElement) Handle(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, request.ElementContents)
}

// AttemptToServeRequests - attempts to handle incoming requests via data provided in request
func (request RequestElement) AttemptToServeRequests() error {
	fmt.Println("atttempting to serve requests")
	if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.GET(fullPath, request.Handle)

		fasthttp.ListenAndServe(":8080", router.Handler)

		return nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[1])) {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.POST(fullPath, request.Handle)

		fasthttp.ListenAndServe(":8080", router.Handler)

		return nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[2])) {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.DELETE(fullPath, request.Handle)

		fasthttp.ListenAndServe(":8080", router.Handler)

		return nil
	}

	return errors.New("invalid request")
}

// NewRequestServer - checks values of request, returns requestelement
func NewRequestServer(name string, location string, requestType string, requestContents interface{}) (RequestElement, error) {
	tempRequest := RequestElement{}
	if name != "" && common.StringInSlice(requestType, AvailableRequestTypes) {
		json, err := json.Marshal(requestContents)

		if err != nil {
			return tempRequest, err
		}

		request := RequestElement{ElementName: name, BaseElementLocation: location, ElementRequestType: requestType, ElementContents: string(json)}
		return request, nil
	}

	return tempRequest, errors.New("invalid request")
}
