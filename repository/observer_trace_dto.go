package repository

type ObserverTraceDto struct {
	ChatID      int64
	URL         *string
	FileName    *string
	Data        *[]byte
	CssSelector *string
	IsActive    string
}
