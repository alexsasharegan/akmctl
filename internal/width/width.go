// Package width measures the display width of a string the way Prettier 3 does
// when it pads Markdown tables: ASCII counts its length, East Asian wide and
// fullwidth characters count 2, an emoji sequence counts 2 (1 for a bare
// text-style emoji), and control characters, combining marks U+0300 to U+036F
// and variation selectors count 0.
//
// The tables in tables.go are generated from Unicode data (see gen/main.go);
// testdata holds the cell widths Prettier chose, which the tests compare with.
package width

import "sort"

//go:generate go run gen/main.go data/EastAsianWidth.txt data/emoji-data.txt

func inRanges(t [][2]rune, r rune) bool {
	i := sort.Search(len(t), func(i int) bool { return t[i][1] >= r })
	return i < len(t) && t[i][0] <= r
}

func isRI(r rune) bool { return r >= 0x1F1E6 && r <= 0x1F1FF }

// runeWidth is the width of one code point outside an emoji sequence.
func runeWidth(r rune) int {
	switch {
	case r <= 0x1F || r >= 0x7F && r <= 0x9F, r >= 0x300 && r <= 0x36F, r >= 0xFE00 && r <= 0xFE0F:
		return 0
	case r < 0x1100:
		return 1
	case inRanges(wideRanges, r):
		return 2
	}
	return 1
}

// emojiCluster recognises an emoji sequence starting at rs[i]: the number of
// runes it uses (0 when there is none) and its width.
func emojiCluster(rs []rune, i int) (n, w int) {
	r := rs[i]
	if (r >= '0' && r <= '9' || r == '#' || r == '*') && i+1 < len(rs) {
		j := i + 1
		if rs[j] == 0xFE0F {
			j++
		}
		if j < len(rs) && rs[j] == 0x20E3 {
			return j + 1 - i, 2 // keycap
		}
		return 0, 0
	}
	if isRI(r) {
		if i+1 < len(rs) && isRI(rs[i+1]) {
			return 2, 2 // flag
		}
		return 0, 0
	}
	if !inRanges(pictoRanges, r) {
		return 0, 0
	}
	// element := base (skin tone | FE0F [20E3] | tag sequence)? ; elements join with a ZWJ
	j := i + 1
	for {
		switch {
		case j < len(rs) && rs[j] >= 0x1F3FB && rs[j] <= 0x1F3FF:
			j++
		case j < len(rs) && rs[j] == 0xFE0F:
			j++
			if j < len(rs) && rs[j] == 0x20E3 {
				j++
			}
		case j < len(rs) && rs[j] >= 0xE0020 && rs[j] <= 0xE007E:
			k := j
			for k < len(rs) && rs[k] >= 0xE0020 && rs[k] <= 0xE007E {
				k++
			}
			if k < len(rs) && rs[k] == 0xE007F {
				j = k + 1
			}
		}
		if j+1 < len(rs) && rs[j] == 0x200D && inRanges(pictoRanges, rs[j+1]) {
			j += 2
			continue
		}
		break
	}
	if j-i == 1 && !inRanges(emojiPresRanges, r) && !inRanges(wideRanges, r) {
		return 1, 1 // a bare text-style emoji is narrow
	}
	return j - i, 2
}

// String returns the display width of s.
func String(s string) int {
	ascii := true
	for i := 0; i < len(s); i++ {
		if s[i] < 0x20 || s[i] > 0x7F {
			ascii = false
			break
		}
	}
	if ascii {
		return len(s)
	}
	rs := []rune(s)
	w := 0
	for i := 0; i < len(rs); {
		if n, cw := emojiCluster(rs, i); n > 0 {
			w += cw
			i += n
			continue
		}
		w += runeWidth(rs[i])
		i++
	}
	return w
}
