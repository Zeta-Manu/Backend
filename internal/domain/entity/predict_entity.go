package entity

type ProcessedAvg struct {
	Key     string
	Average float64
	Sum     float64
	Count   int
}

type PredictResponse struct {
	Class      string  `json:"class"`
	Translated string  `json:"translated"`
	Average    float64 `json:"average"`
	Sum        float64 `json:"sum"`
	Count      int     `json:"count"`
}
