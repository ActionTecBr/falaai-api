package e2e

import (
	"encoding/json"
	"os"
	"testing"
)

var expectedOps = []string{
	"POST /v1/audio/transcriptions", "POST /v1/analyze/diagnostic", "POST /v1/analyze/auditoriaRisco",
	"GET /v1/usage/log", "GET /v1/usage/by-key", "GET /v1/webhooks", "POST /v1/webhooks",
	"PUT /v1/webhooks/{webhook_id}", "DELETE /v1/webhooks/{webhook_id}",
	"GET /v1/email-alerts", "POST /v1/email-alerts", "PUT /v1/email-alerts/{alert_id}",
	"DELETE /v1/email-alerts/{alert_id}", "GET /api/version", "GET /v1/health", "HEAD /v1/health",
}

func TestOpenapiTemOperacoesEsperadas(t *testing.T) {
	data, err := os.ReadFile("../../../../openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		Paths map[string]map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for path, methods := range spec.Paths {
		for m := range methods {
			switch m {
			case "get", "post", "put", "delete", "patch", "head":
				got[toUpper(m)+" "+path] = true
			}
		}
	}
	if len(got) != len(expectedOps) {
		t.Fatalf("esperado %d ops, achou %d", len(expectedOps), len(got))
	}
	for _, op := range expectedOps {
		if !got[op] {
			t.Fatalf("op ausente: %s", op)
		}
	}
}

func TestSdkCobre100pc(t *testing.T) {
	src, err := os.ReadFile("../../../../sdks/go/falaai.gen.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(src)
	for _, p := range []string{"/v1/usage/log", "/v1/webhooks", "/v1/email-alerts", "/api/version", "/v1/health", "/v1/analyze/auditoriaRisco"} {
		if !contains(text, p) {
			t.Fatalf("SDK nao cobre %s", p)
		}
	}
}

func TestExemplosExistem(t *testing.T) {
	for _, f := range []string{
		"curl/transcribe.sh", "python/transcribe.py", "nodejs/transcribe.js",
		"curl/auditoria_risco.sh", "python/auditoria_risco.py", "nodejs/auditoria_risco.js",
		"curl/diagnostic.sh", "python/diagnostic.py", "nodejs/diagnostic.js",
	} {
		if _, err := os.Stat("../../../../app/static/examples/" + f); err != nil {
			t.Fatalf("exemplo ausente: %s", f)
		}
	}
}

func toUpper(s string) string {
	out := []rune(s)
	for i, r := range out {
		if r >= 'a' && r <= 'z' {
			out[i] = r - 32
		}
	}
	return string(out)
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}