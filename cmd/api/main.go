package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/config"
	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/handlers"
	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/httpserver"
	"github.com/VieiraGabrielAlexandre/cityscope-api/internal/ibge"
)

func main() {
	// Carrega .env se existir (opcional e simples, sem libs):
	loadDotEnvIfPresent(".env")

	cfg := config.Load()

	ibgeClient := ibge.NewClient(cfg.IBGEBaseURL, time.Duration(cfg.IBGETimeoutSecond)*time.Second)

	cached := ibge.NewCachedClient(ibgeClient, time.Duration(cfg.IBGECacheTTLSeconds)*time.Second)

	health := handlers.NewHealthHandler()
	locations := handlers.NewLocationsHandler(cached)
	cities := handlers.NewCitiesHandler(cached)

	router := httpserver.NewRouter(httpserver.RouterDeps{
		APIToken:  cfg.APIToken,
		Health:    health,
		Locations: locations,
		Cities:    cities,
	})

	addr := ":" + cfg.Port
	log.Printf("CityScope API running on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

// Carregador simples de .env (KEY=VALUE por linha). Sem dependência externa.
func loadDotEnvIfPresent(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := splitLines(string(b))
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || startsWith(line, "#") {
			continue
		}
		k, v, ok := splitKV(line)
		if !ok {
			continue
		}
		// não sobrescreve env já setada
		if os.Getenv(k) == "" {
			_ = os.Setenv(k, v)
		}
	}
}

// Helpers (pra não puxar libs)
func splitLines(s string) []string {
	out := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start <= len(s)-1 {
		out = append(out, s[start:])
	}
	return out
}

func trimSpace(s string) string {
	// minimalista
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

func startsWith(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	return s[:len(prefix)] == prefix
}

func splitKV(line string) (string, string, bool) {
	// KEY=VALUE (não lida com quotes complexos, mas resolve 95% dos .env)
	for i := 0; i < len(line); i++ {
		if line[i] == '=' {
			k := trimSpace(line[:i])
			v := trimSpace(line[i+1:])
			if k == "" {
				return "", "", false
			}
			return k, v, true
		}
	}
	return "", "", false
}
