package testgrp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string
	}{
		Status: "Test OK",
	}
	json.NewEncoder(w).Encode(status)
	statusCode := http.StatusOK
	h.Log.Infow("test", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}
