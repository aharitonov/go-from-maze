package global

import "testing"

func TestRoute_SerializeAndUnSerialize(t *testing.T) {

	type testCase struct {
		in         [][2]int
		isPositive bool
	}

	newCase := func(in [][2]int, validOptional ...bool) testCase {
		valid := true
		if len(validOptional) > 0 {
			valid = validOptional[0]
		}
		return testCase{in, valid}
	}

	testCases := []testCase{
		newCase([][2]int{}),
		newCase([][2]int{{-3, -2}}),
		newCase([][2]int{{3, 4}, {0, 4}}),
	}

	for _, tc := range testCases {
		t.Run("TestSerializeAndUnSerialize()", func(t *testing.T) {

			r1 := NewRouteFromSlice(tc.in)
			ser := r1.Serialize()

			r2 := (&Route{}).Unserialize(ser)
			result := r1.Eq(r2)

			if tc.isPositive != result {
				t.Errorf("Failure on %#v", tc.in)
			}
		})
	}
}

func TestRouteFrame_Eq(t *testing.T) {

	type testCase struct {
		in, out    RouteFrame
		isPositive bool
	}

	newCase := func(in, out RouteFrame, validOptional ...bool) testCase {
		valid := true
		if len(validOptional) > 0 {
			valid = validOptional[0]
		}
		return testCase{in, out, valid}
	}

	testCases := []testCase{
		newCase(RouteFrame{}, RouteFrame{}, true),
		newCase(RouteFrameNew(0, 0, 0, 0), RouteFrameNew(0, 0, 0, 0)),
		newCase(RouteFrameNew(3, 4, 2, 4), RouteFrameNew(3, 4, 2, 4)),
		newCase(RouteFrameNew(-1, -2, -3, -4), RouteFrameNew(-1, -2, -3, -4)),

		newCase(RouteFrameNew(0, 0, 0, 0), RouteFrameNew(0, 0, 0, 1), false),
		newCase(RouteFrameNew(0, 0, 0, 0), RouteFrameNew(0, 0, 1, 0), false),
		newCase(RouteFrameNew(0, 0, 0, 0), RouteFrameNew(0, 1, 0, 0), false),
		newCase(RouteFrameNew(0, 0, 0, 0), RouteFrameNew(1, 0, 0, 0), false),
	}

	for _, tc := range testCases {
		t.Run("TestRouteFrame_Eq()", func(t *testing.T) {
			result := tc.in.Eq(tc.out)
			if result != tc.isPositive {
				t.Errorf("Failure on %#v", tc.in)
			}
		})
	}
}

func TestRoute_ToFrames(t *testing.T) {

	type testCase struct {
		in         *Route
		out        FrameList
		isPositive bool
	}

	newCase := func(in *Route, out []RouteFrame, validOptional ...bool) testCase {
		valid := true
		if len(validOptional) > 0 {
			valid = validOptional[0]
		}
		return testCase{in, out, valid}
	}

	testCases := []testCase{
		newCase(&Route{}, []RouteFrame{}),
		newCase((&Route{}).Unserialize("[-3 2]"), []RouteFrame{}),
		newCase((&Route{}).Unserialize("[3 4] [0 4] [2 4]"), []RouteFrame{
			RouteFrameNew(3, 4, 0, 4),
			RouteFrameNew(0, 4, 2, 4),
		}),
	}

	for _, tc := range testCases {
		t.Run("TestRoute_ToFrames()", func(t *testing.T) {
			result := tc.in.ToFrames().Eq(tc.out)
			if result != tc.isPositive {
				t.Errorf("Failure on %#v", tc.in)
			}
		})
	}
}
