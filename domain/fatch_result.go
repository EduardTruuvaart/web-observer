package domain

type FetchResult string

const (
	Updated   FetchResult = "UPDATED"
	Unchanged FetchResult = "UNCHANGED"
)
