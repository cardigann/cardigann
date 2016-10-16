package releaseinfo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSceneChecker(t *testing.T) {
	for idx, test := range []struct {
		postTitle    string
		isSceneTitle bool
	}{
		// scene releases
		{"South.Park.S04E13.Helen.Keller.The.Musical.720p.WEBRip.AAC2.0.H.264-GC", true},
		{"Robot.Chicken.S07E02.720p.WEB-DL.DD5.1.H.264-pcsyndicate", true},
		{"Archer.2009.S05E06.Baby.Shower.720p.WEB-DL.DD5.1.H.264-iT00NZ", true},
		{"30.Rock.S04E17.720p.HDTV.X264-DIMENSION", true},
		{"30.Rock.S04.720p.HDTV.X264-DIMENSION", true},

		// not scene releases
		{"S08E05 - Virtual In-Stanity [WEBDL-720p]", false},
		{"S08E05 - Virtual In-Stanity.With.Dots [WEBDL-720p]", false},
		{"Something", false},
		{"86de66b7ef385e2fa56a3e41b98481ea1658bfab", false},
		{"30.Rock.S04E17.720p.HDTV.X264", false},   // no group
		{"S04E17.720p.HDTV.X264-DIMENSION", false}, // no series title
		{"30.Rock.S04E17-DIMENSION", false},        // no quality
	} {
		result, err := IsSceneTitle(test.postTitle)
		require.NoError(t, err)

		if test.isSceneTitle {
			require.True(t, result,
				fmt.Sprintf("Row %d should be a scene release with no error", idx+1))
		} else {
			require.False(t, result,
				fmt.Sprintf("Row %d should not be a scene release", idx+1))
		}
	}
}
