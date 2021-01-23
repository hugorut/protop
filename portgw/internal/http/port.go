package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/hugorut/protop/portgw/internal/processor"
)

// ProcessorProvider is an interface that fetches a given
// FileProcessor based on a lookup string.
type ProcessorProvider interface {
	Get(string) (processor.FileProcessor, error)
}

// errorResponse is a simple struct to marshal a json error message to the client.
type errorResponse struct {
	Error string `json:"error"`
}

// Handler is struct that coordinate http interaction within the portgw.
// It holds reusable fields across handler methods.
type Handler struct {
	Logger logrus.FieldLogger

	ProcessorProvider ProcessorProvider
}

// processFileRequest represents a json struct used to call the
// process file endpoint.
type processFileRequest struct {
	Location string `json:"location"`
}

// processFileResponse represents a json struct used to return information
// from the process file endpoint.
type processFileResponse struct {
	ID int `json:"id"`
}

// ProcessFile ingests a request to initiate a port file consumption.
// It returns an id for the started process so that the status
// of the consumption can be checked in further requests.
//
// TODO: the returned id for the request simply returns a static int and there
// is no logic to check the status of the file consumption.
func (h Handler) ProcessFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	provider := vars["provider"]
	processor, err := h.ProcessorProvider.Get(provider)
	if err != nil {
		msg := fmt.Sprintf("invalid provider %s given", provider)
		h.writeError(w, msg)
		return
	}

	var req processFileRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.writeError(w, "malformed request given")
		return
	}

	id, err := processor.Process(req.Location)
	if err != nil {
		h.Logger.Errorf("process file error, %s", err)
		h.writeError(w, "unable to start processing given file")
		return
	}

	err = json.NewEncoder(w).Encode(processFileResponse{ID: id})
	if err != nil {
		h.Logger.Errorf("could not write process file response to client, %s", err)
	}
}

func (h Handler) writeError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(errorResponse{Error: msg})

	if err != nil {
		h.Logger.Errorf("could not write error to client, %s", err)
	}
}
