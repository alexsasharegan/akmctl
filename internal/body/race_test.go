//go:build race

package body

// raceBuild is true when the tests run under the race detector, which is 20
// times slower here: the generated batteries run an eighth of their documents
// (work item 261108-rwqj, finding F13). CI runs the full count without -race.
const raceBuild = true
