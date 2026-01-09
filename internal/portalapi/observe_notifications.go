package portalapi

type (
	ObserveNotificationsRequest struct {
		PoolName string `json:"poolName"`
	}

	ObserveNotificationsReleaseConnectionResponse struct {
		Status string `json:"status"`
	}
)
