package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"strconv"
	"time"
)

type UserRequest struct {
	Username string
	Request  *fasthttp.RequestCtx
}

var ErrorModels = map[int]string{
	1:  "User is not authenticated",
	2:  "User is not allowed to access this resource",
	3:  "Resource not found",
	5:  "Invalid request body",
	6:  "Invalid request parameters",
	7:  "Server side error",
	8:  "You can vouch/modify user reputation just 1 time",
	9:  "Cannot vouch/modify your own reputation",
	10: "You have not vouched for user",
	11: "Could not verify signature. IP has been banned",
	12: "You are banned from the forum until %s. Reason: %s",
}

// RandomString generates a random string of the given length
func RandomString(n int) string {
	var buf = make([]byte, n*2)
	rand.Read(buf)

	return hex.EncodeToString(buf)[:n]
}

// ParseParam is used to get a specific parameter of a request
// if it's not provided, it returns false and writes an error to
// the client
func ParseParam(name string, ctx *fasthttp.RequestCtx, params fasthttprouter.Params) (string, bool) {
	var result string

	result = params.ByName(name)
	if result == "" {
		Error(ctx, 6, "Invalid request")
		return result, false
	}

	return result, true
}

// ParseBody parses the JSON data from the request body
// and returns the given object
func ParseBody(ctx *UserRequest, object interface{}) interface{} {
	if string(ctx.Request.Request.Header.Peek("Content-Type")) != "application/json" {
		Error(ctx.Request, 5, "Invalid request body")
		return nil
	}

	var body = ctx.Request.Request.Body()
	var err = json.Unmarshal(body, object)
	if err != nil {
		if ctx.Username == "" {
			logger.Logf(logger.INFO, "Invalid body used in request to endpoint '%s'",
				string(ctx.Request.Path()))
		} else {
			logger.Logf(logger.INFO, "Invalid body used in request by user '%s' to endpoint '%s'",
				ctx.Username, string(ctx.Request.Path()))
		}

		Error(ctx.Request, 5, "Invalid request body")

		return nil
	}

	return object
}

// Error is used to return an error to the client
// it makes use of the error codes defined above
// and also takes a custom message
func Error(ctx *fasthttp.RequestCtx, code int, message string) {
	JSON(ctx, models.ErrorResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  false,
			Message: message,
		},
		Error: models.ErrorModel{
			Code:    code,
			Message: ErrorModels[code],
		},
	})
}

// JSON is used to turn an interface into bytes
// that we can send as a response
func JSON(ctx *fasthttp.RequestCtx, data interface{}) {
	var encoded []byte

	ctx.Response.Header.Set("Content-Type", "application/json")

	encoded, _ = json.Marshal(data)

	_, _ = ctx.Write(encoded)
}

// ConvertToInteger is used to convert a user-provided
// variable into an integer, if it fails we return an error
// to the user, we take the endpoint as an argument for logging
func ConvertToInteger(x string, endpoint string, ctx *fasthttp.RequestCtx) (int, bool) {
	var result int
	var err error

	result, err = strconv.Atoi(x)
	if err != nil {
		logger.Logf(logger.WARNING, "Invalid usage of %s endpoint. User may attempt automation.", endpoint)
		Error(ctx, 6, "Invalid request")
		return 0, false
	}

	return result, true
}

// ParseMySQLJson decodes the base64 string MySQL uses to store JSOn
// and returns an interface which can be converted to whatever
// is needed
func ParseMySQLJson(data []byte) interface{} {
	var response interface{}

	var err = json.Unmarshal(data, &response)
	if err != nil {
		logger.Logf(logger.ERROR, "Could not parse MySQL JSON data. Error: %s", err.Error())
		return nil
	}

	return response
}

// IntfSliceContains is used to check if a given item
// is part of a string slice
func IntfSliceContains(s []interface{}, v string) bool {
	for _, item := range s {
		if item == v {
			return true
		}
	}
	return false
}

// FormatTime returns the given time.Time object
// as a string format
func FormatTime(t time.Time) string {
	return t.Format("02-01-2006 15:04:05")
}
