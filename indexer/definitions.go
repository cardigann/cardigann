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
		size:    1252,
		modtime: 1475754160,
		compressed: `
H4sIAAAJbogA/6RSTYvbMBC951eIHEoDa7tftMWwLUuPpYfSsIeWpSjSJB6qSK5mlJBd9r9XcmxHLllS
qA9G8+bN09OTiqKYCUHIUAu5Q2J5H2srt7G+6ervsTbSboLcRAxsESghaH9RHRdCFKJhbqmuql6gZFfN
YkvJtqcoybBxHqGv0zdf3hbfGrcXS+c9WJ7XYnl76n5xO4Ss19WU9QOhyvo3QaPr2lunT/sQSK+aWvz4
fdcjvCtO4FUikLNXAtq75Nm4DdrjcCs5cioZuKk6uEPXzm/r7t+VaNvA426wlWh+BgJ/THD+8CDKT86u
cVMOqHh8nPf0VhLtnddT4oCORPDe+XqcOecqnhqI0g2NCREYUBwHRSkNeBYfRFtbbgrVoNHPXy06JgPx
39JKuWA5pdEHlcUxr/gYOX1Ee/3i2ZFxnex/DeAP5Wc4JPM0mPdun13GYInlykDZ/TNPrxfRJK+cPkyM
JtAfw0cwepTrX9Xh3JlZZxIvkwROgHFEMntchfT+GdnA8ErS+rywKPsIijUamFbpgvshDRwfA/2XxsRe
42E9aLu9NU7qS+IDr0Dl7AVNwvsnDpwF93YITsfsL7LfpNyplWe3ztMmAA3+iawywXfD9gZANf8y8X4h
Zn8CAAD////Fua3kBAAA
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

	"/definitions/bitmetv.yml": {
		local:   "definitions/bitmetv.yml",
		size:    1342,
		modtime: 1475754160,
		compressed: `
H4sIAAAJbogA/6RUwW7bMAy99ysI79IAtdNsGIYZyA7rcdht2CXwQbYYW6gseRITL+v676Pk2E6AdM2w
Xlw+8T2ST4zSNL0B8Iowh1JRi7Tn2IiW48+KvuK37xxrYeqdqBlDk+58QJR59Dn/A5BCQ9Tly2Xf99lR
I7OuXt4EZSRSpp5SB+XK2keFEQKgQ8cQ4U86AlqUqHN4GJIYrER3FKgEYW2dwmMMsMqBW4xBa+WMexSu
anLY/CjGOvt0Bu9CgrfmDrArAl3bWpmBzAM0Vp51qUy3o0l7OMgheXqC7MGaraqzAYPn5yQmEXoa0ztB
XDOY5C+41KL3bK3PuqYLjThByg7UvxNLZ3uPkTYMrLEi63LwnTCZb4XWW2sIPkH4bCqrrVsnb+7jX1IM
txP9+Pdi536MtkY/vuCht076yYqQXWmJgi1NVgMW1Gb22DiJUuOmQq07ISWvzTp5nxQ8AJVWHsLX5Y3w
tyQzLYLDiyixVajlfDnDihzGGOJq5bAa10CRxvnwpLrMDTVp1Sgtb98uuJyYsgSRU+Uu/Eoi/3ggkYTS
/j/UGofbUcz2Rlt26TW1d9eoefXr9Sk/LKYMh63dh0fAje2wj9cNFtcr1Fvz9RYvK/JNbZUmdCeGzW+C
77SiE5wHc/xuwAaS38kdrKC4QOLJ2wucJEum5UJ5VvCFQT6OTmjEqrmGsbpf3PwJAAD//3ggy4s+BQAA
`,
	},

	"/definitions/cinemaz.yml": {
		local:   "definitions/cinemaz.yml",
		size:    1252,
		modtime: 1475754160,
		compressed: `
