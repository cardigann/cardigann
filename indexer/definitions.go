package indexer

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDir struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	data []byte
	once sync.Once
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDir) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDir{fs: _escLocal, name: name}
	}
	return _escDir{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/definitions/alpharatio.yml": {
		local:   "definitions/alpharatio.yml",
		size:    1941,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/6xUXU/DNhR976+wyjRRjSTjGyIxicG0SVO1DiZWiVaVG98mFqmd2Q6FVfnvu3FiJ2VF
7GE8FN/jc8/9jIMgGBCiuYGY0LzIqKKGS4QEXSN0W0MPLZRTkZY0RRhEUOoa4eJFx3ggJCCZMYWOo6iT
CZMkGuBtQouWlVADqVQcWpuQ45iQP56ix/vWPmnsX5x96u6J+ztA+/7p/oEXLeVsL2VCkxeveu5U/0Xx
gS6QMpavmFqXzGWHed5Vn9dKWXsn4HXf8yPLax1/H5PpdEpILy3LQtBRsEG3JeMyeuIMZEPBo7vGfk3u
op9xWtpB2LI7KbTMIZou5ZuDzzp48njq0PMO/ZNzh17Uqs64tCHGNHHAtQUmmRQQ/GYyUJjTbVHosVzy
HNwgm9qchWX8KOWLjn5a4j+Hnrji7G//5rS9cTZmb0NFY66TgUXXknV7pIGqJIvJ81/zFjGvQQce1QQt
xRGBYl672wVtnAtqkBNxweAtLLJmrTTkkBipYjI80IYavbAeJKyNYa2Qy5SLHQWLeIWVVOvY/lqTi6I0
Pt1Sg2o+seF2S0IcwYqnoUNJVQ1bYkG13kjFdokO9URQCnP1PnvSwX6B1vXnO3DL1tUYbqgSXKSxkOYw
zDhjIEaWZ0CbD7pc4FZZ3YHve78NqKhAGO1D75ZOvlF0Y8tR+KAA1uQfhapa8dyAWuA78Ty/2W7Dqvp2
uwXBqqoJpI26qfvwewnqPfwV3usuaN8GJTe9jXDVGbrM4aBNa2EtYlTYAs2wOOTMu7bP1Pu+VhkWJvU6
JDIn1N9TYxRflvUzmilYebypp6ud1C9lM3gFKbwVvQtUUanGznRNmM2eZzP23Ww2vzmsDyO3FoabfO8k
Mb0lTxdcrHBXUyXLojn/QGgsTBYkGc/Z4cmodWVgKM/1/6D0WQ+Y3IhcUvZVCF1QsaN8/JWy5n9/0oKe
yoWvFGf6Jft8RELD17Avsm253y1gO2P9RO/KRc8Bkuy/eFyPBoN/AgAA///JzNEdlQcAAA==
`,
	},

	"/definitions/avistaz.yml": {
		local:   "definitions/avistaz.yml",
		size:    1240,
		modtime: 1475379448,
		compressed: `
H4sIAAAJbogA/6RSTYsUMRC9z68Igwfn0N34gUqDwuJRPIjDHhSRTFLTHcwkbaqyzeyy/90k/bFp6WEE
+9Ck3qt6qXqVoig2jKEiqBm/U0j8PsSGn0J8k+JvIdbcNJ43AQNTeIyIMr+wDgfGCtYSdVhX1ShQkq02
gRK8G1MEJ2isUzDG8dvub4uvre3Z3joHhrY1298+sZ/tnYKMSzFmvEclMv7GS2UTfbLy6R4E7kRbs++/
f8SWtG2UGbiOU8Ar7qmtEpzQo3WnOv1TqEznaRaDE1f6p0dwg0HbhwdWfrTmqJpyQtnj43ZM7zhib51c
Jk7onAjOWVfPNWtdhaEAMS5gNgBBg6BQyEquwRH7wLraUFuIVmn5/OUuZRIg/S0thPWGohujObkdNBiK
K9OzZ473YZah6n0c6YsHdy4/wTkOhPNEzvbZAqY+iR80lOmfNfpqFzqng5XnRfcRdMNGFGg5y40v6bxm
BMlM4kWUUAtgLuFETh18fPOkSMNIpPO6MCtHY4qj0rCM4tbHIgkUXgj+l8aivdbBcdK2vdGWy2viU16h
hDVXNFHdXxg4M+7NZJwM3l/Nfh19x46vXp27jQAS3AWvMsG30/UaQLT/UvFuxzZ/AgAA///2fUby2AQA
AA==
`,
	},

	"/definitions/beyondhd.yml": {
		local:   "definitions/beyondhd.yml",
		size:    2359,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/5SV327yNhjGz7kKi50MjRAo4V/O1tKJ6fuQKsrYJ32qJpO8EKuOndpOIe16vqvYxe1K
5vxxAoMGdpY8+T2PX/t1bMuyGghJosBFa0g48wNfCwyHWrjNhNlUCxSzbYy3WgRmxTJVCHuWrn5AyEKB
UpF0bdtEdPbJm/3ZtxDsRjooKEXYtszIx/Q4fyaQSQipJNKSgr0qBIrXQF10l0Na9HBUBHhYwZYLAsU7
Qv2Ri+b8VSv2LY0XONHaD7mCbKQlS+CkYEe9ku1Pc+mA7U8LbNw/h/0i8Bwel3xRgaNq7NlpXq877kY2
KdiJU8c6X0zkuA4b3XSjAhzU5s0fnILr1da4gDDem8TuOfJXpkAwTDVcrUA2NVPx4GrbQfnO5GrX4n7+
27fCNuxdsM2m344M4+EFw+OdszA16QVYriouZ5cr9Bjw3clmcpx6+p553Adp6HE9PZsuV6bkST36gL3n
MnZYzx721xnUs7/f31rTr+a30sE/xz7h9lcuJQUps20TS+Jp1GimKxMDzx/6ZbKBtWZ22MBwK+ID/w93
/L8MhnXs4U66qQNNzzI4TBtiDg4JWHiBi76/PJmD6NWqxHYKSM7aCKKn1E75lrDcHIIKuH90jBEWxarM
zj+4qPn+jjp3nG3ItpNr6OOjmUEKpDJ4hJUe0w71iurTV3aiIEpHFFgR7hb0Xum4f/76u5mfqlmZjcq8
FnwnIXOelmNmlZXzBZIdF74sK0mdFUnBU1zo8xivKXT0owCmipf1Hz4oTChSwmUqsLyAUP9H9tNNK7Nv
CFC/WoP8qE7M+1G4f+DvtRAuGayUIOs4vakCAZtS3xCqf1lZhVWXyUsMIpHaxrYHX3WU0NdOWoZpL1EU
LlbjVNXkk5X/x/LZBHy+Y5Rj/2JW/3KWJG+XpzFulYSAkL+mN78wxejOXDMpGWHmUixVLp0EknaVWdsh
GVGizvTmO2r+2WyjHno6Y9JTD894mp1muVfBPxrws/3VNWtBAbzgKkuv1Wg0/g0AAP//TebINDcJAAA=
`,
	},

	"/definitions/bithdtv.yml": {
		local:   "definitions/bithdtv.yml",
		size:    1559,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/5xUT0/bThC951OMwiFYP5w//Aq0lnqAcqCqeqEIISEOG+/EXmHvmt1xQor47h2v107S
glL14mTevPdmdnbsOI4HAE4RJjBXlEtacqxFyfHF15v46vLmlgGJLrWqImV0AlemRDALuFJZDpe4UFo1
CfDMQuisFhnLUce1axClH13CfwBiyIkql0wmq9VqzPXipuA4NeVkwIRUVIGYCsLMWIUhBpglwAUm51qV
GKBjhr6bJZMmF0V9LdYB/9BSL01al6hJ2C5x0iZ+VMZSgE4ZOq+lMiE+6z0D8LEjTG6VxI42mybtcX3A
vd3d3XXRcZOCA3AoHE+lEumjG/hkaeTmQJy2aZ7A/dNDQGgZb8CjoD8CrB4aeWEypVtxJYg5E4+Mq7zy
4MLYMvFPHypd1dQXqx3a9lKHLy8w/mL0QmXjDoXX12EgVsK5lbFyl9ihPRGtNTbpNb4dEo+42xKfGJ1r
liGEzakLTIm1QGJe4FgiCVXAmPCZepLF0iy5V0frAo9g7hOEjn6rWK59KQat4A3cGY7SEp/7TjZlhweO
BNVuLiwPS1OS8kMo7Q5H143JKIL/fGLY+Ib72DJmE8tL5Xrr3UF3F+in9w3XzdTc1nx5sROY+sCa1ZZs
Zyz3KyUp/zw6O5mOHoDmRq6BbKINHTY/lMdprgp5OIui9u4VFrI3Cy/P+s2pyy35cQSi5wgiq+Z18yHI
LS56fKEKQus2Zs1r3O7SU4127Vims60sW9nMJU0b3WIrKt7ege1u/n+nG68OiXZd3D97bZ1MmpUujJB7
vWYRVHvdnPq5/4Rn0R8rPrddO3xpew1Oo35fUO5cyjuCT52gQEzzv1HMptFgMPgVAAD//8qu79gXBgAA
`,
	},

	"/definitions/cinemaz.yml": {
		local:   "definitions/cinemaz.yml",
		size:    1240,
		modtime: 1475379448,
		compressed: `
H4sIAAAJbogA/6RSTYsUMRC9z68Igwfn0N34gUqDguxRPIjDHhSRTFLTHcwkbapiM7vsfzdJf2xaehhh
+9CkXlW91HuVoig2jKEiqJlQBk78LsSGn0J8k+JvIdbcNJ43AQNTeIyIMr+wDgfGCtYSdVhX1UhQkq02
ISV4N5YITtBYp2CM47fd3xZfW9uzvXUODG1rtr99zH62fxRkuRRjlveoRJb/6KWyKX2y8vEeBO5EW7Pv
v3/EkbRtlBlyHaeAV9xTWyU4oUfrTnX6p1CZztNMFsQp/dMjuMGg7f09K2+sOaqmnFD28LAdyzuO2Fsn
l4UTOheCc9bVc8/aVEEUIMYFzAYgaBAUGlnJNThiH1hXG2oL0Sotn7/cpUoCpH+phbDeUHRjNCe3gwZD
cUU9e+Z4H7QMXe+jpC8e3Ln8BOcoCGdFzvbZAqY5iR80lOmfDfpqFyang5XnxfQRdMNGFGg5040v6bxm
BMmM4kWkUAtgbuFETh18fPOkSMOYSOd1YlaOxhRHpWEZxa2PTRIovBB8EsdivNbBceK2vdGWy2vkU12h
hDVXOFHdXRCcGfdmMk4G769Wv46+Y8dXr87dRgAJ7oJXGeHb6XoNINr/6Xi3Y5u/AQAA//90eDmR2AQA
AA==
`,
	},

	"/definitions/demonoid.yml": {
		local:   "definitions/demonoid.yml",
		size:    1969,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/6xUTW/jNhC9+1cQRgHHSGTF2e62KyAN0hTooShatMVeAjegyLFEhCIVchRvutj/3iFN
fXjrNDnsxfJ7Gg4f5z0qy7IZY14hFExCY41VkgjDGyJ+GgkJXjjVorLmgNfcVB2vqBhM1vnAKHPvC/rD
WMZqxLbI891ut+q7r4TIZ/Ra8DaVCY5QWacgYcbWBWO/2kdiEnFBxHUnlU34DeG/PiTwLYHfb/KfSXNf
/zZSCbwj8BvW4BL+/gv8PjbLr41qoBdwXrAfrb3v+63XCec3tlFioElHlJV/UBJ6cevvejb+lrRuFl81
Vo5n9MCdqAt2+7BJDD5mI3kWCrw1ZwzaTViubaXMfnHLkWpyLoTtDN7V3EgNbtXW7X4bwNrKgrXWYySU
aTscNjZK3O/9nX/6xFY31mxVteo8uMCyz5/nqbDl3u+sk4eFPTspdICdM0lV4v7sykZhkZ6J3CmsHwq2
jhCcs26chgaBhNnWuuZWyctFWd4JGthdIBYbdsqwtPKJ3hu8FZpUXC4cyMWmMFhnolZanqyXqR0lwYdU
JvjV+4czN/aRZujxScMZK4NFjtP9OLAoTJV2Meisvmu5AT24NCricfhMNVV8wRGdKrtwI1Gh3kdyqzSC
GzzM0g11UMHHdhDFXeXJrT+iDnayOl1+Mw/CUqymyqhhui2H6XjowD0lx3+Bp+C0n1id7ipVnPfOdeVI
Xms99OFa4QHlASRQmi4Sho8YMqdHplOyb+zszv83HCuBvNQQRwoUfWNpqpL9wCIdntFDerrCWDy5JRGV
2SwjiEaSsZjczMzpm2XvKN+SmD6a9L2jI9XA5WTmUx3zfePLhYYtLjbzoeQLn/7Hq4lff19LGgw7uSqs
Ycsr8u0qGLfvB1oO7YZBH5GE8iCoNAM+5mKMVO1g+xq1MQeelpnqiOReSP/lCjl9UdTFVJSwTUMOHp3u
86ueO4q0O6Mtly92e/uabl798/Jp3i2nsT4ek4MF78PW4euSCjWAqF+zcH0+rJzN/g0AAP//KDOd5LEH
AAA=
`,
	},

	"/definitions/eztv.yml": {
		local:   "definitions/eztv.yml",
		size:    751,
		modtime: 1475316103,
		compressed: `
H4sIAAAJbogA/5SST2vjMBDF7/kUQy7ZQGyT/XMx7HFPe1pYemgIYWKNbVFZcqVx0sT4u3ec2EmgBdOb
nt7Tj/mjKIpmAEEzpUBnPoiwWIn48/z/SYRBWzRY9K6NmtDfaPsSUjkARFAy1yFNkv5pjEUyk/sM68HP
kKlwXtOgAdYpCPYiKqfu93yIAqHPyhQ2ryuQc3B2BVRv+/BgXbI1soTmbatziP815E/xXzodnVeh667B
pG0/Wm1LJlDXiUlWdd38QvPueKshkKGMnRc6495QnDvfVLuSUJHf7QVDHthv+vn8XpTuQH6x/Sx0Reea
jLrBh1GcRi0t0xunsB4HoNnQ3bzXwiq1XEZZqY369n0JOGQUMWoTvvIEAJm93jf9sktP+YhyR2scqknW
D2HFFRaWeAIZ9Hm6m5/LsQAZzmT61/KWyLVh8g+991/x+m2xrmW9D4ZU6IsgOwUs3Pw9AAD//6HMSPvv
AgAA
`,
	},

	"/definitions/filelist.yml": {
		local:   "definitions/filelist.yml",
		size:    1766,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/6yU8WrbMBDG/89TiG7QmtVxnXTrahija+kK22C0oxRKB4p9sUUUyZXkdFnpu+xZ9mQ7
K5Jid3RbYP+Y+Hffff7uLCeO4wEhmhnIyJRx4EwbBILOEZwi+LgCnIqyoSVCEHGjW8LETGf4g5CYVMbU
WZJ4h6GSyQBLOa2dJKcGSqkYuHtCRvsZXr9cJkeCzcHBNG3hUVMw6WUvW/JJLrA1GZ94uteh73hzTpe+
0hZ85VQqYKVwpXG3dHLpvfa7+MzT9DB7yiftNlyEBhvpWAotuR/HepDPx8l73Kf2woPMQh941Lt97Vp8
NXVbOjshz4gGqqUgNc1n3m00doKLJwRpENRSGQ/d8t2tTUSurq4GFsxlsX5PaKnyKiPXtzeOmEW8hrvu
mbsE6pu2ncuSiVVzTQ1qEkuGdVVbOJVqntnrNc0Nk+LNtqEzCKLt1WOYqBsTQjQa1OpMbt3fkyFuecrK
oafk4WHLCWuq9Z1URV/oaUfYCC7zGcrSrTa1ohill3q+DJE1cMiNVBkZakNNoydUkYItMmGqOK8YL3ZG
0SOwH62mZdyACnPE7stSUMK32kFCqCo1Rvl6blP8/LEzfPE2em6DuU13kk2UvNMQwvUX5V+Mnf4DLNup
dWds/A4zsmdvWpt1W5gQrwqEwaobAHgRdO4zXmYh+u+dhk44dDaRRkTXVBC6HtcYxSZN+59TKZgG/mhZ
3YXdNqCWGttE2an6zWEsfzaZwYdvkG70p3TWzRUKMJRx/d+8O5MX8g5PIy028R7/o7dm3zfax0GE36Yw
oUPBXC6gPXU+LB6BTQxf9Qw1QNF7xX83OIycmAPk1YbN6V40GAx+BQAA//+adzOF5gYAAA==
`,
	},

	"/definitions/freshon.yml": {
		local:   "definitions/freshon.yml",
		size:    1347,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/6xTTW/UMBC991eM0ksrNVmtuPkABwQXDghR9YJQ5Y1n1xaObTyTllD1v2N787VooQc4
7Cbz5s3z+M2krusLADKMAvYRSXuXYie7FL/P8Ud3e5cQhdRGE9h4J6C6vYPPGA0SSDd4h2+qRLHSHXp5
SIXo6p4yYtw3EukFoAbNHEhsNuMpDT9sLlKqlWGktJLx4LPqMQbYCkiHl6DzasH5oSaUsdUCvny/gfRO
3t0Ahq+ZbP3BuCM1SE6cTUGaoEMB9z52ovyX0LjQ8yzdE8bj7aunJ2jeerc3h2ZC4fm5GolBEj36qE6J
EzoTMUYfJ/E6tWqx5YSA3orWO5bG0VX1rrCq65GXrotE2coZgFVpdVkuBKFapSN2/iG1TTxYvIFdSTES
i7nh4kWIfm8sFjdSJso00hOztO9w9mp9KLFkgiAc67rVxqqrV9dAQboVsr2usuo4nbXsLvpHWoRPXZ+m
Waz81GMcmg84ZCNpdjLXL/ypLZY7i5faHLRNP4bXwDuvhvyMwnm+alpvNUp19HZv0KpZZty4YbGZ8QcL
2E57ZtiuZrCcKpv0iOj4Pm/Ffd7zmSWZo9n1+YMq9WNCYZq1pX9Q0xH3U+O+6xL/P6kp/+isl+qcGqum
WFx00ud+MusXdMn8POverJkJF78v8C5ObaXp/LVcKoXqz/WEqDCe9Wjp4MgBWRZ5ZFrEVr9UOZFA/goA
AP//ODev5kMFAAA=
`,
	},

	"/definitions/hdarea.yml": {
		local:   "definitions/hdarea.yml",
		size:    1915,
		modtime: 1475379482,
		compressed: `
H4sIAAAJbogA/7RUXU/jOBR976+4ovtAtUrS8rkbiV3BhFGlATEDVXmYjpAb3yYWbpyxnSkdxH8fJ06c
poBgHual6j3nnuP74djzvB6AYhpDSCmRSEyYkaUJx9GpDTnJkoIkBsLMK1SJsOxeheYPgAep1nkYBKvV
yrcOfiyCniFjktdJMdGYCMmwjgEOhnshwGQK0Dc/N1hyjtrfpFKxcsxoaJhL8cMkB+PIJHwsOB9HHXoy
fZEabSm3qEa1Ae8/U0ymHbLVtMSwe07/jBfXZN1hrWyLGR1Y5iaq4n40jTpUbWng5/xhy1tT6Nehmd81
LosHl3rUpu63BnYVwe35mRddWLB/i3OP8k5KLTRpJms74dh6XE3G59fQZ58J7VC19kqnKLv0sO48EnGx
xEwT2Y7r0FKnGVuiA8seTgvKRDBlFIVptlAshipo71Bdz00upHbgvwasSggumYod/I9zvBBKcVSqV3FL
Qdsbq5DIOA3h6/dvJctFwjLL5UQbPNDkHivUz9Pc6lGngoaQC2VrYFleaOdYKJT2U9t5fAT/g8gWLPEb
FJ6edurEnCi1EpJ2ExvUJaKUQoZO80pRZVlKlZ9zHZatcYy10YKmvsYHW6xGpbfcaqWqzHpuJp0pCCnN
FpU7r9sz/CXJyvRhlSdlO18KlGv/E67LZpTrRprvvp29K5DMOfrNIaDngq5BSz8TmdIsvr8ThUZ5N08q
5YIhp86lfoXWLzceZjr14pRxujsawH9AOoCTEK0lmxfVeylx4fAF4+bgtk8on0a7XYkJPuQbhHGRiTJT
mM3+N0Wd7M5m9O9Bs23NNH9lORsl7Q26s6iuTD0OCgTmtQFFTRhXf6BnKlYZF4S+aX3w29aK/Xx7AoeN
CzV7fV8RKifZm3VU83cXD2lnq6+4HzVOHDFO36M4HvwKAAD//+xOZ2d7BwAA
`,
	},

	"/definitions/hdtorrents.yml": {
		local:   "definitions/hdtorrents.yml",
		size:    1844,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/5RUWU8bPRR9z6+w+D4VIpglK3QkKtEihFRVVSlCSIgHZ3yZsXDswfYkTaP893oce5YU
FJqHKOf43O3YN0EQ9BBSVEOCcqKFlMC1MhTHc0NdXwa3Dccwz0qcGR54UFqG8meVmB8IBSjXulBJFOUk
8IlCIbOoZ85TXDhdijVkQlJwGKFBgtA3sTBM9JmVN3jl+GHNO2LUCK8vHTduuJ+em7yimxrhX9yo5r7r
HKQP/5ig27uWLLa4LRnFXcnozOK6g+mkez421S9KQoU/H+zgYRdPTrt4PNnRT3f0pv79/X10eVdX3DIe
nbbR1CEL54I0V6EAyzRP0MPLo2P0ImjIk0qgBD9BUDxW4UxklG+DC6yNJrJMWOSFJZ+EnCf220LKi1LX
xUpKEnSwXqPwi+BPNAtLBbJ6d2izOXCaYrmjKbBSSyFJrQEphWz6Z5Cax5cgjWcMwnTGRPqsTAdcP6SC
CXl++N/VVWw+h9sJNSjto90IS6pyRpW2U5gjiTUVnSF73Vohx4tSIU0SrvMgzSkjR6P+1gDKNMh65MDt
lYQMfhWORAjLTJkxb2whdBQe968ufhxUtZ33reL1anmPu6b+L/HSOibNsoKxrV63zcat3urh8Xy9Djeb
D+s1cNKY7S/a+v0VVpXPqnUZONV0YbofWCzFUr3l+xxTbq1PjfGmW62Rli13+LHzx6SBuaiSmvMcqyNN
wsp7rLx/wEhdxg+Q1M61CrfdH/TRJ4Q7RGO21pLOSvuPJ+Gp5neuqn1dLyXIlTJhPGud+nvzXfmNoZrB
3g5HfTRD2KkIaEyZ+regtyYhYsmZwGRvtkl/by5Ff+8f5cx7S4wRe9Wn/fc4XqUqsFTwit+DSRKPk3iC
UDyM4kE0jONp/Q6BdDK+9T5i46QTMYA0f1fQoArq9f4EAAD//xDXkPY0BwAA
`,
	},

	"/definitions/iptorrents.yml": {
		local:   "definitions/iptorrents.yml",
		size:    2463,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/5RVbW/bNhD+3l9BKMUQo5FfVL/IArah8zZs2ApkWDAEaLKCFs8WEYpUSdqp5/m/7ySL
erO7uF9k33P33B3vkY6+778ixHALEeGZVVqDtAYhSVOEfr29qyEGJtYYxJWMiHc1Ireab6kFcqdp/ATa
wxhB5XpD10gF6W9yluDyyUT4hxCfJNZmJhoM6lL9WKWDV+iOaVaGxZhzrTSH0iZkFkTkvdoi4oCZAwZ/
/lhi4fwUmw/PYNNTDKETbBxWcb84bDLuNDINTpO9rYk/Kw18LV1w2GEHw9MSszOtzN5G5O4vR5rmRu2c
TNr2LGzbQUGuCwTjjr/Dn3byh0HHP2k0EzY7m807mYN25a4fG2mf5GhX8fN5h48B7zaMq8ZYmnY4LO3B
78oYAcbNed4JnA6LxO8kT6GERlj7duH8WPh2MXhPYweMXebiuVTqyYmNPf2AZlVqXNqDhUp57OAQ89/f
3zsLizWsScsXNq15mzc+WoWZKlZ/IQaojpOIfPj0WCJ269fgTR5glLwhkD3mdKHWXB7JGbUYMyiQfpZk
BbhSOo2KZ2FymW1sVWxjQB/3g7ffk/5CyRVf9x1KDgevDMyoMc9Ks3agQ6tA0Fppl9zHVgXEuB8islRs
F8VKWsqlufYWNLNxQskWNF9xXBO4iciKcgHM65V0nAvqnm+gCiCdjA2HhlRt8RxLfaZ6MmrU/qno8dIy
3lUxTpJ5Z4oZuxNwQ5aFy4KxUTWuQol0V8iAoM6P2FKJSwafcXtWQtU1+/HHIj4nlso3mN7AemekfK3p
c6GOxsUNKFG1ew+H/b5/OHyz34Nkh8Onb3MF/9iA3vV/g12un6kE1Oq58SZWQ7B0KeDKbXryHbH5+PNf
HUmb+HHCBbuWb4JelFBzbVnffhR0CaJ3TLviIFiVuLwWdvXE61KWNRKOeliCVlHUWs2Xm/yGSzSsKnzF
hQVtmgL65bWnYQ2fs4YDs+i1wSP9/fDw/fXDA3vTe+2ktdwKeLGpoNkUA3yrhPk60pdOghdoms/3xWyT
S7Ix9SyFouzFbONLshn+z8ujmbqPiqHCl8wEX5PY1qv7f6U0meD2jJIfiPevd0P8EXn8OhZZ7ggShxXP
ALBW8S80HrpjCoA4uYQx7/0XAAD//ydNe46fCQAA
`,
	},

	"/definitions/morethantv.yml": {
		local:   "definitions/morethantv.yml",
		size:    1281,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/6RTwY7TMBC99ytG5UKljUuLQJArR4S4rLigVeUm08Qisc14krCs9s5X8HF8CbbjpC10
VSR6qOw3773xPDtZli0AnGLMoTWEXEvNvYe0bD30wUO3HhK3nzzWSF11svI46qxzAVH6i8v9AiCDmtm6
fL0ehkFMXoL79cLXC2kTr5CMlSGFaQ+wCX16D6T9NofYLvxe5vCRa6RF3LemPMocSirqHD5/vUsI99kR
vAkEZ/QNoL0L8sZUSo9iK9lz1hERtrYRPBhq8/gvZMf1LqxiQWnb8dy2c0hjOMuHBxDvjD6oSkwoPD4u
E9FK5wZD5TlxQmciEhnKZ82Fg/m50bmQe9qG2Rss2OtgeTwsiEGSVrqCAjUjjf6Mjv+wZ2M3L6K9x0my
MnmifmPv+OvHz2WopCwXZ0oi7+3ms51nMyocU5r5Pd6HWd08LJnhhDzNwHLf4LPkvYs7YBIJGC9HYVPO
0vSG7i8FwqXwZbcrTAOl6mdGId1JgOCDi7Q2vrxlDpu/a9x7fDs9LsXNxRvwDfeq2il9MCAqMp0d1zLX
XGdFrZry+XaVhCWyVI37bx8AyUxq34UPtyY8TP5m0I2R5bUGzkp95ry55uzU9yfGP3F5Pc/pb+gq+9UK
BKsWL3WOcc8vBUukJ0I78XszdW8Qi/pfFG9XvwMAAP//W1yqZwEFAAA=
`,
	},

	"/definitions/ncore.yml": {
		local:   "definitions/ncore.yml",
		size:    2101,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/7RU3crjNhC9z1OItGRvEn80uxRqcKHsUgrt9oeW3oQQFHtii8iSV5KdH9fv3pEtK0rW
+XrTvTHSzDlHM8cjrVarGSGaGYiJSKUC3Ala4u7X98OOU5HXNMcIiFWtbYSJo45xQciKFMZUOn556clR
ms4wntLK5VNqIJeKgdsTcm5YtitqEfe7j7LB3MuPyGW5CCBDOoD8+cFls8by4/vsvQBCbvwpge/+owIL
eV5BETTwRKAICxghP30IGtSgnMpff0/0j2mngOmw9Bvvc+KQHk8OiUXI+5xYhDyb9qWW1dug2R/qjMkH
LiKCXgfEx9/fuiyXWnPQehCZ4o+IOOD/4mIOknJWxQ9HeEiPKWV2GzENVKVFTDafti5imtUtuLQALcWS
QLW1dC5zJgZyRQ1iXvpIVBVVHzxIVcb9t98yUdXGHyagicm8bUn0XooDy6MarbQ3iHTd3GEqatsLQTZy
kirzIFBKqtjjJ6rAJrFdexFnoxMaOKQGeWT+VcH2tL6CAENMtmkoR4eTNyXLMg5vtsMhBrR5OKNS8sA4
9KfMvHWhF6ivQBjtC7nv/2tFTzH5x9fUtgofDMBO/dXvuiPDgqi+GmnMzrCq1ptt0rZR1y3aFkTWdSVT
kFiD/qhBXaKf4WLt0ejPomR7EIm1dNFTEcYOz/XpcXeUV1ygMtfQdZTznTwJd9BC1/uSmeicvHs3ri/J
N9+6dfLbcWForpPBMSVPwViNbkd7ed45X3YoT76/Cw1Dw4BnnuvewcvUv+upWH61Y2WOUtRjqDGK7Wv7
NhcKDj6Of8yA0jcx+xAPr/Yna59GmsiDLEqpHEew92+8E8zwyVmKxs7M2dhylneR9bMKe73xIQJDGdf/
o3rQf4Y/k0uafRn1V91VUHGawoSzmzlNDZMica3Pl8RHXL3z8TXS7DptvB2EEhSY9dgqjs1T5AG4kdw0
sPYIBaVssE62JHvlxxayu3YeZPTae8IB0uI1LLfY2ezfAAAA//85UY/FNQgAAA==
`,
	},

	"/definitions/privatehd.yml": {
		local:   "definitions/privatehd.yml",
		size:    1250,
		modtime: 1475379448,
		compressed: `
H4sIAAAJbogA/6RSTYsUMRC9768Igwfn0D34gUqDguhBEEFx2IuIZJKa7mAmaasq24zL/neT/tq09DKC
fWhSr169VL1KURRXQpBhqESL5kYyNDoiTp4i8nlAPryPiJWuDrKOKLgiUEKM+0lVPAhRiIa5pWq367qu
nIVK9rurSFCyHYkq4rVHA2Ocvs3+uvja+E7sPSI43lRif32f/eRvDGS5PqYsH8ioLP82aOP79Mnr+3sI
JKqmEt9+fU8tWV8bN+RayRHfycDNrod79OjxVPX/PjSuDTyLwUka+yMQ4GDU5vZWlO+8O5q6nFBxd7cZ
6a0k6jzqJXFCZyIgeqzmmrWu4lBAlNYwG0BgQXEsFKW0gCzeiLZy3BSqMVY/frrtmQzEf0sr5YPj5MZo
Tm4HD4bSyvTiEcouzjJUvU4jfQmA5/IjnNNANE+EvssWMPXJ8mCh7P9Zo8+2sXM+eH1edJ9AHDZiwOpZ
bnxJ5zUjWGcST5KEWQBziWRGcwjp9bNhC2OiP68Li3I0pjgaC8sobX0s0sDxhdB/aSzaaxCOk7bvnPVS
XxKfeIVR3l3QJPP7gYEz415Mxuno/UX28+Q7tXL16txtAtCAD3iVCb6crrcAqvmXildbcfUnAAD//+aC
bfviBAAA
`,
	},

	"/definitions/speedcd.yml": {
		local:   "definitions/speedcd.yml",
		size:    1762,
		modtime: 1474957311,
		compressed: `
H4sIAAAJbogA/4xU0U7bMBR971dY5YVKJGnTFkYeNgEFMSFGR1E3adqDG982FsEOtgPrEP++Gyd22grE
XsA+59zjc6/dBEHQIURzAwnRBQBLGe4FfXD70AIMdKp4YbgUCenezckMFAdNqFhLAV+6KMmpWJV0hXUg
glJXCBf3OsEFIQHJjCmSKHKeUQfxlBYNn1IDK1lZ1ntCxnFC7ubR5aTa7OHqNLila0cOLDnz5GQ+CW4b
Mt7gLHlecC2xg4YfHe84X06c7dgyNyYDZZkrznzVYIeb0vTekeN+fWYhlbGkXTl2WLMngj9AA6HZtXzC
dmtDd0bsYFe5C4yGvtCnHh3tiOJPXnSal+3URi0+9MV9j+EMnQHOYXoWfZ3duBw4sjMptMwh+sG5qx23
6HQ2dNoWdciw1f1cyD9keNh3Fof2pGmGryjYHEU8skSfbd05Itc0dRLs+1RKfwkxep2UjEvnPdrex8fN
PppzBrJj4YfqYbg3Z54CDVSlWUJ+PR4QXGspDggUvytxLldc1NKCGtREht7DuYXDIitqPzCZZAkppDYW
4KIojT+h1KDqH1f35YWEOJUlX4UOJa+v3UZYUK2fpWLbQod6ISgllTMPMHEOqUGEmIVk6ySVwlAu9H73
gudAvklDLmQpWLdnSwxok/gDbUtVFPxtVt+EehKb/S6UfNbge91uzU3O5v1eglqHV7Cu0moft6pv9S5r
dw//KRDmji4wZYhPBPs1CJDPxFjMtkOMqm2WHHLW3ho3ObjNpu8iESYL0oznbH/QawTNl2b9VoFhWxV4
OuNPuxDdAOKed6HGKL4oq+9opmDp8SXPDSjdHlfdU/0EHqshaSwTqw0WrdRKJ1XQBmXyWeSSsg8jD3fz
DT7Kx3Aab9uG+dLUA8C/uqAirLTN+i1Xew3ubvnfd2w3wo17/iUA2xrROwWHriAHSLP/qTjqdTr/AgAA
//+nmNn+4gYAAA==
`,
	},

	"/definitions/thepiratebay.yml": {
		local:   "definitions/thepiratebay.yml",
		size:    1411,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/7xTy27bMBC85ysW8iVBLSmPviygKFr4EKAN+kpzCdKCFtcWUYZUyZUdV/C/dynLtJsH
klN90szODkcjOk3TPQCvCAugCmvlBOFELJk04prJ8wrhc8fC+47WwswaMeMRmrTxgVHmly/4ASCFiqj2
RZ7vmmXWzfKH5j4I6sWdee3sDQ+czfd4VIq6P6LknZl1CnsMcHx4VMCZnTMVmeMNk48vxpE9KeBdI5XN
L5RECzCAs8arEuYBbpefx+VPHNMFWYBQalVvVS+4nIuIXgbU68NvAKfCyAq1jJJX0fZ0G+l1t7dDjKLq
JJKj0b+J9rrBtZXbEmieehSurAq4/D0EfvbWDAHrqyDuR522FsSipG3VFLIvDbpl9gGXC+ukX63Wwrxt
747yw3w0yg/bFrVHhg5LNMTQyNUq6aydXcRAHjWWZB0fNVi7fkXfaAKaWLkEckUl/D7JbI6OzquDtcNU
cWPRo//Wyw3edSVZGKrSslJa7h8dgCi08LTGUS6InJo04XZXDqeRnypN6PzWN1y99YX3tVa0w7OJm3lu
FZI8GUJ6BFebyhVpvC9aJpE+8p+iHzESSvsnKB8KLO3CaCvkoz0ccw+XXa4fb5Jxv5VcPWLv1Z973+O2
9dQaConH6MunVOlwhjf1PV0m3/hE2M+evT0YJpt35G/9f0N8r0M7KG8FGfDhKHfcBrfyZCQmGn/2KhB8
ZYSJWo1YVo9vb2Qg/gYAAP//HuvaLYMFAAA=
`,
	},

	"/definitions/torrentday.yml": {
		local:   "definitions/torrentday.yml",
		size:    2975,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/5SW3W7bNhTH7/sUhFMMMRpZseS4qYBtaONlBTYP3hJkBZpsoEXGIkKTqkg5dj2/+44o
UaIay0l1kYhHP/7PF2nS87xXCCmmaYS0zDIqNMEbMAm8BNN1aZoYE8dikeMFmKnwclVYmHhQEbwg5KFE
6zTyfU28PCVY00Esl/6rQp1qzcSiBkvpWMoHRo0JIb1JiwDoWlcGjueUR+iihMAY47QSiEF7ITNGqzFC
7yKEpnIFFtQ8R5WpYoKzqDL4V5NvGX90fppW4HBYgx+fgh94npliFM9Z7bcb9S5zzis+bPg9MUxuJt5f
NtrhoWins5HlgkPcH1J4v4gFZyqxyYWH+BmOH2y9RqND5NXEXwdjG8Xo/FDJADyz7g8W4NOKTSoQBNH1
DWo/R2CyeUN01zeOiP3udjIMDPTxKdTqYjjs0nI7EoYHKWadjsZd3FTOGbcLfmgS2JNgqwlBp1i7AW9N
ufZxDhR21szpUNCp5HQneGeU3gu2pA5kxrYMoDO78H+Fra4cxIxtCc6L/S2U5NSfXYUWgVdnETTArAFm
VuG0UfibMQvAq/PDYIFPc7kOx6cFULx68F5RUGT0PidMtnOe5orF1hNU2CCw98JvEN98sKWpwd+lUpwq
5YCXHFvBILTcpcwoWwhX8Om+HQ2jzhBbK2Y4tro3jFDZBo3JKgYWNH/n8ENremicfIBRrQjg7AK1HwDT
9KvN5byTaK9maJdRbnOus/DUrKyJjPMlnDs42xjEGdvww8KnP8VxuyBQYEMsJWnOB0VxFicR+vzlzp43
K68xnhSAkuIE0fSumM7lgoly8pLqRJLWacVEmutau/wQod52iwaw1u7ZYlDa0G7XM5CmSls8xRp8+rAL
FJykapAmaeExw5rJkqkIJghdm89lBpzGcD6DnyMca7ZiejNhK6RSLCKhEy9OGCfHQb9XHrgmNUeu5+ve
nuBfZ/jRhJ7B0U4h/vpo3e2228Fu98N2SwXZ7b78WKT3Z06zzeA3unmUGVF1fpl8dEptA9V4zulRdae4
LgboJ6TnkmyK/5kTtXgT9KMEq2NNBvpfc/L3jdw9o5w0lS6D29hxyxlx9IZ98IBrCmudsXle3HCSjN7X
9nvGNc1UI9fcTTK6oOv0lbtOcQYXGNT75/b25+PbW/Km/7pn1xLTnD4bVOAGRajGjKvvm9SVCVyzir3x
vNroJWpEPgouMXlWLXyJmmJfny/NWd/6hg6/pCawSmAP1GfNwVaqlDO9p5OfUe+/3gnyhuju+2ah+QbB
xNN6nqKUtJx3BD62aXJK4+QlM972/w8AAP//0xltjp8LAAA=
`,
	},

	"/definitions/torrentleech.yml": {
		local:   "definitions/torrentleech.yml",
		size:    2314,
		modtime: 1475375088,
		compressed: `
H4sIAAAJbogA/5xW/W/iNhj+nb/Coqfp0Jq4QHsfkTbp2mof2m7Xjeo2qVSTiV+C1cTO2U6Bofzvs42d
BI4ONH4o9fM+74efx3GIoqiHkGIaEqSFlMB1DpAuDMhJYcD7LfirB3PCs4pkJgA8qpRFGH9SifkHoQgt
tC5VgvFyuYy75WIhM9wzpJSUnpwSDZmQDPwaoXcJQh/Fs0Hw5Bb5zxm6IYUnvD9MuJ/4+PAiORSfpAAc
ZCANO6SAjQ5g4wb7qcEuD2BXAUOo7YmuxUqBVp40MqPff8a3Iq0KowmR6xB44wJN19Fbu0a7n6/KjUcu
q53BZN3d4B+NY4EyNGreCK5EDvivmVgF+P0uPH5zEToP28jdZBzQURe9C2in+J+MhakuWvS328YUs8MP
FWUCf2YUROAOPRrWl25HHzgrIEBXh6W4IVIL08bTDMuoI54aqcdOi58nnwJw6YCPJA3AlQPuFoJD9Ekv
mqMx3qZeULLuOagQtD2eCohMFwl6+PLoEf0cteC5JSjBzxGUjzY9Fxnj2+SSaMPBlQKJSZqKimvswtiF
50IWifvrloyXlW7a2qTto9jfbFBsBJ6zLA4oquu+J5ZEqaWQdJcY0IYIUgoZikdm6BxS86iarDM30g9m
jOtcpE8orliktHlMI5cT+hhZQCl7B/RaX04vg8p+J09CIZ7N3lRJjHJKS8EzF9agdNLszMnnLxSFZ1Is
FfT2O8cFFDOQMyL/Jrm2FkiimdixoLc/7DbnmvjtzVmuQapWn630EjJYlc3cRGbKJP/hysevp1P67XQa
269B3/b1p6LTuL8/Pd5s2BzFv1cg1/EvsLYeqbpmnMIKf7Eoti7uxzcb4DRYuXtQXkmydN5Lc0eDOQDN
BVvX/rJdPzx+t9nEdf3NTh07UPKVnJrMcjjzY7sF+h7pmaBr+y29XJDTJjd0aQ9GpxpNuF5E6YLl9PVw
YEqQVk6tJZtV9h20kDBv8D0z/sOQjinTKXaGDF6Fc6aZeQkdnWk0QLFjNoNR0ITl6n9kvrSlVBT2+j9e
8vIUgahY8lwQerTaeHC0lmL/HJfoahA6G6NPkOUUI1WZM33AxwfUR4Kj/jmKhuixOZxAd+q80PptaO1+
fJyS8W7wbwAAAP//P+lCKAoJAAA=
`,
	},

	"/": {
		isDir: true,
		local: "templates",
	},

	"/definitions": {
		isDir: true,
		local: "templates/definitions",
	},
}
