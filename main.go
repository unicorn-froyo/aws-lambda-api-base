package api_base

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Request struct {
	body       interface{}
	headers    interface{}
	pathParams interface{}
}

type Response struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

type ApiHandlerFunc func(r *Request) Response
type LambdaHandlerFunc func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
type RequestOption func(*Request)

func NewApi(opts ...RequestOption) *Request {
	ao := Request{}
	for _, o := range opts {
		o(&ao)
	}
	return &ao
}

func WithBody(b interface{}) RequestOption {
	return func(r *Request) {
		r.body = b
	}
}

func WithHeaders(h interface{}) RequestOption {
	return func(r *Request) {
		r.headers = h
	}
}

func WithPathParams(p interface{}) RequestOption {
	return func(r *Request) {
		r.pathParams = p
	}
}

func parseMap(p map[string]string, v interface{}) error {
	pj, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(pj, &v)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *Request) Run(f ApiHandlerFunc) LambdaHandlerFunc {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if r.body != nil {
			err := json.Unmarshal([]byte(request.Body), &r.body)
			if err != nil {
				panic(err)
			}
		}
		if r.headers != nil {
			err := parseMap(request.Headers, r.headers)
			if err != nil {
				panic(err)
			}
		}

		if r.pathParams != nil {
			err := parseMap(request.PathParameters, r.pathParams)
			if err != nil {
				panic(err)
			}
		}
		resp := f(r)
		return transformResponse(resp), nil
	}

}

func transformResponse(r Response) events.APIGatewayProxyResponse {

	if len(r.Headers) == 0 {
		r.Headers = map[string]string{
			"Content-Type": "application/json",
		}
	} else {
		r.Headers["Content-Type"] = "application/json"
	}
	return events.APIGatewayProxyResponse{
		StatusCode:      r.StatusCode,
		Body:            r.Body,
		Headers:         r.Headers,
		IsBase64Encoded: false,
	}
}
