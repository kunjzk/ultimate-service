package checkgrp

import (
	"encoding/json"
	"net/http"
	"os"

	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints
type Handlers struct {
	Build string
	Log   *zap.SugaredLogger
}

// we want useful information for our own logs.
// we want to log, we want to connect to a database
// Better to use methods than functions, because we can't change the function signature to pass these things in.
func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}
	statusCode := http.StatusOK
	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("readiness", "ERROR", err)
	}
	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

func (h Handlers) Liveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     h.Build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	statusCode := http.StatusOK
	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("liveness", "ERROR", err)
	}

	h.Log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

func response(w http.ResponseWriter, statusCode int, data interface{}) error {
	// Take our data which is of type any and turn it into json
	// If we can't do that, there's a problem
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set a response header that tells clients this is json data
	// so that clients can handle it properly
	w.Header().Set("Content-Type", "application/json")

	// write status code
	w.WriteHeader(statusCode)

	// Send the data back to the client by calling Write. Write expects json data
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
