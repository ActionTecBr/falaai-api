package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func postJSON(t *testing.T, path string, body any) (int, any) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", BaseURL()+path, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+TestKey())
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func Test401InvalidKey(t *testing.T) {
	req, _ := http.NewRequest("GET", BaseURL()+"/v1/usage/log?page=1&limit=1", nil)
	req.Header.Set("Authorization", "Bearer fai_chave_invalida_000")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	Log("errors_401", "GET", "/v1/usage/log", nil, out, fmt.Sprintf("HTTP %d", resp.StatusCode), resp.StatusCode)
	if resp.StatusCode != 401 {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func Test422DiagnosticMissing(t *testing.T) {
	body := map[string]any{"language": "pt-BR", "dialog": "Speaker 1: ola"}
	st, out := postJSON(t, "/v1/analyze/diagnostic", body)
	Log("errors_422_diag", "POST", "/v1/analyze/diagnostic", body, out, fmt.Sprintf("HTTP %d", st), st)
	if st != 422 {
		t.Fatalf("status %d", st)
	}
}

func Test422AuditoriaMissing(t *testing.T) {
	body := map[string]any{"language": "pt-BR", "dialog": "Speaker 1: ola"}
	st, out := postJSON(t, "/v1/analyze/auditoriaRisco", body)
	Log("errors_422_aud", "POST", "/v1/analyze/auditoriaRisco", body, out, fmt.Sprintf("HTTP %d", st), st)
	if st != 422 {
		t.Fatalf("status %d", st)
	}
}

func Test422ExtraForbidden(t *testing.T) {
	body := map[string]any{"dialog": "Speaker 1: ola", "language": "pt-BR", "response_language": "pt-BR", "duration_seconds": 10, "threshold_multiplier": 1}
	st, out := postJSON(t, "/v1/analyze/auditoriaRisco", body)
	Log("errors_422_extra", "POST", "/v1/analyze/auditoriaRisco", body, out, fmt.Sprintf("HTTP %d", st), st)
	if st != 422 {
		t.Fatalf("status %d", st)
	}
}

func Test400AuditoriaLanguage(t *testing.T) {
	body := map[string]any{"dialog": "Speaker 1: ola", "language": "xx", "response_language": "pt-BR", "duration_seconds": 10}
	st, out := postJSON(t, "/v1/analyze/auditoriaRisco", body)
	Log("errors_400_aud", "POST", "/v1/analyze/auditoriaRisco", body, out, fmt.Sprintf("HTTP %d", st), st)
	if st != 400 {
		t.Fatalf("status %d", st)
	}
}

func Test400DiagnosticLanguage(t *testing.T) {
	body := map[string]any{"dialog": "Speaker 1: ola", "language": "xx", "duration_seconds": 10}
	st, out := postJSON(t, "/v1/analyze/diagnostic", body)
	Log("errors_400_diag", "POST", "/v1/analyze/diagnostic", body, out, fmt.Sprintf("HTTP %d", st), st)
	if st != 400 {
		t.Fatalf("status %d", st)
	}
}