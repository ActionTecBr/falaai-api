package e2e

import (
	"context"
	"fmt"
	"testing"

	falaai "github.com/actiontecbr/falaai-api"
)

const (
	wName  = "E2E Test Webhook"
	wURL   = "https://e2e-falaai.invalid/hook"
	aName  = "E2E Test Alert"
	aEmail = "e2e-test@falaai.invalid"
)

func sptr(s string) *string { return &s }
func bptr(b bool) *bool     { return &b }

func assertWebhookFull(t *testing.T, w falaai.WebhookItem, wid, name string, active bool) {
	t.Helper()
	if w.Id != wid {
		t.Fatalf("id %q != %q", w.Id, wid)
	}
	mustStr(t, "user_id", w.UserId)
	if w.Name != name {
		t.Fatalf("name %q != %q", w.Name, name)
	}
	if w.Url != wURL {
		t.Fatalf("url %q != %q", w.Url, wURL)
	}
	mustStr(t, "secret", w.Secret)
	if len(w.Events) != 2 {
		t.Fatalf("events len %d != 2", len(w.Events))
	}
	if w.Active != active {
		t.Fatalf("active %v != %v", w.Active, active)
	}
	mustStr(t, "created_at", w.CreatedAt)
	mustStr(t, "updated_at", w.UpdatedAt)
}

func assertAlertFull(t *testing.T, a falaai.EmailAlertItem, aid, name string, active bool) {
	t.Helper()
	if a.Id != aid {
		t.Fatalf("id %q != %q", a.Id, aid)
	}
	mustStr(t, "user_id", a.UserId)
	if a.Name != name {
		t.Fatalf("name %q != %q", a.Name, name)
	}
	if a.Email != aEmail {
		t.Fatalf("email %q != %q", a.Email, aEmail)
	}
	if len(a.Events) != 2 {
		t.Fatalf("events len %d != 2", len(a.Events))
	}
	if a.Active != active {
		t.Fatalf("active %v != %v", a.Active, active)
	}
	mustStr(t, "created_at", a.CreatedAt)
	mustStr(t, "updated_at", a.UpdatedAt)
}

func TestWebhooksCrud(t *testing.T) {
	c, _ := NewClient(BaseURL(), TestKey())
	ctx := context.Background()
	page, limit := 1, 100
	lst, _ := c.ListWebhooksV1WebhooksGetWithResponse(ctx, &falaai.ListWebhooksV1WebhooksGetParams{Page: &page, Limit: &limit})
	if lst.JSON200 != nil {
		for _, w := range lst.JSON200.Data {
			if w.Url == wURL {
				_, _ = c.DeleteWebhookV1WebhooksWebhookIdDeleteWithResponse(ctx, w.Id)
			}
		}
	}

	body := falaai.CreateWebhookRequest{Name: wName, Url: wURL, Events: []falaai.WebhookEvent{falaai.WebhookEventCreditsLow, falaai.WebhookEventPaymentFailed}}
	cr, err := c.CreateWebhookV1WebhooksPostWithResponse(ctx, body)
	if err != nil {
		t.Fatal(err)
	}
	Log("webhooks_create", "POST", "/v1/webhooks", body, cr.JSON200, fmt.Sprintf("HTTP %d", cr.StatusCode()), cr.StatusCode())
	if cr.StatusCode() != 200 {
		t.Fatalf("create status %d", cr.StatusCode())
	}
	wid := cr.JSON200.Id
	assertWebhookFull(t, *cr.JSON200, wid, wName, true)

	ub := falaai.UpdateWebhookRequest{Name: sptr(wName + " (updated)"), Active: bptr(false)}
	ur, err := c.UpdateWebhookV1WebhooksWebhookIdPutWithResponse(ctx, wid, ub)
	if err != nil {
		t.Fatal(err)
	}
	Log("webhooks_update", "PUT", "/v1/webhooks/"+wid, ub, ur.JSON200, fmt.Sprintf("HTTP %d", ur.StatusCode()), ur.StatusCode())
	if ur.StatusCode() != 200 {
		t.Fatalf("update status %d", ur.StatusCode())
	}
	if ur.JSON200.Message != "updated" {
		t.Fatalf("message %q != updated", ur.JSON200.Message)
	}

	lst2, err := c.ListWebhooksV1WebhooksGetWithResponse(ctx, &falaai.ListWebhooksV1WebhooksGetParams{Page: &page, Limit: &limit})
	if err != nil {
		t.Fatal(err)
	}
	var row *falaai.WebhookItem
	for i := range lst2.JSON200.Data {
		if lst2.JSON200.Data[i].Id == wid {
			row = &lst2.JSON200.Data[i]
		}
	}
	if row == nil {
		t.Fatalf("webhook %s nao encontrado na lista", wid)
	}
	assertWebhookFull(t, *row, wid, wName+" (updated)", false)

	dr, err := c.DeleteWebhookV1WebhooksWebhookIdDeleteWithResponse(ctx, wid)
	if err != nil {
		t.Fatal(err)
	}
	Log("webhooks_delete", "DELETE", "/v1/webhooks/"+wid, nil, dr.JSON200, fmt.Sprintf("HTTP %d", dr.StatusCode()), dr.StatusCode())
	if dr.StatusCode() != 200 {
		t.Fatalf("delete status %d", dr.StatusCode())
	}
	if dr.JSON200.Message != "deleted" {
		t.Fatalf("message %q != deleted", dr.JSON200.Message)
	}

	lst3, err := c.ListWebhooksV1WebhooksGetWithResponse(ctx, &falaai.ListWebhooksV1WebhooksGetParams{Page: &page, Limit: &limit})
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range lst3.JSON200.Data {
		if w.Id == wid {
			t.Fatalf("webhook %s ainda existe apos delete", wid)
		}
	}
}

