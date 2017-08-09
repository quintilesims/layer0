package controllers

import (
	"fmt"
	"net/http"

	"github.com/zpatrick/fireball"
)

func newJobResponse(jobID string) fireball.ResponseFunc {
	return fireball.ResponseFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", fmt.Sprintf("/job/%s", jobID))
		w.Header().Set("X-JobID", jobID)
		w.WriteHeader(http.StatusAccepted)
		w.Write(nil)
	})
}
