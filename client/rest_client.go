// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"fmt"
	glmClient "github.com/hewlettpackard/hpegl-metal-client/v1/pkg/client"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	"github.com/hpe-hcss/lh-cdc-singularity/model"
	log "github.com/hpe-storage/common-host-libs/logger"
)

type TokenType string

const (
	TokenTypeIAM TokenType = "IAM"
	TokenTypeQ   TokenType = "Q"
)

type SessionInfo struct {
	TokenMethod TokenType `json:"token_method"` // access token method, choices are "IAM" or "Q"
	TokenFile   string    `json:"token_file"`   // path and filename for token file which contains url, token, and user id
	RESTURL     string    `json:"rest_url"`     // rest_url parsed from tokenfile
	User        string    `json:"user"`         // user id parsed from tokenfile
	Token       string    `json:"token"`        // token string parsed from tokenfile
	DebugMode   bool      `json:"debugmode"`    // debug mode for logging http requests and responses
	Space       string    `json:"space"`        // space name
}

func GetREST(url string, userName string, membershipID string) (context.Context,
	*glmClient.APIClient, error) {
	log.Infof("GetREST function")
	var session SessionInfo
	// create REST Client context

	ctx := context.Background()
	sessionToken, err := model.GetSessionToken()
	if err != nil {
		log.Errorln(err)
		return ctx, nil, err
	}
	session.RESTURL = url
	session.Token = sessionToken
	session.User = userName
	// set up authentication credendtials by method

	// session object is loaded from saved file by before function run by cli command at start
	if session.RESTURL == "" {
		return ctx, nil, fmt.Errorf("no valid REST server address found")
	}
	if session.Token == "" {
		return ctx, nil, fmt.Errorf("no valid access token found")
	}

	if session.User == "" {
		return ctx, nil, fmt.Errorf("no valid memberid found for Q access token")
	}
	// add access token for auth to Client context as required by the Client API
	ctx = context.WithValue(ctx, glmClient.ContextAccessToken, session.Token)

	// Get a new Client configuration with basepath set to Quake portal URL and add base version path /rest/v1
	cfg := glmClient.NewConfiguration()
	cfg.BasePath = session.RESTURL + constants.REST_API_VERSION

	// set 'apikey' in the context to pass MembershipID or ProjectID in the header
	cfg.AddDefaultHeader("Membership", membershipID)
	// Client API debug mode flag
	if session.DebugMode {
		cfg.Debug = true
	}
	log.Infof("cfg value: %v", cfg)
	// get new API Client with basepath and auth credentials setup in configuration and context
	r := glmClient.NewAPIClient(cfg)
	log.Infof("NewAPIClient response: %+v", r)
	return ctx, r, nil
}
