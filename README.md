# utopia-client

Existing Kino Cat UTOPIA integration. Research, verified site behavior and the
strict series subscription design: [docs/research.md](docs/research.md).

## Runtime

Go 1.23.2. Existing environment variables remain `UTP_USER`, `UTP_PASSWORD`,
`BROKER_IP`, `BROKER_PORT`. Configure credentials only on the deployment host.
The Cloud Browser login used during research is not transferred to this service.
Automated Go login through the site's protection has not been verified.

Existing subjects:

* `UTOPIA_GET_TORRENT_INFO`: one positive torrent ID in `arguments`.
* `UTOPIA_GET_TORRENT_FILE`: one positive torrent ID; saves the response in
  `/root/torrent_files`. This is a live action, not a dry-run command.

Commands accept `user_id` and legacy `chat_id`; conflicting IDs are rejected.
Outgoing messages still use the pinned core helper and its `user_id` contract.
Check deployed bot/text versions before upgrading: accessible master versions
use different envelopes and do not contain UTOPIA command routes.

`ParseSimilar` and `GetSimilarPageByID` now support grouped similar pages,
but there is no new NATS listener or Telegram route for them in this patch.
No subscriptions worker or automatic download behavior is enabled.

## Verification without accounts or downloads

```sh
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

Tests use reduced sanitized HTML and local HTTP test servers only.
