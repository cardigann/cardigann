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
		size:    1955,
		modtime: 1476329823,
		compressed: `
H4sIAAAJbogA/6xU72vrNhT9nr9CpGM07Nle+uP1PUMGXTs2GGFZO7pAE4Ji3diijuRJctMu+H/ftWzJ
TpfSfVg/pLpH5x7de3StIAgGhGhuICY0LzKqqOESIUG3CF3X0F0L5VSkJU0RBhGUuka4eNIxLggJSGZM
oeMo6mTCJIkGuJvQomUl1EAqFYc2JmQcE/LHQ3R/28ZnTfyLi8/dPnF/JxjfPtze8aKlXBylzGjy5FUv
neq/KP6gz0iZymcsrSvmqsM870uf10rZ+ODAr/3MtyyvNf4+JvP5nJBeWZaFoKOgQdcl4zJ64AxkQ8Gl
20a/ZjfRz3hb2kFo2Y0UWuYQzdfyxcEXHTy7P3foZYf+yblDP9eqLriyR0xp4oCvFphlUkDwm8lAYU3X
RaGncs1zcBfZ9OYibONHKZ909NMa/zn0zDVnf/s75+2Oi7F6e1Q05ToZWHQrWTdHGqhKspg8/rVsEfMc
dOCnmqCl+ESgWNbpdkCb5IIa5ERcMHgJi6wZKw05JEaqmAxPtKFGr2wGCetgWCvkMuXiQMEiXmEj1Ta2
vzbkoiiNL7fUoJpPbLjfkxCvYMPT0KGkqoYtsaBa76Rih0SHeiIohbX6nCPloF+gdf35DtywdT2GO6oE
F2kspDkNM84YiJHlGdDmjS4XOFVWd+B979uAigqE0f7ow9bJN4rusJ0mU5frLTeT8bf7vcL3BbBF/0ZU
1YbnBtQKn43H5WS/D6sKeSBYVbXZRk1qW34vQb2Gv8JrbYr2rii56w2Ia9bQdQ4nbZUrGxGjwhZo7o5D
znxq+2q9HnPOsDCppyOROaF+nxqj+LqsX9VMwcbjTT+dFaR+OJs5UJDCS9HbQBWVajSqM2GxeFws2HeL
xXJyWi9GbkoMN/nRi8Xy1jxdcbHB0U2VLItm/QOhsTBZkGQ8Z6dnozaVgaE81/+D0nseMLkTuaTsoyN0
QcWB8vgjZc3/fseCnsql7xTv9EP2xYiEhm/h2MnWcj9bwA6u9R29K3d6DpBk/yXjy2jwTwAAAP//nIrg
vqMHAAA=
`,
	},

	"/definitions/avistaz.yml": {
		local:   "definitions/avistaz.yml",
		size:    1252,
		modtime: 1476327553,
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
		modtime: 1476327553,
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
		modtime: 1476327553,
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
		modtime: 1476327553,
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
		modtime: 1476327553,
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
		size:    1966,
		modtime: 1476327553,
		compressed: `
H4sIAAAJbogA/6xUTW/jNhC9+1cQRgHHSGTF2e62KyAN3BTooShatMVeDNegxLFEhCIVchSvu9j/3iFN
fXjrNDnsxfI8zgwf5z0ySZIJY04iZExAbbSRggDNawJ+GgABrrCyQWn0Ca64LlteUjLopHUekfrBZfSH
sYRViE2Wpvv9fiF8xaKGdEJrBW9iTsERSmMlxJixZcbYr+aJkAjcELBqhTQxfkPxXx9i8C0Fv9+nPxPh
Lv9tgGLwjoLfsAIb4++/iN+HZulKyxo6AtcZ+9GYh67fchnj9N7Usuhh4hFopR+kgI7c8rsODb851U3C
Um3EcEYH3BZVxtaPm4jgUzKAVz7BGX3FoNn4cmVKqY/FDUfKSXlRmFbjtuJaKLCLpmqO2wBWRmSsMQ4D
IHXTYr+xlsXDUdzpp09scW/0TpaL1oH1KPv8eRoTG+7c3lhxmtiho0QL2FodWUXszzavJWbxG8G9xOox
Y8sQgrXGDtNQUCDFbGdsvZbidpbn24IGtvXAbMMuGeZGHGhd47pQxOJ2ZkHMNpnGKikqqcTFch7bkROc
t2QMv3p/f+baPNEMHR4UXLHcS2Q5XY4TifxUaReN1qhtwzWoXqWBEQ/DZ7IuwwJHtDJv/XVEiepoyZ1U
CLbXMInX00IJH5ueFLelI7X+CDzYxeJy/s3UE4u2GjOjhvG2nLrjsQV7iIr/AgevtBtJHe8qZVx3yrX5
AK6U6vtwJfEEcgACyE03MYaP6D2nBqSVomtszd791xyLAnmuIIwUyPra0FQF+4EF2H+DhvS1mTZ4sSYS
pd7MQxCEJGExqpnoyzfzTlG+IzKdNemxoyNVwMVo5mMe02Pj25mCHc420z7lC53+R6uRXn+vBA2GXdxl
RrP5Hel254U79gMl+nb9oM9QQnFiVJoBH3wxWKqysHsN2+ADR2W6PEO5I9K9XN6nL5K6GZMqTF2Tgmen
+3zVc0cRZq+V4eLFbm9f083Jf14+zbv52NbnbXJS8N5v7V+XmKgAiuo1hcvrvnIy+TcAAP//btL78q4H
AAA=
`,
	},

	"/definitions/eztv.yml": {
		local:   "definitions/eztv.yml",
		size:    820,
		modtime: 1476339286,
		compressed: `
H4sIAAAJbogA/5SSza7aMBCF9zzFiM1tJJKI/i0iddlVV5WqLooQGuJJYtWxU3sChSjv3gkkgaq9Qux8
fI4/j2ccx/ECIGimDOjMBxEWaxGff3z7LsKgLVssB9fGbRh2tP0ZMlkAxFAxNyFL0+FogmW6kP0cm9HP
kal0XtOoAdYZCPYiaqdu+3yIA6HPqww2v1Yg6+DsCqjZDuHRumQbZAktu04XkHxtyZ+SL3Q6Oq9C31+D
adf9a3UdmUB9LyZZ1ffLC82741xDIEM5Oy90xr2hpHC+rXcVoSK/2wuGPLDfDP359FK5A/mX7f9CV3Sh
yagZPrbiNGl5Mv3mDNZTAzQbupm3Wlhllqs4r7RRb95GgGNGEaM24ZkjAMjs9b4dhl15KiaUO1rjUD1k
vRNWUmNpiVdP3xL0+fED30dTTdKvh+kP0ZwotGHyd+0Yfuf1J2PTyMTvDKnQl0HGDFi65Tx9Un8BXrnz
Y7T4EwAA//8NsQ5ENAMAAA==
`,
	},

	"/definitions/filelist.yml": {
		local:   "definitions/filelist.yml",
		size:    1867,
		modtime: 1476704031,
		compressed: `
H4sIAAAJbogA/6yU707bPBTGv/cqLN5Xgmqk+dMCw9I0MRBD2iZNMCEkxCQ3OU2sunawnbIOcS+7ll3Z
HNd2UzZglfYlqn/nnMfPObYbRVEPIUU1YDShDBhV2gBOZgacGvBxCRjhZUNKA4FHjWoJ5VOFzQ+EIlRp
XeM49goDKeKeCeWkdik50VAKScGtEcpG2Hy/XMZHnM7AwTRt4VFTUOHT9lryScxNaTw88TTp0HesOScL
H2kDPnIqJNCSu9CwGzq59FqjLj7zND3ET+mk3YKLUGAtHQuuBPPtWA30+Th+b+apfOIBttAbztaWr12J
j6ZuSmcn6D+kgCjBUU3yqVfLhi7h4omENCTUQmoP3fDd0jpCV1dXPQtmolidk5GUeYXR9e2NI3oereCu
23MXQX3TljNRUr4srok2ObElg7qqLZwIOcP2e01yTQV/s63JFELS9nIbyutGBxONArm8k1v392hgpjyh
5cBT9PCw5RJrotSdkMV6oqedxIYzkU9NWrrVupbEWFlzPVsEywoY5FpIjAZKE92oMZGooHPMdRXlFWXF
TtZ/BEb9ZbeUaZChj8i9LAklfKsdRIjIUhkrX8+ti58/dgav3vb/t8bcpDvOxlLcKQjm1gflD8Z2/wEW
bdeq07Z5hxgldtHKrMpCh+YrgWsTdQ0AK0Kee8YLHKz/XqnJmEFnEmkfqZpwRFbtai3puGn/cyoJk8Af
Das7sNsG5EKZMl52on5yxpa/m1SbzTdwlz3nzqq5QAGaUKb+mXan80LcmdtIik20h3+prej3jeZx0Ddv
k+tQIWEm5tDeOm/WXIFNBPdfEnz23NvdaiIV/OHU0z2cjHCyl2RxksZZkuyH2wzFmt7LLg/7LpkB5NWG
xWnS7/V6vwIAAP//TrgXHUsHAAA=
`,
	},

	"/definitions/freshon.yml": {
		local:   "definitions/freshon.yml",
		size:    1347,
		modtime: 1476327553,
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
		modtime: 1476327553,
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
		modtime: 1476327553,
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
		modtime: 1476327553,
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

	"/definitions/immortalseed.yml": {
		local:   "definitions/immortalseed.yml",
		size:    2354,
		modtime: 1476327553,
		compressed: `
H4sIAAAJbogA/4xV72/bNhD9XqD/w0H5EmORlTjOj3HIhswulmHr1i1BMcA2Alq8SIRlkiUpp27R/32k
fiu2s+WDI57fPb17dzyHYfj2DYDhFgnw9VpqSzODyHxU0LWL/voiytDEmivLpSAQUEhQoKYZWE3jFerA
YzIqkpwmLhtFmJsixMXKEP8EEEJqrTIkirqvHK4xevvGI2KqamhMLSZSc6wDAOcE/rQpauj8HcFtlsEf
+arUWMBGBB4+RreCr7ED88caMnJUHyYAL5iUMg3JBYHbnHEZFZ9LKVceUT+3uDMCP/tAj2mS8oxpFCaa
UG2lFA1+NKrw0buCs8S/63GOLry66BfXBdNyFscacnZJ4L3cOHeiu2kDKSPh3bSBXTWw+x3Y7/IZpvjU
YMevYO8byvPTyphewe9zw+MaMr7YbdRRGakhV0WL7pWbgQ6kOBt42DRMFW7aY3r46AZpfH2qath1gbqb
wi7sjiepr5IL7ge3Ib4+QOwiIXweXY5r5Pe7wA7ynw1vnBkfUnGP1EgBH9w1MV5Sk3F5kLuf4d0vc9aS
dS6EQarjlMDs06IO2U3YRk88whGdAKpFSaGps6EiUNQ6WLTU8tngUKWVnwYzjK3UBBjfHFmp4Ef/NFSU
MS4SdzKKCiJsGsZ+zo9HgzLxiWcWdSsvrBaJxgQ/N80CoDoxboP8XUiB49l8zubz4cniu0FQisxkwkVf
pKUrLMKtzjXaVDICShpbRrhQuW3fnxvUpYLg61cYTqR44smwjsK3b0GNVNSYZ6lZH1lHW6RFY0mbVCjz
fLFqZb00MEXK3FUoPfyCKs4cbVln1akD3fiJyZsSsbe6FW69OFNp/itHvR3+VgW7xXERZznDR+aEPDpZ
bi85Gti226TatVsCp/3RerRb5eyzj96wega9QLJbq6XLDI+M3+ruofhwRdulZFv/X5OUmmPnwaww4CZw
azGzXIWxFNYpChbNFGHG2hdYbjMk7fD0zD3IVdrdmdGzwU6x+zgtG+airmHkWGhnbK3VfJn7H8xUt3tz
d+67s//Jt8W4RJF0v64vQS2mESfX66I9+8TRWWHGTTCpUMHiP9Ux+SwySdmBajsOnQ98uXuu9Svshn/Z
35se88WgnRZkfacOpFw1KRlinP6vnOsmhzlX9+KDXoIYVGPiJ9PvtJmxW+fvPHhyllm3uNyPh/1h7kbT
XaZ/AwAA//+uZYmOMgkAAA==
`,
	},

	"/definitions/iptorrents.yml": {
		local:   "definitions/iptorrents.yml",
		size:    2218,
		modtime: 1476587223,
		compressed: `
H4sIAAAJbogA/5RVf2/bNhD9v5+CUIshRmPLUf1DFrANnYdhw1Ygw4IhQJMVtHiRiFCkStJJPE/ffSfJ
1C+njftPknv33rtHSrqMx+NXhBhuISI8t0prkNYgJGmG0G+XVy3EwMQaSVzJiHivL8il5g/UArnSNL4H
7SFHUJlsaYJSkONtqRJc3psI/yBkTFJrcxP5fjtqEqvMf1VmAGu5TBpqHSBW6p5DBRFidzlCFp7sARB0
AyIi65qEYEzzg0GMwRKlORxqQpZBRD6oB0QcsHSA/9fPByxcHWOr6TPY4hhD6AibhQ3vV4fNZ4Mgi+DY
7F0r/EVp4Il05HCgDqbHI5bPRFm+i8jV3060KIu2OZ/362XYr4NK3A4IZoP+QL8Y+IfBoD/vhAm7yZar
gXPQnzzsY5D+Seq64a9WAz0S3m8ZV51r6dbh9FD7fyhjBBh3z6sBcTGtjN9LnrkX9AJnX65dHwdfrv0P
NHbAzDlXPzf41rqHjZl+wrIZNTvU/lplPHZwiP7X19euwmGdat7rhd1q1dfN6qoqM8XaL8QA1XEakY+f
b90X9zBuwfOSYJQ8J5DflnKhEi5rcQY2Vaz3vXKZb23jXTdwb+z3ZLJW8o4nkxojReFVJAvGOnpOLc70
s90kT/Nylqa4d+ruocclgydcJBWjji8gxr0SkUn8qeLXe6XK31F6vvWeSfhG08cqn8YdBhiy2SBFsd9P
iuK7/R4kK4rP35dn+HMLejf5HXaPSjPTHEKrx859ukCepRsBr93SIz8Qu1FsV/7WkbTpOE65YGfybTCK
UmrOLJvYT9V2G9W2dxwEay+zjrZzdXeUZR3DixGOoA2LWqv5Zlsu+1TDXYPfcWFBm9auXcAaEnjKOw10
0biliffPzc2PZzc37O3ojedeF24FvBgq6IZiYCkX5ttEXzoJ/i/Jyvt90W1+ihtTj1Ioyl50m53iZvi/
L1/NYuRm4xM+5U7wNYltu4C++ihNLrh95kl+JN5/3jkZX5Dbb1ORzY6gcNroDADrDf9C8NAdUwDE6SmK
1ej/AAAA///I50+kqggAAA==
`,
	},

	"/definitions/morethantv.yml": {
		local:   "definitions/morethantv.yml",
		size:    1281,
		modtime: 1476327553,
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
		modtime: 1476327553,
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
		modtime: 1476327553,
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
		modtime: 1476327553,
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

	"/definitions/sceneaccess.yml": {
		local:   "definitions/sceneaccess.yml",
		size:    2622,
		modtime: 1476587223,
		compressed: `
H4sIAAAJbogA/6SVbW/bNhDH3+dTEPYwrF1t2ZKTZgKCIYjRpdiCaMuQeagzgxXPNlFZ1EgqiVf0u+8o
UaRkOwnQ5oUj3f93jyTFwWBwRIjiGmKiUsiBpikohbacbvZsGc1XJV2hHfJBWVl4/knF+EDIgKy1LlQc
BC2nIZTBEcopLSyWUg0rITnYd/N3GhNyJe7RFtxMCenb5+ntdPCHg8LQUY54DE8mDni7r8/uOXP6ZF9P
aPpJOWBsIvx5W9dQ/R90EoRWv7T65Y4+Hld6LXZShyctpZs0QoUkF8EvOG8sjFT/g+TCAccxuRC5EhkE
yU1E+ha4iXz0URtJiGcSz5x65i/OHYPPnok8M/soHqOTUYPZV4/+FB+qutNZxdh4zyEY5oqmTxAT7I2c
l4yL4DehVIZ7ClevVDwN3mU09bOPHHeVRLXNcvjusWOH3XIGArNa6Pb91K/JaUORTqSd7WJX7v3NNYY5
LwplHn1vVh4xum300fT8b+9vNuS1XoMMrrjCAUyvZ957vCOaXz+VsT8xuBv7SZgE7rW7KSdh53C10d39
PYl2zmGL7Z6kSeugGKh69vJx65xY+bI13Ukje5M5H7PZDHn87Y4petvRdk7PcS12vwTj0Fqx6laS1tDe
CQl8lZN+83DwkxOFL3t0Ekejlx26NdlZ7MFo60aOngariBW5Ecx/VxVQma5j8uHfO2vR9wNvfGMAJfI3
BIo7456JFc9r54JqZILKUgcGvRYsJkshN5WB50WpXapSgawvjN7nz2SIh37JV8PGSr586VmwoEo9CMm6
YGN1IEgpZBMcPwrKXDpuFgoySDUCpNevSlzg12mBnae1twalY5ew6oTniJgmbf/tLmmWHWiJfCfpA2ao
Hc5Mtb+XILfDX2FralVY7Pf1WM7COq8UD63hNzVqOdR6gVqlLDlkzFH2Mtwe6k0z9JMLvS2A0A9rCct/
zno/o8dZ787hFBH+sTTXtyGcfckzDbJ1xZr7uV4hCSt4LFoCRpErha3O51X4H+Zz9uOrZsU019nB2dv6
qgWmhwqqPK3AQFOeqa+K02osFZsN5Pqb4zDxkGeCsqfisAXLXgyi+H/PDcbI37YeZiFeD83PfK6Gw1fD
182iMNw3z6SmjAFzGxFYJ/N+oTVhgQwgXT/v0CBH/wcAAP//+zpVrz4KAAA=
`,
	},

	"/definitions/speedcd.yml": {
		local:   "definitions/speedcd.yml",
		size:    1961,
		modtime: 1476327553,
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
		modtime: 1476327553,
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

	"/definitions/torrentbytes.yml": {
		local:   "definitions/torrentbytes.yml",
		size:    2621,
		modtime: 1476329823,
		compressed: `
H4sIAAAJbogA/5xWa2/bNhf+7l9B1C9e2FtkJb6kiYBiaJKlKTZvWVykBupgoKVji4hMaiQVx/P833dI
ibqkdlMsHyLxeZ5zJ0V7ntciRDENAdFCSuB6vtGgEOR0heCnHLwowITyZUaXSAD3Mosw/qgCfCHEI7HW
aeD76/W6V/fW46D9FmpCmhbakGpYCsmgWBMyPAkI+XRPSBv/3dLwscAHA4v7k6ucmlw54iwnbgoCn47p
N0z8q/urO5a6OMdIjsUThibt4lkLNxzW6N8En4TAwZ8+sagusHnuY0/OS/M8fPleCEYVf1PjbxzfP97v
AGvw7lyQfqW5SLI7uill11mSkBxz+Y7yfKt0L2rNGO2rti4w3X+fRUzYZXucKRb69X6dOYH/q1AqAaWc
6jqhYaE6RVGhGt8OXKmj0vSeRSBIm+SGZl5ulMMDGrtSLonTlzU+90+HdfKrGusC087bS/8D7niUEPts
1GgEl4IrkcBBHh2MabiXNsV/7f/S0W8D59y/nQxIu+Anrk+D40ownYtnpzDvrpU1H58ZI06C706BXShz
sKPMFULHIJsaUwb5lsa14qAGj/J0OrU9UDVNrSfmmKCElH9tY1Cd7fNgD10dkoK26GeYO3iYx/24wk+U
Al2YfRx/mPzsjo6dlNP3g9oCdzp+K95ztgIHnVnoWkhgS14D3WFrEAP0/bvpgj9mKmxZdCWi6gungMow
DsiXvx4KRD95FXhkBErwIwLpgzFPxJLx3DilGjW+RXppnJ/OFehYRAFZCLmyAONppstwmQKZf8TfbLek
h0NbsGXPoWS3e1MIU6rUWsioKXRoKQQphQxKG5uQpo/QTMqkpZS5IlpudgoSCPE+wCsm6ml41pbBe0G/
8MY47mjrqVV2q17+XIq1gjJUs1zyP0nXWEJu985U8kcGctP7BTamDlWrGA9hod5uJV5pgFWXF9Jut2CJ
Bvkn3lFfHtBRb7f7/3YLPNrtcudKy2/5N1lWM3e1x3qVkLmINiRiT23GOci1pGkK0gIryvhCmsGYVSi4
xtvTvBNN5wkEXMeeWHh6k0Kn3yXaetIy4EJ3ggWTSnthzJKoawMvGCRRmURx2W72j8S6tradky6hRw0E
Q9HSimot2TwzPxZiCYsSz/tVTYKY3wL51pOwhOe0RqAXuVTY+XyYs5kZ52z2E+b4rjObRT92/+n0fui6
SWmmkwNb6UWWjWVhEIGmLFH/2f5QyZFY80TQ6Lscv+JLsb9fL/BtWRGO8lX1qFvuPogagzlgcO4MEoAw
/h6Lk+Nu698AAAD//9bN7vg9CgAA
`,
	},

	"/definitions/torrentday.yml": {
		local:   "definitions/torrentday.yml",
		size:    2975,
		modtime: 1476327553,
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
		modtime: 1476327553,
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
