package ibge

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type PopulationEstimate struct {
	Year  int   `json:"year"`
	Value int64 `json:"value"`
}

// A API de agregados com view=flat retorna uma lista de mapas (linhas) com chaves tipo:
// V (valor), D1C (localidade), D2C (ano), etc.
type flatRow map[string]string

func (c *Client) GetPopulationEstimateLast(ctx context.Context, municipalityIBGEID string) (PopulationEstimate, error) {
	municipalityIBGEID = strings.TrimSpace(municipalityIBGEID)
	if municipalityIBGEID == "" {
		return PopulationEstimate{}, fmt.Errorf("empty municipality id")
	}

	// /api/v3/agregados/6579/periodos/last/variaveis/9324?localidades=N6[3550308]&view=flat
	path := "/v3/agregados/6579/periodos/last/variaveis/9324"
	q := url.Values{}
	q.Set("localidades", "N6["+municipalityIBGEID+"]")
	q.Set("view", "flat")

	var rows []flatRow
	if err := c.getJSON(ctx, path, q, &rows); err != nil {
		return PopulationEstimate{}, err
	}

	// Normalmente vem 1 linha quando periodos=last e 1 localidade.
	for _, r := range rows {
		v := strings.TrimSpace(r["V"])
		yearStr := strings.TrimSpace(r["D2C"]) // ano (código)
		if v == "" || v == "..." || yearStr == "" {
			continue
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			continue
		}

		// alguns valores podem vir com separador; a maioria vem "123456"
		vClean := strings.ReplaceAll(v, ".", "")
		vClean = strings.ReplaceAll(vClean, ",", "")
		val, err := strconv.ParseInt(vClean, 10, 64)
		if err != nil {
			continue
		}

		return PopulationEstimate{Year: year, Value: val}, nil
	}

	return PopulationEstimate{}, fmt.Errorf("no population estimate found for %s", municipalityIBGEID)
}
