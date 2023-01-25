package domain

type FetchResult string

const (
	Updated           FetchResult = "UPDATED"
	Unchanged         FetchResult = "UNCHANGED"
	NewContentIsAdded FetchResult = "NEW_CONTENT_IS_ADDED"
)
