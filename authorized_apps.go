package gochimp3

import (
	"context"
	"fmt"
	"net/http"
)

const (
	authorizedAppsPath      = "/authorized-apps"
	singleAuthorizedAppPath = authorizedAppsPath + "/%s"
)

type ListOfAuthorizedApps struct {
	baseList `json:""`
	Apps     []AuthorizedApp `json:""`
}

type AuthorizedAppRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AuthorizedApp struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Users       []string `json:"users "`
	withLinks
}

type AuthorizedAppCreateResponse struct {
	AccessToken string `json:"access_token"`
	ViewerToken string `json:"viewer_token"`
}

func (api *API) GetAuthorizedApps(ctx context.Context, params *ExtendedQueryParams) (*ListOfAuthorizedApps, error) {
	response := new(ListOfAuthorizedApps)

	err := api.Request(ctx, http.MethodGet, authorizedAppsPath, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *API) CreateAuthorizedApp(ctx context.Context, body *AuthorizedAppRequest) (*AuthorizedAppCreateResponse, error) {
	response := new(AuthorizedAppCreateResponse)

	err := api.Request(ctx, http.MethodGet, authorizedAppsPath, nil, body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *API) GetAuthorizedApp(ctx context.Context, id string, params *BasicQueryParams) (*AuthorizedApp, error) {
	response := new(AuthorizedApp)
	endpoint := fmt.Sprintf(singleAuthorizedAppPath, id)

	err := api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
