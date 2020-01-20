// Command authentication contains logic needed to make Alexa token auth compatible with the Hue API. It takes in the Alexa
// authentication request and forwards it to either the Hue /token or /refresh endpoint depending on the grant_type
// then proxies back the response
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	HueAPITarget   string = "https://api.meethue.com/"
	HueTokenPath   string = "/oauth2/token"
	HueRefreshPath string = "/oauth2/refresh"

	AuthGrantType    string = "authorization_code"
	RefreshGrantType string = "refresh_token"
)

var (
	// ErrUnrecognizedGrantType unsupported grant_type found in token request
	ErrUnrecognizedGrantType = errors.New("Unrecognized `grant_type` in Alexa token request")
)

// getHueForwardURI gets the correct URI to pass the authorization request on to based on the grant_type of the request
func getHueForwardPath(req *http.Request) (string, error) {
	if err := req.ParseForm(); err != nil {
		return "", fmt.Errorf("Error parsing form: %s", err)
	}

	gt := req.Form.Get("grant_type")

	log.Debug().Str("Handler", "Authentication").Str("Function", "handler").
		Msgf("Extracted %s grant type from request", gt)

	switch gt {
	case AuthGrantType:
		log.Debug().Str("Handler", "Authentication").Str("Function", "handler").
			Msgf("Handling %s grant type...", AuthGrantType)

		return HueTokenPath, nil
	case RefreshGrantType:
		log.Debug().Str("Handler", "Authentication").Str("Function", "handler").
			Msgf("Handling %s grant type...", RefreshGrantType)

		return HueRefreshPath, nil
	default:
		return "", ErrUnrecognizedGrantType
	}
}

// serveReverseProxy sends a proxy request to the target, appending the req path and oreq host as the forwarded host
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) error {
	baseURL, err := url.Parse(target)
	if err != nil {
		return err
	}

	rp := httputil.NewSingleHostReverseProxy(baseURL) // Set up reverse proxy to Hue API

	// Update headers to allow for SSL redirection
	req.Header.Set("X-Forwarded-Host", req.URL.Host)
	req.Host = baseURL.Host

	log.Debug().Str("Handler", "Authentication").Str("Function", "handler").
		Msgf("Sending proxied request to %v...", req.URL)

	rp.ServeHTTP(res, req)

	return nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Debug().Str("Handler", "Authentication").Str("Function", "handler").Msg("Got authentication request...")

	// Convert incoming ProxyRequest to http.Request
	ra := core.RequestAccessor{}
	or, err := ra.ProxyEventToHTTPRequest(request)
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %s", err)
	}

	// Clone incoming authentication request
	cr := or.Clone(context.TODO())
	cr.Body, err = or.GetBody()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while cloning auth request %s", err)
	}

	// Preserve received request, changing path to matching Hue authentication path
	cr.URL.Path, err = getHueForwardPath(or)
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while getting Hue forward URI %s", err)
	}

	w := core.NewProxyResponseWriter()

	// Proxy request to Hue resource server, swapping out host to Hue API host
	if err = serveReverseProxy(HueAPITarget, w, cr); err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while proxying request to Hue: %s", err)
	}

	// Convert http.Response to lambda gateway response
	resp, err := w.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %s", err)
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
