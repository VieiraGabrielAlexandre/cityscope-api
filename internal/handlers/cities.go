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
	// extrai id do path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "missing city id in path", http.StatusBadRequest)
		return
	}
	ibgeID := parts[2]

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	m, err := h.ibge.GetMunicipality(ctx, ibgeID)
	if err != nil {
		http.Error(w, "failed to fetch city", http.StatusBadGateway)
		return
	}

	// Snapshot inicial (localidades). Depois a gente pluga Agregados aqui.
	resp := map[string]any{
		"data": map[string]any{
			"ibge_id": m.ID,
			"name":    m.Nome,
			"state": map[string]any{
				"sigla": m.Microrregiao.Mesorregiao.UF.Sigla,
				"name":  m.Microrregiao.Mesorregiao.UF.Nome,
				"id":    m.Microrregiao.Mesorregiao.UF.ID,
			},
		},
	}

	writeJSON(w, http.StatusOK, resp)
}
