package domain

type FetchState string

const (
	Updated           FetchState = "UPDATED"
	Unchanged         FetchState = "UNCHANGED"
	NewContentIsAdded FetchState = "NEW_CONTENT_IS_ADDED"
)
