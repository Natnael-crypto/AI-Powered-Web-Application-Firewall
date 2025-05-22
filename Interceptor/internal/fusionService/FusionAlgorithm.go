package fusionService

func FusionAlgorithm(ruleResult bool, mlResult bool, Normal float64, Anomaly float64) bool {
	if ruleResult {
		return true
	}
	if mlResult {
		return true
	}
	return false
}
