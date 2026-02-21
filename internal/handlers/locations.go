package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/ibge"
)

type LocationsHandler struct {
	ibge ibge.API
}

func NewLocationsHandler(ibgeClient ibge.API) *LocationsHandler {
	return &LocationsHandler{ibge: ibgeClient}
}

func (h *LocationsHandler) ListStates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	states, err := h.ibge.ListStates(ctx)
	if err != nil {
		WriteError(w, r, "BAD_GATEWAY", "failed to fetch states", http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": states,
	})
}

func (h *LocationsHandler) ListMunicipalities(w http.ResponseWriter, r *http.Request) {
	uf := strings.TrimSpace(r.URL.Query().Get("state"))
	if uf == "" {
		WriteError(w, r, "BAD_REQUEST", "missing query param: state (e.g. SP)", http.StatusBadRequest)
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q")) // filtro local

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	munis, err := h.ibge.ListMunicipalitiesByUF(ctx, uf)
	if err != nil {
		WriteError(w, r, "BAD_GATEWAY", "failed to fetch municipalities", http.StatusBadGateway)
		return
	}

	if q != "" {
		qLower := strings.ToLower(q)
		filtered := make([]ibge.Municipality, 0, len(munis))
		for _, m := range munis {
			if strings.Contains(strings.ToLower(m.Nome), qLower) {
				filtered = append(filtered, m)
			}
		}
		munis = filtered
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": munis,
	})
}
