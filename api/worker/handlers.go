package worker

import (
	"encoding/base64"
	"encoding/json"
	"github.com/ufcg-lsd/arrebol-pb/api"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker"
	"github.com/ufcg-lsd/arrebol-pb/arrebol/worker/auth/token"
	"net/http"
)

const SignatureHeader string = "SIGNATURE";

func (a *WorkerApi) AddWorker(w http.ResponseWriter, r *http.Request) {
	var err error
	signatureEncoded := r.Header.Get(SignatureHeader)
	signature, _ := base64.StdEncoding.DecodeString(signatureEncoded)


	if string(signature) == "" {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "signature header was not found",
			Status:  http.StatusBadRequest,
		})
		return
	}

	var worker *worker.Worker
	if err = json.NewDecoder(r.Body).Decode(&worker); err != nil {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "Maybe the body has a wrong shape",
			Status:  http.StatusBadRequest,
		})
		return
	}

	var tempToken *token.Token
	if tempToken, err = a.auth.Authenticate(signature, worker); err != nil {
		api.Write(w, http.StatusUnauthorized, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusUnauthorized,
		})
		return
	}

	var queueId uint
	if queueId, err = a.manager.Join(*worker); err != nil {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "Maybe the body has a wrong shape",
			Status:  http.StatusBadRequest,
		})
		return
	}

	token, err := (*tempToken).SetPayloadField("QueueID", queueId)
	queue, err := a.storage.RetrieveQueue(queueId)
	queue.Workers = append(queue.Workers, worker)
	err = a.storage.SaveQueue(queue)
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	}

	if err != nil {
		api.Write(w, http.StatusBadRequest, api.ErrorResponse{
			Message: "Maybe the body has a wrong shape",
			Status:  http.StatusBadRequest,
		})
		return
	}

	w.Header().Set("Token", token.String())
	api.Write(w, http.StatusOK, nil)
}

