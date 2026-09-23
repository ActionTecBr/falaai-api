package e2e

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"

	falaai "github.com/actiontecbr/falaai-api"
)

func TestAiChain(t *testing.T) {
	c, _ := NewClient(ProdURL(), TestKey())
	ctx := context.Background()

	f, err := os.Open(AudioPath())
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="analise_25s.mp3"`)
	h.Set("Content-Type", "audio/mpeg")
	fw, _ := mw.CreatePart(h)
	if _, err := io.Copy(fw, f); err != nil {
		t.Fatal(err)
	}
	f.Close()
	_ = mw.WriteField("model", "falaai-transcribe-1")
	_ = mw.WriteField("language", "pt")
	_ = mw.WriteField("client_reference_id", "e2e-call-2026-09-22-001")
	mw.Close()

	tr, err := c.CreateTranscriptionV1AudioTranscriptionsPostWithBodyWithResponse(ctx, mw.FormDataContentType(), &buf)
	if err != nil {
		t.Fatal(err)
	}
	Log("transcriptions", "POST", "/v1/audio/transcriptions",
		map[string]string{"file": "analise_25s.mp3", "model": "falaai-transcribe-1", "language": "pt", "client_reference_id": "e2e-call-2026-09-22-001"},
		tr.JSON200, fmt.Sprintf("HTTP %d", tr.StatusCode()), tr.StatusCode())
	if tr.StatusCode() != 200 {
		t.Fatalf("transcribe status %d | body: %s", tr.StatusCode(), string(tr.Body))
	}
	tp := tr.JSON200
	mustStr(t, "id", tp.Id)
	mustStr(t, "object", tp.Object)
	mustStr(t, "model", tp.Model)
	mustStr(t, "filename", tp.Filename)
	mustStr(t, "processed_at", tp.ProcessedAt)
	mustGt0(t, "usage.audio_seconds", tp.Usage.AudioSeconds)
	_ = tp.Usage.CreditsConsumed
	_ = tp.Usage.ProcessingMs
	mustStr(t, "language", tp.Language)
	mustGt0(t, "duration_seconds", tp.DurationSeconds)
	mustStr(t, "text", tp.Text)
	mustStr(t, "dialog", tp.Dialog)
	for _, ev := range tp.AudioEvents {
		mustStr(t, "event", ev.Event)
		_ = ev.StartS
		_ = ev.EndS
		_ = ev.DurationS
		mustStr(t, "formatted_timestamp", ev.FormattedTimestamp)
	}
	_ = tp.EventTypes
	mustGt0Int(t, "word_count", tp.WordCount)
	_ = tp.Input.DurationS
	mustStr(t, "input.original_format", tp.Input.OriginalFormat)
	mustStr(t, "input.codec", tp.Input.Codec)
	_ = tp.Input.SampleRate
	_ = tp.Input.Channels

	dialog := tp.Dialog
	dur := tp.DurationSeconds
	txt := tp.Text

	dbody := falaai.DiagnosticRequest{Dialog: &dialog, DurationSeconds: dur, Language: "pt-BR", Text: &txt, ClientReferenceId: sptr("e2e-diag-2026-09-22-001")}
	d, err := c.CreateDiagnosticV1AnalyzeDiagnosticPostWithResponse(ctx, dbody)
	if err != nil {
		t.Fatal(err)
	}
	Log("diagnostic", "POST", "/v1/analyze/diagnostic", dbody, d.JSON200, fmt.Sprintf("HTTP %d", d.StatusCode()), d.StatusCode())
	if d.StatusCode() != 200 {
		t.Fatalf("diagnostic status %d", d.StatusCode())
	}
	dp := d.JSON200
	mustStr(t, "id", dp.Id)
	mustStr(t, "response_language", dp.ResponseLanguage)
	if dp.Object != "analysis" {
		t.Fatalf("object %q != analysis", dp.Object)
	}
	if dp.Analysis.DialogueSummary.Explanation == nil {
		t.Fatalf("analysis.dialogue_summary vazio")
	}
	if dp.Analysis.ContactReason.Explanation == nil {
		t.Fatalf("analysis.contact_reason vazio")
	}
	_ = dp.Analysis.IdentifiedAction
	_ = dp.Analysis.IdentifiedLabel
	_ = dp.Analysis.Sentiment
	_ = dp.Usage.Characters
	_ = dp.Usage.CreditsConsumed
	_ = dp.Usage.ProcessingMs

	cd := falaai.Inbound
	abody := falaai.AuditoriaRiscoRequest{
		Dialog:           &dialog,
		DurationSeconds:  dur,
		Language:         "pt-BR",
		ResponseLanguage: "pt-BR",
		Text:             &txt,
		CallDirection:    &cd,
		Participants: &[]falaai.Participant{
			{Interlocutor: "Speaker 1", Name: sptr("Mateus"), Role: falaai.ParticipantRoleAgent},
			{Interlocutor: "Speaker 2", Name: sptr("Cliente"), Role: falaai.ParticipantRoleClient},
		},
		ClientReferenceId: sptr("e2e-aud-2026-09-22-001"),
	}
	a, err := c.CreateAuditoriaRiscoV1AnalyzeAuditoriaRiscoPostWithResponse(ctx, abody)
	if err != nil {
		t.Fatal(err)
	}
	Log("auditoriaRisco", "POST", "/v1/analyze/auditoriaRisco", abody, a.JSON200, fmt.Sprintf("HTTP %d", a.StatusCode()), a.StatusCode())
	if a.StatusCode() != 200 {
		t.Fatalf("auditoria status %d", a.StatusCode())
	}
	pub := a.JSON200.Response
	mustStr(t, "meta.id", pub.Meta.Id)
	_ = pub.Meta.Usage.Characters
	_ = pub.Meta.Usage.CreditsConsumed
	_ = pub.Meta.Usage.ProcessingMs
	_ = pub.Participants
	_ = pub.Verdict
	_ = pub.Scores
	_ = pub.Detections
	_ = pub.Analysis
	_ = pub.Timeline
	_ = pub.AudioEventModel
	if pub.CategoriesSummary == nil {
		t.Fatalf("categories_summary ausente")
	}
	_ = pub.Indexer
	_ = pub.Summary
	if pub.AcoesI18n == nil {
		t.Fatalf("acoes_i18n ausente")
	}
	_ = pub.AuditDecisions
	_ = pub.ScoringExplanation
	_ = pub.HtmlReport
}