H4sIAAAJbogA/6RSTWsbMRC951cIH0oN0W6/aMtCWkqOpYdSk0NLKLI09g6Vpa1mFOOE/PdK6921XBxc
6B4WzZs3T09PklJeCEHI0AiNDjbqPtVObVJ93dffU22VW0e1Thg4GSkj6H5RkxZCSNEyd9TU9SBQsa8v
UkurbqBoxbD2AWGo8zdb3Mhvrd+KhQ8BHM8asbg5dL/4O4Si19dU9COhLvqfokHftzfeHPYhUEG3jfjx
+3ZA+E4ewMtMIO8uBXS32bP1a3T74U5x4tQqclv3cI+ufNg0/b8v0XWRp93S6dH+jARhn+Ds4UFU196t
cF2NqHh8nA30ThFtfTDHxBGdiBCCD800c8pVOjUQ5RuaEiKwoDkNikpZCCw+iK5x3ErdojXPX817JgPx
39Ja++g4pzEEVcQxq3kfOX1Ed/Xi2Z5xle1/jRB21WfYZfM0mg9+W1zGaInV0kLV/wtPr+fJJC+92R0Z
zWDYh49gzSQ3vKrdqTOzKSReZgk8AqYRxRxwGfP7Z2QL4yvJ69PCohoikCu0cFzlCx6GDHB6DPRfGkf2
2gCrUdtvnfXKnBMfeRK1d2c0Ce+fOHAR3NsxOJOyP8t+k3OnTp3cukybAAyEJ7IqBN+N21sA3f7LxPu5
uPgTAAD//7b7vjLkBAAA
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
		size:    801,
		modtime: 1475396593,
		compressed: `
H4sIAAAJbogA/5SSy86bMBCF93mKUTZ/kQIovS2Quuyqq0pVF42iaIIHsGpsag9JE8S7d0iApGor1B3H
5/hjLo7jeAUQNFMGdOWTCIu1iI/fvnwVYdCWLZaDa+M2DCfafg+ZfADEUDE3IUvT4WqCZbqS8xyb0c+R
qXRe06gBthkI9iZqpx7nfIoDoc+rDHY/NiDfwdkNULMfwqN1yzbIElp3nS4g+dySvySf6HJ2XoW+vwfT
rvvT6joygfpeTLKq79c3mnfnuYZAhnJ2XuiMR0NJ4XxbHypCRf5wFAx5YL8b5vPhpXIn8i/7v4Xu6EKT
UTN8HMVl0tIy/eQMttMANBt6mI9aWGWWqzivtFGvXkeAY0YRozbhf64AILPXx3ZYduWpmFDubI1Dtch6
I6ykxtISLyCDvi538zaaCpDhLKbfRXOi0IbJP/U+PMX7s8WmkfU+GVKhL4PsFLB063nVpH4D/OOf76PV
rwAAAP//KY7dniEDAAA=
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
		size:    1926,
		modtime: 1475754160,
		compressed: `
