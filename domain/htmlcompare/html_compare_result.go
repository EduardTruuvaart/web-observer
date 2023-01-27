package htmlcompare

type HtmlCompareResult struct {
	State       HtmlCompareState
	Differences []string
	DiffSize    int
}
