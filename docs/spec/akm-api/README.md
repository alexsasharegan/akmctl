# AKM API schemas and fixtures

The JSON Schemas (`schemas/`), the validator (`validate.js`) and the conformance fixtures (`.fixtures/`) of the AKM API. The contract they belong to is [akm-api-spec](../akm-api-spec.md); the standard is [agent-ready-repo-spec](../agent-ready-repo-spec.md).

`.fixtures/make_fixtures.py` regenerates the fixture repositories and the expected responses. The directory starts with a dot so a notes vault does not index the sample records; Prettier ignores it, because the fixtures are byte-exact.

`testdata/api` and `testdata/yaml-lite.json` at the repository root are a copy of `.fixtures/cases`, `.fixtures/repos` and `.fixtures/yaml-lite.json`, kept in step by hand. The Go tests read the copy. Change both together until the work item [261029-x4ws](../../work/261029-x4ws-single-source-api-fixtures.md) makes one the source.
