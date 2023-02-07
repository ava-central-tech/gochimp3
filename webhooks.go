package gochimp3

import (
	"context"
	"fmt"
	"net/http"
)

const (
	webhooksPath      = "/lists/%s/webhooks"
	singleWebhookPath = webhooksPath + "/%s"
)

type ListOfWebHooks struct {
	baseList
	ListID   string    `json:"list_id"`
	WebHooks []WebHook `json:"webhooks"`
}

type WebHookRequest struct {
	URL     string      `json:"url"`
	Events  HookEvents  `json:"events"`
	Sources HookSources `json:"sources"`
}

type WebHook struct {
	WebHookRequest
	ID     string `json:"id"`
	ListID string `json:"list_id"`
	withLinks
}

type HookSources struct {
	User  bool `json:"user"`
	Admin bool `json:"admin"`
	API   bool `json:"api"`
}

type HookEvents struct {
	Subscribe   bool `json:"subscribe"`
	Unsubscribe bool `json:"unsubscribe"`
	Profile     bool `json:"profile"`
	Cleaned     bool `json:"cleaned"`
	Upemail     bool `json:"upemail"`
	Campaign    bool `json:"campaign"`
}

func (list *ListResponse) CreateWebHooks(ctx context.Context, body *WebHookRequest) (*WebHook, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(webhooksPath, list.ID)
	response := new(WebHook)

	return response, list.api.Request(ctx, http.MethodPost, endpoint, nil, &body, response)
}

func (list *ListResponse) UpdateWebHook(ctx context.Context, id string, body *WebHookRequest) (*WebHook, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(singleWebhookPath, list.ID, id)
	response := new(WebHook)

	return response, list.api.Request(ctx, http.MethodPatch, endpoint, nil, &body, response)
}

// TODO - does this take filters? undocumented

func (list *ListResponse) GetWebHooks(ctx context.Context) (*ListOfWebHooks, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(webhooksPath, list.ID)
	response := new(ListOfWebHooks)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, nil, nil, response)
}

func (list *ListResponse) GetWebHook(ctx context.Context, id string) (*WebHook, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(singleWebhookPath, list.ID, id)
	response := new(WebHook)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, nil, nil, response)
}

func (list *ListResponse) DeleteWebHook(ctx context.Context, id string) (bool, error) {
	if err := list.CanMakeRequest(); err != nil {
		return false, err
	}

	endpoint := fmt.Sprintf(singleWebhookPath, list.ID, id)
	return list.api.RequestOk(ctx, http.MethodDelete, endpoint)
}
