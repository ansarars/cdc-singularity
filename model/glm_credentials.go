// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package model

import (
	"errors"
	"fmt"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	log "github.com/hpe-storage/common-host-libs/logger"
	"gopkg.in/ini.v1"
)

type GLMCredDetails struct{}

type AuthSvcInfo struct {
	AuthURL            string   `json:"auth_url"`
	AuthClientID       string   `json:"auth_client_id"`
	AuthCLIClientID    string   `json:"auth_cli_client_id"`
	AuthAudience       string   `json:"auth_audience"`
	AuthSAMLConnection string   `json:"auth_saml_connection"`
	NoPasswordDomains  []string `json:"no_password_domains"`
}

type AuthRequest struct {
	GrantType string `json:"grant_type"`
	ClientID  string `json:"client_id"`
	Audience  string `json:"audience"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Scope     string `json:"scope"`
	Realm     string `json:"realm"`
}

type AuthResponse struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

var TokenError = errors.New("session token is not present in plugin.conf file")

func (glm *GLMCredDetails) GetGLMCredDetails() (map[string]string, error) {
	cfg, err := ini.Load(constants.PLUGIN_CONF_FILE)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file '%s'", constants.PLUGIN_CONF_FILE)
	}

	glmSection, err := cfg.GetSection(constants.GLM_CREDENTIALS)
	if err != nil {
		return nil, fmt.Errorf("GLM section %s not configured in plugin.conf file", constants.GLM_CREDENTIALS)
	}
	glmCredDetails := map[string]string{}
	keys := glmSection.Keys()
	for _, key := range keys {
		glmCredDetails[key.Name()] = key.Value()
	}
	return glmCredDetails, nil
}

func (glm *GLMCredDetails) UpdateGLMCredentials(key string, value string) error {
	cfg, err := ini.Load(constants.PLUGIN_CONF_FILE)
	if err != nil {
		return fmt.Errorf("failed to read configuration file '%s'", constants.PLUGIN_CONF_FILE)
	}
	cfg.Section(constants.GLM_CREDENTIALS).Key(key).SetValue(value)
	err = cfg.SaveTo(constants.PLUGIN_CONF_FILE)
	if err != nil {
		return fmt.Errorf("failed to save contents in configuration file '%s'", constants.PLUGIN_CONF_FILE)
	}
	return nil
}

func GetSessionToken() (string, error) {
	glm := &GLMCredDetails{}
	glmCredDetails, err := glm.GetGLMCredDetails()
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	if glmCredDetails[constants.SESSION_TOKEN] == "" {
		log.Errorln("session token is not present in plugin.conf file")
		return "", TokenError
	}
	return glmCredDetails[constants.SESSION_TOKEN], nil
}
