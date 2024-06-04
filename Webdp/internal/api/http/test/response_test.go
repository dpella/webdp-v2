package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"webdp/internal/api/http/response"
)

func TestHttpResponse(t *testing.T) {
	var temp response.HttpResponse[string]
	temp = mockServiceFunction(10)
	js, err := json.Marshal(temp)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("json = %s\n", string(js))

	temp = mockServiceFunction(-10)
	js, err = json.Marshal(temp)

	if err != nil {
		t.Error(err)
	}
	fmt.Printf("json = %s\n", string(js))

	temp = mockServiceFunction(0)
	js, err = json.Marshal(temp)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("json = %s\n", string(js))
}

func TestHttpNoContent(t *testing.T) {
	resp := mockServiceNoContent("michael", "1")
	if resp.GetStatusCode() == 204 {
		t.Error("didnt expect no content")
	}

	js, err := json.Marshal(resp)

	if err != nil {
		t.Error(err)
	}

	if string(js) == "" {
		t.Error("didnt expect no content")
	}

	resp = mockServiceNoContent("adam", "loooong pwd")
	if resp.GetStatusCode() == 204 {
		t.Error("didnt expect no content")
	}

	js, err = json.Marshal(resp)

	if err != nil {
		t.Error(err)
	}

	if string(js) == "" {
		t.Error("didn expect no content")
	}

	resp = mockServiceNoContent("success", "1")

	if resp.GetStatusCode() != 204 {
		t.Errorf("expected no content but got: %d", resp.GetStatusCode())
	}

	_, err = json.Marshal(resp)

	if err != nil {
		t.Error(err)
	}

	resp = mockServiceNoContent("hubbabubba", "1")

	if resp.GetStatusCode() != 204 {
		t.Errorf("expected no content but got: %d", resp.GetStatusCode())
	}

	_, err = json.Marshal(resp)

	if err != nil {
		t.Error(err)
	}

}

func mockServiceFunction(input int) response.HttpResponse[string] {
	if input > 0 {
		return response.NewFail[string](http.StatusBadRequest).
			WithDetail("that number is too darn high!").
			WithTitle("yikes, that's one high number!").
			WithType("Huge Jacked Man Number")
	} else if input < 0 {
		return response.NewSuccess(http.StatusAccepted, "nice low number")
	} else {
		return response.NewFail[string](418).
			WithDetail("but I am a teapot?")
	}
}

func mockServiceNoContent(username string, pwd string) response.HttpResponse[response.Void] {
	if username == "michael" {
		return response.NewFail[response.Void](http.StatusBadRequest).
			WithDetail("you cannot be called michael!")
	} else if len(pwd) > 3 {
		return response.NewFail[response.Void](http.StatusBadRequest).
			WithDetail("too long to read").
			WithTitle("gigantic password!")

	} else {
		return response.NoContent()
	}
}
