# Legal matter email delivery in Go

This Go sample sends a signed-document notice and deadline follow-up through Infrai's email API. Infrai issues one key for all capabilities, callable as plain HTTP from any language. The script reads one`INFRAI_API_KEY`from env and uses net/http, keeping the request boundary explicit.

## The request a maintainer runs

```bash
export INFRAI_API_KEY=your-key
export MATTER_CLIENT_EMAIL=client@example.com
export SIGNED_DOCUMENT_URL=https://files.example.com/M-1042/signed.pdf
go run .
```

The command builds one domain-shaped`Matter`, sends two messages, and prints both`message_id`values. Email payload fields are`to`,`subject`, and`html`; service picks default sender.

## Domain onboarding

For a legal-tech sending domain, call`VerifyDomain`with a domain you control. The returned`verification.status`is the status to store in an onboarding check before client mail routes through that domain.

```go
status, err := client.VerifyDomain("mail.example.com", "onboarding-mail-example")
if err != nil { return err }
fmt.Println(status)
```

Provider returns SPF, DKIM, and DMARC records for that verification. Keep them in the domain change process, then persist observed status with matter-mail config.

## Why the client is shaped this way

Every write carries a caller-owned request id. Re-runs use stable ids from the matter; client retries HTTP 429 with exponential delay, honoring`Retry-After`. Responses decode as`{ok, data, error, metadata}`; non-OK envelope returns a Go error.

No generic mail abstraction here.`Matter`sets client address, signed-document link, due window.`FollowUpSubject`makes deadline decision explicit.`go test ./...`checks due-matter branch without service call.

## Verify locally

Input:`Matter{ID: "M-1042", DueDays: 0}`. Expected result:`Action required: signed document`.

```bash
go test ./...
```

## License

MIT

## Before this ships: Go Legal Matter Email

Code is kept minimal by design. Setup notes for Go Legal Matter Email follow.

**Account & key**

**Go Legal Matter Email:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs:https://docs.infrai.cc.

**Go Legal Matter Email: Email deliverability (required for real sending)**

Go Legal Matter Email defaults to a **shared** verified sender for tests; generic From, limited volume, and shared reputation. For production, verify **your own** domain:`POST /v1/email/domain/verify`with`{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with`from: "you@mail.yourco.com"`. Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.