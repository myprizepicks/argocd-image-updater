package tag

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/argoproj-labs/argocd-image-updater/pkg/log"
)

// semverCollection is a replacement for semver.Collection that breaks version
// comparison ties through a lexical comparison of the original version strings.
// Using this, instead of semver.Collection, when sorting will yield
// deterministic results that semver.Collection will not yield.
type modifiedSemverCollection []*semver.Version

// Len returns the length of a collection. The number of Version instances
// on the slice.
func (s modifiedSemverCollection) Len() int {
	return len(s)
}

// Less is needed for the sort interface to compare two Version objects on the
// slice. If checks if one is less than the other.
func (s modifiedSemverCollection) Less(i, j int) bool {
	baseI, suffixNumI := extractBaseAndSuffix(s[i].Original())
	baseJ, suffixNumJ := extractBaseAndSuffix(s[j].Original())

	log.Debugf("Base I %s | Base J %s", baseI, baseJ)

	// Compare base versions without suffixes
	baseVersionI, _ := semver.NewVersion(baseI)
	baseVersionJ, _ := semver.NewVersion(baseJ)
	comp := baseVersionI.Compare(baseVersionJ)
	if comp != 0 {
		return comp < 0
	}

	return suffixNumI < suffixNumJ
}

// Swap is needed for the sort interface to replace the Version objects
// at two different positions in the slice.
func (s modifiedSemverCollection) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func extractBaseAndSuffix(version string) (string, int) {
	// Regex to match a base version and an optional suffix with .n
	re := regexp.MustCompile(`^([^-]+)(?:-(.*))?`)
	matches := re.FindStringSubmatch(version)

	baseVersion := matches[1]
	suffixNum := 0

	if len(matches) > 2 && matches[2] != "" {
		// Check for a trailing .n in the suffix, and extract the number if present
		suffixParts := strings.Split(matches[2], ".")
		if len(suffixParts) > 1 {
			num, err := strconv.Atoi(suffixParts[len(suffixParts)-1])
			if err == nil {
				suffixNum = num
			}
		}
	}

	return baseVersion, suffixNum
}
