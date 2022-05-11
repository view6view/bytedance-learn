package benchstruct

func EmptyStructMap(n int) {
	m := make(map[int]struct{})

	for i := 0; i < n; i++ {
		m[i] = struct{}{}
	}
}

func BoolMap(n int) {
	m := make(map[int]bool)

	for i := 0; i < n; i++ {
		m[i] = false
	}
}
