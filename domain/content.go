package domain

type Content struct {
	URL         string
	Data        []byte
	CssSelector string
	IsActive    bool
}
