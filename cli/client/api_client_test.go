package client

import (
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
	"net/http"
	"testing"
)

func TestExecuteErrors(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		se := models.ServerError{Message: "msg", ErrorCode: 1}
		MarshalAndWrite(t, w, se, 500)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	err := client.Execute(client.Sling("").Get(""), nil)
	if err == nil {
		t.Fatalf("Error was nil!")
	}

	se, ok := err.(*ServerError)
	if !ok {
		t.Fatalf("Error was not of type *ServerError")
	}

	testutils.AssertEqual(t, se.Message, "msg")
	testutils.AssertEqual(t, se.ErrorCode, int64(1))
}

func TestExecuteWithJob(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		headers := map[string]string{
			"Location": "/job/jobid",
			"X-JobID":  "jobid",
		}

		MarshalAndWriteHeader(t, w, "", headers, 202)
	}

	client, server := newClientAndServer(handler)
	defer server.Close()

	jobID, err := client.ExecuteWithJob(client.Sling("").Get(""))
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertEqual(t, jobID, "jobid")
}

func TestExecuteWithJobErrors(t *testing.T) {
	cases := map[string]Handler{
		"Should error if no headers returned": func(w http.ResponseWriter, r *http.Request) {
			MarshalAndWrite(t, w, "", 202)
		},
		"Should error if non-202 status returned": func(w http.ResponseWriter, r *http.Request) {
			headers := map[string]string{
				"Location": "/job/jobid",
				"X-JobID":  "jobid",
			}

			MarshalAndWriteHeader(t, w, "", headers, 200)
		},
	}

	for name, handler := range cases {
		client, server := newClientAndServer(handler)
		defer server.Close()

		if _, err := client.ExecuteWithJob(client.Sling("").Get("")); err == nil {
			t.Fatalf("%s: Error was nil!", name)
		}
	}
}
