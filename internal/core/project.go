package core

const (
	ProjectAggregateType = "project"
	ProjectEventCreated  = "created"
)

type Project struct {
	Id          *int64
	Name        string
	DataSources []DataSource
	Storage     Storage
}
