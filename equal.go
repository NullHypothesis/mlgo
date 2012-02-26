package mlgo

import "math"

func ApproximatelyEqual(a, b, epsilon float64) bool {
	var t float64
	if math.Fabs(a) < math.Fabs(b) {
		t = math.Fabs(b) * epsilon
	} else {
		t = math.Fabs(a) * epsilon
	}
	return math.Fabs(a - b) <= t
}

func EssentiallyEqual(a, b, epsilon float64) bool {
	var t float64
	if math.Fabs(a) > math.Fabs(b) {
		t = math.Fabs(b) * epsilon
	} else {
		t = math.Fabs(a) * epsilon
	}
	return math.Fabs(a - b) <= t
}

func DefinitelyGreaterThan(a, b, epsilon float64) bool {
	var t float64
	if math.Fabs(a) < math.Fabs(b) {
		t = math.Fabs(b) * epsilon
	} else {
		t = math.Fabs(a) * epsilon
	}
	return (a - b) > t
}

func DefinitelyLessThan(a, b, epsilon float64) bool {
	var t float64
	if math.Fabs(a) < math.Fabs(b) {
		t = math.Fabs(b) * epsilon
	} else {
		t = math.Fabs(a) * epsilon
	}
	return (a - b) < t
}

