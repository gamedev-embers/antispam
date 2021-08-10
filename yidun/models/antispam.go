package models

type AntiSpam struct {
	TaskID       string           `json:"taskId"`
	Action       int              `json:"action"`
	CensorType   int              `json:"censorType"`
	Lang         []string         `json:"lang"`
	IsRelatedHit bool             `json:"isRelatedHit"`
	Labels       []*AntiSpamLabel `json:"labels"`
}

type AntiSpamLabel struct {
	Label   int `json:"label"`
	Level   int `json:"level"`
	Details struct {
		Hints []*Hint `json:"hints"`
		// HitInfos []string
	} `json:"details"`
}

type Hint struct {
	Hint      string          `json:"hint"`
	Positions []*HintPosition `json:"positions"`
}

type HintPosition struct {
	PositionType int `json:"positionType"`
	StartPos     int `json:"startPos"`
	EndPos       int `json:"endPos"`
}
