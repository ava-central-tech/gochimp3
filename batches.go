package gochimp3

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	batchesPath     = "/batches"
	singleBatchPath = batchesPath + "/%s"
)

func (api *API) GetBatchOperations(ctx context.Context, params *ListQueryParams) (*ListOfBatchOperations, error) {
	response := new(ListOfBatchOperations)

	err := api.Request(ctx, http.MethodGet, batchesPath, params, nil, response)
	if err != nil {
		return nil, err
	}

	for _, l := range response.BatchOperations {
		l.api = api
	}

	return response, nil
}

type ListOfBatchOperations struct {
	baseList
	BatchOperations []BatchOperationResponse `json:"batches"`
}

func (api *API) GetBatchOperation(ctx context.Context, id string, params *BasicQueryParams) (*BatchOperationResponse, error) {
	endpoint := fmt.Sprintf(singleBatchPath, id)
	response := new(BatchOperationResponse)
	response.api = api

	return response, api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (api *API) CreateBatchOperation(ctx context.Context, body *BatchOperationCreationRequest) (*BatchOperationResponse, error) {
	response := new(BatchOperationResponse)
	response.api = api
	return response, api.Request(ctx, http.MethodPost, batchesPath, nil, body, response)
}

type BatchOperationCreationRequest struct {
	Operations []BatchOperation `json:"operations"`
}

type BatchOperationResponse struct {
	Links []Link `json:"_links,omitempty"`

	ID                 string `json:"id"`
	Status             string `json:"status"`
	TotalOperations    int    `json:"total_operations"`
	FinishedOperations int    `json:"finished_operations"`
	ErroredOperations  int    `json:"errored_operations"`
	SubmittedAt        string `json:"submitted_at,omitempty"`
	CompletedAt        string `json:"completed_at,omitempty"`
	ResponseBodyUrl    string `json:"response_body_url"`

	api *API
}

type BatchOperation struct {
	Method      string     `json:"method"`
	Path        string     `json:"path"`
	Params      url.Values `json:"params,omitempty"`
	Body        string     `json:"body"`
	OperationID string     `json:"operation_id,omitempty"`

	api *API
}
