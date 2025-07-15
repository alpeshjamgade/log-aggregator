package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log-aggregator/internal/constants"
	"log-aggregator/internal/logger"
	"log-aggregator/internal/models"
	"net/http"
	"net/url"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

func WriteUnimplemented(w http.ResponseWriter, message ...string) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	var errMessage string
	if len(message) > 0 {
		errMessage = message[0]
	} else {
		errMessage = "API unimplemented!!"
	}

	_, err := w.Write([]byte(errMessage))

	if err != nil {
		return err
	}

	return nil
}

func WriteInvalidSession(w http.ResponseWriter, res ...*models.HTTPResponse) error {

	defaultRes := &models.HTTPResponse{Data: map[string]any{}, Status: "error", Message: constants.Empty}
	if len(res) != 0 {
		return WriteJSON(w, http.StatusForbidden, defaultRes)
	}

	return WriteJSON(w, http.StatusForbidden, res)

}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload models.HTTPResponse
	payload.Error = true
	payload.Status = "error"
	payload.Message = err.Error()
	payload.Data = map[string]any{}

	return WriteJSON(w, statusCode, payload)
}

func BuildURL(parsedURL *url.URL, path string, queryParams map[string]string) string {
	u := &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
		Path:   path,
	}

	if len(queryParams) > 0 {
		qp := url.Values{}
		for k, v := range queryParams {
			qp.Set(k, v)
		}
		u.RawQuery = qp.Encode()
	}

	return u.String()
}

func MakeHttpRequest(ctx context.Context, client *http.Client, method, url string, reqBody any, headers map[string]string) (int, json.RawMessage, error) {
	Logger := logger.CreateFileLoggerWithCtx(ctx)

	//create request body
	var (
		reqBodyBytes []byte
		err          error
	)
	if reqBody != nil {
		reqBodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			Logger.Errorw("error json marshalling request body", "err", err)
			return http.StatusInternalServerError, nil, err
		}
	}

	//create the http request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqBodyBytes))
	if err != nil {
		Logger.Errorw("error creating http request", "err", err)
		return http.StatusInternalServerError, nil, err
	}
	if len(headers) != 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	//call the api
	resp, err := client.Do(req)
	if err != nil {
		Logger.Errorw("could not place http request", "err", err)
		return http.StatusInternalServerError, nil, err
	}
	defer resp.Body.Close()

	var respBodyBytes json.RawMessage
	respBodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		Logger.Errorw("error reading response body", "err", err)
		return http.StatusInternalServerError, nil, err
	}
	return resp.StatusCode, respBodyBytes, nil
}
