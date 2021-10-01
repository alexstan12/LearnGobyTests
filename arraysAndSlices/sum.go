package arraysAndSlices

func Sum(numbers [5]int) int {
	var sum int
	for _, value := range numbers {
		sum += value
	}
	return sum
}

func SumSlices(numbers []int) int {
	var sum int
	for _, value := range numbers {
		sum += value
	}
	return sum
}

func SumAll(slices ...[]int) (sums []int) {

	for _, slice := range slices {
		sums = append(sums, SumSlices(slice))
	}
	return sums
}

func SumAllTails(slices ...[]int) (sums []int) {
	for _, slice := range slices {
		if len(slice) == 0 {
			sums = append(sums, 0)
		} else {
			sums = append(sums, SumSlices(slice[1:]))
		}
	}
	return sums
}
