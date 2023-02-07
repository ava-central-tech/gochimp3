package gochimp3

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
)

const (
	automationsPath      = "/automations"
	singleAutomationPath = automationsPath + "/%s"
	pauseAllEmailsPath   = singleAutomationPath + "/actions/pause-all-emails"
	startAllEmailsPath   = singleAutomationPath + "/actions/start-all-emails"

	automationEmailPath       = singleAutomationPath + "/emails"
	singleAutomationEmailPath = automationEmailPath + "/%s"

	pauseSingleEmailPath = singleAutomationEmailPath + "/actions/pause"
	startSingleEmailPath = singleAutomationEmailPath + "/actions/start"

	automationQueuesPath      = singleAutomationEmailPath + "/queue"
	singleAutomationQueuePath = automationQueuesPath + "/%s"

	removedSubscribersAutomationPath = singleAutomationPath + "/removed-subscribers"
)

type ListOfAutomations struct {
	baseList
	Automations []Automation `json:"automations"`
}

type Automation struct {
	ID              string                  `json:"id"`
	CreateTime      string                  `json:"create_time"`
	StartTime       string                  `json:"start_time"`
	Status          string                  `json:"status"`
	EmailsSent      int                     `json:"emails_sent"`
	Recipients      AutomationRecipient     `json:"recipients"`
	Settings        AutomationSettingsShort `json:"settings"`
	Tracking        AutomationTracking      `json:"tracking"`
	TriggerSettings WorkflowType            `json:"trigger_settings"`
	ReportSummary   ReportSummary           `json:"report_summary"`

	withLinks
	api *API
}

type AutomationRecipient struct {
	ListID         string            `json:"list_id"`
	SegmentOptions AutomationOptions `json:"segment_options"`
}

type AutomationOptions struct {
	SavedSegmentID int                  `json:"saved_segment_id"`
	Match          string               `json:"match"`
	Conditions     []SegmentConditional `json:"conditions"`
}

type AutomationSettingsShort struct {
	UseConversation bool   `json:"use_conversation"`
	ToName          string `json:"to_name"`
	Title           string `json:"title"`
	FromName        string `json:"from_name"`
	ReplyTo         string `json:"reply_to"`
	Authenticate    bool   `json:"authenticate"`
	AutoFooter      bool   `json:"auto_footer"`
	InlineCSS       bool   `json:"inline_css"`
}

type AutomationSettingsLong struct {
	Title        string   `json:"title"`
	FromName     string   `json:"from_name"`
	ReplyTo      string   `json:"reply_to"`
	Authenticate bool     `json:"authenticate"`
	AutoFooter   bool     `json:"auto_footer"`
	InlineCSS    bool     `json:"inline_css"`
	SubjectLine  string   `json:"subject_line"`
	AutoTweet    bool     `json:"auto_tweet"`
	AutoFBPost   []string `json:"auto_fb_post"`
	FBComments   bool     `json:"fb_comments"`
	TemplateID   int      `json:"template_id"`
	DragAndDrop  bool     `json:"drag_and_drop"`
}

type AutomationTracking struct {
	Opens           bool       `json:"opens"`
	HTMLClicks      bool       `json:"html_clicks"`
	TextClicks      bool       `json:"text_clicks"`
	GoalTracking    bool       `json:"goal_tracking"`
	Ecomm360        bool       `json:"ecomm360"`
	GoogleAnalytics string     `json:"google_analytics"`
	Clicktale       string     `json:"clicktale"`
	Salesforce      Salesforce `json:"salesforce"`
	Highrise        Highrise   `json:"highrise"`
	Capsule         Capsule    `json:"capsule"`
}

type Salesforce struct {
	Campaign bool `json:"campaign"`
	Notes    bool `json:"notes"`
}

type Highrise struct {
	Campaign bool `json:"campaign"`
	Notes    bool `json:"notes"`
}

type Capsule struct {
	Notes bool `json:"notes"`
}

type ReportSummary struct {
	Opens            int     `json:"opens"`
	UniqueOpens      int     `json:"unique_opens"`
	OpenRate         float64 `json:"open_rate"`
	Clicks           int     `json:"clicks"`
	SubscriberClicks int     `json:"subscriber_clicks"`
	ClickRate        float64 `json:"click_rate"`
}

func (auto *Automation) CanMakeRequest() error {
	if auto.ID == "" {
		return errors.New("no ID provided")
	}

	return nil
}

func (api *API) GetAutomations(ctx context.Context, params *BasicQueryParams) (*ListOfAutomations, error) {
	response := new(ListOfAutomations)

	err := api.Request(ctx, http.MethodGet, automationsPath, params, nil, response)
	if err != nil {
		return nil, err
	}

	for _, l := range response.Automations {
		l.api = api
	}

	return response, nil
}

