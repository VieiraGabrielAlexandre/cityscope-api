package ibge

import (
	"context"
	"net/url"
	"strings"
)

type State struct {
	ID    int    `json:"id"`
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
}

type Municipality struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`

	Microrregiao struct {
		Mesorregiao struct {
			UF struct {
				ID    int    `json:"id"`
				Sigla string `json:"sigla"`
				Nome  string `json:"nome"`
			} `json:"UF"`
		} `json:"mesorregiao"`
	} `json:"microrregiao"`
}

func (c *Client) ListStates(ctx context.Context) ([]State, error) {
	var out []State
	// docs: /api/v1/localidades/estados
	err := c.getJSON(ctx, "/v1/localidades/estados", nil, &out)
	return out, err
}

func (c *Client) ListMunicipalitiesByUF(ctx context.Context, uf string) ([]Municipality, error) {
	uf = strings.TrimSpace(strings.ToUpper(uf))
	var out []Municipality
	// docs: /api/v1/localidades/estados/{UF}/municipios
	err := c.getJSON(ctx, "/v1/localidades/estados/"+url.PathEscape(uf)+"/municipios", nil, &out)
	return out, err
}

func (c *Client) GetMunicipality(ctx context.Context, ibgeID string) (Municipality, error) {
	var out Municipality
	// docs: /api/v1/localidades/municipios/{id}
	err := c.getJSON(ctx, "/v1/localidades/municipios/"+url.PathEscape(strings.TrimSpace(ibgeID)), nil, &out)
	return out, err
}
