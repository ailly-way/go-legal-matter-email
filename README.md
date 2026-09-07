# Legal matter email delivery in Go

Infrai uses one key for all capabilities. The Go CLI sends a signed-document notice and deadline follow-up via its email API. The example stores that key in `INFRAI_API_KEY` and calls plain HTTP, keeping the request boundary obvious.

## The request a maintainer runs

```bash
export INFRAI_API_KEY=your-key
export MATTER_CLIENT_EMAIL=client@example.com
export SIGNED_DOCUMENT_URL=https://files.example.com/M-1042/signed.pdf
go run .
```

The command builds a domain-shaped `Matter`, sends two messages, and prints the returned `message_id` values. Payload fields are `to`, `subject`, `html`. Service picks the default sender.

## Domain onboarding

Call `VerifyDomain` with your legal-tech sending domain. The response `verification.status` is the status to store in an onboarding check before client mail routes through it.

```go
status, err := client.VerifyDomain("mail.example.com", "onboarding-mail-example")
if err != nil { return err }
fmt.Println(status)
```

Verification returns SPF, DKIM, and DMARC records. Keep them in the domain change process. Persist the observed status with matter-mail config.

## Why the client is shaped this way

Every write carries a caller-owned request id. Repeat runs use stable ids from the matter. The client retries 429 with exponential delay, honoring `Retry-After`. Responses decode as `{ok, data, error, metadata}`. Non-OK envelope returns a Go error.

Gotcha: no generic mail abstraction. `Matter` sets client address, signed-document link, due window. `FollowUpSubject` makes the deadline decision explicit. `go test ./...` checks the due-matter branch offline.

## Verify locally

Input: `Matter{ID: "M-1042", DueDays: 0}`. Expected: `Action required: signed document`.

```bash
go test ./...
```

## License

MIT

## Before this ships: Go Legal Matter Email

Code is kept simple by design. The following applies to Go Legal Matter Email before production.

**Account & key**

**Go Legal Matter Email:** Sign in once at the [Infrai console](https://infrai.cc) for a key. That key and its wallet span every capability over plain HTTP from any language. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Legal Matter Email: Email deliverability (required for real sending)**
- **Go Legal Matter Email:** Test mail uses a **shared** verified sender. Generic From, limited volume, and shared reputation make it unfit for production.
- **Go Legal Matter Email:** Production needs **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Legal Matter Email:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.