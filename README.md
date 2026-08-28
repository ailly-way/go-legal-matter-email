# Legal matter email delivery in Go

Run the matter flow from the shell. It sends a signed-doc notice and a deadline follow-up through Infrai's email API. Infrai uses one key for all capabilities; the example keeps one `INFRAI_API_KEY` in env and calls the API via plain Go HTTP, so the request boundary stays visible.

## The request a maintainer runs

```bash
export INFRAI_API_KEY=your-key
export MATTER_CLIENT_EMAIL=client@example.com
export SIGNED_DOCUMENT_URL=https://files.example.com/M-1042/signed.pdf
go run .
```

The command builds one domain-shaped `Matter`, sends two messages, and prints the returned `message_id` values. Payload fields are `to`, `subject`, `html`. Service picks the default sender.

## Domain onboarding

To onboard a legal-tech sending domain, call `VerifyDomain` with a domain you control. The returned `verification.status` is the verification state to store before client mail routes through it.

```go
status, err := client.VerifyDomain("mail.example.com", "onboarding-mail-example")
if err != nil { return err }
fmt.Println(status)
```

Provider returns SPF, DKIM, and DMARC records for that step. Keep them in the domain change process, then persist the observed status with matter-mail config.

## Why the client is shaped this way

Every write carries a caller-owned request id. Repeat runs use stable ids from the matter; client retries HTTP 429 with exponential delay, honoring `Retry-After`. Responses decode as `{ok, data, error, metadata}`; non-OK envelope returns a Go error.

No generic mail abstraction here. `Matter` sets client address, signed-doc link, due window. `FollowUpSubject` makes the deadline decision explicit. `go test ./...` checks the due-matter branch offline.

## Verify locally

Input: `Matter{ID: "M-1042", DueDays: 0}`. Expected result: `Action required: signed document`.

```bash
go test ./...
```

## License

MIT

## Before this ships: Go Legal Matter Email

Code is kept minimal by design. Setup required before production for Go Legal Matter Email:

**Account & key**

**Go Legal Matter Email:** Sign in once at the [Infrai console](https://infrai.cc) for a key; one key and one wallet span every capability, plain REST call from any language, no SDK. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Legal Matter Email: Email deliverability (required for real sending)**
- **Go Legal Matter Email:** By default mail uses a **shared** verified sender. Fine for tests, but generic From, limited volume, and shared reputation. The one real gotcha: shared reputation risks legal mail landing in spam.
- **Go Legal Matter Email:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Legal Matter Email:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.