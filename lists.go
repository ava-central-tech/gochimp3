package gochimp3

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
)

const (
	listsPath      = "/lists"
	singleListPath = listsPath + "/%s"

	abuseReportsPath      = "/lists/%s/abuse_reports"
	singleAbuseReportPath = abuseReportsPath + "/%s"

	activityPath = "/lists/%s/activity"
	clientsPath  = "/lists/%s/clients"

	historyPath       = "/lists/%s/growth-history"
	singleHistoryPath = historyPath + "/%s"

	interestCategoriesPath     = "/lists/%s/interest-categories"
	singleInterestCategoryPath = interestCategoriesPath + "/%s"

	interestsPath      = "/lists/%s/interest-categories/%s/interests"
	singleInterestPath = interestsPath + "/%s"

	listsBatchSubscribeMembers = "/lists/%s"

	mergeFieldsPath = "/lists/%s/merge-fields"
	mergeFieldPath  = mergeFieldsPath + "/%s"
)

type ListQueryParams struct {
	ExtendedQueryParams

	BeforeDateCreated      string
	SinceDateCreated       string
	BeforeCampaignLastSent string
	SinceCampaignLastSent  string
	Email                  string
}

func (q ListQueryParams) Params() map[string]string {
	m := q.ExtendedQueryParams.Params()
	m["before_date_created"] = q.BeforeDateCreated
	m["since_date_created"] = q.SinceDateCreated
	m["before_campaign_last_sent"] = q.BeforeCampaignLastSent
	m["since_campaign_last_sent"] = q.SinceCampaignLastSent
	m["email"] = q.Email
	return m
}

type ListOfLists struct {
	baseList
	Lists []ListResponse `json:"lists"`
}

type ListCreationRequest struct {
	Name                string           `json:"name"`
	Contact             Contact          `json:"contact"`
	PermissionReminder  string           `json:"permission_reminder"`
	UseArchiveBar       bool             `json:"use_archive_bar"`
	CampaignDefaults    CampaignDefaults `json:"campaign_defaults"`
	NotifyOnSubscribe   string           `json:"notify_on_subscribe"`
	NotifyOnUnsubscribe string           `json:"notify_on_unsubscribe"`
	EmailTypeOption     bool             `json:"email_type_option"`
	Visibility          string           `json:"visibility"`
}

type ListResponse struct {
	ListCreationRequest
	withLinks

	ID                string   `json:"id"`
	DateCreated       string   `json:"date_created"`
	ListRating        int      `json:"list_rating"`
	SubscribeURLShort string   `json:"subscribe_url_short"`
	SubscribeURLLong  string   `json:"subscribe_url_long"`
	BeamerAddress     string   `json:"beamer_address"`
	Modules           []string `json:"modules"`
	Stats             Stats    `json:"stats"`

	api *API
}

func (list *ListResponse) CanMakeRequest() error {
	if list.ID == "" {
		return errors.New("No ID provided on list")
	}

	return nil
}

type Stats struct {
	MemberCount               int     `json:"member_count"`
	UnsubscribeCount          int     `json:"unsubscribe_count"`
	CleanedCount              int     `json:"cleaned_count"`
	MemberCountSinceSend      int     `json:"member_count_since_send"`
	UnsubscribeCountSinceSend int     `json:"unsubscribe_count_since_send"`
	CleanedCountSinceSend     int     `json:"cleaned_count_since_send"`
	CampaignCount             int     `json:"campaign_count"`
	CampaignLastSent          string  `json:"campaign_last_sent"`
	MergeFieldCount           int     `json:"merge_field_count"`
	AvgSubRate                float64 `json:"avg_sub_rate"`
	AvgUnsubRate              float64 `json:"avg_unsub_rate"`
	TargetSubRate             float64 `json:"target_sub_rate"`
	OpenRate                  float64 `json:"open_rate"`
	ClickRate                 float64 `json:"click_rate"`
	LastSubDate               string  `json:"last_sub_date"`
	LastUnsubDate             string  `json:"last_unsub_date"`
}

type CampaignDefaults struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	Subject   string `json:"subject"`
	Language  string `json:"language"`
}

func (api *API) GetLists(ctx context.Context, params *ListQueryParams) (*ListOfLists, error) {
	response := new(ListOfLists)

	err := api.Request(ctx, http.MethodGet, listsPath, params, nil, response)
	if err != nil {
		return nil, err
	}

	for i := range response.Lists {
		response.Lists[i].api = api
	}

	return response, nil
}

