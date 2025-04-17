package fox

import (
	. "maze/internal/global"
	"testing"
)

func TestUnloopingLine(t *testing.T) {

	type testCase struct {
		in, out    Route
		isPositive bool
	}

	newCase := func(in, out [][2]int, validOptional ...bool) testCase {
		valid := true
		if len(validOptional) > 0 {
			valid = validOptional[0]
		}
		return testCase{NewRouteFromSlice(in), NewRouteFromSlice(out), valid}
	}

	in1 := [][2]int{
		{7, 6},
		{7, 3},
		{4, 3}, // x = redundant point
		{10, 3},
		{10, 1},
		{12, 1},
		{12, 23},
		{10, 23},
		{10, 21},
		{4, 21}, // x
		{7, 21},
		{7, 18},
	}
	out1 := [][2]int{
		{7, 6},
		{7, 3},
		{10, 3},
		{10, 1},
		{12, 1},
		{12, 23},
		{10, 23},
		{10, 21},
		{7, 21},
		{7, 18},
	}

	in2 := [][2]int{
		{0, 0},
		{9, 0},
		{9, 21},
		{4, 21}, // x
		{7, 21}, // x
		{0, 21},
		{0, -1},
		{-1, -1},
	}
	out2 := [][2]int{
		{0, 0},
		{9, 0},
		{9, 21},
		{0, 21},
		{0, -1},
		{-1, -1},
	}

	inTail := [][2]int{
		{0, 0},
		{9, 0},
		{9, 21},
		{4, 21}, // x
		{7, 21},
		{-1, 21},
	}
	outTail := [][2]int{
		{0, 0},
		{9, 0},
		{9, 21},
		{-1, 21},
	}

	testCases := []testCase{
		newCase([][2]int{}, [][2]int{}),
		newCase([][2]int{{7, 6}, {7, 3}}, [][2]int{{7, 6}, {7, 3}}),
		newCase(in1, out1),
		newCase(in2, out2),
		newCase(inTail, outTail),
	}

	for _, tc := range testCases {
		t.Run("unloopingLine()", func(t *testing.T) {
			resultRoute := unloopingLine(tc.in)
			if tc.isPositive != tc.out.Eq(&resultRoute) {
				f := "Failure (F(IN) = OUT), OUT ≠ RESULT):\n    IN: %v\n   OUT: %v\nRESULT: %v"
				t.Errorf(f, tc.in, tc.out, resultRoute)
			}
		})
	}
}

func TestUnlooping(t *testing.T) {

	type testCase struct {
		in, out    Route
		isPositive bool
	}

	newCase := func(in, out [][2]int, validOptional ...bool) testCase {
		valid := true
		if len(validOptional) > 0 {
			valid = validOptional[0]
		}
		return testCase{NewRouteFromSlice(in), NewRouteFromSlice(out), valid}
	}

	testCases := []testCase{
		newCase([][2]int{}, [][2]int{}),
		newCase([][2]int{{-4, -3}}, [][2]int{{-4, -3}}),
		newCase([][2]int{{1, 1}, {1, 7}, {1, 1}}, [][2]int{{1, 1}}),
		newCase(
			[][2]int{{7, 6}, {7, 3}, {8, 3}, {5, 3}, {5, 2}, {2, 5}, {2, 7}, {7, 6}, {1, 6}, {1, 1}},
			[][2]int{{7, 6}, {1, 6}, {1, 1}},
		),
	}

	for _, tc := range testCases {
		t.Run("unlooping()", func(t *testing.T) {
			resultRoute := unlooping(tc.in)
			if tc.isPositive != tc.out.Eq(&resultRoute) {
				f := "Failure (F(IN) = OUT), OUT ≠ RESULT):\n    IN: %v\n   OUT: %v\nRESULT: %v"
				t.Errorf(f, tc.in, tc.out, resultRoute)
			}
		})
	}
}

func TestShortByIntersection(t *testing.T) {

	rA := &Route{}
	rA.Unserialize("[8 9] [10 9] [10 8] [10 10] [10 7] [10 11] [16 11] [16 6] [24 6] [24 5] [25 5] [25 6] [26 6] [26 8] [24 8] [24 9] [22 9] [22 8] [19 8] [19 9]")

	rB := &Route{}
	rB.Unserialize("[8 9] [8 11] [26 11] [26 8] [24 8] [24 9] [22 9] [22 8] [20 8] [19 8] [18 8] [18 9] [19 9]")

	rC := &Route{}
	rC.Unserialize("[8 9] [8 11] [26 11] [26 8] [24 8] [24 9] [22 9] [22 8] [19 8] [19 9]")

	t.Run("shortByIntersection()", func(t *testing.T) {
		result := shortByIntersection(*rA, *rB)
		if !rC.Eq(&result) {
			f := "Failure (F(r1,r2) = OUT), OUT ≠ RESULT):\n   IN1: %v\n   IN2: %v\n   OUT: %v\nRESULT: %v"
			t.Errorf(f, *rA, *rB, *rC, result)
		}
	})
}
