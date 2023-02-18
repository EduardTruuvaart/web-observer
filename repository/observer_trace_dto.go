package repository

type ObserverTraceDto struct {
	ChatID      int64
	URL         *string
	FileName    *string
	CssSelector *string
	IsActive    string
}
