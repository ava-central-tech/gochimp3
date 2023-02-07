package gochimp3

import (
	"context"
	"net/http"
)

const (
	campaignFoldersPath = "/campaign-folders"
	// single folder endpoint not implemented
)

type CampaignFolderQueryParams struct {
	ExtendedQueryParams
}

type ListOfCampaignFolders struct {
	baseList
	Folders []CampaignFolder `json:"folders"`
}

type CampaignFolder struct {
	withLinks

	Name  string `json:"name"`
	ID    string `json:"id"`
	Count uint   `json:"count"`

	api *API
}

type CampaignFolderCreationRequest struct {
	Name string `json:"name"`
}

func (api *API) GetCampaignFolders(ctx context.Context, params *CampaignFolderQueryParams) (*ListOfCampaignFolders, error) {
	response := new(ListOfCampaignFolders)

	err := api.Request(ctx, http.MethodGet, campaignFoldersPath, params, nil, response)
	if err != nil {
		return nil, err
	}

	for _, l := range response.Folders {
		l.api = api
	}

	return response, nil
}

func (api *API) CreateCampaignFolder(ctx context.Context, body *CampaignFolderCreationRequest) (*CampaignFolder, error) {
	response := new(CampaignFolder)
	response.api = api
	return response, api.Request(ctx, http.MethodPost, campaignFoldersPath, nil, body, response)
}
