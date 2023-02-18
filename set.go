package gopkg

func union[K comparable, V any](m ...map[K]V) map[K]V {

	ret := make(map[K]V)
	for _, n := range m {
		for key, value := range n {
			ret[key] = value
		}
	}
	return ret
}

func complement[K comparable, V any](a, b map[K]V) map[K]V {

	aComplementB := make(map[K]V)
	for aK, aV := range a {
		if _, ok := b[aK]; !ok {
			aComplementB[aK] = aV
		}
	}
	return aComplementB
}
