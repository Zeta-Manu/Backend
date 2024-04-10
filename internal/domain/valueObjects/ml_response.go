package valueobjects

type MlResponse struct {
	Results struct {
		Raw []struct {
			Class string  `json:"class"`
			Conf  float64 `json:"conf"`
		} `json:"raw"`
		Avg map[string]struct {
			Average float64 `json:"average"`
			Sum     float64 `json:"sum"`
			Count   int     `json:"count"`
		} `json:"avg"`
	} `json:"results"`
}
