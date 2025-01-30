package fusionService

func FusionAlgorithm(ruleResult bool, ruleReason string, mlResult bool, mlReason string) (bool, string) {
	if ruleResult {
		return true, ruleReason
	}
	if mlResult {
		return true, mlReason
	}
	return false, ""
}
