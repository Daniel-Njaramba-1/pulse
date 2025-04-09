package interfaces

type Store interface {
	FeedGetId() *int
	FeedCreateQuery() string
	FeedGetByIdQuery() string
	FeedGetAllQuery() string
	FeedUpdateDetailsQuery() string
	FeedDeleteQuery() string
	FeedDeactivateQuery() string
	FeedReactivateQuery() string
}