package domain

type FetchResult struct {
	State      FetchState
	Difference string
	DiffSize   int
}
