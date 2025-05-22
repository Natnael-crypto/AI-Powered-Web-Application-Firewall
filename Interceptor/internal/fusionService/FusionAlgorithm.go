package fusionService

import (
	"math"
)

func FusionAlgorithm(ruleResult bool, normalScore, anomalyScore float64) bool {
	var rule float64
	if ruleResult {
		rule = 1.0
	} else {
		rule = 0.0
	}

	weightRule := 4.80215439
	weightNormal := -15.02618438
	weightAnomaly := 14.62489455
	bias := -0.77542936

	logit := rule*weightRule + normalScore*weightNormal + anomalyScore*weightAnomaly + bias

	prob := 1.0 / (1.0 + math.Exp(-logit))

	return prob > 0.5
}
