# FalaAI API Go SDK - Speech-to-Text, Call Analytics & Compliance Audit

[![version](https://img.shields.io/github/v/tag/ActionTecBr/falaai-api?label=go&color=00ADD8)](https://pkg.go.dev/github.com/actiontecbr/falaai-api)
[![license](https://img.shields.io/badge/license-MIT-green)](https://github.com/ActionTecBr/falaai-api/blob/main/LICENSE)
[![build](https://github.com/ActionTecBr/falaai-api/actions/workflows/ci.yml/badge.svg)](https://github.com/ActionTecBr/falaai-api/actions/workflows/ci.yml)

Official Go SDK for the **FalaAI API** - call transcription, conversation intelligence and compliance auditing (**COPC CX**, **ISO 18295-1**).

## What is FalaAI API?

FalaAI API turns conversations into auditable business intelligence, in three steps:

1. **Transcribe** - audio (calls, voice notes, meetings) to text, with speaker separation.
2. **Diagnose** - summary, reason, recommended action, topic and sentiment per conversation.
3. **Audit compliance** - risk score and violations against **COPC CX** and **ISO 18295-1**.

It works with **phone calls and call recordings** (PABX IP, Asterisk, FreePBX, contact center),
**messaging** (WhatsApp, Telegram, web chat, SMS) and email - anything that can be turned into
text. Three core endpoints (plus health, usage, webhooks and email alerts), one API key, no setup.

## Who it's for

| Role | What they get |
| --- | --- |
| **Contact Center / Quality** | Audit 100% of conversations instead of a sample |
| **Compliance / Legal** | Forensic, auditable evidence for audits and disputes |
| **CX / Operations** | Risk score, sentiment and reason for every conversation |
| **Developers** | One typed SDK, three core endpoints, one API key |
| **Data / BI** | Clean, typed JSON ready for your database or BI tool |

## Install

```bash
go get github.com/actiontecbr/falaai-api
```

## Quick start

```go
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	falaai "github.com/actiontecbr/falaai-api"
)

func authEditor(key string) falaai.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+key)
		return nil
	}
}

func main() {
	c, err := falaai.NewClientWithResponses(
		"https://api01-falaai.action.tec.br",
		falaai.WithRequestEditorFn(authEditor("fai_xxxxxx")),
	)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	f, err := os.Open("call.mp3")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("file", "call.mp3")
	if _, err := io.Copy(fw, f); err != nil {
		panic(err)
	}
	_ = mw.WriteField("model", "falaai-transcribe-1")
	_ = mw.WriteField("language", "pt")
	mw.Close()

	res, err := c.CreateTranscriptionV1AudioTranscriptionsPostWithBodyWithResponse(
		ctx, mw.FormDataContentType(), &buf)
	if err != nil {
		panic(err)
	}
	if res.JSON200 == nil {
		panic(fmt.Sprintf("HTTP %d: %s", res.StatusCode(), string(res.Body)))
	}
	fmt.Println(res.JSON200.Text)
}
```

## Request and response examples

Real requests and responses, captured against the live API.

### Health

**Request**

```bash
curl https://api01-falaai.action.tec.br/v1/health \
  -H "Authorization: Bearer fai_xxxxxx"
```

**Response**

```json
{
  "status": "ok",
  "version": "api_v1.21.47",
  "uptime_seconds": 92855,
  "database": true,
  "phase": "production",
  "launch_date": "2026-08-01"
}
```

### Transcribe

**Response** (the request is in [Quick start](#quick-start))

```json
{
  "id": "tr-a38a97d0-452e-4519-b6bb-3c63694bf7ae",
  "object": "transcription",
  "model": "falaai-transcribe-1",
  "filename": "analise_25s.mp3",
  "processed_at": "2026-09-23T17:42:54.242449+00:00",
  "usage": {
    "audio_seconds": 25.0,
    "credits_consumed": 25,
    "processing_ms": 951
  },
  "language": "por",
  "duration_seconds": 25.0,
  "text": "Novatechno tudo, bom dia. Bom dia, Mateus. Tudo bem? Bem e você? Tudo bem. Não tem nada no sistema. Eu queria ver se de  ...",
  "dialog": "Speaker 1: [00:00:01.639 - 00:00:02.560] Novatechno tudo, bom dia.\nSpeaker 2: [00:00:02.980 - 00:00:03.980] Bom dia, Mateus. Tudo bem?\nSpeaker 1: [00:00:04.700 - 00:00:05.179] Bem e você?\nSpeaker 2: [ ...",
  "audio_events": [
    {
      "event": "[tosse]",
      "start_s": 22.06,
      "end_s": 22.44,
      "duration_s": 0.38,
      "formatted_timestamp": "00:00:22.059"
    }
  ],
  "event_types": [
    "[tosse]"
  ],
  "word_count": 129,
  "input": {
    "duration_s": 25.2,
    "original_format": "mp3",
    "codec": "mp3",
    "sample_rate": 8000,
    "channels": 1
  },
  "language_confidence": 1.0,
  "client_reference_id": "readme-example"
}
```

### Diagnose

**Request**

```bash
curl https://api01-falaai.action.tec.br/v1/analyze/diagnostic \
  -H "Authorization: Bearer fai_xxxxxx" \
  -H "Content-Type: application/json" \
  -d '{"dialog": "<transcription dialog>", "language": "pt-BR", "duration_seconds": 25.0}'
```

**Response**

```json
{
  "id": "di-15c79d7c-7245-4af4-9a4e-ead237c1a6ed",
  "response_language": "pt-BR",
  "object": "analysis",
  "analysis": {
    "dialogue_summary": {
      "explanation": "O cliente (Speaker 2) contatou a Novatechno (Speaker 1) para obter informações sobre um produto (código 4675005) que saiu de linha. O cliente já realizou uma pe ..."
    },
    "contact_reason": {
      "explanation": "O cliente busca informação sobre qual produto substituiu o código 4675005, que saiu de linha."
    },
    "identified_action": {
      "list_choice": "Sem Ação",
      "justification": "O cliente buscou informações sobre um produto que saiu de linha, mas a conversa terminou sem uma ação definida ou resolução.",
      "evidence_phrases": [
        "Novatechno tudo, bom dia.",
        "Bom dia, Mateus. Tudo bem?",
        "Tudo bem. Não tem nada no sistema. Eu queria ver se de repente você tem alguma informação.",
        "Quatro meia sete cinco zero zero cinco saiu de linha. Ééé, eu procurei no Google ver se o que entrou no lugar. Aí ele fala"
      ]
    },
    "identified_label": {
      "list_choice": "Pedido de Informação",
      "justification": "O cliente busca informações sobre um produto que saiu de linha, perguntando o que o substituiu, indicando uma necessidade de dados.",
      "evidence_phrases": [
        "Não tem nada no sistema. Eu queria ver se de repente você tem alguma informação.",
        "É, o produto, posso te falar o código?",
  ... (truncated)
```

### Audit compliance

**Request**

```bash
curl https://api01-falaai.action.tec.br/v1/analyze/auditoriaRisco \
  -H "Authorization: Bearer fai_xxxxxx" \
  -H "Content-Type: application/json" \
  -d '{"dialog": "<transcription dialog>", "language": "pt-BR", "duration_seconds": 25.0, "response_format": "v2"}'
```

**Response**

```json
{
  "response": {
    "meta": {
      "id": "ar-3439fd5f-1932-4ebe-bfb5-944ec1046e60",
      "usage": {
        "characters": 378,
        "credits_consumed": 90,
        "processing_ms": 6534
      },
      "object": "auditoria_risco",
      "call_duration_s": 25.0,
      "analyzed_at": "2026-09-23T17:43:02.092299",
      "client_reference_id": "readme-example"
    },
    "participants": {
      "identified": [
        {
          "confidence": "high",
          "source": "input",
          "evidence": "fornecido pelo input",
          "interlocutor": "Speaker 1",
          "name": "Mateus",
          "role": "agent"
        },
        {
          "confidence": "high",
          "source": "input",
          "evidence": "fornecido pelo input",
          "interlocutor": "Speaker 2",
          "name": "Cliente",
          "role": "client"
        }
      ],
      "call_direction": "inbound",
      "role_inference_reliable": true,
      "identification_status": "input",
      "unidentified_items_count": 0
    },
    "verdict": {
      "label": "Limpa",
      "level_code": "LIMPA",
      "color": "#4fff4d",
      "icon": "mdi:check-circle",
      "risk_matrix": {
        "severity": "NOTE",
        "likelihood_avg": 0.0,
        "impact_avg": 0.0,
        "likelihood_level": "LOW",
        "impact_level": "LOW"
      },
      "applied_actions": [
        {
          "action_type": "nenhuma_acao",
          "label": "Nenhuma Acao",
          "description": "Nenhuma acao necessaria. A chamada nao apresentou deteccoes que exijam intervencao. Continuar monitoramento de rotina.",
          "priority": "BAIXO",
          "color": "#6B7280",
          "icon": "mdi:check-circle-outline",
          "condition": "sempre",
          "reason": "condição=sempre"
        }
      ],
  ... (truncated)
```

### The forensic HTML report

The audit response carries **`response.html_report`** - the complete **forensic report as HTML**, encoded as a **base64 gzip** string (~19 KB). Decode it and open in a browser (or convert to PDF) to get the auditable evidence: timeline, audio events (MAC/MVAD/MOD), violations against **COPC CX** and **ISO 18295-1**, and the executive verdict.

```python
import base64, gzip, json, httpx

res = httpx.post(
    "https://api01-falaai.action.tec.br/v1/analyze/auditoriaRisco",
    headers={"Authorization": "Bearer fai_xxxxxx"},
    json={"dialog": dialog, "language": "pt-BR", "duration_seconds": 25.0, "response_format": "v2"},
    timeout=900,
).json()

html = gzip.decompress(base64.b64decode(res["response"]["html_report"])).decode("utf-8")
open("report.html", "w", encoding="utf-8").write(html)
```

## Use cases

- Call and voice-note **transcription** with speaker separation
- **Contact center quality assurance (QA)** automation
- **Compliance auditing** against **COPC CX** and **ISO 18295-1**
- **Risk detection** - churn risk, legal threats, escalation
- **WhatsApp, Telegram and chat** conversation analysis
- **CRM and help desk** enrichment
- **Call analytics** and **speech-to-text** at scale - every call transcribed and scored
- **Conversation intelligence** - summary, reason, action, topic and **sentiment analysis** per conversation
- **Speaker diarization** (speaker separation) on stereo or mono audio
- **Quality monitoring** and **agent performance** - audit 100% instead of a sample
- Built from a single **OpenAPI** contract, so all SDKs stay in sync
- **LGPD**-aware handling of customer conversations

## Where it fits

Plugs into **omnichannel service platforms**, **help desks**, **chatbots** and **unified messaging** - anywhere a conversation becomes text.

Common Go stacks in contact center, CRM and help desk - if you build on any of these, the SDK drops in:

WAHA - NATS / RabbitMQ - high-performance microservices

> Product names are trademarks of their respective owners, listed as common stacks in this ecosystem. No partnership is implied.

## Endpoints

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/v1/audio/transcriptions` | Audio to text, with speaker separation |
| `POST` | `/v1/analyze/diagnostic` | Conversation analysis - summary, reason, action, topic, sentiment |
| `POST` | `/v1/analyze/auditoriaRisco` | Compliance audit - risk score and violations |
| `GET` | `/v1/health` | Service health check (public) |

All endpoints require `Authorization: Bearer fai_xxxxxx`.
Full reference: <https://api01-falaai.action.tec.br/docs>
Other operations: usage logs, webhooks and email alerts - see the full reference.

## Links

- **Product:** <https://falaai.action.tec.br/api>
- **API reference:** <https://api01-falaai.action.tec.br/docs>
- **Get an API key:** <https://falaai.action.tec.br/api/auth>
- **Package (pkg.go.dev):** <https://pkg.go.dev/github.com/actiontecbr/falaai-api>
- **Source:** <https://github.com/ActionTecBr/falaai-api>

## Other official SDKs

| Language | Install | Package |
| --- | --- | --- |
| **Python** | `pip install falaai-api` | [PyPI](https://pypi.org/project/falaai-api/) |
| **Node.js / TypeScript** | `npm install falaai-api` | [npm](https://www.npmjs.com/package/falaai-api) |
| **PHP** | `composer require actiontecbr/falaai-api` | [Packagist](https://packagist.org/packages/actiontecbr/falaai-api) |
| **Go** | `go get github.com/actiontecbr/falaai-api` | [pkg.go.dev](https://pkg.go.dev/github.com/actiontecbr/falaai-api) |
| **Ruby** | `gem install falaai-api` | [RubyGems](https://rubygems.org/gems/falaai-api) |
| **Java** | `io.github.actiontecbr:falaai-api` | [Maven Central](https://central.sonatype.com/artifact/io.github.actiontecbr/falaai-api) |
| **.NET / C#** | `dotnet add package FalaAI.Api` | [NuGet](https://www.nuget.org/packages/FalaAI.Api) |

## Requirements, support and more

- **Requirements:** an API key (`fai_...`) and a runtime able to upload audio files.
- **Documentation:** <https://api01-falaai.action.tec.br/docs>
- **Support:** open an issue in this repository.
- **Security:** responsible disclosure policy in [SECURITY.md](SECURITY.md).
- **Contributing:** see [CONTRIBUTING.md](CONTRIBUTING.md) · [Code of Conduct](CODE_OF_CONDUCT.md).
- **Changelog:** see [CHANGELOG.md](CHANGELOG.md).

## License

MIT (c) 2026 Action Tec Br - see [LICENSE](LICENSE).
