package ibge

import "context"

type API interface {
	ListStates(ctx context.Context) ([]State, error)
	ListMunicipalitiesByUF(ctx context.Context, uf string) ([]Municipality, error)
	GetMunicipality(ctx context.Context, ibgeID string) (Municipality, error)
	GetPopulationEstimateLast(ctx context.Context, municipalityIBGEID string) (PopulationEstimate, error)
	GetUrbanIndicators4714Last(ctx context.Context, municipalityIBGEID string) (UrbanIndicators4714, error)
}
