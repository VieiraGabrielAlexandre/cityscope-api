package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/ibge"
)

type CitiesHandler struct {
	ibge *ibge.Client
}

func NewCitiesHandler(ibgeClient *ibge.Client) *CitiesHandler {
	return &CitiesHandler{ibge: ibgeClient}
}

// GET /v1/cities/{ibge_id}
func (h *CitiesHandler) GetCitySnapshot(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "missing city id in path", http.StatusBadRequest)
		return
	}
	ibgeID := parts[2]

	ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
	defer cancel()

	m, err := h.ibge.GetMunicipality(ctx, ibgeID)
	if err != nil {
		http.Error(w, "failed to fetch city", http.StatusBadGateway)
		return
	}

	// Agregados: população estimada (last)
	pop, popErr := h.ibge.GetPopulationEstimateLast(ctx, ibgeID)

	inds, indsErr := h.ibge.GetUrbanIndicators4714Last(ctx, ibgeID)

	area, areaErr := h.ibge.GetTerritorialArea(ctx, ibgeID)
	if areaErr != nil {
		// fallback: Cidades e Estados (IBGE)
		area, areaErr = h.ibge.GetTerritorialAreaFromCidadesEstados(
			ctx,
			m.Microrregiao.Mesorregiao.UF.Sigla,
			m.Nome,
		)
	}

	data := map[string]any{
		"ibge_id": m.ID,
		"name":    m.Nome,
		"state": map[string]any{
			"sigla": m.Microrregiao.Mesorregiao.UF.Sigla,
			"name":  m.Microrregiao.Mesorregiao.UF.Nome,
			"id":    m.Microrregiao.Mesorregiao.UF.ID,
		},
	}

	if indsErr == nil {
		data["indicators"] = inds
	} else {
		data["indicators"] = map[string]any{
			"available": false,
			"error":     indsErr.Error(),
		}
	}

	// densidade (calculada pelo CityScope)
	if popErr == nil && areaErr == nil && area.ValueKm2 > 0 {
		density := float64(pop.Value) / area.ValueKm2
		data["density_per_km2"] = density
	}

	// Se falhar, não derruba o snapshot — só marca como indisponível.
	if popErr == nil {
		data["population_estimate"] = pop
	} else {
		data["population_estimate"] = map[string]any{
			"available": false,
			"error":     popErr.Error(),
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}
