# UTOPIA repair plan — 2026-11-15

Goal: repair the existing HTML integration using observed UNIT3D 9.2.0 pages;
preserve subjects and JSON fields, and document a fail-closed subscription design.
The user's brief authorizes implementation after research; no deployment or live
downloads are part of this plan. Work is in a fresh checkout on a separate branch.

1. Add reduced, sanitized fixtures based on observed episode and pack markup.
   Verify old parser fails for release title, complete file list and language.
   Repair selectors; keep raw dimensions separate from the site's quality tier.
   Retain video codec/HDR and avoid inventing studios from generic track labels.
2. Implement ParseSimilar from grouped table rows, including closed details,
   inherited row-group type and season/episode context. Preserve old signatures
   with wrappers and add a string-ID navigator for category.TMDB IDs.
3. Guard listener IDs and HTML formatting; accept both existing user_id and
   legacy chat_id envelopes, reject conflicting identities. Keep subjects and
   core output behavior unchanged. Do not claim missing cross-repo routing works.
4. Add bounded HTTP timeouts and reject login pages; propagate initialization
   failures instead of terminating the process. Validate downloaded bytes locally
   before writing; never exercise the live download route.
5. Run all tests, race tests, vet and build; review diff. Record research, API
   verification limits and subscription algorithm in docs/research.md.

Review focus: absent MediaInfo; non-media files; subtitles becoming audio;
rowspan inheritance; negative/overflow IDs; chat identity conflicts; expired
sessions; HTML returned as .torrent; source names containing markup.

API replacement and a durable subscription worker are intentionally deferred
until the actual authorized API response and downloader acknowledgment contract
are verified. No automatic download hook will be added.

## Execution record

All five repair tasks completed. The baseline full suite failed vet on the
existing Printf `%w`; repaired. New parser tests were first observed failing
on release title, files, language, quality and false studio extraction. Missing
MediaInfo initially reproduced a nil panic in formatting. Similar-page and
HTTP-login tests were also observed failing before their fixes.

Final verification: 15 tests; `go test -race ./...`, `go vet ./...` and
`go build` succeeded with Go 1.23.2. No live NATS, tracker API, torrent download,
container deployment or Transmission request was exercised.

Independent code review found a shared-directory permission regression from
CreateTemp's 0600 mode. A regression test reproduced it; explicit 0644 before
rename fixed it, and the full race suite/vet/build passed afterward.

Implementation decisions: retain HTML until actual API authorization can be
verified; keep outgoing core envelope and subjects; do not add unverified
Telegram routing or a download worker. Their limits are recorded in research.md.
Use `[skip ci]` when saving the branch because its existing push workflow would
otherwise publish `beta-latest`. No deployment is authorized by this repair.