// TODO query params?
func (api *API) GetAutomation(ctx context.Context, id string) (*Automation, error) {
	endpoint := fmt.Sprintf(singleAutomationPath, id)

	response := new(Automation)
	response.api = api

	return response, api.Request(ctx, http.MethodGet, endpoint, nil, nil, response)
}

// ------------------------------------------------------------------------------------------------
// Actions for Sending Emails
// ------------------------------------------------------------------------------------------------

func (auto *Automation) PauseSendingAll(ctx context.Context) (bool, error) {
	if err := auto.CanMakeRequest(); err != nil {
		return false, err
	}
	return auto.api.PauseSendingAll(ctx, auto.ID)
}

func (api *API) PauseSendingAll(ctx context.Context, id string) (bool, error) {
	endpoint := fmt.Sprintf(pauseAllEmailsPath, id)
	return api.RequestOk(ctx, http.MethodPost, endpoint)
}

func (auto *Automation) StartSendingAll(ctx context.Context) (bool, error) {
	if err := auto.CanMakeRequest(); err != nil {
		return false, err
	}
	return auto.api.StartSendingAll(ctx, auto.ID)
}

func (api *API) StartSendingAll(ctx context.Context, id string) (bool, error) {
	endpoint := fmt.Sprintf(startAllEmailsPath, id)
	return api.RequestOk(ctx, http.MethodPost, endpoint)
}

func (email *AutomationEmail) PauseSending(ctx context.Context) (bool, error) {
	return email.api.PauseSending(ctx, email.WorkflowID, email.ID)
}

func (api *API) PauseSending(ctx context.Context, workflowID, emailID string) (bool, error) {
	endpoint := fmt.Sprintf(pauseSingleEmailPath, workflowID, emailID)
	return api.RequestOk(ctx, http.MethodPost, endpoint)
}

func (email *AutomationEmail) StartSending(ctx context.Context) (bool, error) {
	return email.api.StartSending(ctx, email.WorkflowID, email.ID)
}

func (api *API) StartSending(ctx context.Context, workflowID, emailID string) (bool, error) {
	endpoint := fmt.Sprintf(startSingleEmailPath, workflowID, emailID)
	return api.RequestOk(ctx, http.MethodPost, endpoint)
}

// ------------------------------------------------------------------------------------------------
// Automation Emails
// ------------------------------------------------------------------------------------------------

type ListOfEmails struct {
	baseList
	Emails []AutomationEmail `json:"emails"`
}

type AutomationEmail struct {
	ID              string                 `json:"id"`
	WorkflowID      string                 `json:"workflow_id"`
	Position        int                    `json:"position"`
	Delay           AutomationDelay        `json:"delay"`
	CreateTime      string                 `json:"create_time"`
	StartTime       string                 `json:"start_time"`
	ArchiveURL      string                 `json:"archive_url"`
	Status          string                 `json:"status"`
	EmailsSent      int                    `json:"emails_sent"`
	SendTime        string                 `json:"send_time"`
	ContentType     string                 `json:"content_type"`
	Recipients      AutomationRecipient    `json:"recipients"`
	Settings        AutomationSettingsLong `json:"settings"`
	Tracking        AutomationTracking     `json:"tracking"`
	SocialCard      SocialCard             `json:"social_card"`
	TriggerSettings WorkflowType           `json:"trigger_settings"`
	ReportSummary   ReportSummary          `json:"report_summary"`

	withLinks
	api *API
}

