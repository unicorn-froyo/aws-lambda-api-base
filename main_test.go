package api_base

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func Test_HandleRequest(t *testing.T) {
	type testModel struct {
		Message string
	}
	type Tests struct {
		name      string
		event     events.APIGatewayProxyRequest
		response  events.APIGatewayProxyResponse
		handlerFn ApiHandlerFunc
		options   []RequestOption
	}

	tests := []Tests{
		{
			name: "WithBody/WithHeaders/WithPathParams",
			event: events.APIGatewayProxyRequest{
				Body:    "{\"Message\": \"boogers\"}",
				Headers: map[string]string{},
			},
			response: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
					"Message":      "In A Bottle",
				},
				MultiValueHeaders: map[string][]string(nil),
				Body:              "{\"Message\":\"boogers\"}",
				IsBase64Encoded:   false},
			handlerFn: func(r *Request) Response {
				b, _ := json.Marshal(r.body)
				return Response{
					StatusCode: 200,
					Body:       string(b),
					Headers:    map[string]string{"Message": "In A Bottle"},
				}
			},
			options: []RequestOption{
				WithBody(testModel{}),
				WithHeaders(testModel{}),
				WithPathParams(testModel{}),
			},
		},
		{
			name:  "Nothing",
			event: events.APIGatewayProxyRequest{},
			response: events.APIGatewayProxyResponse{
				StatusCode:        400,
				Headers:           map[string]string{"Content-Type": "application/json"},
				MultiValueHeaders: map[string][]string(nil),
				Body:              "",
				IsBase64Encoded:   false,
			},
			handlerFn: func(r *Request) Response {
				return Response{StatusCode: 400, Body: ""}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn := NewApi(test.options...).Run(test.handlerFn)
			resp, err := fn(test.event)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resp, test.response)

		})
	}

}
