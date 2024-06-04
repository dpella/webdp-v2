package utils

func Map[A, B any](xs []A, f func(A) B) []B {
	out := make([]B, 0)
	for _, val := range xs {
		out = append(out, f(val))
	}
	return out
}

func Reduce[A, B any](identity B, xs []A, fun func(A, B) B) B {
	res := identity
	for _, x := range xs {
		res = fun(x, res)
	}
	return res
}

func Filter[A any](xs []A, fun func(A) bool) []A {
	out := make([]A, 0)
	for _, x := range xs {
		if fun(x) {
			out = append(out, x)
		}
	}
	return out
}
