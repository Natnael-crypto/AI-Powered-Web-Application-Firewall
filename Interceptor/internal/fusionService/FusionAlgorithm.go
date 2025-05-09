package fusionService

func FusionAlgorithm(ruleResult bool, mlResult bool, percent float64) bool {
	if ruleResult {
		return true
	}
	if mlResult {
		return true
	}
	return false
}
