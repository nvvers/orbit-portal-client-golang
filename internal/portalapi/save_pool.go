package portalapi

type (
	SavePoolRequest struct {
		PoolName                  string   `json:"poolName"`
		SubjectPattern            string   `json:"subjectPattern"`
		EventTypes                []string `json:"eventTypes"`
		NotificationRetryInterval int      `json:"notificationRetryInterval"`
		NotificationTTL           int      `json:"notificationTTL"`
	}

	SavePoolResponse struct{}
)
