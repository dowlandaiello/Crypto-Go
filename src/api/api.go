package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mitsukomegumi/Crypto-Go/src/pairs"

	"github.com/mitsukomegumi/Crypto-Go/src/orders"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/buaazp/fasthttprouter"
	"github.com/mitsukomegumi/Crypto-Go/src/accounts"
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

	ElementDb *mgo.Database `json:"database"`

	Dynamics string `json:"dynamics"` // Dynamic data
}

// SetupRoutes - setup all necessary routes for operation
func SetupRoutes(database *mgo.Database) error {
	router, err := SetupAccountRoutes(database)

	if err != nil {
		return err
	}

	SetupOrderRoutes(router, database)

	err = fasthttp.ListenAndServe(":8080", router.Handler)

	if err != nil {
		return err
	}

	return nil
}

// Handle - attempt to serve specified data
func (request RequestElement) Handle(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, request.ElementContents)
}

// HandleDel - attempt to delete
func (request RequestElement) HandleDel(ctx *fasthttp.RequestCtx) {
	keys := strings.Split(request.Dynamics, "/:")
	keys = append(keys[:0], keys[0+1:]...)

	values := []string{}

	x := 0

	for x != len(keys) {
		values = append(values, ctx.UserValue(keys[x]).(string))
		x++
	}

	if common.StringInSlice("username", keys) && !common.StringInSlice("pair", keys) {
		acc, err := findAccount(request.ElementDb, values[0])

		if err != nil {
			fmt.Fprintf(ctx, err.Error())
		} else {
			if common.ComparePasswords(acc.PassHash, []byte(values[1])) {
				err := removeAccount(request.ElementDb, acc)

				if err != nil {
					fmt.Fprintf(ctx, err.Error())
				} else {
					fmt.Fprintf(ctx, "removed")
				}
			}
		}
	} else if common.StringInSlice("pair", keys) {
		acc, err := findAccount(request.ElementDb, values[3])

		if err == nil {
			if common.ComparePasswords(acc.PassHash, []byte(values[4])) {
				split := strings.Split(values[0], "-")
				pair := pairs.NewPair(split[0], split[1])
				amount, _ := strconv.ParseFloat(values[2], 64)
				order, _ := orders.NewOrder(acc, values[1], pair, amount)

				err = addOrder(request.ElementDb, &order)

				if err != nil {
					fmt.Fprintf(ctx, err.Error())
				} else {
					json, err := json.MarshalIndent(order, "", "  ")

					if err != nil {
						fmt.Fprintf(ctx, err.Error())
					} else {
						fmt.Fprintf(ctx, string(json[:]))
					}
				}
			}
		}
	}
}

// HandleVar - handle request, with dynamics
func (request RequestElement) HandleVar(ctx *fasthttp.RequestCtx) {
	key := strings.Split(common.TrimLeftChar(request.ElementName), "/:")[0]
	value := ctx.UserValue(key).(string)

	collection := strings.Split(request.BaseElementLocation, "/")[2]

	if strings.Contains(request.BaseElementLocation, "orders") {
		collection = value

		key = strings.Split(common.TrimLeftChar(request.ElementName), "/:")[1]
		value = ctx.UserValue(key).(string)
	}

	val, err := findValue(request.ElementDb, collection, strings.ToLower(key), value)

	if err != nil {
		fmt.Fprint(ctx, err.Error())
	} else {
		json, err := json.MarshalIndent(val, "", "  ")

		if err != nil {
			fmt.Fprint(ctx, err.Error())
		} else {
			fmt.Fprint(ctx, string(json[:]))
		}
	}
}

