package ibge

import (
	"context"
	"fmt"
	"time"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/cache"
)

type CachedClient struct {
	inner API
	ttl   time.Duration

	statesCache *cache.TTLCache[[]State]
	munisCache  *cache.TTLCache[[]Municipality]
	muniCache   *cache.TTLCache[Municipality]
	popCache    *cache.TTLCache[PopulationEstimate]
	indCache    *cache.TTLCache[UrbanIndicators4714]
}

func NewCachedClient(inner API, ttl time.Duration) *CachedClient {
	return &CachedClient{
		inner:       inner,
		ttl:         ttl,
		statesCache: cache.NewTTLCache[[]State](),
		munisCache:  cache.NewTTLCache[[]Municipality](),
		muniCache:   cache.NewTTLCache[Municipality](),
		popCache:    cache.NewTTLCache[PopulationEstimate](),
		indCache:    cache.NewTTLCache[UrbanIndicators4714](),
	}
}

func (c *CachedClient) ListStates(ctx context.Context) ([]State, error) {
	key := "states"
	return c.statesCache.GetOrSet(key, c.ttl, func() ([]State, error) {
		return c.inner.ListStates(ctx)
	})
}

func (c *CachedClient) ListMunicipalitiesByUF(ctx context.Context, uf string) ([]Municipality, error) {
	key := fmt.Sprintf("munis:%s", uf)
	return c.munisCache.GetOrSet(key, c.ttl, func() ([]Municipality, error) {
		return c.inner.ListMunicipalitiesByUF(ctx, uf)
	})
}

func (c *CachedClient) GetMunicipality(ctx context.Context, ibgeID string) (Municipality, error) {
	key := fmt.Sprintf("muni:%s", ibgeID)
	return c.muniCache.GetOrSet(key, c.ttl, func() (Municipality, error) {
		return c.inner.GetMunicipality(ctx, ibgeID)
	})
}

func (c *CachedClient) GetPopulationEstimateLast(ctx context.Context, municipalityIBGEID string) (PopulationEstimate, error) {
	key := fmt.Sprintf("pop:last:%s", municipalityIBGEID)
	return c.popCache.GetOrSet(key, c.ttl, func() (PopulationEstimate, error) {
		return c.inner.GetPopulationEstimateLast(ctx, municipalityIBGEID)
	})
}

func (c *CachedClient) GetUrbanIndicators4714Last(ctx context.Context, municipalityIBGEID string) (UrbanIndicators4714, error) {
	key := fmt.Sprintf("ind4714:last:%s", municipalityIBGEID)
	return c.indCache.GetOrSet(key, c.ttl, func() (UrbanIndicators4714, error) {
		return c.inner.GetUrbanIndicators4714Last(ctx, municipalityIBGEID)
	})
}
