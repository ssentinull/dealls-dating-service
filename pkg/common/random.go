package common

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	alfanums     = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	LeadIDPrefix = "520"
)

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alfanums[rand.Intn(len(alfanums))]
	}
	return string(b)
}

func RandSplit[T any](a []T, chunks int) [][]T {
	res := make([][]T, 0)

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	rng.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})

	chunkSize := len(a) / chunks
	remainder := len(a) % chunks // Get the remainder

	currIdx := 0
	for i := 0; i < chunks; i++ {
		endIdx := currIdx + chunkSize

		if i < remainder { // If there's a remainder, distribute it evenly
			endIdx++
		}

		res = append(res, a[currIdx:endIdx])

		currIdx = endIdx
	}

	return res
}