type SocialCard struct {
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type AutomationDelay struct {
	Amount    int    `json:"amount"`
	Type      string `json:"type"`
	Direction string `json:"direction"`
	Action    string `json:"action"`
}

func (email *AutomationEmail) CanMakeRequest() error {
	if email.ID == "" {
		return errors.New("no ID provided")
	}

	return nil
}

func (auto *Automation) GetEmails(ctx context.Context) (*ListOfEmails, error) {
	if err := auto.CanMakeRequest(); err != nil {
		return nil, err
	}

	return auto.api.GetAutomationEmails(ctx, auto.ID)
}

func (api *API) GetAutomationEmails(ctx context.Context, automationID string) (*ListOfEmails, error) {
	endpoint := fmt.Sprintf(automationEmailPath, automationID)
	response := new(ListOfEmails)

	for _, l := range response.Emails {
		l.api = api
	}

	return response, api.Request(ctx, http.MethodGet, endpoint, nil, nil, response)
}

func (auto *Automation) GetEmail(ctx context.Context, id string) (*AutomationEmail, error) {
	if err := auto.CanMakeRequest(); err != nil {
		return nil, err
	}

	return auto.api.GetAutomationEmail(ctx, auto.ID, id)
}

func (api *API) GetAutomationEmail(ctx context.Context, automationID, emailID string) (*AutomationEmail, error) {
	endpoint := fmt.Sprintf(singleAutomationEmailPath, automationID, emailID)
	response := new(AutomationEmail)
	response.api = api

	return response, api.Request(ctx, http.MethodGet, endpoint, nil, nil, response)
}

// ------------------------------------------------------------------------------------------------
// Queues
// ------------------------------------------------------------------------------------------------

type AutomationQueueRequest struct {
	EmailAddress string `json:"email_address"`
}

type ListOfAutomationQueues struct {
	baseList
	WorkflowID string            `json:"workflow_id"`
	EmailID    string            `json:"email_id"`
	Queues     []AutomationQueue `json:"queue"`
}

type AutomationQueue struct {
	ID           string `json:"id"`
	WorkflowID   string `json:"workflow_id"`
	EmailID      string `json:"email_id"`
	ListID       string `json:"list_id"`
	EmailAddress string `json:"email_address"`
	NextSend     string `json:"next_send"`
	withLinks

	api *API
}

func (email *AutomationEmail) GetQueues(ctx context.Context) (*ListOfAutomationQueues, error) {
	if err := email.CanMakeRequest(); err != nil {
		return nil, err
	}

	return email.api.GetAutomationQueues(ctx, email.WorkflowID, email.ID)
}

func (api *API) GetAutomationQueues(ctx context.Context, workflowID, emailID string) (*ListOfAutomationQueues, error) {
	endpoint := fmt.Sprintf(automationQueuesPath, workflowID, emailID)

	response := new(ListOfAutomationQueues)
	for _, l := range response.Queues {
		l.api = api
	}

	return response, api.Request(ctx, http.MethodGet, endpoint, nil, nil, response)
}

func (email *AutomationEmail) GetQueue(ctx context.Context, id string) (*AutomationQueue, error) {
	if err := email.CanMakeRequest(); err != nil {
		return nil, err
	}

	return email.api.GetAutomationQueue(ctx, email.WorkflowID, email.ID, id)
}

func (api *API) GetAutomationQueue(ctx context.Context, workflowID, emailID, subsID string) (*AutomationQueue, error) {
	endpoint := fmt.Sprintf(singleAutomationQueuePath, workflowID, emailID, subsID)

	response := new(AutomationQueue)
	response.api = api

	return response, api.Request(ctx, http.MethodGet, endpoint, nil, nil, response)
}

func (email *AutomationEmail) CreateQueue(ctx context.Context, emailAddress string) (*AutomationQueue, error) {
	if err := email.CanMakeRequest(); err != nil {
		return nil, err
	}

	return email.api.CreateAutomationEmailQueue(ctx, email.WorkflowID, email.ID, emailAddress)
}

func (api *API) CreateAutomationEmailQueue(ctx context.Context, workflowID, emailID, emailAddress string) (*AutomationQueue, error) {
	endpoint := fmt.Sprintf(automationQueuesPath, workflowID, emailID)
	response := new(AutomationQueue)

	body := &AutomationQueueRequest{
		EmailAddress: emailAddress,
	}

	err := api.Request(ctx, http.MethodPost, endpoint, nil, body, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ------------------------------------------------------------------------------------------------
// Removed Subscribers
// ------------------------------------------------------------------------------------------------

type RemovedSubscriberRequest struct {
	EmailAddress string `json:"email_address"`
}

type ListOfRemovedSubscribers struct {
	baseList
	WorkflowID  string              `json:"workflow_id"`
	Subscribers []RemovedSubscriber `json:"subscribers"`
}

type RemovedSubscriber struct {
	ID           string `json:"id"`
	WorkflowID   string `json:"workflow_id"`
	ListID       string `json:"list_id"`
	EmailAddress string `json:"email_address"`

	withLinks
}

func (auto *Automation) GetRemovedSubscribers(ctx context.Context) (*ListOfRemovedSubscribers, error) {
	if err := auto.CanMakeRequest(); err != nil {
		return nil, err
	}

	return auto.api.GetAutomationRemovedSubscribers(ctx, auto.ID)
}

func (api *API) GetAutomationRemovedSubscribers(ctx context.Context, workflowID string) (*ListOfRemovedSubscribers, error) {
	endpoint := fmt.Sprintf(removedSubscribersAutomationPath, workflowID)

	response := new(ListOfRemovedSubscribers)

	return response, api.Request(ctx, http.MethodGet, endpoint, nil, nil, response)
}

func (auto *Automation) CreateRemovedSubscribers(ctx context.Context, emailAddress string) (*RemovedSubscriber, error) {
	if err := auto.CanMakeRequest(); err != nil {
		return nil, err
	}

	return auto.api.CreateAutomationRemovedSubscribers(ctx, auto.ID, emailAddress)
}

func (api *API) CreateAutomationRemovedSubscribers(ctx context.Context, workflowID, emailAddress string) (*RemovedSubscriber, error) {
	endpoint := fmt.Sprintf(removedSubscribersAutomationPath, workflowID)

	response := new(RemovedSubscriber)
	body := &RemovedSubscriberRequest{
		EmailAddress: emailAddress,
	}

	return response, api.Request(ctx, http.MethodPost, endpoint, nil, body, response)
}