H4sIAAAJbogA/7RU3U77NhS/71NY9IZqJGn53CIxBAtTpYHYoCoXK0JufJpYuHZmO5QO8e5zYsdpCgh2
8b+pen5fOT4ncRAEPYQU1RCjnGAJ2JQcL005Ts5tyTDPSpwZCHhQqgqh/EnF5g9CAcq1LuIoWq1WoU0I
UxH1DJniwolSrCETkoKrEToc7scITaYI9c3PHVScpw42qVysPDMaGuZaPBtxNE6M4PeSMfOnw0+mH3Oj
Le846VCNbQM+eOeYTDtk62mJYfc5/QtW3uJ1h7W2LWZ0aJm7pK77yTTpUC7SwO/5o5a3oajvSjPBW1iW
L1563EoP2gC7jOj+8iJIrizYv4d5QFhH4oxGZlTbghObcTMZX96iPv0Tkw7lvDc6B9mlh+7kiUjLJXCN
ZTuuI0udc7oED1ZnOC8JFdGUEhDmsKWiKaqL9i1y/dwVQmoP/mLAuoXomqrUwz/7xCuhFAOlejW3FKR9
ZxVgmeYx+vufB4fo56AF9yqBEnwPQfFQ2ZnIKLfmAmujiTR+ghoNi7ywDwCdCxKjQijbJOVFqf0jSwXS
fo07r68o/E3wBc3CBkVvbztOWGClVkKSrrBBvRCkFDL2nk+aqtpSqvriXVmdnUGqjRdpEmp4sc1qUHor
zTlVHdbzQ9uYwk5kcqTZc605s4LTquu/SpDr8A9YVz2rpmdpLoB2Bb4NPGcQNklIzwVZIy1DLrjSNH16
FKUG+TjPaueCAiM+xV1H64+PF3OdB2lOGdkdDdCvCHcAb8FaSzov64tTwsLjC8rMg1WbXd2RdocSMngp
NgiTIjNlRjKbnZmmTndnM/LToNmpppp9soKNlvYH3VnUL4YbB0EYzV0AAY0pUz/gzESsOBOYfBl9+L+j
Ff336wkcNSnE7PV7TagC8y/7qOfvXzwgna1+kn7cJDGANP+O42SAev8FAAD//4ozdm2GBwAA
`,
	},

	"/definitions/hdme.yml": {
		local:   "definitions/hdme.yml",
		size:    2540,
		modtime: 1475754160,
		compressed: `
H4sIAAAJbogA/6RV/U/rNhT9vX+FFbaJavT7gxKJTQzexrTXN1YGqsTQkxNfEo/EzrMdSl/V/32OEzst
CwNp/EDj63Ovzz3HH51Op4WQpAp8FJMU9IDhVA8uL+Yf9CDBLMpxpAPAOrksIpQ9Sl9/INRBsVKZ9Hu9
IrULeUuHQ5xV0yFWEHFBoRojNBz76M/b3hmjKZQRdIDMyAImPro67/UJXqMakGVfq/nxsY/m/EmX7F1e
uPnb88sLW2DqAD8l+aIoc4D0l8DrCjEZN5QgweLD/GZpVxk0QC54mKfAFC4asrX6DcCfeS5+cYTGTctV
4hrCMx+d5YTy3kcuZQJSFvN/oHkuaWhrzBpq0HPMiG1qPGlA/PpJfbSrnPjodxWDQKjWdU6lW+HE6H4V
cwadEqjneUAT68yoqVUTQIP+rE8trEm6GpZZ2PBV2PGwRo2sNreUAK9QhTDIBKwNI7Ot0O7fAbrC4aPz
adSw3ALS/NkCpk2Am6Vta9Kk7zX9RK8VLCyLqdnc1xkXyrEwI8fz2CBshRKhiV9Dsaf2dZy9CXUqnfwn
dEfPiTs9O8gbnTgcTB1orI1eLpf7YupAy8ynnNTHWQIWYeyjuy/3VUQ9dergUQGQnB0hyO6L9IRHlJXJ
GVYa01P4EUy0m8UlgRRUzImPMi6VCTxwkfrmvxlSluXKMcgliPK+8jYb1D3n7IFGXRtF261XATMs5YoL
sg+0UQcEIbiwxTuafwKh0hGkSBfSAAgB0rKypPqwFldjq1ZqL0HBc9mCAql8R8Q0nq5NxzoosKJ8T5Rq
br+ezDDryhQnyQNnCv2Aip9SIJooELJmXQoiIILnzHHDIpK6+YVZ7bD7/Y/tmyzhWPfj/SsxS3AILzLv
vCPvCHmeMbLyeJd0IPhKgiO+b9M3Aq+M8kI/J6Dldy/DdhtuNt3t9nTw3WYDjNSO2X1k/PoN1oVPcsdR
ysKEAC4cHdhYkOAio1+OC0b1VnXG4CCBuxUlKj71Bv3+t969VlMFnKyLX+HHWB5q9wKuFE/vAn2NRILn
jJx6n2mq/Za9IPqsSwn9GHT/ziLvvl35AAlxC1aP37reHbt7w2cq7oQxTcjhsI1wrbVSgga5eY0FPLj4
C4937fqSg1hLncainVlrm6ZhjyZVCbzJZtTWIjTyMfm2OZ4WT6H8H+V22iN8xYqt+Ga1weA95ST9+naf
07ZDCEj5ky4RCMtHW/eOzt5jzotDWPvinRVXiY/MSfyL1ZseyF6pV1af2dUTgDB+T8ZJu/VPAAAA//8N
SFJc7AkAAA==
`,
	},

	"/definitions/hdtorrents.yml": {
		local:   "definitions/hdtorrents.yml",
		size:    1869,
		modtime: 1475399794,
		compressed: `
