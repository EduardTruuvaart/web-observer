package domain

type ObserverTrace struct {
	ChatID      int64
	URL         *string
	FileName    *string
	Data        *[]byte
	CssSelector *string
	IsActive    bool
}
