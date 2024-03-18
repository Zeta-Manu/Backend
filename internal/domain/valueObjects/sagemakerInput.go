package valueobjects

type SageMakerInput struct {
	Instance []Instance `json:"instances"`
}

type Instance struct {
	Data map[string]string `json:"data"`
}