// HandlePost - handle POST request, with dynamics
func (request RequestElement) HandlePost(ctx *fasthttp.RequestCtx) {
	keys := strings.Split(request.Dynamics, "/:")
	keys = append(keys[:0], keys[0+1:]...)

	values := []string{}

	x := 0

	for x != len(keys) {
		values = append(values, ctx.UserValue(keys[x]).(string))
		x++
	}

	if common.StringInSlice("username", keys) && !common.StringInSlice("pair", keys) {
		acc := accounts.NewAccount(values[0], values[1], values[2])

		err := addAccount(request.ElementDb, &acc)

		if err != nil {
			fmt.Fprintf(ctx, err.Error())
		} else {
			json, err := json.MarshalIndent(acc, "", "  ")

			if err != nil {
				fmt.Fprintf(ctx, err.Error())
			} else {
				fmt.Fprintf(ctx, string(json[:]))
			}
		}
	} else if common.StringInSlice("pair", keys) {
		acc, err := findAccount(request.ElementDb, values[3])

		if err == nil {
			if common.ComparePasswords(acc.PassHash, []byte(values[4])) {
				split := strings.Split(values[0], "-")
				pair := pairs.NewPair(split[0], split[1])
				amount, _ := strconv.ParseFloat(values[2], 64)
				order, _ := orders.NewOrder(acc, values[1], pair, amount)

				err = addOrder(request.ElementDb, &order)

				if err != nil {
					fmt.Fprintf(ctx, err.Error())
				} else {
					json, err := json.MarshalIndent(order, "", "  ")

					if err != nil {
						fmt.Fprintf(ctx, err.Error())
					} else {
						fmt.Fprintf(ctx, string(json[:]))
					}
				}
			}
		}
	}
}

// AttemptToServeRequestsWithRouter - attempts to handle incoming requests via data provided in request
func (request RequestElement) AttemptToServeRequestsWithRouter(router *fasthttprouter.Router) (*fasthttprouter.Router, error) {
	fmt.Println("atttempting to serve requests with handler: " + request.ElementName)

	if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics == "" {
		fullPath := request.BaseElementLocation + "/" + request.ElementName

		router.GET(fullPath, request.Handle)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics != "" {
		fullPath := request.BaseElementLocation + "/" + request.ElementName

		router.GET(fullPath, request.HandleVar)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[1])) {
		fullPath := request.BaseElementLocation + request.Dynamics

		router.POST(fullPath, request.HandlePost)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[2])) {
		fullPath := request.BaseElementLocation + request.Dynamics

		router.DELETE(fullPath, request.HandleDel)

		return router, nil
	}

	return nil, errors.New("invalid request")
}

// AttemptToServeRequests - attempts to handle incoming requests via data provided in request
func (request RequestElement) AttemptToServeRequests() (*fasthttprouter.Router, error) {
	fmt.Println("atttempting to serve requests")
	if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics == "" {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.GET(fullPath, request.Handle)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics != "" {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.GET(fullPath, request.HandleVar)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[1])) {
		fullPath := request.BaseElementLocation + request.Dynamics

		router := fasthttprouter.New()

		router.POST(fullPath, request.HandlePost)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[2])) {
		fullPath := request.BaseElementLocation + "/" + request.ElementName
		router := fasthttprouter.New()

		router.DELETE(fullPath, request.HandleDel)

		return router, nil
	}

	return nil, errors.New("invalid request")
}

// NewRequestServer - checks values of request, returns requestelement
func NewRequestServer(name string, location string, requestType string, requestContents interface{}, db *mgo.Database, dynamics string) (RequestElement, error) {
	tempRequest := RequestElement{}
	if name != "" && common.StringInSlice(requestType, AvailableRequestTypes) {
		if requestContents != "" {
			json, err := json.MarshalIndent(requestContents, "", "  ")

			if err != nil {
				return tempRequest, err
			}

			if dynamics == "" {
				request := RequestElement{ElementName: name, BaseElementLocation: location, ElementRequestType: requestType, ElementContents: string(json)}

				return request, nil
			}

			request := RequestElement{ElementName: name, BaseElementLocation: location, ElementRequestType: requestType, ElementContents: string(json), ElementDb: db, Dynamics: dynamics}

			return request, nil
		}
		request := RequestElement{ElementName: name, BaseElementLocation: location, ElementRequestType: requestType, ElementDb: db, Dynamics: dynamics}
		return request, nil
	}

	return tempRequest, errors.New("invalid request")
}

func findValue(database *mgo.Database, collection string, key string, value string) (interface{}, error) {
	c := database.C(collection)

	result := make(map[string]interface{})

	err := c.Find(bson.M{key: value}).One(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
