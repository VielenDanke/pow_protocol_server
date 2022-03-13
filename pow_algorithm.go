package main

func fib(n int) int {
	f, s := 0, 1
	for i := 0; i < n; i++ {
		f, s = s, f+s
	}
	return f
}
