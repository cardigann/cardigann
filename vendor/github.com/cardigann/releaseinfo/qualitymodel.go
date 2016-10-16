package releaseinfo

const (
	RevisionDefault int = iota
	RevisionProper
)

type QualityModel struct {
	Quality       Quality
	Revision      int
	QualitySource string
}

func (qm QualityModel) String() string {
	var revision string

	switch qm.Revision {
	case RevisionProper:
		revision = " proper"
	}

	return qm.Quality.String() + revision
}
