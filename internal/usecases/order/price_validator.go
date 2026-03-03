package order

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"yego/internal/domain"
)

// normalizeKey lowercases a string and removes unicode accents (NFD decomposition).
// "Descripción" → "descripcion", "Código" → "codigo"
func normalizeKey(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, strings.ToLower(s))
	return result
}

// findColValue searches the import data map for a key that contains any of the
// given patterns (accent-insensitive, case-insensitive) and returns its string value.
func findColValue(data map[string]any, patterns []string) (string, bool) {
	for k, v := range data {
		normK := normalizeKey(k)
		for _, p := range patterns {
			if strings.Contains(normK, normalizeKey(p)) {
				return strings.TrimSpace(strings.ReplaceAll(fmt.Sprintf("%v", v), "\u00a0", " ")), true
			}
		}
	}
	return "", false
}

// findImportByCode returns the first ImportRecord whose code column matches the given code.
func findImportByCode(records []*domain.ImportRecord, code string) *domain.ImportRecord {
	normCode := normalizeKey(strings.TrimSpace(code))
	log.Printf("[PriceValidator] looking for code=%q (normalized=%q) in %d import records", code, normCode, len(records))
	for _, r := range records {
		val, ok := findColValue(r.Data, []string{"codigo", "code", "sku", "ref", "referencia"})
		log.Printf("[PriceValidator]   record id=%s code_col_found=%v code_val=%q", r.ID, ok, val)
		if ok && normalizeKey(val) == normCode {
			log.Printf("[PriceValidator]   MATCH found record id=%s", r.ID)
			return r
		}
	}
	log.Printf("[PriceValidator] NO match for code=%q", code)
	return nil
}

// importName extracts the product name from an import record.
func importName(data map[string]any) string {
	val, ok := findColValue(data, []string{"descripcion", "nombre", "name", "producto", "description"})
	if !ok {
		return ""
	}
	return val
}

// importPrice extracts the unit price from an import record, rounded to 2 decimal places.
func importPrice(data map[string]any) (float64, bool) {
	val, ok := findColValue(data, []string{"precio unitario", "precio", "price", "costo", "valor", "importe"})
	if !ok {
		log.Printf("[PriceValidator] price column not found in data keys: %v", mapKeys(data))
		return 0, false
	}
	log.Printf("[PriceValidator] price raw value=%q", val)
	cleaned := strings.NewReplacer("$", "", " ", "", "\u00a0", "").Replace(val)
	cleaned = strings.ReplaceAll(cleaned, ",", ".")
	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		log.Printf("[PriceValidator] price parse error: %v (cleaned=%q)", err, cleaned)
		return 0, false
	}
	return math.Round(price*100) / 100, true
}

func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// findImportByName returns the best matching ImportRecord for the given name.
// Strategy: exact match first, then partial match (import contains item name or vice versa).
func findImportByName(records []*domain.ImportRecord, name string) *domain.ImportRecord {
	normName := normalizeKey(strings.TrimSpace(name))
	log.Printf("[PriceValidator] looking for name=%q (normalized=%q) in %d import records", name, normName, len(records))

	// Pass 1: exact match
	for _, r := range records {
		val, ok := findColValue(r.Data, []string{"descripcion", "nombre", "name", "producto", "description"})
		if ok && normalizeKey(val) == normName {
			log.Printf("[PriceValidator]   EXACT match found record id=%s import_name=%q", r.ID, val)
			return r
		}
	}

	// Pass 2: partial match — import name contains item name, or item name contains import name
	for _, r := range records {
		val, ok := findColValue(r.Data, []string{"descripcion", "nombre", "name", "producto", "description"})
		if !ok {
			continue
		}
		normVal := normalizeKey(val)
		if strings.Contains(normVal, normName) || strings.Contains(normName, normVal) {
			log.Printf("[PriceValidator]   PARTIAL match found record id=%s import_name=%q", r.ID, val)
			return r
		}
	}

	log.Printf("[PriceValidator] NO match for name=%q", name)
	return nil
}

// correctItemPrices looks up each item by code (or name as fallback) against the
// import records and corrects Name and Price to match. Returns the (possibly
// corrected) slice and a boolean indicating whether any changes were made.
func correctItemPrices(items []domain.OrderItem, records []*domain.ImportRecord) ([]domain.OrderItem, bool) {
	hasChanges := false
	corrected := make([]domain.OrderItem, len(items))
	for i, item := range items {
		corrected[i] = item
		log.Printf("[PriceValidator] item[%d] code=%q name=%q price=%.2f", i, item.Code, item.Name, item.Price)

		if item.Code == "" && item.Name == "" {
			log.Printf("[PriceValidator] item[%d] WARNING: no code and no name, skipping", i)
			continue
		}

		var matched *domain.ImportRecord
		if item.Code != "" {
			matched = findImportByCode(records, item.Code)
			if matched == nil {
				log.Printf("[PriceValidator] item[%d] no import match for code=%q, trying by name", i, item.Code)
			}
		}
		if matched == nil && item.Name != "" {
			matched = findImportByName(records, item.Name)
		}
		if matched == nil {
			log.Printf("[PriceValidator] item[%d] no match found, skipping", i)
			continue
		}

		log.Printf("[PriceValidator] item[%d] matched import id=%s", i, matched.ID)
		if name := importName(matched.Data); name != "" && name != item.Name {
			log.Printf("[PriceValidator] item[%d] correcting name: %q → %q", i, item.Name, name)
			corrected[i].Name = name
			hasChanges = true
		}
		if price, ok := importPrice(matched.Data); ok && price != item.Price {
			log.Printf("[PriceValidator] item[%d] correcting price: %.2f → %.2f", i, item.Price, price)
			corrected[i].Price = price
			hasChanges = true
		}
	}
	log.Printf("[PriceValidator] hasChanges=%v", hasChanges)
	return corrected, hasChanges
}
