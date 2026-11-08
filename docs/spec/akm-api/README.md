# AKM API schemas and fixtures

The JSON Schemas (`schemas/`), the validator (`validate.js`) and the conformance fixtures (`.fixtures/`) of the AKM API. The contract they belong to is [akm-api-spec](../akm-api-spec.md); the standard is [agent-ready-repo-spec](../agent-ready-repo-spec.md).

`.fixtures/` is the one source of the conformance fixtures. Edit the cases and sample repositories there. The directory starts with a dot so a notes vault does not index the sample records; Prettier ignores it, because the fixtures are byte-exact.

The Go tests read a copy: `testdata/api` and `testdata/yaml-lite.json` at the repository root. After a change run `scripts/sync-fixtures.sh` and commit both; `TestFixtureCopiesAgree` (`go test ./internal/api`) fails, naming the files, when the copy differs from the source in any file or byte (work item [261029-x4ws](../../work/done/261029-x4ws-single-source-api-fixtures.md)). `make_fixtures.py` stays in the source only.

`.fixtures/make_fixtures.py` is a frozen reference for how the two sample repositories were built, not a generator to run: it no longer knows every case (`bases`, `show-3`, `ls-7`, the newer `close` and `init` cases), and running it rewrites cases that are kept by hand. Add or change a case by editing its JSON file (write only `name`, `repo`, `request` and `"expected": {"operation": "<op>"}` for a new one), run `scripts/sync-fixtures.sh`, then `go test ./internal/api -run TestConformance -update`: it rewrites `expected` (and `after`, `would_be`) from the tool's own output in both copies, and the diff is the review. Check the new response by eye; the test passes whatever the tool says.
