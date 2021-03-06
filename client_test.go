package prom

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestCase struct {
	ApiKey      string
	Response    DummyResponse
	GetParams   map[string]string
	PostRequest DummyPostRequest
	IsError     bool
}

type DummyResponse struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

type DummyPostRequest struct {
}

func ServerDummy(w http.ResponseWriter, r *http.Request) {
	authKey := r.Header.Get("Authorization")

	if authKey != "Bearer authorizedTestKey" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`<html>
 <head>
  <title>401 Not Authenticated</title>
 </head>
 <body>
  <h1>401 Not Authenticated</h1>
  This server could not verify that you are authorized to access the document you requested. Either you supplied the wrong credentials (e.g., bad password), or your browser does not understand how to supply the credentials required.<br /><br />
Not Authenticated


 </body>
</html>`))
		return
	}

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		switch r.URL.Query().Get("test_param") {
		case "1":
			w.Write([]byte("wrong string non-json"))
			return
		case "2":
			w.Write([]byte("{\"error\": \"some errors happen\"}"))
			return
		case "3":
			w.Write([]byte("{\"data\": \"a new data not error\"}"))
			return
		}

	} else if r.Method == http.MethodPost {

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`<html>
 <head>
  <title>405 Method Not Allowed</title>
 </head>
 <body>
  <h1>405 Method Not Allowed</h1>
  The method PUT is not allowed for this resource. <br /><br />

 </body>
</html>`))
	}
}

func TestClient(t *testing.T) {
	c := NewClient("testApiKey")
	cTest := &Client{
		apiUrl: defaultApiUrl,
		apiKey: "testApiKey",
	}
	if !reflect.DeepEqual(cTest, c) {
		t.Errorf("[%d] wrong result, expected %#v, got %#v", 0, cTest, c)
	}
}

func TestGet(t *testing.T) {
	cases := []TestCase{
		TestCase{
			ApiKey:  "unauthorizedKey",
			IsError: true,
		},
		TestCase{
			ApiKey:    "authorizedTestKey",
			GetParams: map[string]string{"test_param": "1"},
			Response:  DummyResponse{},
			IsError:   true,
		},
		TestCase{
			ApiKey:    "authorizedTestKey",
			GetParams: map[string]string{"test_param": "2"},
			Response: DummyResponse{
				Error: "some errors happen",
			},
			IsError: false,
		},
		TestCase{
			ApiKey:    "authorizedTestKey",
			GetParams: map[string]string{"test_param": "3"},
			Response: DummyResponse{
				Data: "a new data not error",
			},
			IsError: false,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(ServerDummy))

	for caseNum, item := range cases {
		c := &Client{
			apiUrl: ts.URL,
			apiKey: item.ApiKey,
		}
		var response DummyResponse

		err := c.Get("/test/get", item.GetParams, &response)
		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(item.Response, response) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Response, response)
		}
	}
	ts.Close()
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ServerDummy))
	c := &Client{
		apiUrl: ts.URL,
		apiKey: "testApiKey",
	}
	var response DummyResponse
	request := &DummyPostRequest{}
	c.Post("/test/post", request, &response)
}
