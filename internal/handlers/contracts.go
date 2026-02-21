package handlers

import "github.com/VieiraGabrielAlexandre/cityscope-api/internal/ibge"

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
