package entity

type LinkStatus string

func (l LinkStatus) String() string { return string(l) }

var (
	ActiveStatus   = LinkStatus("active")
	ExpiredStatus  = LinkStatus("expired")
	DisabledStatus = LinkStatus("disabled")
	BlockedStatus  = LinkStatus("blocked")
)
