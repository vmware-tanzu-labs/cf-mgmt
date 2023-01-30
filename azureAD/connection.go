package azureAD

import (
	"encoding/json"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Token struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

func getAccessToken(config *config.AzureADConfig) (Token, error) {

	token := Token{}
	reqBody := url.Values{}
	reqBody.Set("client_id", config.ClientId)
	reqBody.Set("client_secret", config.Secret)
	reqBody.Set("scope", SCOPE)
	reqBody.Set("grant_type", GRANT_TYPE)

	tokenUrl := TOKEN_LOGIN_URL + config.TenantID + TOKEN_OAUTH_URI
	req, err := http.NewRequest(http.MethodPost, tokenUrl, strings.NewReader(reqBody.Encode()))
	if err != nil {
		return token, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return token, err
	}

	defer response.Body.Close()

	// Result processing
	body, _ := io.ReadAll(response.Body)

	if err = json.Unmarshal(body, &token); err != nil {
		lo.G.Debugf("Reading body into Token struct failed: %s", err)
		return token, nil
	}

	lo.G.Debugf("Result Token: %s", token.AccessToken)
	return token, nil
}
