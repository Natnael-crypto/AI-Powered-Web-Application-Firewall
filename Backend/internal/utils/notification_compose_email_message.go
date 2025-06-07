package utils

import (
	"fmt"
	"time"
)

func ComposeEmailMessage(username, ruleName, hostName string, clientIPs []string, count int64, startTimestamp int64) string {
	seconds := startTimestamp / 1000
	nanoseconds := (startTimestamp % 1000) * int64(time.Millisecond)
	startTime := time.Unix(seconds, nanoseconds).Format("2006-01-02 15:04:05 MST")

	endTime := time.Now().Format("2006-01-02 15:04:05 MST")

	return fmt.Sprintf(
		`Hello %s,

Your Web Application Firewall has detected suspicious activity.

🛡️ Rule Triggered: %s
📌 Application: %s
🌐 Source IP(s): %v
🔢 Occurrence Count: %d
🕒 Time: %s to %s

Recommended Action: Please review the related logs and ensure appropriate mitigation steps are in place.

Best regards,
Gasha WAF Security Monitoring System`,
		username,
		ruleName,
		hostName,
		clientIPs,
		count,
		startTime,
		endTime,
	)
}
