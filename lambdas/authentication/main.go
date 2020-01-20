package main

import (
	"context"
	"errors"
	"net/http/httputil"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Hue token URLs
const (
	HueAPIBase    string = "https://api.meethue.com/"
	HueTokenURL   string = "https://api.meethue.com/oauth2/token"
	HueRefreshURL string = "https://api.meethue.com/oauth2/refresh"

	AuthGrantType    string = "authorization_code"
	RefreshGrantType string = "refresh_token"
)

var (
	// ErrUnrecognizedGrantType unsupported grant_type found in token request
	ErrUnrecognizedGrantType = errors.New("Unrecognized `grant_type` in Alexa token request")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Debug().Str("Handler", "Authentication").Str("Function", "handler").Msg("Got authentication request...")

	// Convert incoming ProxyRequest to http.Request
	ra := core.RequestAccessor{}
	httpReq, err := ra.ProxyEventToHTTPRequest(request)
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	if err := httpReq.ParseForm(); err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error parsing form: %s", err)
	}

	gt := httpReq.Form.Get("grant_type")
	httpReq.PostForm.Get("grant_type")

	log.Debug().Str("Handler", "Authentication").Str("Function", "handler").
		Msgf("Extracted %s grant type from request", gt)

	// Determine if Alexa is trying to get a new token or refreshing one
	var forwardPath string
	switch gt {
	case AuthGrantType:
		log.Debug().Str("Handler", "Authentication").Str("Function", "handler").
			Msgf("Handling %s grant type...", AuthGrantType)
		forwardPath = HueTokenURL
	case RefreshGrantType:
		log.Debug().Str("Handler", "Authentication").Str("Function", "handler").
			Msgf("Handling %s grant type...", RefreshGrantType)
		forwardPath = HueRefreshURL
	default:
		return events.APIGatewayProxyResponse{}, ErrUnrecognizedGrantType
	}

	// Clone incoming authentication request
	clone := httpReq.Clone(context.TODO())
	clone.Body, err = httpReq.GetBody()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while cloning auth request %v", err)
	}

	// Preserve received request, changing URL to correct Hue API URL
	forwardURL, err := url.ParseRequestURI(forwardPath)
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while creating new forward URL %v", err)
	}
	clone.URL = forwardURL

	// Proxy request to Hue resource server
	w := core.NewProxyResponseWriter()
	baseURL, _ := url.Parse(HueAPIBase)
	rp := httputil.NewSingleHostReverseProxy(baseURL) // Set up reverse proxy to Hue API

	// Update headers to allow for SSL redirection
	clone.Header.Set("X-Forwarded-Host", httpReq.Header.Get("Host"))
	clone.Host = baseURL.Host

	log.Debug().Str("Handler", "Authentication").Str("Function", "handler").
		Msg("Sending authentication request to Hue...")
	rp.ServeHTTP(w, clone)

	resp, err := w.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	if resp.StatusCode != 200 {
		return core.GatewayTimeout(), core.NewLoggedError(
			"Non 200 Response found getting token. Resp code: %d, Resp body: %s", resp.StatusCode, resp.Body,
		)
	}

	log.Debug().Str("Handler", "Authentication").Str("Function", "handler").Msg("Success. Passing response back.")

	return resp, nil
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	lambda.Start(handler)
}
