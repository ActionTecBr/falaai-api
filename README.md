# FalaAI API — Go SDK

Official Go SDK for the **FalaAI API** — AI-powered call transcription, diagnosis and compliance auditing.

## Install

```bash
go get github.com/actiontecbr/falaai-api
```

## Quick start

```go
package main

import (
    "context"
    "fmt"

    falaai "github.com/actiontecbr/falaai-api"
)

func main() {
    // Create an HTTP client with the API base URL and your API key.
    client, err := falaai.NewClient("https://api01-falaai.action.tec.br")
    if err != nil {
        panic(err)
    }
    client.RequestEditorFn = falaai.BearerAuth("fai_xxx")

    res, err := client.HealthCheckWithResponse(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Println(res.JSON200.Status)
}
```

## Endpoints

| Method | Path | Description |
|---|---|---|
| POST | `/v1/audio/transcriptions` | Audio to text (diarization, audio events) |
| POST | `/v1/analyze/diagnostic` | Conversation analysis |
| POST | `/v1/analyze/auditoriaRisco` | Compliance audit (risk) |
| GET | `/v1/usage/log` | Usage log |
| GET | `/v1/usage/by-key` | Usage grouped by API key |
| GET | `/v1/webhooks` | List webhooks |
| POST | `/v1/webhooks` | Create webhook |
| PUT | `/v1/webhooks/{webhook_id}` | Update webhook |
| DELETE | `/v1/webhooks/{webhook_id}` | Delete webhook |
| GET | `/v1/email-alerts` | List email alerts |
| POST | `/v1/email-alerts` | Create email alert |
| PUT | `/v1/email-alerts/{alert_id}` | Update email alert |
| DELETE | `/v1/email-alerts/{alert_id}` | Delete email alert |
| GET | `/api/version` | API version |
| GET | `/v1/health` | Health check |
| HEAD | `/v1/health` | Health check (HEAD) |

## Authentication

All authenticated endpoints require an API key in the `Authorization` header:

```
Authorization: Bearer fai_xxx
```

Get your API key at [falaai.action.tec.br/api](https://falaai.action.tec.br/api).

## Tests

An end-to-end suite (19 tests) lives in `tests/e2e/` and runs against the live API:

```bash
FALAAI_E2E_BASE=https://api01-falaai.action.tec.br \
FALAAI_TEST_KEY=fai_xxx \
FALAAI_E2E_AUDIO=/path/to/audio.mp3 \
go test ./tests/e2e/ -count=1
```

## License

[MIT](LICENSE) © Action Tec Br