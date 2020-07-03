package metrics

type FloatArray []float64

// Max find the max item in FloatArray
func (array FloatArray) Max() float64 {
	if len(array) == 0 {
		return 0
	}
	max := array[0]
	for _, v := range array[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Min find the min item in FloatArray
func (array FloatArray) Min() float64 {
	if len(array) == 0 {
		return 0
	}
	min := array[0]
	for _, v := range array[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Avg calc the average value of FloatArray
func (array FloatArray) Avg() float64 {
	if len(array) == 0 {
		return 0
	}
	var sum float64 = 0
	for _, v := range array[1:] {
		sum += v
	}
	return sum / float64(len(array))
}
