package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dowlandaiello/Crypto-Go/src/database"
	"github.com/dowlandaiello/Crypto-Go/src/market"

	"github.com/dowlandaiello/Crypto-Go/src/pairs"

	"github.com/dowlandaiello/Crypto-Go/src/orders"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/buaazp/fasthttprouter"
	"github.com/dowlandaiello/Crypto-Go/src/accounts"
	"github.com/valyala/fasthttp"

	"github.com/dowlandaiello/Crypto-Go/src/common"
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

	SetupMarketRoutes(router, database)

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
	keys := strings.Split(request.Dynamics, "?")
	keys = append(keys[:0], keys[0+1:]...)

	values := []string{}

	x := 0

	for x != len(keys) {
		peekVal := string(ctx.PostArgs().Peek(keys[x]))

		if peekVal == "" {
			break
		}

		values = append(values, peekVal)

		x++
	}

	if len(values) < 3 {
		values, _ = request.GetUserValues(keys, ctx)
	}

	if common.StringInSlice("username", keys) && !common.StringInSlice("pair", keys) {
		acc, err := findAccount(request.ElementDb, values[0])

		if err != nil {
			fmt.Fprintf(ctx, err.Error())
		} else {
			if common.ComparePasswords(acc.PassHash, []byte(values[1])) {
				err := removeAccount(request.ElementDb, &acc)

				if err != nil {
					fmt.Fprintf(ctx, err.Error())
				} else {
					fmt.Fprintf(ctx, "removed")
				}
			}
		}
	} else if common.StringInSlice("pair", keys) {
		acc, err := findAccount(request.ElementDb, values[2])

		if err == nil {
			if common.ComparePasswords(acc.PassHash, []byte(values[3])) {
				split := strings.Split(values[0], "-")
				pair := pairs.NewPair(split[0], split[1])
				order, err := findOrder(request.ElementDb, values[1], pair)

				if err != nil {
					fmt.Fprintf(ctx, err.Error())
				} else {
					err = removeOrder(request.ElementDb, order)

					if err != nil {
						fmt.Fprintf(ctx, err.Error())
					} else {
						fmt.Fprintf(ctx, "order removed")
					}
				}
			}
		}
	}
}

