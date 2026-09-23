package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func Log(name, method, path string, payload, response any, result string, status int) string {
	re := regexp.MustCompile(`[^A-Za-z0-9]+`)
	safe := strings.ToUpper(re.ReplaceAllString(name, "_"))
	dir := filepath.Join("logs", safe)
	_ = os.MkdirAll(dir, 0o777)
	ts := time.Now().Format("20060102_150405")
	file := filepath.Join(dir, fmt.Sprintf("%s_%s.log", safe, ts))

	enc := func(v any) string {
		if v == nil {
			return "(sem dados)"
		}
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
	lines := []string{
		strings.Repeat("=", 70),
		fmt.Sprintf("TESTE: %s %s", method, path),
		"DATA: " + time.Now().Format(time.RFC3339),
		strings.Repeat("=", 70),
		"",
		"--- PAYLOAD (enviado) ---",
		enc(payload),
		"",
		"--- RESPOSTA (saida do SDK) ---",
		fmt.Sprintf("HTTP: %d", status),
		enc(response),
		"",
		"--- RESULTADO ---",
		result,
		"",
	}
	_ = os.WriteFile(file, []byte(strings.Join(lines, "\n")), 0o644)
	return file
}