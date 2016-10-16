package releaseinfo

import "encoding/json"

type Quality struct {
	Id   int
	Name string
}

func QualityFromString(q string) Quality {
	for _, quality := range AllQualities {
		if quality.Name == q {
			return quality
		}
	}
	return QualityUnknown
}

func (q Quality) String() string {
	return q.Name
}

func (q Quality) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.Name)
}

func (q Quality) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}
	found := QualityFromString(name)
	q.Name = found.Name
	q.Id = found.Id
	return nil
}

var (
	QualityUnknown     = Quality{0, "Unknown"}
	QualitySDTV        = Quality{1, "SDTV"}
	QualityDVD         = Quality{2, "DVD"}
	QualityWEBDL1080p  = Quality{3, "WEBDL-1080p"}
	QualityHDTV720p    = Quality{4, "HDTV-720p"}
	QualityWEBDL720p   = Quality{5, "WEBDL-720p"}
	QualityBluray720p  = Quality{6, "Bluray-720p"}
	QualityBluray1080p = Quality{7, "Bluray-1080p"}
	QualityWEBDL480p   = Quality{8, "WEBDL-480p"}
	QualityHDTV1080p   = Quality{9, "HDTV-1080p"}
	QualityRAWHD       = Quality{10, "Raw-HD"}
	QualityHDTV2160p   = Quality{16, "HDTV-2160p"}
	QualityWEBDL2160p  = Quality{18, "WEBDL-2160p"}
	QualityBluray2160p = Quality{19, "Bluray-2160p"}
)

var AllQualities = []Quality{
	QualityUnknown,
	QualitySDTV,
	QualityDVD,
	QualityWEBDL1080p,
	QualityHDTV720p,
	QualityWEBDL720p,
	QualityBluray720p,
	QualityBluray1080p,
	QualityWEBDL480p,
	QualityHDTV1080p,
	QualityRAWHD,
	QualityHDTV2160p,
	QualityWEBDL2160p,
	QualityBluray2160p,
}
