# FalaAI API - Go SDK

[![version](https://img.shields.io/badge/version-1.21.47-blue)](https://pkg.go.dev/github.com/actiontecbr/falaai-api)
[![license](https://img.shields.io/badge/license-MIT-green)](https://github.com/ActionTecBr/falaai-api/blob/main/LICENSE)
[![build](https://github.com/ActionTecBr/falaai-api/actions/workflows/ci.yml/badge.svg)](https://github.com/ActionTecBr/falaai-api/actions/workflows/ci.yml)

Official Go SDK for the **FalaAI API**.

## What is FalaAI API?

FalaAI API turns conversations into auditable business intelligence, in three steps:

1. **Transcribe** - audio (calls, voice notes, meetings) to text, with speaker separation.
2. **Diagnose** - summary, reason, recommended action, topic and sentiment per conversation.
3. **Audit compliance** - risk score and violations against **COPC CX** and **ISO 18295-1**.

It works with phone calls, WhatsApp, Telegram, chat, email, PDF and images.
Three REST endpoints, one API key, no setup.

## Who it's for

| Role | What they get |
| --- | --- |
| **Contact Center / Quality** | Audit 100% of conversations instead of a sample |
| **Compliance / Legal** | Forensic, auditable evidence for audits and disputes |
| **CX / Operations** | Risk score, sentiment and reason for every conversation |
| **Developers** | One typed SDK, three REST endpoints, one API key |
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
	"os"

	falaai "github.com/actiontecbr/falaai-api"
)

func main() {
	c, err := falaai.NewClientWithResponses("https://api01-falaai.action.tec.br")
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
	fmt.Println(res.JSON200.Text)
}
```

## Use cases

- Call and voice-note **transcription** with speaker separation
- **Contact center quality assurance (QA)** automation
- **Compliance auditing** against **COPC CX** and **ISO 18295-1**
- **Risk detection** - churn risk, legal threats, escalation
- **WhatsApp, Telegram and chat** conversation analysis
- **CRM and help desk** enrichment
- **LGPD**-aware handling of customer conversations

## Where it fits

Common Go stacks in contact center, CRM and help desk - if you build on any of these, the SDK drops in:

WAHA - NATS / RabbitMQ - high-performance microservices

> Product names are trademarks of their respective owners, listed as common stacks in this ecosystem. No partnership is implied.

## Endpoints

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/v1/audio/transcriptions` | Audio to text, with speaker separation |
| `POST` | `/v1/analyze/diagnostic` | Conversation analysis - summary, reason, action, topic, sentiment |
| `POST` | `/v1/analyze/auditoriaRisco` | Compliance audit - risk score and violations |

All endpoints require `Authorization: Bearer fai_xxxxxx`.
Full reference: <https://api01-falaai.action.tec.br/docs>

## Links

- **Product:** <https://falaai.action.tec.br/api>
- **API reference:** <https://api01-falaai.action.tec.br/docs>
- **Get an API key:** <https://falaai.action.tec.br/api/auth>
- **Package (pkg.go.dev):** <https://pkg.go.dev/github.com/actiontecbr/falaai-api>
- **Source:** <https://github.com/ActionTecBr/falaai-api>

## License

MIT (c) 2026 Action Tec Br - see [LICENSE](LICENSE).