H4sIAAAJbogA/5RUXU8bOxB9z6+wuFcXVrAf+djAXYlKtAghVVVVihAS4sFZD7sWjr3Y3qRplP9er+P9
CqDQPEQ5x2c8M2c88X1/gJCiGhKUEy2kBK6VoTieG+r60r9tOYZ5VuLM8MD90jKUP6vE/EDIR7nWhUrC
MCd+fVEgZBYOzHmKC6dLsYZMSAoOIzRMEPomFoYJP7PyBq8cP2p4R4xb4fWl4yYt97Pm4jd0UyN8xY0b
7rvOQdbh/yfo9q4jiyzuSsZRXzI+s7ipYBr3zycm+0VJqKjPhzt41MfxaR9P4h39dEdv8t/f34eXd03G
LVOj0y6aOmThXJB2FAqwTPMEPbw8OkYv/JY8qQRK8BMExWMVzkRG+Ta4wNpoQssERV5Y8knIeWK/LaS8
KHWTrKQkQQfrNQq+CP5Es6BUIKt3hzabA6cpljuaAiu1FJI0GpBSyLZ+Bql5fAnSeMYgSGdMpM/KVMD1
QyqYkOeH/1xdReZzuO1Qg9J1tGthSVXOqNK2C3Mksaai1+SgnyvgeFEqpEnCde6nOWXkaOxtDaBMg2xa
9t1eScjgV+FIhLDMlGnzxiZCR8Gxd3Xx46DK7bzvJG9Wq/a4b+q/Ei+tY9IsKxjbmnXbbNzqrR4ez9fr
YLP5b70GTlqz60Fbv7/CqvJZdYaBU00XpvqhxVIs1Xu+zzHl1vrUGG+q1Rpp2XGHH4+9JMfqSJNXWs/d
KWEuqmwm0CmroWBVGwuMNPnrzpLG0k5F3bEMPfQJ4R7RTkFrSWel/SuU8NTwOzPszvGlBLlSJoxnndN6
oHVV9SpRzWBvhWMPzRB2KgIaU6b+Lui9TohYciYw2Xtb7O29S9Hf+1s5q70lxoi96lPvI45XVxVYKnjD
72GcRJMkihGKRmE0DEdRNG0eKJDeje+9j8g46UQMIM0/FDSsggaDPwEAAP//DbAV+U0HAAA=
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

	"/definitions/norbits.yml": {
		local:   "definitions/norbits.yml",
		size:    1734,
		modtime: 1475790994,
		compressed: `
H4sIAAAJbogA/5RUTW/bOBC9+1cMlMUmwcbS2vkEgSywDZoe2hQpGuQSuAElji3CNKmSlB3X0H8vKVOK
3VpwepE5M+8NHx+H7vf7PQDDLRKQSqfcGhdLOnPx5zYWVE5KOvGYtC+Vz3A5NcQtAPqQW1uQJAn8WKJN
eq6U0SJAMmpxojTHEAMMCNypuUuEeEjg4TGsTwnc34T1mV8nH5yeBnlO4P+ScRXCCwLvlJqa5H3qfkLy
MmCS+rtRuGoKj5yh6tXZmWKvugxSneUEnr6PQsbO+6/JEw8wSp4AFiNP19RytSYX1DpMwiXDl7jIi966
n8DMKk0gOqixkWcJNeFyi1VnWtZY6Rmpv3XIZVHaVmJpUK8vKFqtIL5RcswncZOFqooCsKDGLJRm28Am
2wJRa6ev5eyQs3mMOFPSorSpeoGU+IByaY4Ob5GLw+MAd7dl/LiEcJPfwbFo7C8a/IkYOqgwtZJeezub
xqVaLQy2Uretgr80XdTH126C0XnQDmJVzZyIZzeZT6Pr1Squqr9XK5SsqtabXHvLvpSol/FHXHrDTOuY
35L85oylqcADt9TOngcfgHOJLcFqIm3ez3Iu2JH8Z7g+8ZijYG2b8EKWuyyzbIM/OIb/gPF5i8uo2TAa
IHLvb+HaPdfnu+VihjoiMOiEPDy68rCzfK/VRNPZustpJ+xrwYVwiLNOxF1p+HTqIOfdUjgzZqr52DrY
RSfs05Kl01rP5Z7d5v6Z18irXvOeuRU7J3PL5qG3mW753lKotZqnpf/LrJuFQpjVN7Z+yjWOv11HSaBF
o139PahprxZSKMr+sH9D29vf8B/7bblsXGBuXveiz4/bR4IM9X5rBv82DIGY5W+iuIvp/QwAAP//yrvr
qsYGAAA=
`,
	},

	"/definitions/privatehd.yml": {
		local:   "definitions/privatehd.yml",
		size:    1256,
		modtime: 1475754160,
		compressed: `
H4sIAAAJbogA/6RTXYsTMRR9768IfRALO1O/UBlYRfRBEEGx7Isskk5uOxfTZMy9aanL/neTmck0I10q
2IchOffcc09O0qIoZkIQMlSidbiXDI0KiJG7gHzpkY8fAqKl2Xq5DSiYwlNE0PykKiyEKETD3FK1XI4i
JdtZqNWyHTh1gLfWIQz7+JuvbopvjT2IlXUODM8rsbo5VT/bPUJW6/aU1T1hndXfeYXdVCF2Vp0GEUhX
N5X4/ut2QHhfnMCrSCBrrgS0t7Fd2y2avrmVHDhL6blZdnCHbqzbVd2326JpPY/TYCdR//AErk9xfncn
yvfWbHBbJlTc388HeiuJDtapKTGhIxGcs64ae865CqcGonhHY0QEGmoOjaKUGhyLN6KtDDdF3aBWj58t
OiYD8d/SdW294ZjGEFQWx3zJfeb0Fs31k0c94zra/+rBHctPcIzmKZl39pBdRrLEcq2h7L6Zp+eLYJLX
Vh0nRiPo+vARtBrlhmd1HM+c6atM4WlUwAmQOiSzw7WP/wFG1pDeSFyfi5KVKIcAig1qmO7i9Q5NCjg8
BfovjYm9xsEmaduD0VaqS+KJV2BtzQVNwt8PHDjL7WXKTYXkL7JfxNiplWdH52kTgAL3QFaZ4Ks0XgPU
zb90vF6I2Z8AAAD//9NWYd7oBAAA
`,
	},

	"/definitions/speedcd.yml": {
		local:   "definitions/speedcd.yml",
		size:    1961,
		modtime: 1475754160,
		compressed: `
H4sIAAAJbogA/4xVXW/bNhR9z68g1JcYiCRbttNVDxuaukWGoqtXB96AZQ+0eC0RkUmVpJJ6Qf77LimR
st2vPMS6POfcy/slJY7jM0I0N5AT3QCwguFZ0J0/Jw5goAvFG8OlyEl0syYrUBw0oWIvBfwWoaSmomxp
iX4g4lZbhIs7naNBSEwqY5o8TX3M9AzxgjY9X1ADpbQhuzMh8ywnN+v0emEPL9C6ij/RvScnjlwFcrFe
xJ96MjvgHPm24VpiBT0/e3US+Xrhw84d89FUoBzznrPgNTnhlrS48+R83N3ZSGUc6SzPTjv2teA76CEM
9kHeY7ldQH9H5mHveQrMpsExZD17eSLKfgmiq7odujYb8GlwHgcMe+gDYB+Wb9LfVx99HtiyN1JoWUP6
F+fedz6gy9XUawfUI9NB9/dGfplejn2ES3fRssIlig87kc0cMWZHI0fkAy28BMu+kjLMIMNYr1vGpY89
Oz5nr/pzuuYM5JmDd3Yv/MqZ+1gDVUWVk38+XxC0tRQXBJp/rVhR3P5O2lCDmnSj5IOGpKkah2qooTBS
5WQj2Z78Shi/T4proAxtQzc1JNpoa/e8UfaH5cJUcVHxmp1PRp3bAZSNXPAtrw2okGrcv6IKSvjS9CAh
VJU6J+e35JY9Ti6yp9sEDXyQka2gliUXRxUYegdvHRyq2IGpJMtJI7VxABdNa8LFrQbVXR09PpIEx7rl
ZeJR8vQU9cKGav0gFTsWejQIQSnsWKhq6KFrUl5IYSgX+jx6x2sgf0hD3slWsKjrigFt8nChK8mmgh8X
+1HrZvmDiR2X5mfv8v2zBbVP3sPeZqtDutZ/0Ptcoxf4UCDMjZ0ySXDHsV6DgJ98P3Ojon6aULNh77ip
IQ9DPOgBS+qt6VYCfyn+bXpZ/8Hcf9vtxyvlIPr1jrkNMkbxTWv/HVQKtgE/2b/DHfxsW6XRTZQHrF9G
TLRHmXwQtaTspylPT/Ob/Cw/ht14RgN1Q0Vitb39rahuGH7C/L/vhD1Ibj4K+wDsqEXfcbj0DjVAUT3H
4yW+vv8HAAD//4W/SnepBwAA
`,
	},

	"/definitions/thepiratebay.yml": {
		local:   "definitions/thepiratebay.yml",
		size:    1391,
		modtime: 1475396077,
		compressed: `
H4sIAAAJbogA/7xTX2/aPhR976e4Sl9AP5Lwp79tRJqmTTxU2qr963ipmGTiS2LNjTP7BsoivvvsEAxt
QeVpPOWce+7xyYkJw/ACwAjCBCjHUmhGOGdrSxbs3pK3OcKXhoUPDS1ZkVUssyMswso4RhS/TGIfAELI
iUqTxPGhWaR0Fp+aGycoV8/mpVYPdqBVfGFHKSvbI1K7kyktsMUAw/4ggRu1tJRnhjsmnkwnnh0l8L7i
QsVTwVEBXMJNZUQKSwf3y1d++bONqZ3MQUilKPeq/205U49eOdTq3e8SrlnBc5TcS1572+t9pDfN3gEx
9qqRJ8fjx4kumsG94vsSaBkaZDrNE7j73QP7bFTRAyxnTtyOGm3JyIqCuhYLiL5WqNfRR1yvlOZms9kK
47p+Por78Xgc9+sapUELNaZYkIUF32yCxlqrlQ9kUGJKStujLreu39BUkoDmiq+BdJIz0yEeLVHTbd7d
OiyEbcx7tN96vcOHrsSTgvIwzYXknUEXWCKZoS32ckakxbxytzvXuPD8QkhCbfa+7uptL7wppaAD3pro
zNhWIYiDHoQDmO0qFyTxWLSII32yf4p2ZBET0pyhPBWYq1UhFeMv9jC0Pdw1uX6+DSbtVjB7wd6IP0ff
46n1QhXkEk/QpOdUqTHDh/JIl8F3eyJ0ov/edXvB7h3tt/63IX6Urh3kT4IYRP7I60SWUbdVSMQ0P2fj
qvs3AAD//zrY6S9vBQAA
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
