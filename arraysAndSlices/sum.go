package arraysAndSlices

func Sum(numbers [5]int) int{
	var sum int
	for _, value := range numbers{
		sum += value
	}
	return sum
}
