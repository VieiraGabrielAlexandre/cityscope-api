

# `CityScope API 🇧🇷`

API em **Go** que fornece um *snapshot urbano* de cidades brasileiras usando dados oficiais do IBGE.

O objetivo do projeto é transformar o enorme volume de dados estatísticos públicos do IBGE em uma API simples de consumir para aplicações, dashboards, chatbots e automações.

> Em vez de lidar com o modelo complexo do SIDRA/IBGE, o CityScope entrega informações prontas sobre uma cidade.

---

## ✨ O que a API já faz

Atualmente o CityScope:

* Lista estados brasileiros
* Lista municípios por UF
* Busca dados de uma cidade pelo código IBGE
* Retorna um snapshot padronizado
* Obtém população estimada oficial (IBGE – SIDRA Agregados)
* Protege endpoints com autenticação por token

---

## 📊 Exemplo de resposta

### `GET /v1/cities/3550308`

```json
{
  "data": {
    "ibge_id": 3550308,
    "name": "São Paulo",
    "state": {
      "sigla": "SP",
      "name": "São Paulo",
      "id": 35
    },
    "population_estimate": {
      "year": 2024,
      "value": 12345678
    }
  }
}
```

---

## 🔐 Autenticação

Todos os endpoints `/v1/*` usam **Bearer Token**.

Header:

```
Authorization: Bearer SEU_TOKEN
```

O token é definido no `.env`.

---

## ⚙️ Configuração

Crie um arquivo `.env`:

```env
PORT=8080
API_TOKEN=changeme-super-secret
IBGE_BASE_URL=https://servicodados.ibge.gov.br/api
IBGE_TIMEOUT_SECONDS=12
```

---

## ▶️ Executando localmente

```bash
go run ./cmd/api
```

Servidor:

```
http://localhost:8080
```

---

## 🔎 Testando

Healthcheck:

```bash
curl http://localhost:8080/health
```

Listar estados:

```bash
curl -H "Authorization: Bearer $TOKEN" \
http://localhost:8080/v1/locations/states
```

Buscar municípios:

```bash
curl -H "Authorization: Bearer $TOKEN" \
"http://localhost:8080/v1/locations/municipalities?state=SP&q=camp"
```

Snapshot da cidade:

```bash
curl -H "Authorization: Bearer $TOKEN" \
http://localhost:8080/v1/cities/3550308
```

---

## 🧠 Como funciona (resumo técnico)

O CityScope consome duas partes da API do IBGE:

### Localidades

Divisões administrativas oficiais:

* estados
* municípios

### Agregados (SIDRA)

Tabelas estatísticas do IBGE.

Exemplo usado:

| Agregado | Variável | Descrição                    |
| -------- | -------- | ---------------------------- |
| 6579     | 9324     | População residente estimada |

Isso equivale a:

> “População oficial estimada do município no último ano disponível”

---

## 🏗️ Estrutura do projeto

```
cmd/api
internal/config
internal/httpserver
internal/handlers
internal/ibge
```

* `ibge` → client HTTP e integração
* `handlers` → endpoints REST
* `httpserver` → router e middleware de auth
* `config` → carregamento do .env

---

## 📌 Próximos passos

Planejados:

* área territorial (km²)
* densidade demográfica
* PIB municipal
* cache Redis
* documentação OpenAPI (Swagger)
* ranking de cidades

---

## 🎯 Objetivo do projeto

Criar uma API pública e simples para responder perguntas como:

* “Qual cidade é maior?”
* “Onde abrir um negócio?”
* “Qual município cresce mais?”
* “Qual a densidade populacional?”

Usando **dados oficiais do Brasil**, mas com ergonomia de API moderna.

---

## Licença

MIT