// HandleVar - handle request, with dynamics
func (request RequestElement) HandleVar(ctx *fasthttp.RequestCtx) {
	key := strings.Split(common.TrimLeftChar(request.ElementName), "?")[0]

	value := request.GetUserValue(key, ctx)

	collection := strings.Split(request.BaseElementLocation, "/")[2]

	if strings.Contains(request.BaseElementLocation, "orders") {
		collection = value

		key = strings.Split(common.TrimLeftChar(request.ElementName), "?")[1]
		value = request.GetUserValue(key, ctx)
	}

	if strings.Contains(collection, "?") {
		collection = strings.Split(collection, "?")[0]
	}

	if strings.Contains(request.Dynamics, "password") && !strings.Contains(request.BaseElementLocation, "orders") {
		passKey := strings.Split(common.TrimLeftChar(request.ElementName), "?")[1]

		passVal := request.GetUserValue(strings.ToLower(passKey), ctx)

		accVal, err := findAccount(request.ElementDb, request.GetUserValue("username", ctx))

		if err != nil {
			fmt.Fprintf(ctx, err.Error())
		} else {
			if common.ComparePasswords(accVal.PassHash, []byte(passVal)) {
				val, err := accounts.DecryptPrivateKeys(accVal.WalletRawHashedKeys, string(passVal))

				if err != nil {
					fmt.Fprintf(ctx, err.Error())
				}

				json, err := json.MarshalIndent(val, "", "  ")

				if err != nil {
					fmt.Fprint(ctx, err.Error())
				} else {
					fmt.Fprint(ctx, string(json[:]))
				}
			} else {
				fmt.Fprintf(ctx, "incorrect password")
			}
		}
	} else if strings.Contains(request.Dynamics, "pair") {
		strVal := strings.ToUpper(request.GetUserValue(common.TrimLeftChar(request.Dynamics), ctx))
		split := strings.Split(strVal, "-")
		if !strings.Contains(request.BaseElementLocation, "volume") {
			currentPrice, err := market.CheckPrice(pairs.NewPair(split[0], split[1]))

			if err != nil {
				fmt.Fprintf(ctx, err.Error())
			} else {
				fmt.Fprintf(ctx, common.FloatToString(currentPrice))
			}
		} else {
			currentVolume := market.CheckVolume(pairs.NewPair(split[0], split[1]))

			fmt.Fprintf(ctx, common.FloatToString(currentVolume))
		}
	} else {
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
}

// HandlePost - handle POST request, with dynamics
func (request RequestElement) HandlePost(ctx *fasthttp.RequestCtx) {
	keys := strings.Split(request.Dynamics, "?")
	keys = append(keys[:0], keys[0+1:]...)

	values := []string{}

	x := 0

	for x != len(keys) {
		peekVal := string(ctx.PostArgs().Peek(keys[x]))

		if peekVal == "" {
			break
		}

		values = append(values, peekVal)

		x++
	}

	if len(values) < 2 {
		values, _ = request.GetUserValues(keys, ctx)
	}

	if common.StringInSlice("username", keys) && !common.StringInSlice("pair", keys) && !common.StringInSlice("symbol", keys) {
		fAcc, err := findAccount(request.ElementDb, values[0])

		if err != nil {
			fmt.Println(values)
			acc := accounts.NewAccount(values[0], values[1], values[2])

			err = addAccount(request.ElementDb, &acc)

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
		} else {
			if common.ComparePasswords(fAcc.PassHash, []byte(values[3])) {
				update := fAcc

				update.Username = values[0]
				update.Email = values[1]
				update.PassHash = common.HashAndSalt([]byte(values[2]))

				err = updateAccount(request.ElementDb, fAcc, &update)

				json, err := json.MarshalIndent(update, "", "  ")

				if err != nil {
					fmt.Fprintf(ctx, err.Error())
				} else {
					fmt.Fprintf(ctx, string(json[:]))
				}
			} else {
				fmt.Fprint(ctx, errors.New("invalid password").Error())
			}
		}
	} else if common.StringInSlice("pair", keys) && !strings.Contains(request.BaseElementLocation, "fill") {
		if !strings.Contains(request.BaseElementLocation, "update") {
			acc, err := findAccount(request.ElementDb, values[4])

			if err == nil {
				if common.ComparePasswords(acc.PassHash, []byte(values[5])) {
					split := strings.Split(values[0], "-")
					pair := pairs.NewPair(split[0], split[1])
					amount, _ := strconv.ParseFloat(values[2], 64)
					fillprice, _ := strconv.ParseFloat(values[3], 64)

					acc.Deposit(pair.StartingSymbol, request.ElementDb)
					acc.Deposit(pair.EndingSymbol, request.ElementDb)

					order, err := orders.NewOrder(&acc, values[1], pair, amount, fillprice)

					if err != nil {
						fmt.Fprintln(ctx, err.Error())
					} else {
						err = addOrder(request.ElementDb, &order)

						fAcc, _ := findAccount(request.ElementDb, values[4])

						updateAccount(request.ElementDb, fAcc, &acc)

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
		} else {
			fAcc, err := findAccount(request.ElementDb, values[2])

			if err == nil {
				if common.ComparePasswords(fAcc.PassHash, []byte(values[3])) {
					split := strings.Split(values[0], "-")

					pair := pairs.NewPair(split[0], split[1])
					fOrder, err := findOrder(request.ElementDb, values[1], pair)

					if err != nil {
						fmt.Fprintf(ctx, errors.New("invalid orderid").Error())
					} else {
						floatValueAmount, err := strconv.ParseFloat(values[5], 64)

						if err != nil {
							fmt.Fprintf(ctx, err.Error())
						} else {
							floatFill, err := strconv.ParseFloat(values[4], 64)

							if err != nil {
								fmt.Fprint(ctx, err.Error())
							}

							update, err := orders.NewOrder(&fAcc, fOrder.OrderType, pair, floatValueAmount, floatFill)

							update.OrderID = fOrder.OrderID

							if err != nil {
								fmt.Fprintf(ctx, err.Error())
							}

							err = database.UpdateOrder(request.ElementDb, *fOrder, update)

							if err != nil {
								fmt.Fprintf(ctx, err.Error())
							} else {
								json, err := json.MarshalIndent(update, "", "  ")

								if err != nil {
									fmt.Fprintf(ctx, err.Error())
								} else {
									fmt.Fprintf(ctx, string(json[:]))
								}
							}
						}
					}

				} else {
					fmt.Fprint(ctx, "incorrect password")
				}
			} else {
				fmt.Fprintf(ctx, "invalid account "+"'"+values[2]+"'")
			}
		}
	} else if common.StringInSlice("symbol", keys) {
		acc, err := findAccount(request.ElementDb, values[0])

		if err == nil {
			fmt.Fprintf(ctx, "waiting for deposit")

			if !common.CheckSafeSlice(acc.WalletBalances) {
				acc.WalletBalances = []float64{float64(0), float64(0), float64(0)}
			}

			fAcc, fErr := findAccount(request.ElementDb, values[0])

			if fErr != nil {
				fmt.Fprint(ctx, fErr.Error())
			}

			updateAccount(request.ElementDb, fAcc, &acc)

			go acc.Deposit(values[1], request.ElementDb)
		} else {
			fmt.Fprintf(ctx, err.Error())
		}
	} else if strings.Contains(request.BaseElementLocation, "fill") {
		split := strings.Split(values[0], "-")
		pair := pairs.NewPair(split[0], split[1])

		order, _ := findOrder(request.ElementDb, values[1], pair)

		if *order != (orders.Order{}) {
			fAcc, err := findAccount(request.ElementDb, order.Issuer.Username)

			if err != nil {
				fmt.Fprintf(ctx, err.Error())
			} else {
				err = orders.FillOrder(order, values[2])

				order.OrderPair.Volume += order.Amount

				if err != nil {
					fmt.Fprintf(ctx, err.Error())
				} else {
					updateAccount(request.ElementDb, fAcc, order.Issuer)

					removeOrder(request.ElementDb, order)

					fmt.Fprint(ctx, "order filled")
				}
			}
		} else {
			fmt.Fprintf(ctx, errors.New("invalid order").Error())
		}
	}
}

// HandleGETCollection - handle GET requests for collections
func (request RequestElement) HandleGETCollection(ctx *fasthttp.RequestCtx) {
	var collection interface{}
	var collectionKey string

	if strings.Contains(request.BaseElementLocation, "?") {
		collectionKey = strings.Split(request.BaseElementLocation, "?")[1]

		collection = string(ctx.FormValue(collectionKey))
	} else {
		collection = strings.Split(request.BaseElementLocation, "api/")[1]
	}

	var results []interface{}

	c := request.ElementDb.C(collection.(string))

	err := c.Find(nil).All(&results)
	if err != nil {
		fmt.Fprintf(ctx, err.Error())
	} else {
		json, err := json.MarshalIndent(results, "", "  ")

		if err != nil {
			fmt.Fprintf(ctx, err.Error())
		} else {
			fmt.Fprintf(ctx, string(json[:]))
		}
	}
}

// AttemptToServeRequestsWithRouter - attempts to handle incoming requests via data provided in request
func (request RequestElement) AttemptToServeRequestsWithRouter(router *fasthttprouter.Router) (*fasthttprouter.Router, error) {
	fmt.Println("atttempting to serve requests with handler: " + request.ElementName)

	if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics == "" {
		router.GET(request.BaseElementLocation, request.Handle)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics != "" {
		router.GET(request.BaseElementLocation, request.HandleVar)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[1])) {
		router.POST(request.BaseElementLocation, request.HandlePost)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[2])) {
		router.DELETE(request.BaseElementLocation, request.HandleDel)

		return router, nil
	}

	return nil, errors.New("invalid request")
}

// AttemptToServeRequests - attempts to handle incoming requests via data provided in request
func (request RequestElement) AttemptToServeRequests() (*fasthttprouter.Router, error) {
	fmt.Println("atttempting to serve requests")
	if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics == "" {
		router := fasthttprouter.New()

		router.GET(request.BaseElementLocation, request.Handle)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[0])) && request.Dynamics != "" {
		router := fasthttprouter.New()

		router.GET(request.BaseElementLocation, request.HandleVar)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[1])) {
		router := fasthttprouter.New()

		router.POST(request.BaseElementLocation, request.HandlePost)

		return router, nil
	} else if strings.Contains(strings.ToLower(request.ElementRequestType), strings.ToLower(AvailableRequestTypes[2])) {

		router := fasthttprouter.New()

		router.DELETE(request.BaseElementLocation, request.HandleDel)

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
				request := RequestElement{ElementName: name, BaseElementLocation: location, ElementRequestType: requestType, ElementContents: string(json), ElementDb: db}

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

// GetUserValues - attempts to fetch user values from specified request
func (request RequestElement) GetUserValues(keys []string, ctx *fasthttp.RequestCtx) ([]string, error) {
	x := 0

	values := []string{}

	if len(keys) == 0 {
		return []string{}, errors.New("invalid keys")
	}

	params := strings.Split(string(ctx.RequestURI()), request.BaseElementLocation)[1] // All user parameters

	for x != len(keys) {
		key := "?" + keys[x] + "=" // Key to search for in user params

		userVal := strings.Split(params, key)

		if len(userVal) == 1 {
			return values, nil
		}

		formattedVal := strings.Split(userVal[1], "?")[0]

		values = append(values, formattedVal)
		x++
	}

	return values, nil
}

// GetUserValue - attempts to fetch specified user value from request
func (request RequestElement) GetUserValue(key string, ctx *fasthttp.RequestCtx) string {
	initVal := string(ctx.PostArgs().Peek(key))

	if initVal == "" {
		initVal = string(ctx.QueryArgs().Peek(key))

		if initVal == "" || strings.Contains(initVal, "?") {
			params := strings.Split(string(ctx.RequestURI()), request.BaseElementLocation)[1] // All user parameters
			formattedKey := "?" + key + "="                                                   // Key to search for in user params

			userVal := strings.Split(params, formattedKey)[1]

			initVal = strings.Split(userVal, "?")[0]
		}
	}
	return initVal
}
