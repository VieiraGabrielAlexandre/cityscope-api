package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/ibge"
)

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, id)
}

func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(ctxKeyRequestID).(string)
	return id
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func WriteError(w http.ResponseWriter, r *http.Request, code string, message string, status int) {
	requestID, _ := r.Context().Value(ctxKeyRequestID).(string)

	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	}

	writeJSON(w, status, resp)
}

type ctxKey string

const ctxKeyRequestID ctxKey = "request_id"

type AvailabilityError struct {
	Available bool   `json:"available"`
	Error     string `json:"error,omitempty"`
}

type CitySnapshotResponse struct {
	Data CitySnapshot `json:"data"`
}

type CitySnapshot struct {
	IBGEID int       `json:"ibge_id"`
	Name   string    `json:"name"`
	State  CityState `json:"state"`
	// Um dos dois: ou valor, ou erro
	PopulationEstimate any `json:"population_estimate,omitempty"`
	Indicators         any `json:"indicators,omitempty"`
}

type CityState struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Sigla string `json:"sigla"`
}

// Se você quiser ainda mais tipado:
// PopulationEstimateResult { Data *ibge.PopulationEstimate; Error *AvailabilityError }
// IndicatorsResult { Data *ibge.UrbanIndicators4714; Error *AvailabilityError }
// (eu te passo essa versão quando você re-upar os arquivos)
var _ = ibge.PopulationEstimate{}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
