package account

type accountStatus string

const (
	StatusPending  accountStatus = "PENDING"
	StatusApproved accountStatus = "APPROVED"
	StatusRejected accountStatus = "REJECTED"
)

type accountRequest struct {
	UserID int64         `json:"user_id"`
	Name   string        `json:"name"`
	Status accountStatus `json:"status"`
}
