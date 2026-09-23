package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	falaai "github.com/actiontecbr/falaai-api"
)

func mustStr(t *testing.T, label, v string) {
	t.Helper()
	if v == "" {
		t.Fatalf("%s: string vazia", label)
	}
}

func mustGt0(t *testing.T, label string, v float32) {
	t.Helper()
	if v <= 0 {
		t.Fatalf("%s: %v <= 0", label, v)
	}
}

func mustGt0Int(t *testing.T, label string, v int) {
	t.Helper()
	if v <= 0 {
		t.Fatalf("%s: %d <= 0", label, v)
	}
}

func TestHealthGet(t *testing.T) {
	c, _ := NewClient(BaseURL(), TestKey())
	r, err := c.HealthCheckWithResponse(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Log("health_get", "GET", "/v1/health", nil, r.JSON200, fmt.Sprintf("HTTP %d", r.StatusCode()), r.StatusCode())
	if r.StatusCode() != 200 {
		t.Fatalf("status %d", r.StatusCode())
	}
	p := r.JSON200
	if p.Status != "ok" {
		t.Fatalf("status != ok")
	}
	mustStr(t, "version", p.Version)
	if p.UptimeSeconds < 0 {
		t.Fatalf("uptime_seconds < 0")
	}
	_ = p.Database
	mustStr(t, "phase", p.Phase)
	mustStr(t, "launch_date", p.LaunchDate)
}

func TestHealthHead(t *testing.T) {
	req, _ := http.NewRequest("HEAD", BaseURL()+"/v1/health", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	Log("health_head", "HEAD", "/v1/health", nil, map[string]int{"status": resp.StatusCode}, fmt.Sprintf("HTTP %d", resp.StatusCode), resp.StatusCode)
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestVersion(t *testing.T) {
	c, _ := NewClient(BaseURL(), TestKey())
	r, err := c.GetVersionApiVersionGetWithResponse(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	Log("version", "GET", "/api/version", nil, r.JSON200, fmt.Sprintf("HTTP %d", r.StatusCode()), r.StatusCode())
	if r.StatusCode() != 200 {
		t.Fatalf("status %d", r.StatusCode())
	}
	p := r.JSON200
	if p.Service != "FalaAI API" {
		t.Fatalf("service != FalaAI API")
	}
	mustStr(t, "version", p.Version)
	mustStr(t, "deployDate", p.DeployDate)
}

func TestUsageLog(t *testing.T) {
	c, _ := NewClient(BaseURL(), TestKey())
	page, limit := 1, 5
	p := falaai.GetUsageLogV1UsageLogGetParams{Page: &page, Limit: &limit}
	r, err := c.GetUsageLogV1UsageLogGetWithResponse(context.Background(), &p)
	if err != nil {
		t.Fatal(err)
	}
	Log("usage_log", "GET", "/v1/usage/log", map[string]int{"page": 1, "limit": 5}, r.JSON200, fmt.Sprintf("HTTP %d", r.StatusCode()), r.StatusCode())
	if r.StatusCode() != 200 {
		t.Fatalf("status %d", r.StatusCode())
	}
	res := r.JSON200
	if res.Page != 1 || res.Limit != 5 {
		t.Fatalf("page/limit %d/%d", res.Page, res.Limit)
	}
	for _, it := range res.Data {
		mustStr(t, "id", it.Id)
		mustStr(t, "endpoint", it.Endpoint)
		mustStr(t, "status", it.Status)
		mustStr(t, "created_at", it.CreatedAt)
		_ = it.CreditsCost
		_ = it.ErrorsCount
	}
}

func TestUsageByKey(t *testing.T) {
	c, _ := NewClient(BaseURL(), TestKey())
	r, err := c.GetUsageByKeyV1UsageByKeyGetWithResponse(context.Background(), &falaai.GetUsageByKeyV1UsageByKeyGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	Log("usage_by_key", "GET", "/v1/usage/by-key", nil, r.JSON200, fmt.Sprintf("HTTP %d", r.StatusCode()), r.StatusCode())
	if r.StatusCode() != 200 {
		t.Fatalf("status %d", r.StatusCode())
	}
	for _, it := range *r.JSON200 {
		mustStr(t, "key_id", it.KeyId)
		_ = it.KeyName
		_ = it.TotalCredits
		_ = it.RequestCount
	}
}

func TestWebhooksList(t *testing.T) {
	c, _ := NewClient(BaseURL(), TestKey())
	page, limit := 1, 5
	p := falaai.ListWebhooksV1WebhooksGetParams{Page: &page, Limit: &limit}
	r, err := c.ListWebhooksV1WebhooksGetWithResponse(context.Background(), &p)
	if err != nil {
		t.Fatal(err)
	}
	Log("webhooks_list", "GET", "/v1/webhooks", map[string]int{"page": 1, "limit": 5}, r.JSON200, fmt.Sprintf("HTTP %d", r.StatusCode()), r.StatusCode())
	if r.StatusCode() != 200 {
		t.Fatalf("status %d", r.StatusCode())
	}
	res := r.JSON200
	if res.Page != 1 || res.Limit != 5 {
		t.Fatalf("page/limit %d/%d", res.Page, res.Limit)
	}
	for _, w := range res.Data {
		mustStr(t, "id", w.Id)
		mustStr(t, "user_id", w.UserId)
		_ = w.Name
		mustStr(t, "url", w.Url)
		_ = w.Secret
		_ = w.Events
		_ = w.Active
		_ = w.RetryEnabled
		_ = w.FailureCount
		mustStr(t, "created_at", w.CreatedAt)
		mustStr(t, "updated_at", w.UpdatedAt)
	}
}

func TestEmailAlertsList(t *testing.T) {
	c, _ := NewClient(BaseURL(), TestKey())
	page, limit := 1, 5
	p := falaai.ListEmailAlertsV1EmailAlertsGetParams{Page: &page, Limit: &limit}
	r, err := c.ListEmailAlertsV1EmailAlertsGetWithResponse(context.Background(), &p)
	if err != nil {
		t.Fatal(err)
	}
	Log("email_alerts_list", "GET", "/v1/email-alerts", map[string]int{"page": 1, "limit": 5}, r.JSON200, fmt.Sprintf("HTTP %d", r.StatusCode()), r.StatusCode())
	if r.StatusCode() != 200 {
		t.Fatalf("status %d", r.StatusCode())
	}
	res := r.JSON200
	if res.Page != 1 || res.Limit != 5 {
		t.Fatalf("page/limit %d/%d", res.Page, res.Limit)
	}
	for _, a := range res.Data {
		mustStr(t, "id", a.Id)
		mustStr(t, "user_id", a.UserId)
		_ = a.Name
		mustStr(t, "email", a.Email)
		_ = a.Events
		_ = a.Active
		mustStr(t, "created_at", a.CreatedAt)
		mustStr(t, "updated_at", a.UpdatedAt)
	}
}