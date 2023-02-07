package gochimp3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

var testServer = "http://localhost:9999"
var delegate func(http.ResponseWriter, *http.Request)

func fatalIf(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Shouldn't have gotten an error %s", err)
	}
}

func TestMain(m *testing.M) {
	http.HandleFunc("/somewhere", func(w http.ResponseWriter, r *http.Request) {
		delegate(w, r)
	})
	go func() { _ = http.ListenAndServe(":9999", nil) }()
	os.Exit(m.Run())
}

func testAPI() *API {
	api := New("apikey")
	api.endpoint = testServer
	api.Debug = true
	return api
}

func TestGoodGet(t *testing.T) {
	ctx := context.Background()

	expected := map[string]any{
		"one": "thing",
		"two": "thing",
	}
	delegate = func(w http.ResponseWriter, r *http.Request) {
		// check the headers
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// check auth
		_, pass, ok := r.BasicAuth()
		assert.Equal(t, "apikey", pass)
		assert.True(t, ok)

		// check the query params
		assert.Empty(t, r.URL.Query())

		// check we sent an empty body
		defer func() { _ = r.Body.Close() }()
		body, err := io.ReadAll(r.Body)
		fatalIf(t, err)
		assert.Equal(t, 0, len(body))

		// return something
		data, _ := json.Marshal(expected)
		_, _ = fmt.Fprintf(w, string(data))
	}

	api := testAPI()

	actual := make(map[string]any)
	err := api.Request(ctx, http.MethodGet, "/somewhere", nil, nil, &actual)
	fatalIf(t, err)

	assert.EqualValues(t, expected, actual)
}

func TestGetWithParams(t *testing.T) {
	ctx := context.Background()

	expected := map[string]any{
		"one": "thing",
		"two": "thing",
	}
	params := BasicQueryParams{
		Fields:        []string{"marp", "parm"},
		ExcludeFields: []string{"red", "fish"},
	}

	delegate = func(w http.ResponseWriter, r *http.Request) {
		// check the query params
		for k, v := range r.URL.Query() {
			switch {
			case k == "fields":
				assert.EqualValues(t, []string{strings.Join(params.Fields, ",")}, v)
			case k == "exclude_fields":
				assert.EqualValues(t, []string{strings.Join(params.ExcludeFields, ",")}, v)
			default:
				t.Fail()
			}
		}
		data, _ := json.Marshal(expected)
		_, _ = fmt.Fprintf(w, string(data))
	}

	api := testAPI()

	actual := make(map[string]any)
	err := api.Request(ctx, http.MethodGet, "/somewhere", &params, nil, &actual)
	fatalIf(t, err)

	assert.EqualValues(t, expected, actual)
}

func TestGetEmptyResponse(t *testing.T) {
	ctx := context.Background()

	delegate = func(w http.ResponseWriter, r *http.Request) {}
	api := testAPI()
	err := api.Request(ctx, http.MethodGet, "/somewhere", nil, nil, nil)
	fatalIf(t, err)
	result, err := api.RequestOk(ctx, http.MethodGet, "/somewhere")
	fatalIf(t, err)
	assert.True(t, result)
}

func TestGetWithBody(t *testing.T) {
	ctx := context.Background()

	s := struct {
		A string
		B string
	}{"string1", "string2"}

	delegate = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		defer func() { _ = r.Body.Close() }()

		body, err := io.ReadAll(r.Body)
		fatalIf(t, err)

		parsed := struct {
			A string
			B string
		}{}
		err = json.Unmarshal(body, &parsed)
		fatalIf(t, err)
		assert.EqualValues(t, s, parsed)
	}

	api := testAPI()
	err := api.Request(ctx, http.MethodPost, "/somewhere", nil, &s, nil)
	fatalIf(t, err)
}

func TestGetWithNon200Response(t *testing.T) {
	ctx := context.Background()

	delegate = func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(&APIError{
			Type:     "some type",
			Title:    "a title",
			Status:   500,
			Detail:   "you done screwed up",
			Instance: "123123123",
		})
		fatalIf(t, err)
		http.Error(w, string(data), 500)
	}

	api := testAPI()
	ok, err := api.RequestOk(ctx, http.MethodGet, "/somewhere")
	assert.False(t, ok)
	assert.NotNil(t, err)
}

func TestMissingEndpoint(t *testing.T) {
	ctx := context.Background()

	api := testAPI()
	ok, err := api.RequestOk(ctx, http.MethodGet, "/nowhere")
	assert.False(t, ok)
	assert.NotNil(t, err)
}