func TestEmailAlertsCrud(t *testing.T) {
	c, _ := NewClient(BaseURL(), TestKey())
	ctx := context.Background()
	page, limit := 1, 100
	lst, _ := c.ListEmailAlertsV1EmailAlertsGetWithResponse(ctx, &falaai.ListEmailAlertsV1EmailAlertsGetParams{Page: &page, Limit: &limit})
	if lst.JSON200 != nil {
		for _, a := range lst.JSON200.Data {
			if a.Email == aEmail {
				_, _ = c.DeleteEmailAlertV1EmailAlertsAlertIdDeleteWithResponse(ctx, a.Id)
			}
		}
	}

	body := falaai.CreateEmailAlertRequest{Name: aName, Email: aEmail, Events: []falaai.EmailEvent{falaai.EmailEventCreditsLow, falaai.EmailEventPaymentFailed}}
	cr, err := c.CreateEmailAlertV1EmailAlertsPostWithResponse(ctx, body)
	if err != nil {
		t.Fatal(err)
	}
	Log("email_alerts_create", "POST", "/v1/email-alerts", body, cr.JSON200, fmt.Sprintf("HTTP %d", cr.StatusCode()), cr.StatusCode())
	if cr.StatusCode() != 200 {
		t.Fatalf("create status %d", cr.StatusCode())
	}
	aid := cr.JSON200.Id
	assertAlertFull(t, *cr.JSON200, aid, aName, true)

	ub := falaai.UpdateEmailAlertRequest{Name: sptr(aName + " (updated)"), Active: bptr(false)}
	ur, err := c.UpdateEmailAlertV1EmailAlertsAlertIdPutWithResponse(ctx, aid, ub)
	if err != nil {
		t.Fatal(err)
	}
	Log("email_alerts_update", "PUT", "/v1/email-alerts/"+aid, ub, ur.JSON200, fmt.Sprintf("HTTP %d", ur.StatusCode()), ur.StatusCode())
	if ur.StatusCode() != 200 {
		t.Fatalf("update status %d", ur.StatusCode())
	}
	if ur.JSON200.Message != "updated" {
		t.Fatalf("message %q != updated", ur.JSON200.Message)
	}

	lst2, err := c.ListEmailAlertsV1EmailAlertsGetWithResponse(ctx, &falaai.ListEmailAlertsV1EmailAlertsGetParams{Page: &page, Limit: &limit})
	if err != nil {
		t.Fatal(err)
	}
	var row *falaai.EmailAlertItem
	for i := range lst2.JSON200.Data {
		if lst2.JSON200.Data[i].Id == aid {
			row = &lst2.JSON200.Data[i]
		}
	}
	if row == nil {
		t.Fatalf("alert %s nao encontrado na lista", aid)
	}
	assertAlertFull(t, *row, aid, aName+" (updated)", false)

	dr, err := c.DeleteEmailAlertV1EmailAlertsAlertIdDeleteWithResponse(ctx, aid)
	if err != nil {
		t.Fatal(err)
	}
	Log("email_alerts_delete", "DELETE", "/v1/email-alerts/"+aid, nil, dr.JSON200, fmt.Sprintf("HTTP %d", dr.StatusCode()), dr.StatusCode())
	if dr.StatusCode() != 200 {
		t.Fatalf("delete status %d", dr.StatusCode())
	}
	if dr.JSON200.Message != "deleted" {
		t.Fatalf("message %q != deleted", dr.JSON200.Message)
	}

	lst3, err := c.ListEmailAlertsV1EmailAlertsGetWithResponse(ctx, &falaai.ListEmailAlertsV1EmailAlertsGetParams{Page: &page, Limit: &limit})
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range lst3.JSON200.Data {
		if a.Id == aid {
			t.Fatalf("alert %s ainda existe apos delete", aid)
		}
	}
}