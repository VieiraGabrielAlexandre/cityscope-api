package ibge

import (
	"sort"
	"strconv"
	"strings"
)

func parseYear(raw string) (int, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "..." {
		return 0, false
	}
	y, err := strconv.Atoi(raw)
	return y, err == nil
}

// Para série do tipo {"2010":"...", "2022":"..."} -> pega o ano mais recente
func lastYear(serie map[string]string) (year int, val string, ok bool) {
	if len(serie) == 0 {
		return 0, "", false
	}
	years := make([]int, 0, len(serie))
	for k := range serie {
		y, err := strconv.Atoi(k)
		if err == nil {
			years = append(years, y)
		}
	}
	if len(years) == 0 {
		return 0, "", false
	}
	sort.Ints(years)
	year = years[len(years)-1]
	return year, serie[strconv.Itoa(year)], true
}

func parseInt64BR(raw string) (int64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "..." {
		return 0, false
	}
	raw = strings.ReplaceAll(raw, ".", "")
	raw = strings.ReplaceAll(raw, ",", "")
	v, err := strconv.ParseInt(raw, 10, 64)
	return v, err == nil
}

// Float "flex": aceita "1.23" e "1,23" e também com milhares "1.234,56"
func parseFloat64Flex(raw string) (float64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "..." {
		return 0, false
	}

	// Caso BR com milhares e vírgula decimal: "1.234,56"
	if strings.Contains(raw, ",") {
		raw = strings.ReplaceAll(raw, ".", "")
		raw = strings.ReplaceAll(raw, ",", ".")
	}

	// Caso US/normal: "8510417.771" (já ok)
	v, err := strconv.ParseFloat(raw, 64)
	return v, err == nil
}
