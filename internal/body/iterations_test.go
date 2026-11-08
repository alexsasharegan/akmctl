package body

// iterations is the number of generated documents a battery runs: n, or a
// eighth of n under the race detector.
func iterations(n int) int {
	if raceBuild {
		return n / 8
	}
	return n
}