// NewListResponse returns a *ListResponse that is minimally viable for making
// API requests. This is useful for such API requests that depend on a
// ListResponse for its ID (e.g. CreateMember) without having to make a second
// network request to get the list itself.
func (api *API) NewListResponse(id string) *ListResponse {
	return &ListResponse{
		ID:  id,
		api: api,
	}
}

func (api *API) GetList(ctx context.Context, id string, params *BasicQueryParams) (*ListResponse, error) {
	endpoint := fmt.Sprintf(singleListPath, id)

	response := new(ListResponse)
	response.api = api

	return response, api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (api *API) CreateList(ctx context.Context, body *ListCreationRequest) (*ListResponse, error) {
	response := new(ListResponse)
	response.api = api
	return response, api.Request(ctx, http.MethodPost, listsPath, nil, body, response)
}

func (api *API) UpdateList(ctx context.Context, id string, body *ListCreationRequest) (*ListResponse, error) {
	endpoint := fmt.Sprintf(singleListPath, id)

	response := new(ListResponse)
	response.api = api

	return response, api.Request(ctx, http.MethodPatch, endpoint, nil, body, response)
}

func (api *API) DeleteList(ctx context.Context, id string) (bool, error) {
	endpoint := fmt.Sprintf(singleListPath, id)
	return api.RequestOk(ctx, http.MethodDelete, endpoint)
}

// ------------------------------------------------------------------------------------------------
// Abuse Reports
// ------------------------------------------------------------------------------------------------

type ListOfAbuseReports struct {
	baseList

	ListID  string        `json:"list_id"`
	Reports []AbuseReport `json:"abuse_reports"`
}

type AbuseReport struct {
	ID           string `json:"id"`
	CampaignID   string `json:"campaign_id"`
	ListID       string `json:"list_id"`
	EmailID      string `json:"email_id"`
	EmailAddress string `json:"email_address"`
	Date         string `json:"date"`

	withLinks
}

func (list *ListResponse) GetAbuseReports(ctx context.Context, params *ExtendedQueryParams) (*ListOfAbuseReports, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(abuseReportsPath, list.ID)
	response := new(ListOfAbuseReports)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (list *ListResponse) GetAbuseReport(ctx context.Context, id string, params *ExtendedQueryParams) (*AbuseReport, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(singleAbuseReportPath, list.ID, id)
	response := new(AbuseReport)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

// ------------------------------------------------------------------------------------------------
// Activity
// ------------------------------------------------------------------------------------------------

type ListOfActivity struct {
	baseList

	ListID     string     `json:"list_id"`
	Activities []Activity `json:"activity"`
}

type Activity struct {
	Day             string `json:"day"`
	EmailsSent      int    `json:"emails_sent"`
	UniqueOpens     int    `json:"unique_opens"`
	RecipientClicks int    `json:"recipient_clicks"`
	HardBounce      int    `json:"hard_bounce"`
	SoftBounce      int    `json:"soft_bounce"`
	Subs            int    `json:"subs"`
	Unsubs          int    `json:"unsubs"`
	OtherAdds       int    `json:"other_adds"`
	OtherRemoves    int    `json:"other_removes"`

	withLinks
}

func (list *ListResponse) GetActivity(ctx context.Context, params *BasicQueryParams) (*ListOfActivity, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(activityPath, list.ID)
	response := new(ListOfActivity)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

// ------------------------------------------------------------------------------------------------
// Clients
// ------------------------------------------------------------------------------------------------

type ListOfClients struct {
	baseList

	ListID  string   `json:"list_id"`
	Clients []Client `json:"clients"`
}

type Client struct {
	Client  string `json:"client"`
	Members int    `json:"members"`
	ListID  string `json:"list_id"`

	withLinks
}

func (list *ListResponse) GetClients(ctx context.Context, params *BasicQueryParams) (*ListOfClients, error) {
	if list.ID == "" {
		return nil, errors.New("No ID provided on list")
	}

	endpoint := fmt.Sprintf(clientsPath, list.ID)
	response := new(ListOfClients)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

// ------------------------------------------------------------------------------------------------
// Growth History
// ------------------------------------------------------------------------------------------------

type ListOfGrownHistory struct {
	baseList

	ListID  string          `json:"list_id"`
	History []GrowthHistory `json:"history"`
}

type GrowthHistory struct {
	ListID   string `json:"list_id"`
	Month    string `json:"month"`
	Existing int    `json:"existing"`
	Imports  int    `json:"imports"`
	OptIns   int    `json:"optins"`

	withLinks
}

func (list *ListResponse) GetGrowthHistory(ctx context.Context, params *ExtendedQueryParams) (*ListOfGrownHistory, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(historyPath, list.ID)
	response := new(ListOfGrownHistory)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (list *ListResponse) GetGrowthHistoryForMonth(ctx context.Context, month string, params *BasicQueryParams) (*GrowthHistory, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(singleHistoryPath, list.ID, month)
	response := new(GrowthHistory)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

// ------------------------------------------------------------------------------------------------
// Interest Categories
// ------------------------------------------------------------------------------------------------

type ListOfInterestCategories struct {
	baseList
	ListID     string             `json:"list_id"`
	Categories []InterestCategory `json:"categories"`
}

type InterestCategoryRequest struct {
	Title        string `json:"title"`
	DisplayOrder int    `json:"display_order"`
	Type         string `json:"type"`
}

type InterestCategory struct {
	InterestCategoryRequest

	ListID string `json:"list_id"`
	ID     string `json:"id"`

	withLinks
	api *API
}

func (ic *InterestCategory) CanMakeRequest() error {
	if ic.ID == "" {
		return errors.New("no ID provided on interest category")
	}

	return nil
}

type InterestCategoriesQueryParams struct {
	ExtendedQueryParams

	Type string `json:"type"`
}

func (q *InterestCategoriesQueryParams) Params() map[string]string {
	m := q.ExtendedQueryParams.Params()
	m["type"] = q.Type
	return m
}

func (list *ListResponse) GetInterestCategories(ctx context.Context, params *InterestCategoriesQueryParams) (*ListOfInterestCategories, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(interestCategoriesPath, list.ID)
	response := new(ListOfInterestCategories)

	err := list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
	if err != nil {
		return nil, err
	}

	for i := range response.Categories {
		response.Categories[i].api = list.api
	}

	return response, nil
}

func (list *ListResponse) GetInterestCategory(ctx context.Context, id string, params *BasicQueryParams) (*InterestCategory, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(singleInterestCategoryPath, list.ID, id)
	response := new(InterestCategory)
	response.api = list.api

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (list *ListResponse) CreateInterestCategory(ctx context.Context, body *InterestCategoryRequest) (*InterestCategory, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(interestCategoriesPath, list.ID)
	response := new(InterestCategory)
	response.api = list.api

	return response, list.api.Request(ctx, http.MethodPost, endpoint, nil, body, response)
}

func (list *ListResponse) UpdateInterestCategory(ctx context.Context, id string, body *InterestCategoryRequest) (*InterestCategory, error) {
	if list.ID == "" {
		return nil, errors.New("No ID provided on list")
	}

	endpoint := fmt.Sprintf(singleInterestCategoryPath, list.ID, id)
	response := new(InterestCategory)
	response.api = list.api

	return response, list.api.Request(ctx, http.MethodPatch, endpoint, nil, body, response)
}

func (list *ListResponse) DeleteInterestCategory(ctx context.Context, id string) (bool, error) {
	if list.ID == "" {
		return false, errors.New("No ID provided on list")
	}

	endpoint := fmt.Sprintf(singleInterestCategoryPath, list.ID, id)
	return list.api.RequestOk(ctx, http.MethodDelete, endpoint)
}

// ------------------------------------------------------------------------------------------------
// Interest
// ------------------------------------------------------------------------------------------------

type ListOfInterests struct {
	Interests  []Interest `json:"interests"`
	CategoryID string     `json:"category_id"`
	ListID     string     `json:"list_id"`
	TotalItems int        `json:"total_items"`
	withLinks
}

type Interest struct {
	CategoryID   string `json:"category_id"`
	ListID       string `json:"list_id"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
	withLinks
}

type InterestRequest struct {
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
}

func (list *ListResponse) GetInterests(ctx context.Context, interestCategoryID string, params *ExtendedQueryParams) (*ListOfInterests, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(interestsPath, list.ID, interestCategoryID)
	response := new(ListOfInterests)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (list *ListResponse) GetInterest(ctx context.Context, interestCategoryID, interestID string, params *BasicQueryParams) (*Interest, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(singleInterestPath, list.ID, interestCategoryID, interestID)
	response := new(Interest)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (ic *InterestCategory) CreateInterest(ctx context.Context, body *InterestRequest) (*Interest, error) {
	if err := ic.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(interestsPath, ic.ListID, ic.ID)
	response := new(Interest)

	return response, ic.api.Request(ctx, http.MethodPost, endpoint, nil, body, response)
}

// ------------------------------------------------------------------------------------------------
// Batch subscribe list members
// ------------------------------------------------------------------------------------------------
type BatchSubscribeMembersError struct {
	EmailAddress string `json:"email_address"`
	ErrorMessage string ` json:"error"`
}

type BatchSubscribeMembersResponse struct {
	withLinks

	NewMembers     []ListOfMembers              `json:"new_members"`
	UpdatedMembers []ListOfMembers              `json:"updated_members"`
	ErrorMessages  []BatchSubscribeMembersError `json:"errors"`
	TotalCreated   int                          `json:"total_created"`
	TotalUpdated   int                          `json:"total_updated"`
	ErrorCount     int                          `json:"error_count"`

	api *API
}

type BatchSubscribeMembersRequest struct {
	Members        []MemberRequest `json:"members"`
	UpdateExisting bool            `json:"update_existing"`
}

func (list *ListResponse) BatchSubscribeMembers(ctx context.Context, body *BatchSubscribeMembersRequest) (*BatchSubscribeMembersResponse, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(listsBatchSubscribeMembers, list.ID)
	response := new(BatchSubscribeMembersResponse)

	return response, list.api.Request(ctx, http.MethodPost, endpoint, nil, body, response)
}

// ------------------------------------------------------------------------------------------------
// Merge Fields
// ------------------------------------------------------------------------------------------------

type MergeFieldsParams struct {
	ExtendedQueryParams

	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type MergeFieldParams struct {
	BasicQueryParams

	MergeID string `json:"_"`
}

type ListOfMergeFields struct {
	baseList

	ListID      string       `json:"list_id"`
	MergeFields []MergeField `json:"merge_fields"`
}

type MergeField struct {
	MergeID      int               `json:"merge_id"`
	Tag          string            `json:"tag"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Required     bool              `json:"required"`
	DefaultValue string            `json:"default_value"`
	Public       bool              `json:"public"`
	DisplayOrder int               `json:"display_order"`
	Options      MergeFieldOptions `json:"options"`
	HelpText     string            `json:"help_text"`
	ListID       string            `json:"list_id"`

	withLinks
}

type MergeFieldOptions struct {
	DefaultCountry int      `json:"default_Country"`
	PhoneFormat    string   `json:"phone_format"`
	DateFormat     string   `json:"date_format"`
	Choices        []string `json:"choices,omitempty"`
	Size           int      `json:"size"`
}

type MergeFieldRequest struct {
	// The tag used in MailChimp campaigns and for the /members endpoint.
	Tag string `json:"tag"`

	// The name of the merge field.
	Name string `json:"name"`

	// The type for the merge field.
	// Possible Values: text, number, address, phone, date, url, image, url, radio, dropdown, birthday, zip
	Type string `json:"type"`

	// The boolean value if the merge field is required.
	Required bool `json:"required"`

	// The default value for the merge field if null.
	DefaultValue string `json:"default_value"`

	// Whether the merge field is displayed on the signup form.
	Public bool `json:"public"`

	// The order that the merge field displays on the list signup form.
	DisplayOrder int `json:"display_order"`

	// The order that the merge field displays on the list signup form.
	Options MergeFieldOptions `json:"options"`

	// Extra text to help the subscriber fill out the form.
	HelpText string `json:"help_text"`
}

func (list *ListResponse) GetMergeFields(ctx context.Context, params *MergeFieldsParams) (*ListOfMergeFields, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(mergeFieldsPath, list.ID)
	response := new(ListOfMergeFields)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (list *ListResponse) GetMergeField(ctx context.Context, params *MergeFieldParams) (*MergeField, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(mergeFieldPath, list.ID, params.MergeID)
	response := new(MergeField)

	return response, list.api.Request(ctx, http.MethodGet, endpoint, params, nil, response)
}

func (list *ListResponse) CreateMergeField(ctx context.Context, body *MergeFieldRequest) (*MergeField, error) {
	if err := list.CanMakeRequest(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(mergeFieldsPath, list.ID)
	response := new(MergeField)

	return response, list.api.Request(ctx, http.MethodPost, endpoint, nil, body, response)
}
