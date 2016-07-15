package torznab

type Indexer interface {
	Search(query Query) ([]ResultItem, error)
	Capabilities() Capabilities
}
