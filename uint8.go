package main

func add(a, b uint8) uint8 {
	if a > 255-b {
		return 255
	}
	return a + b
}

func sub(a, b uint8) uint8 {
	if a < b {
		return 0
	}
	return a - b
}
