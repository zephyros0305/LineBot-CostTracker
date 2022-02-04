package main

func SumUintSlice(slice []uint) uint {
	var sum uint
	for _, v := range slice {
		sum += v
	}

	return sum
}
