package slices

func Cut[T comparable](from []T, cut []T) []T {
	set := make(map[T]struct{}, len(cut))
	for _, t := range cut {
		set[t] = struct{}{}
	}

	out := make([]T, 0, len(from))
	for _, t := range from {
		if _, ok := set[t]; !ok {
			out = append(out, t)
		}
	}

	return out
}
