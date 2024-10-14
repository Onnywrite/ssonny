package must

func Ok1(err error) {
	if err != nil {
		panic(err)
	}
}

func Ok2[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

func Ok3[T1 any, T2 any](value1 T1, value2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}

	return value1, value2
}
