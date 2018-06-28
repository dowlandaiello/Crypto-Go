package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/mgo.v2"

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

	Dynamics string `json:"dynamics"` // Dynamic data
}

// Handle - attempt to serve specified data
func (request RequestElement) Handle(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, request.ElementContents)
}

// HandleVar - handle request, with dynamics
func (request RequestElement) HandleVar(ctx *fasthttp.RequestCtx) {
	db := mgo.Database{}

	json.Unmarshal([]byte(request.ElementContents), db)

	fmt.Fprintf(ctx, request.ElementName)

	val, err := findValue(&db, strings.SplitAfter(strings.SplitAfter(request.BaseElementLocation, "/")[1], "/")[1], request.ElementName, ctx.UserValue(request.ElementName))
	//fmt.Fprint(ctx, findValue(request.ElementContents, request.BaseElementLocation, ctx.UserValue(request.Dynamics), ctx.UserValue(request.Dynamics)))
}

// AttemptToServeRequests - attempts to handle incoming requests via data provided in request
func (request RequestElement) AttemptToServeRequests() error {
	fmt.Println("atttempting to serve requests")
	if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics == "" {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.GET(fullPath, request.Handle)

		err := fasthttp.ListenAndServe(":8080", router.Handler)

		if err != nil {
			return err
		}

		return nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics != "" {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.GET(fullPath, request.HandleVar)

		err := fasthttp.ListenAndServe(":8080", router.Handler)

		if err != nil {
			return err
		}

		return nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[1])) {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.POST(fullPath, func(ctx *fasthttp.RequestCtx) {})

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
func NewRequestServer(name string, location string, requestType string, requestContents interface{}, dynamics string) (RequestElement, error) {
	tempRequest := RequestElement{}
	if name != "" && common.StringInSlice(requestType, AvailableRequestTypes) {
		if requestContents == "" {
			json, err := json.Marshal(requestContents)

			if err != nil {
				return tempRequest, err
			}

			if dynamics == "" {
				request := RequestElement{ElementName: name, BaseElementLocation: location, ElementRequestType: requestType, ElementContents: string(json)}

				return request, nil
			}

			request := RequestElement{ElementName: name, BaseElementLocation: location, ElementRequestType: requestType, ElementContents: string(json), Dynamics: dynamics}

			return request, nil
		}
		request := RequestElement{ElementName: name, BaseElementLocation: location, ElementRequestType: requestType, Dynamics: dynamics}
		return request, nil
	}

	return tempRequest, errors.New("invalid request")
}
