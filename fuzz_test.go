package heraldhelp

import "testing"

func FuzzMdTable(f *testing.F) {
	f.Add("Flag|Type|Default|Description", "--|--|--|--", "--output|string|stdout|Output destination")

	f.Fuzz(func(t *testing.T, header, sep, row string) {
		// Build a rows slice from the fuzz inputs.
		splitRow := func(s string) []string {
			var cells []string
			current := ""
			for _, c := range s {
				if c == '|' {
					cells = append(cells, current)
					current = ""
				} else {
					current += string(c)
				}
			}
			cells = append(cells, current)
			return cells
		}

		rows := [][]string{splitRow(header), splitRow(row)}
		// mdTable must not panic on any input.
		_ = mdTable(rows)
	})
}

func FuzzMdTableRow(f *testing.F) {
	f.Add("cell1", "cell2", "cell3", 10, 10, 10)

	f.Fuzz(func(t *testing.T, c1, c2, c3 string, w1, w2, w3 int) {
		// Ensure widths are non-negative to avoid pathological repeats.
		if w1 < 0 {
			w1 = 0
		}
		if w2 < 0 {
			w2 = 0
		}
		if w3 < 0 {
			w3 = 0
		}
		// Cap widths to prevent excessive memory allocation.
		if w1 > 1000 {
			w1 = 1000
		}
		if w2 > 1000 {
			w2 = 1000
		}
		if w3 > 1000 {
			w3 = 1000
		}

		row := []string{c1, c2, c3}
		widths := []int{max(w1, len(c1)), max(w2, len(c2)), max(w3, len(c3))}
		// mdTableRow must not panic on any input.
		_ = mdTableRow(row, widths)
	})
}
