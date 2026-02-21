package httpserver

import (
	"net/http"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/handlers"
)

type RouterDeps struct {
	APIToken string

	Health    *handlers.HealthHandler
	Locations *handlers.LocationsHandler
	Cities    *handlers.CitiesHandler
}

func NewRouter(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	// Público
	mux.HandleFunc("/health", deps.Health.Handle)
	mux.HandleFunc("/docs", DocsUIHandler)
	mux.HandleFunc("/openapi.json", OpenAPIJSONHandler)

	// Protegido
	auth := AuthMiddleware(deps.APIToken)
	protected := http.NewServeMux()

	protected.HandleFunc("/v1/locations/states", deps.Locations.ListStates)
	protected.HandleFunc("/v1/locations/municipalities", deps.Locations.ListMunicipalities)
	protected.HandleFunc("/v1/cities/", deps.Cities.GetCitySnapshot) // /v1/cities/{ibge_id}

	// Encapsula o protected no middleware
	mux.Handle("/v1/", auth(protected))

	return mux
}
