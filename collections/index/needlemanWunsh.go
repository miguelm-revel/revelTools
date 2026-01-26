package index

const (
	MATCH    = 0
	GAP      = 1
	MISMATCH = 1
)

type matrix [][]int

func newMatrix(n, m int) matrix {
	nm := make([][]int, n)
	for i := range nm {
		nm[i] = make([]int, m)
	}
	return nm
}

func score(seq1, seq2 string) int {
	n, m := len(seq1), len(seq2)
	nm := newMatrix(n+1, m+1)
	for i := range seq1 {
		nm[i][0] = GAP * i
	}
	for i := range seq2 {
		nm[0][i] = GAP * i
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			hor := nm[i-1][j] + GAP
			ver := nm[i][j-1] + GAP
			var diag int
			if seq1[i-1] == seq2[j-1] {
				diag = MATCH
			} else {
				diag = MISMATCH
			}
			diag += nm[i-1][j-1]
			nm[i][j] = min(hor, ver, diag)
		}
	}
	return nm[n][m]
}
