package core

type (
	DataSourceProvider string
	DataSourceType     string
)

const (
	// Providers
	S3DataSource DataSourceProvider = "s3"

	// Types
	HTML DataSourceType = "html"

	// Events
	DatasourceAggregateType = "datasource"
	DatasourceEventCreated  = "created"
)

type DataSource struct {
	Id          *int64
	Provider    DataSourceProvider
	Credentials map[string]string
	DataType    DataSourceType
}
