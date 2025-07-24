package handler

import (
	"encoding/json"
	"log-aggregator/internal/constants"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/models"
	"log-aggregator/internal/utils"
	"net/http"
)

func (h *Handler) SaveLog(w http.ResponseWriter, r *http.Request) {
	ctx := utils.ContextWithValueIfNotPresent(r.Context(), constants.TraceID, utils.GetUUID())
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	req := &models.FluentBitReq{}
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

	var logDecoded models.RawLog
	if err := json.Unmarshal([]byte(req.Log), &logDecoded); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, res)
		return
	}

	req.LogDecoded = logDecoded

	err = h.Service.SaveLog(ctx, req)
	if err != nil {
		Logger.Errorw("error saving log", "error", err)
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res.Message = "Log saved"
	utils.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) SaveBulkLog(w http.ResponseWriter, r *http.Request) {
	ctx := utils.ContextWithValueIfNotPresent(r.Context(), constants.TraceID, utils.GetUUID())
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	var req []*models.FluentBitReq
	res := &models.HTTPResponse{Data: map[string]any{}, Status: "success", Message: constants.Empty}

	err := utils.ReadJSON(w, r, &req)
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

	for _, fluentBitLog := range req {

		cleanLog := utils.ExtractJSONFromLog(fluentBitLog.Log)

		var rawLog models.RawLog
		if err := json.Unmarshal([]byte(cleanLog), &rawLog); err != nil {
			res.Message = errs[0].Error()
			utils.WriteJSON(w, http.StatusBadRequest, res)
			return
		}

		fluentBitLog.LogDecoded = rawLog
	}

	err = h.Service.SaveBulkLog(ctx, req)
	if err != nil {
		Logger.Errorw("error saving bulk log", "error", err)
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res.Message = "Logs saved"
	utils.WriteJSON(w, http.StatusOK, res)

}

// SaveBulkLogV2 flattens received log and saves into db
func (h *Handler) SaveBulkLogV2(w http.ResponseWriter, r *http.Request) {
	ctx := utils.ContextWithValueIfNotPresent(r.Context(), constants.TraceID, utils.GetUUID())
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	var req []map[string]any
	res := &models.HTTPResponse{Data: map[string]any{}, Status: "success", Message: constants.Empty}

	err := utils.ReadJSON(w, r, &req)
	if err != nil {
		Logger.Errorw("error reading request", "error", err)
		res.Status = "error"
		res.Message = "Invalid Request"
		utils.WriteJSON(w, http.StatusBadRequest, res)
		return
	}

	var reqFlattened []map[string]any
	for _, record := range req {
		recordFlattened := utils.FlattenMap("", record)

		if recordFlattened["log"] != nil {
			jsonEncodedString := utils.ExtractJSONFromLog(recordFlattened["log"].(string))
			var logDecoded map[string]any
			if err := json.Unmarshal([]byte(jsonEncodedString), &logDecoded); err != nil {
				res.Message = err.Error()
				utils.WriteJSON(w, http.StatusBadRequest, res)
				return
			}

			recordFlattened["log"] = utils.FlattenMap("", logDecoded)
		}

		reqFlattened = append(reqFlattened, recordFlattened)

	}

	err = h.Service.SaveBulkLogV2(ctx, reqFlattened)
	if err != nil {
		Logger.Errorw("error saving bulk log", "error", err)
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	res.Message = "Logs saved"
	res.Data = reqFlattened
	utils.WriteJSON(w, http.StatusOK, res)

}
