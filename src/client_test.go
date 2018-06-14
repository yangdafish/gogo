package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestWhisperMessageHandlerShouldFailIfNoParameters(t *testing.T) {
	req, err := http.NewRequest("GET", "/whisper", nil)
    if err != nil {
        t.Fatal(err)
	}

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(WhisperMessageHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
	}
}

func TestRecieveMessageHandlerShouldFailIfNoParameters(t *testing.T) {
	req, err := http.NewRequest("GET", "/recieve", nil)
    if err != nil {
        t.Fatal(err)
	}

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(RecieveMessageHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
	}
}

func TestConnectNodeHandlerShouldFailIfNoParameters(t *testing.T) {
	req, err := http.NewRequest("GET", "/connect", nil)
    if err != nil {
        t.Fatal(err)
	}

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(ConnectNodeHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
	}
}

func TestBuildHttpWhisperRequest(t *testing.T) {
	testSender := "testSender"
	testMessage := "testMessage"
	testTargetUrl := "testTargetUrl"
	whisperRoutePath := "/recieve"

	req := BuildHttpWhisperRequest(testTargetUrl, testSender, testMessage)

	requestNameParam := req.URL.Query().Get("name")
	if requestNameParam != testSender {
		t.Errorf("Request 'name' parameter was invalid: Got %s want %s",
			requestNameParam, testSender)
	}

	requestMessageParam := req.URL.Query().Get("message")
	if requestMessageParam != testMessage {
		t.Errorf("Request 'message' parameter was invalid: Got %s want %s",
			requestMessageParam, testMessage)
	}

	requestUrl := req.URL.Host
	if requestUrl != testTargetUrl {
		t.Errorf("Request 'url' was invalid: Got %s want %s",
			requestUrl, testTargetUrl)
	}

	requestPath := req.URL.Path
	if requestPath != whisperRoutePath {
		t.Errorf("Request 'path' was invalid: Got %s want %s",
			requestUrl, whisperRoutePath)
	}
}

func TestBuildHttpConnectNodeRequest(t *testing.T) {
	testTargetUrl := "testTargetUrl"
	connectRoutePath := "/connect"

	req := BuildHttpConnectNodeRequest(testTargetUrl)

	requestUrl := req.URL.Host
	if requestUrl != testTargetUrl {
		t.Errorf("Request 'url' was invalid: Got %s want %s",
			requestUrl, testTargetUrl)
	}

	requestPath := req.URL.Path
	if requestPath != connectRoutePath {
		t.Errorf("Request 'path' was invalid: Got %s want %s",
			requestUrl, connectRoutePath)
	}
}
