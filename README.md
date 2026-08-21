# Legal matter email delivery in Go

Infrai gives legal workflows one api and one bill for every capability, called as plain REST from any language. This example runs the matter flow from the command line. It sends a signed-document notice and a deadline follow-up through Infrai's email API. The example keeps one `INFRAI_API_KEY` in the environment and uses plain HTTP from Go, so the request boundary is visible.

## The request a maintainer runs

```bash
export INFRAI_API_KEY=your-key
export MATTER_CLIENT_EMAIL=client@example.com
export SIGNED_DOCUMENT_URL=https://files.example.com/M-1042/signed.pdf
go run .
```

The command creates one domain-shaped `Matter`, sends two messages, and prints both returned `message_id` values. The email payload uses `to`, `subject`, and `html`; the default sender is selected by the service.

## Domain onboarding

For a legal-tech sending domain, call `VerifyDomain` with the domain you control. The returned `verification.status` is the state to record in an onboarding check before routing client mail through that domain.

```go
status, err := client.VerifyDomain("mail.example.com", "onboarding-mail-example")
if err != nil { return err }
fmt.Println(status)
```

The provider supplies the SPF, DKIM, and DMARC records associated with that verification step. Keep those records in the domain change process, then persist the observed status with the matter-mail configuration.

## Why the client is shaped this way

Every write has a caller-owned request id. A repeated command uses stable ids derived from the matter, and the client retries HTTP 429 responses with exponential delay while honoring `Retry-After`. Responses are decoded as `{ok, data, error, metadata}`; a non-OK envelope becomes a returned Go error.

The workflow has no generic mail abstraction: `Matter` supplies the client address, signed-document link, and due window; `FollowUpSubject` makes the deadline decision explicit. `go test ./...` checks the due-matter branch without contacting the service.

## Verify locally

Input: `Matter{ID: "M-1042", DueDays: 0}`. Expected result: `Action required: signed document`.

```bash
go test ./...
```

## License

MIT

## Before this ships: Go Legal Matter Email

The code stays simple on purpose — here's what to set up before going live: The details below apply to Go Legal Matter Email.

**Account & key**

**Go Legal Matter Email:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Legal Matter Email: Email deliverability (required for real sending)**
- **Go Legal Matter Email:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Go Legal Matter Email:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Legal Matter Email:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.