package e2e

import (
	"context"
	"net/http"
	"os"
	"strings"

	falaai "github.com/actiontecbr/falaai-api"
)

var envCache map[string]string

func loadEnv() map[string]string {
	if envCache != nil {
		return envCache
	}
	envCache = map[string]string{}
	for _, p := range []string{"../../../.env.e2e", "../.env.e2e", ".env.e2e"} {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			i := strings.Index(line, "=")
			envCache[strings.TrimSpace(line[:i])] = strings.TrimSpace(line[i+1:])
		}
		break
	}
	return envCache
}

func getenv(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return loadEnv()[k]
}

func BaseURL() string {
	if v := getenv("FALAAI_E2E_BASE"); v != "" {
		return v
	}
	if v := getenv("FALAAI_LOCAL_URL"); v != "" {
		return v
	}
	return "http://localhost:8002"
}

func ProdURL() string {
	if v := getenv("FALAAI_PROD_URL"); v != "" {
		return v
	}
	return "https://api01-falaai.action.tec.br"
}

func TestKey() string { return getenv("FALAAI_TEST_KEY") }

func AudioPath() string { return getenv("FALAAI_E2E_AUDIO") }

func authEditor(key string) falaai.RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+key)
		return nil
	}
}

func NewClient(baseURL, key string) (*falaai.ClientWithResponses, error) {
	return falaai.NewClientWithResponses(baseURL, falaai.WithRequestEditorFn(authEditor(key)))
}