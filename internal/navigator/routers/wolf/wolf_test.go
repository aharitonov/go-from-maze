package wolf

import (
	. "maze/internal/global"
	"testing"
)

func TestShortByInterval(t *testing.T) {

	rA := &Route{}
	rA.Unserialize("[0 0] [2 0] [2 3] [9 3] [9 0]")
	//             [9 0]-[12 0]
	//             [9 3]-[9 0]
	//             [2 3]-[9 3]
	//             [2 0]-[2 3]
	//             [0 0]-[2 0]

	rB := &Route{}
	rB.Unserialize("[0 0] [2 0] [2 9] [4 9] [4 12] [7 12] [7 1] [9 1] [9 0]")
	// (0,0) (2,0)
	// (2,0) (2,9)
	// (2,9) (4,9)
	// (4,9) (4,12)
	// (4,12) (7,12)
	// (7,12) (7,1)
	// (7,1) (9,1)
	// (9,1) (9,0)

	rC := &Route{}
	rC.Unserialize("[0 0] [2 0] [2 3] [9 3] [9 0]")

	t.Run("shortByInterval()", func(t *testing.T) {
		result := shortByInterval(rA, rB)
		if !rC.Eq(&result) {
			f := "Failure (EXPECT = F(r1,r2)), EXPECT ≠ RESULT):\n" +
				"   IN1: %v\n   IN2: %v\nEXPECT: %v\nRESULT: %v" +
				"\nFRAMES A: %v" +
				"\nFRAMES B: %v"
			t.Errorf(f, *rA, *rB, *rC, result, rA.ToFrames(), rB.ToFrames())
		}
	})
}
