package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/ibge"
)

type CitiesHandler struct {
	ibge ibge.API
}

func NewCitiesHandler(ibgeClient ibge.API) *CitiesHandler {
	return &CitiesHandler{ibge: ibgeClient}
}

// GET /v1/cities/{ibge_id}
func (h *CitiesHandler) GetCitySnapshot(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		WriteError(w, r, "BAD_REQUEST", "missing city id in path", http.StatusBadRequest)
		return
	}
	ibgeID := parts[2]

	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()

	m, err := h.ibge.GetMunicipality(ctx, ibgeID)
	if err != nil {
		WriteError(w, r, "BAD_GATEWAY", "failed to fetch city", http.StatusBadGateway)
		return
	}

	// Agregados: população estimada (last)
	pop, popErr := h.ibge.GetPopulationEstimateLast(ctx, ibgeID)

	inds, indsErr := h.ibge.GetUrbanIndicators4714Last(ctx, ibgeID)

	resp := CitySnapshotResponse{
		Data: CitySnapshot{
			IBGEID: m.ID,
			Name:   m.Nome,
			State: CityState{
				m.Microrregiao.Mesorregiao.UF.ID,
				m.Microrregiao.Mesorregiao.UF.Nome,
				m.Microrregiao.Mesorregiao.UF.Sigla,
			},
		},
	}

	if popErr == nil {
		resp.Data.PopulationEstimate = pop
	} else {
		resp.Data.PopulationEstimate = AvailabilityError{Available: false, Error: popErr.Error()}
	}

	if indsErr == nil {
		resp.Data.Indicators = inds
	} else {
		resp.Data.Indicators = AvailabilityError{Available: false, Error: indsErr.Error()}
	}

	writeJSON(w, http.StatusOK, resp)
}
