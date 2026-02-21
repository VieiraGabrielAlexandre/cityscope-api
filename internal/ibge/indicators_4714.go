package ibge

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type UrbanIndicators4714 struct {
	ReferenceYear      int     `json:"reference_year"`
	PopulationResident int64   `json:"population_resident"`
	AreaKm2            float64 `json:"area_km2"`
	DensityPerKm2      float64 `json:"density_per_km2"`
}

type agg4714RespItem struct {
	ID         string `json:"id"`
	Variavel   string `json:"variavel"`
	Unidade    string `json:"unidade"`
	Resultados []struct {
		Series []struct {
			Serie map[string]string `json:"serie"`
		} `json:"series"`
	} `json:"resultados"`
}

// 4714 = População Residente, Área territorial e Densidade demográfica
// Variáveis:
// 93   População residente
// 6318 Área da unidade territorial (km²)
// 614  Densidade demográfica
func (c *Client) GetUrbanIndicators4714Last(ctx context.Context, municipalityIBGEID string) (UrbanIndicators4714, error) {
	municipalityIBGEID = strings.TrimSpace(municipalityIBGEID)
	if municipalityIBGEID == "" {
		return UrbanIndicators4714{}, fmt.Errorf("empty municipality id")
	}

	// Uma chamada só com 3 variáveis:
	// /api/v3/agregados/4714/periodos/last/variaveis/93|6318|614?localidades=N6[3550308]
	path := "/v3/agregados/4714/periodos/last/variaveis/93|6318|614"
	q := url.Values{}
	q.Set("localidades", "N6["+municipalityIBGEID+"]")

	var items []agg4714RespItem
	if err := c.getJSON(ctx, path, q, &items); err != nil {
		return UrbanIndicators4714{}, err
	}

	// extrai (ano, valor) de cada item
	type pair struct {
		year int
		val  string
	}

	findLast := func(m map[string]string) (pair, bool) {
		if len(m) == 0 {
			return pair{}, false
		}
		years := make([]int, 0, len(m))
		for k := range m {
			y, err := strconv.Atoi(k)
			if err != nil {
				continue
			}
			years = append(years, y)
		}
		if len(years) == 0 {
			return pair{}, false
		}
		sort.Ints(years)
		last := years[len(years)-1]
		return pair{year: last, val: m[strconv.Itoa(last)]}, true
	}

	var out UrbanIndicators4714
	var havePop, haveArea, haveDen bool

	for _, it := range items {
		if len(it.Resultados) == 0 || len(it.Resultados[0].Series) == 0 {
			continue
		}
		serie := it.Resultados[0].Series[0].Serie
		p, ok := findLast(serie)
		if !ok {
			continue
		}

		// define ano de referência (se vierem diferentes, fica o maior)
		if p.year > out.ReferenceYear {
			out.ReferenceYear = p.year
		}

		raw := strings.TrimSpace(p.val)
		if raw == "" || raw == "..." {
			continue
		}

		switch it.ID {
		case "93": // População residente (inteiro)
			// normalmente vem "203080756"
			rawClean := strings.ReplaceAll(raw, ".", "")
			rawClean = strings.ReplaceAll(rawClean, ",", "")
			v, err := strconv.ParseInt(rawClean, 10, 64)
			if err != nil {
				continue
			}
			out.PopulationResident = v
			havePop = true

		case "6318": // Área km² (float)
			// pode vir "8510417.771" (ponto decimal)
			rawClean := strings.ReplaceAll(raw, ",", ".")
			v, err := strconv.ParseFloat(rawClean, 64)
			if err != nil {
				continue
			}
			out.AreaKm2 = v
			haveArea = true

		case "614": // Densidade (float)
			rawClean := strings.ReplaceAll(raw, ",", ".")
			v, err := strconv.ParseFloat(rawClean, 64)
			if err != nil {
				continue
			}
			out.DensityPerKm2 = v
			haveDen = true
		}
	}

	if !havePop && !haveArea && !haveDen {
		return UrbanIndicators4714{}, fmt.Errorf("no indicators found for municipality %s", municipalityIBGEID)
	}

	// se densidade não veio (ou veio 0), calcula como fallback se tiver pop+area
	if (!haveDen || out.DensityPerKm2 == 0) && havePop && haveArea && out.AreaKm2 > 0 {
		out.DensityPerKm2 = float64(out.PopulationResident) / out.AreaKm2
	}

	return out, nil
}
