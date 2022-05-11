package test

func test1(s string) {
	// bad
	for index := 0; index < len(s); index++ {

	}

	// good
	for i := 0; i < len(s); i++ {

	}
}

type Client struct {
}
