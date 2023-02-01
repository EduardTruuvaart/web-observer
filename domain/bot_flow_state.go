package domain

type BotFlowState string

const (
	NotStarted        BotFlowState = "NOT_STARTED"
	URLRequsted       BotFlowState = "URL_REQUESTED"
	SelectorRequested BotFlowState = "SELECTOR_REQUESTED"
)
