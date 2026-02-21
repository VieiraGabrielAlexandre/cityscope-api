package ibge

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type TerritorialArea struct {
	ValueKm2 float64 `json:"value_km2"`
}

func (c *Client) GetTerritorialArea(ctx context.Context, municipalityIBGEID string) (TerritorialArea, error) {
	municipalityIBGEID = strings.TrimSpace(municipalityIBGEID)
	if municipalityIBGEID == "" {
		return TerritorialArea{}, fmt.Errorf("empty municipality id")
	}

	// agregado 29171 variável 214
	path := "/v3/agregados/29171/periodos/last/variaveis/214"
	q := url.Values{}
	q.Set("localidades", "N6["+municipalityIBGEID+"]")
	q.Set("view", "flat")

	var rows []map[string]string
	if err := c.getJSON(ctx, path, q, &rows); err != nil {
		return TerritorialArea{}, err
	}

	for _, r := range rows {
		v := strings.TrimSpace(r["V"])
		if v == "" || v == "..." {
			continue
		}

		// área vem com vírgula decimal
		v = strings.ReplaceAll(v, ".", "")
		v = strings.ReplaceAll(v, ",", ".")

		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}

		return TerritorialArea{ValueKm2: val}, nil
	}

	return TerritorialArea{}, fmt.Errorf("territorial area not found")
}
