package torznab

type Indexer interface {
	Search(query Query) (*ResultFeed, error)
	Capabilities() Capabilities
}
