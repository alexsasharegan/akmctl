// Package conformance runs the API conformance fixtures
// (docs/spec/akm-api/.fixtures, copied to testdata/api) through each binding of
// the tool: the in-process API, the HTTP handler and, for reads, the CLI with
// --format json. The code is in the test files; ADR 261108-15dw describes the
// fixture format (matchers, replayed git history, bindings).
package conformance
