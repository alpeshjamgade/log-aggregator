package handler

import (
	"log-aggregator/internal/constants"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/models"
	"log-aggregator/internal/utils"
	"net/http"
)

func (h *Handler) SaveEvent(w http.ResponseWriter, r *http.Request) {
	ctx := utils.ContextWithValueIfNotPresent(r.Context(), constants.TraceID, utils.GetUUID())
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	req := &models.Event{}
	res := &models.HTTPResponse{Data: map[string]any{}, Status: "success", Message: constants.Empty}

	err := utils.ReadJSON(w, r, req)
	if err != nil {
		Logger.Errorw("error reading request", "error", err)
		res.Status = "error"
		res.Message = "Invalid Request"
		utils.WriteJSON(w, http.StatusBadRequest, res)
		return
	}

	errs := utils.ValidateParams(req)
	if errs != nil {
		res.Status = "error"
		res.Message = errs[0].Error()
		utils.WriteJSON(w, http.StatusBadRequest, res)
		return
	}

	err = h.Service.SaveEvent(ctx, req)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res.Message = "Event saved"
	utils.WriteJSON(w, http.StatusOK, res)
}
