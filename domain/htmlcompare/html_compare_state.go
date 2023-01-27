package htmlcompare

type HtmlCompareState string

const (
	Identical                 HtmlCompareState = "IDENTICAL"
	Different                 HtmlCompareState = "DIFFERENT"
	SelectionNotFoundInSource HtmlCompareState = "SELCETION_NOT_FOUND_IN_SOURCE"
	SelectionNotFoundInTarget HtmlCompareState = "SELCETION_NOT_FOUND_IN_TARGET"
)
