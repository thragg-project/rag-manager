package core

const (
	StorageAggregateType = "storage"
	StorageEventCreated  = "created"
)

type Storage struct {
	Id             *int64 `json:"-"`
	Collection     string `json:"collection"`
	EmbeddingModel string `json:"embeddingModel"`
	Dimensions     string `json:"dimensions"`
	Metric         string `json:"metric"`
}
