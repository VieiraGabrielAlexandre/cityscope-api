package ibge

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// Extrai algo como: "1.521,202 km²"
var reArea = regexp.MustCompile(`Área Territorial\s+([\d\.\,]+)\s+km²`)

func (c *Client) GetTerritorialAreaFromCidadesEstados(ctx context.Context, ufSigla, cityName string) (TerritorialArea, error) {
	ufSigla = strings.ToLower(strings.TrimSpace(ufSigla))
	if ufSigla == "" || cityName == "" {
		return TerritorialArea{}, fmt.Errorf("missing ufSigla or cityName")
	}

	slug := slugify(cityName)
	// Ex: https://www.ibge.gov.br/cidades-e-estados/sp/sao-paulo.html
	u := fmt.Sprintf("https://www.ibge.gov.br/cidades-e-estados/%s/%s.html", ufSigla, slug)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return TerritorialArea{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return TerritorialArea{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return TerritorialArea{}, fmt.Errorf("cidades-e-estados non-2xx: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TerritorialArea{}, err
	}

	m := reArea.FindStringSubmatch(string(body))
	if len(m) < 2 {
		return TerritorialArea{}, fmt.Errorf("area not found in cidades-e-estados html")
	}

	// "1.521,202" -> "1521.202"
	raw := strings.ReplaceAll(m[1], ".", "")
	raw = strings.ReplaceAll(raw, ",", ".")
	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return TerritorialArea{}, err
	}

	return TerritorialArea{ValueKm2: val}, nil
}

func slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))

	// remove acentos (NFD) e diacríticos
	t := norm.NFD.String(s)
	b := strings.Builder{}
	b.Grow(len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		// mantém letras/números e converte espaços/pontuação para "-"
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			continue
		}
		if r == ' ' || r == '_' || r == '-' {
			b.WriteByte('-')
		}
	}

	out := b.String()
	out = strings.Trim(out, "-")
	out = regexp.MustCompile(`-+`).ReplaceAllString(out, "-")
	return out
}
