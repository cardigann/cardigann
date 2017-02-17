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

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
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

func (dir _escDirectory) Open(name string) (http.File, error) {
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
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
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

	"/definitions/abnormal.yml": {
		local:   "definitions/abnormal.yml",
		size:    2874,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6RVf2/iRhD9n08xcn/oIgVILm2vspSrwECCEgINBFWKULXYg73F3vXtrg9R8Hc/efEP
nJjjSPgjkWfeezNvdryu1+s1AEkVmkDmjIuA+DUARgI0oVUEfMLciLhowkLUFyKJULaUZg0AoA6eUqE0
m81MorGStRqATcIUYhOFLhcU0+fsZwyG035325l2Hg0TBvwrRdnsTDuHQP1RARtXoto/ApoOx5NeUbEK
ctvZfvp4ERZKt5VKt53t5cWfx3GP3cHTP+09k20/eiTrKmjv6f7+GHIy3Y47hY3J9JXRyTTpbR/xsrFU
o/ddfu8gd9Sy7rapdnVyr/grQHfQH4/7w4fxTn+oPBRlUOuhP+jusi1GAyxnO0MrpXa4HQXIFBEvp/k0
7lvb3n3LMkxoRQ7lzXsupY/y5YFr4GB0leMGo6sqiDV8sLqPkxw2pQ7yMnBkbVujZP9G1qvETWvQHetU
84YEL9eu2x4O75J0m/OlbHbnnC9rGhJwp3hvJBJheyY8f5mlEfW1XgTPE4Dk7BwwnCV0n7uU7cghUZ4J
TR1phF6ogwsuAlP/1Y+UhZHKi0USxe4uMDYbaFicLajbyKIQx0YKDImUKy6cMjCL5kAUgotMvF7dkHaM
Uia3zd6EJPpoKy5MaKyIYJS5OqlQKjNvQsspHl5eaLlaPq49/0ZTcSGQKZlgDPjprx3merOhC2j8HaFY
N+5wnTQu4zixU45BEkRfYhz/stkgczJ35dn9LMhKj0MQ5iI0rPwKjGObqOfZ9WbTiONfSxLF+Rrvbkfw
1d7W5NNL3f+ryNxH+Axqzp118l+YjKsPz7ZPpLw2bO57SBxjdrZbE4q+k8ul1/k6P6BCXzkmU17d9qjv
fLg8g89AMhBRStB5lHxrPIGLLLygvkJR+jDU0y/Ql8SpVCI77exHhCtN0FPMXgKq/L19OdDPx/1+HL5i
PifOUdZvJRcHbEj6//H6v5/liO+4Fhj6xMYKx8/GDTfOwXDnxuxE4kATg9OJd5q4LIgOUcedXiUzkyFh
7zXc92ENBJIeTm7dw0igTKgej4R8Gz+jn8z+T9c8B8Mh69NrSwwIZbvuV4jLNytkAifzA0519YAz5Z1e
njDNXiN5w+QJy7iVVBKGyJwKpgHE5cVVik5p4w7s6h/ZW+kj2t6PMD6d1b4FAAD//+rvg/I6CwAA
`,
	},

	"/definitions/alpharatio.yml": {
		local:   "definitions/alpharatio.yml",
		size:    1961,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6xU72vkNhD97r9i2JQjS892Nz8ud4IU0qS0UELTu2O7kA1BtmZtEa/kSnI2qfH/XixL
srdNSD+cPxjP05s3o6ex4jiOADQ3SIBWdUkVNVxGAIJukcBFD312UEVF0dACCaCIG90jXDxoEgEAxFAa
U2uSpqNMkudpFAHktHasnBospOLoYoAFAfi6TL9cufhoiH/18bFfB/8cwNfl1fLqM68d5eRFyg3NH4Lq
qVf9DyUU+kAAruUjRz02czZigfdxynNSNt4r+Gma+W9W0Fr8QGC1WgFM2rKs1WrlKQsCFw3jMl1yhnKg
LDnzy0cEbi7TX+gWtYeOCVxKoWWF6SqTTx4+GeGbL8cePR3RPzn36Ide1QdntsQ1zT3wyQI3pRQY/25K
VHAAF3Wtr2XGK/QHOezNRwsCP0n5oNOfMykfPHrkN2ff05Vjt+LjEwK2VHrNdR5ZdCvZOEcaqcpLArd/
3TnEPMYj+L4naCneA9Z3fbod0CG5pqYkkHLB8Cmpy2GsNFaYG6kIzA60oUbf2wxI+mDWK1Sy4GJPwSJB
YSPVlti3DbmoGxPabTSq4RebtS0kl1JseJF4FLpu5og11XonFdsnejQQUSmpvHj8ckPWM9S6/4WjceTG
nSY7qgQXBRHSHCYlZwzF3DINakNCR86uTA52RcH9qRlGKoXC6FB+3wD4TtEdgdmQqZtsy8354l3bKioK
hOQy3BRdt+GVQXWfU3N7d962Sde9a1sUrOtctlHnvTl/NKiek9/wubdGB2+U3E3GxG/W0KzCA9flvY3A
qMQBwwlyrFhIdXfX8+jdRIwleT8juayAhnVqjOJZ09+tpcJNwIf96OkhxO7CVVjgUz1ZAKCq0ARmownr
9e16zb5fr+/OD/uPuZ8Vw02Fr7SX8eKei42EpFCyqYfvH4ESYco4L3nFDo/mLpWhobzS30DpNQ+Y3IlK
UvZWCV1Tsae8eEtZ879fsWCichp2Ss3b7JM5JIZv8aXK1vIwW8j2jvUVvTNfvULMy/+T8XEe/RMAAP//
U+2CWKkHAAA=
`,
	},

	"/definitions/alphareign.yml": {
		local:   "definitions/alphareign.yml",
		size:    2757,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/9RWS24jNxDd9ykKXMWJW7KziUHAAYwYgYEgCBAp3kgKUG6WuglTZJuslqIY3ucUucgs
5i5zgbnCgP3Tx4IG8MIDbdTNV6+Kr16JRH/+8DFN0wQgaCYJaMoCPencJgAWFyThJkJ/tpCikHldsnZW
griB27sxOGvWwB6zR/IiATBo8wpzkkA2rUICQDZzSttcwl/jX9OryNH2McgEACCFgrkMcjjcbD4IlCQA
GZYtKUOm3HlN7Rrgd7eMq/a5Aw7vbuXmdTc02oRGXWh8L2F83y/q9PqxgUYN1Kf8wQV52TySGls4tREX
CH1WSJg8zVqEl+kGPI+E4Ow5UDmL6cbl2jbJJXIhYVgjTWHiwikJc+cXNRBfmuUEsziJv69FwxfNbtqW
FfdaqkC+maR4fobBL87OdT7oUHh5ES2xxBBWzqtdYof2RPLe+a54+lpv07+hjJ2XoPRygBUXadQLP8d1
8ztAQ55ThTbvLGQKLHsxsWwgZm3zkOyXxUnhad707SoWtYseWTvZlvqHJYhP//0vYqS1ftvgHW9z4gPO
PVXk160bv9E6uhB6G7xbbY270zXwFCrTFJtrMqrnsGZDsncImb1+qOKZU8iYxlm0QUWM2gR5wM2m7Thv
dt6TZTE7VDGSumJuZY1DdbTaAnNLLL9eDPlIC6gUqc4R/e8RZoz21pEiH45wG0IbN0RZcZTfMfb6XzpT
LWiOded9coZhSyaA+F5IEJfdiajKN6e2F9b6SAJOMoMhXE+FwQcyUP+mpdcL9OupmMnMWUZtw3dTsdSK
3FScwQ9wOC1UWUYh7KUt4k13LE3budvLubudirPYS3+DihNQ/dOPF+Up6r68uDpN4aO9v8nom6neiHhv
BaFwq7eervoD4wQEj7YEv/+MW8HN/uLVdVt/f4nkSwAAAP//oep6DsUKAAA=
`,
	},

	"/definitions/apollo.yml": {
		local:   "definitions/apollo.yml",
		size:    2322,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5xU3YobRxO9N/gditGH8X7szGa9f6ZhHTZCJpA1JMHxzXojWt2lUaOerkl3jYQi9O6h
508joxhtbqSZU+ecqqmurjRNX78CCIZRgCzJWorvThYo4KF/1xiUNyUbcgKSByiqYBSwl2qJPokMK11e
yRwFoEurECF0irRxuYA/Pn9M39cs45ZBxCeAFBbMZRAXF03ezJvy9asYU7LsSEoy5uQNdgDApYCHShuC
EXyKZXT4OwG/jmEED2VpjZKx1tDFrgT8RLQMMIJJWj91kevW7aL+nREto0X33NNuBHyilcHG4RGld8bl
8MVopJ50K+DzFxjBmArUmw69a3NfjKkwKjRxo0LzrQAF6cHXBZReLQQ8/fXcQbxK9+h5ZARy54Dlc2Nh
KTeuNSglL0SDZOWibDMgL0gLKClwgxhXVrxPWgX0zZEn2y1kY3Jzk2cdCrtd0jFLGcKavD5kduiQuUQs
LeU5agGXHdiUCskj5WBcy0Xvyfe1pBDQomLyAubki1GtiU/wAUIpXbZumt8IGAOLfXXx65m8R8ehaUCM
+TgNBx065DSN77JaMwosOUxrWZu1MWrP4ftO3zT3f16u63Z56XKEbNyP9G43N5bRT5Xkp+022+2e7y/f
bLfo9KCRTc7Avu35bxX6TfYLbmLHw7Dl5DX66WwjgE2Bh+habkR9jTtYquY2z+TgDrW5qllhuD82T+vh
fHZtYjmzOGo/f1q/wQfgGelN/PdZG2qUc4NW7200rZ0lqXtg6CyfFh7nf94nw97+2BR83ynfGH2fPO/l
ktmbWRX3WFT3iQab61gubVZZ7qkqp3HI9gyPBa1QtGffXkTDFv+zy3nNYpn3+6JdbpujlqwzFYdQkd2H
lQzDAqC2rGn1Rh7ctIOgHKxEAe+OcrDedwKujjv0G1HA9XG9bZfiqt6JAm6O0lS9GwXc/lvUqCDgrm8R
FUWcgBfMyWlzMTcWj7uyFo4XqVoYq99enfWDJPn42R/wr3t+MH+fwL/p+bmXsxMKut0nQNToT5Dc9RKL
qBYnad6ffXtTV2SrAueyJn9nIhPhibg+UmnR80x6ocixNC68/ZrMPWJdxtfkLBGQ/JC8QPzx98nkcTIZ
/3xc/P+IXfZYVb6o7lb9TwAAAP///towFxIJAAA=
`,
	},

	"/definitions/avistaz.yml": {
		local:   "definitions/avistaz.yml",
		size:    1258,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6RTTWsbQQy9768QPpQGMrv9oi0DSQk9lh5KQw4toYx3ZK/oeGY70to4If+9zH55DDEu
1IfFenp6kp52lVIFAJOgBrMlFvNQAHizQQ03ffyjAHDGrzuzRg3oVccJIf+bdQEAoKARaVlX1ShQSqiK
AqA27UipjeA6RMIxTr/F7Z363oQd3IYY0ctCw+3dIfs1bAmzXB9zlu+Y6ix/01kKfXoT7KEPo4l1o+Hn
n/sRka06gJeJwMFfArb3aWYX1uSH4tZIo6EynTRVD/foKsSN7p99SL7tZO6GG0PuV8cYBwcXj49Qfg5+
RetyQuHpaTHSW8O8C9EeEyd0JmKMIU4t1Im5+s2ROV1pBtL+DmsJUUNpHEaBa2i1l0bVDTn78s1FzxVk
0fNQg3xdh85L8mS0KzNlUclgPH8if/XqxcC4Skt86zDuyy+4TyvwtEIMu+wk00hilg7L/pnN9PYCrkGW
we6PBk1gHE5A6OwsN75b+8PWWQObSbxOEnQEzCVGJNKyS1+BkDic3pX0/3lhKEcL1IocHkfpzGORRTHk
+L80jsZrIq4m7bDzLhh7TnziKaqDP6PJ9HBi4cy495Nx1sh59rvkO7fm2da524xoMZ7wKhP8MLV3iHXz
LxUfL6D4GwAA//+d3Gs+6gQAAA==
`,
	},

	"/definitions/beyondhd.yml": {
		local:   "definitions/beyondhd.yml",
		size:    2359,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5SVz27jNhDG73qKgXtp0MiyY9lxeGviFCl2DQROmi6wCApamkhEaFJLjpJot3vvU/Th
+iSF/lCyG0dOjxr9vo8fOSPK930PwApCBmsstIrT2ANQfIMMzqvC1cIDkFwlOU+QASo/t2VFqEfLPAAA
H1KizLIgcBbDl+Jr8Na7DQZeuSgSCZW0HvWakdaPAqsSABUZMiB8oaYg+Rolg4sa8gAinjUGESdMtBHY
PANMThks9ZNAG5zLfMULAPihrkAA5zL3DS8a9nTcspNFXdpiJ4sGm0/2Yb8YvsSbW73qwNNu7avXfuPR
fJQFomHPwj42/OAs533Y6ckoa8Bpr9/yOmy4cW/GFW7yF+c42kf+qgiN4hKCrROotuYST98t24ofnr1b
tbpc/vapkc3GB2RXi087gvnsgODmIly5TCMGt3cdV7O3d3CT6udXwxSG/fSlinSM1tHzfvpqcXvnIp/1
o9c8emxtZ/3sdn/DaT/7++W5v/joPqsZg5/zWOjgo7ZWorXV2ORWRBCAq7munDl4eT3pZqyBl9cTN2FT
x92JGPV/uN3vZTrrY7cn6aQPdD2r4E3ZEHdxWOQmShl8/nLvLqInvysel4DV6hgwuy/lUidC1eINUqrj
nWtMqCyn1rt+wWDw7RsML7R6EMmwrsH374MKIrTk8IxTyiDYoLU8QTvM0qxc0XASmjX0CzEY/PPX34P6
Vq1iep14bfSzxUr5Oo7bVRXnAxbP2sS2TVIqO1JiRNowIL6WOCRtDCpqHtZ/xEhcSCDDFKV+lAoZ/6h+
Ojmq5A8CZdydQX1VF+55xzze0o+PgLcMJzJinZd/qtTgQ1t/EJLQ2M6s+5l8ydEUloxQydZbAG4Sy8oY
rr2CJB5ME3Zp6s3a/yN5awOxflZS8/ig1+SwlxVfD29jftQSBjf6qfzzGxeG07vOwWZcMckt1aVXhuK4
8+ztkM2koD29+QyDPwfHMIb7PSIyYrNHMxgO2lnFeGfBt+Zr5M5CIkbpuyTjI8/z/g0AAP//TebINDcJ
AAA=
`,
	},

	"/definitions/bithdtv.yml": {
		local:   "definitions/bithdtv.yml",
		size:    1567,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5xUz0/bShC++68YhUOwHs4P3gNeV+oByoGq6oUihIQ4rL0Te4S9a3bHCSnif6+8XjtO
JZSql0TzzTff/PQmSRIBOGIUkBIXitcRgJYVCrj6epfcXN/dRwAKXWapZjJawI2pEMwKbigv4BpXpKl1
gGeWUueNzFEA6qRxLUL62YkIACCBgrl2Yj7fbDazlDhpE84yU82jCCCTdSBmkjE3ljDYAEsBcHc/v9RU
YYBOBcB3syZ086uyuZXbgP/XUa9N1lSoWdrecdY5ftTGcoDOBcBlo8gE+2LQDMD/PWF+Twp72nIhuna9
sRTw8PDQW6etC47AoXRGQy2zZxd5Z2XUriGH0maFgMeXp4DwOtmBJyH+BLB+asNLk5PugmvJhYC5R2Z1
UXtwZWwl/K83SdcND8kah7Zb6uTtDWZfjF5RPutReH+fBGItndsYq/aJPToQ0Vpje/GkL4jlM+4X5btG
59qDGIC29xIzNlYAy7TEmUKWVMKM8ZVHNIuVWaMAx9sSTyD1LkbHYqjW5622PmEEYCWT2RsSaYWvQz27
xJMjx5Ibl0oLK6NZZEazJO2Op7etyDSGf7xj0uqGvYyE2ViLmt0gvT/wfpF+it9w207PjeacSRaw8IY1
m1HY3mAeN6S4+Dy9OFtMn4BTo7bAVmjDx+0fF0lWUKmOl3Hc3QBhqQax8BFtd5MfyatR+GkMcuBIZktp
0z4IhcXVgK+oZLRuvMYkPBQvDdqtY0s6H3kBpM2daMvoD5y4xIPV/PtBNT46OLqDcX+tNepMmY0ujVQH
tZYx1AfVHP083OFFHP1+4qnty5F8WOA8Hu4F1d5SPgj41AeUiFnxJxHLRRxF0a8AAAD//7/+jMgfBgAA
`,
	},

	"/definitions/bitmetv.yml": {
		local:   "definitions/bitmetv.yml",
		size:    1616,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6RUTVPbMBC9+1fs6EQKNrgf06IZOJTeOr0xXEIOsrW2NciSK20SUsp/70iO7dAJJJ2e
PHp67+3qaeU0TRMArwg5FIpapFUCYESLHL4q+oG3dwmAFqZeiho5oEmXPiDKPHieAACk0BB1/Px8vV5n
W4/Muvo8Cc5IpEw9Unvn0toHhRECoE2HHAgfaQtoUaDmcNOTEoBSdFuDUhDW1incrgFyDrd3SVy0Vk64
R+HKhsP852Kos0on8CwQvDVngN0iyLWtlenFLVJj5YsulemWNHr3GxzY0xNkN9ZUqs56DJ6fWSQRehro
naCGx5D8npRa9F7U6LOu6UIjTpCyvfRtYeHs2mOU9QfWWJJ1HHwnTOZboXVlDcE1hM+8tNq6K5ZfXlx8
y9kCTiPc31FM5d9LvkxlCDem8h03a+ukHwMJ7FJLFJIDy3ssuE3qoX0ShcZ5iVp3Qkpl6iv2iS3gGqiw
chO+jjfCn5DMtAg5z6JFpVDL6Yr6QdkMa4gDxiEfhkGRxmlzp7rkhpq0bJSWJ+9ncA1iZAkip4pleCtR
v92QSEJp/x9ujcNqMLNro62QB90+HOPm1a/Dp/w8GxkOW7sKvwI3tCPoyJjikIV6Vyxni9cdAapKaUK3
E9j0Z/CdVrSDAwhXew5zYL/ZGeSw2CMip9o9GpaxcbhQvij4ykEuhyQ0Ytkco8gvBkmlNB7mfxzotRPF
YfqX6WreCM1hjY/dvgRO7u/l6Yz9NVsrq5ctViIWHGWl8Ljrzt6x6akCLLvXpf3bCtw/AQAA//9IOq3I
UAYAAA==
`,
	},

	"/definitions/blubits.yml": {
		local:   "definitions/blubits.yml",
		size:    3973,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/7RWf2/TyBb9H4nvcDWRUAt16qRJ2s4joEKEKr2H9BbYLlKdXU08t8mIicfMjJNmge++
Gse/sbOogn+a5p5zf5wzd+x4nvf4EYARFiksZLIQ1rhAxNZI4ZVMvDzC0YRaxFaoiAK5gusZWM3CT6iJ
gyWLlglbIgWMvCTNkCL6ZKj7D8CDlbWxoaeni6xmP1Tr08ePHByyOOeFzOJSaYF5AKAHb9VGoMm/D0Y0
i5y+ksk7toMevEmkTKfVbJfzxiXvegY9uJ55s5tZUWVSR9/hOrkvUsd1cDiY+HGR2QAH/kUFvPweFAV4
UQfPhy6x1DlTYbLGyDJdkTsc/Jjcs8tDgiYHBA3PDggajg4IGo4PC/pwA++xJuacwoebPblNxMgv8ZqA
4UUJ1Oe7rAPFbGeVUs25rmcfbgreuKvApLvAx48fC2svqfvaJWiSo7Wpx4M83Kz8NjEiLGhnFK4SLtTp
jeCouhb9vMmqWTcaN+H6KBdNeD9SZkKB/k8ZI9EYN4RkYTnzaJD69D5W2kIP0s/iwEfDFLyKxBqhB+ln
gY0o/P819OAPEXG1NXAVxyavu1a88gwwyHS4onD7eZ6H7MYroyeOYVR0AhjP9yWkWoooKxAzu6IgIo73
/XgVv4zZEqcpIWuGdqU4hVgZu4+IKE5s2T8RnAL58gX6r1V0J5b9xKB2z0j49o3kpHjbIMXMmK3SvCSh
1koXVT0wKDG0SlOwbCGxL4WxDF6AXSi+c5/a/eFF3MQsujV2J3FKQiWVpr03b3zf9/9DMmMsGls0aOgu
3cy7JrIXsY2rvjdNMytUu2n7YJlKylx4AVLQUEWWicgcBeRdWgYCcpzpvhPSojYV5fsXjImlsHkQgOml
oXBLgsT3mU9OYDBvZmiMJQvx+xzP88gJEJ9k558txwEtjTPOtyk9wd8S1Lv+f3Hnzs/UTjl9yVmlNUbl
nqv0xWgo+HmEhVZssAhota22yl3kYtNfeM46jJyPXWtAV8wcsduVxrs/p6SxyfkwT7K3525K5se58Sh5
2ZirbSQV47T0r5ylKJ+z0g6CT8m8Yre1WiwS93PBsYvbKKzEw1Xbh/Y4WiakedLdJ61dIs1dOrgd1Q25
EbiFrB0Ftyxlw9y5h0io+f6vTh0W8NntnbFaRMs2EXmjYmy1dj8YzC9yvjq3EX+3H3D8E46muPCkOkk1
eYn3cVsuCYKvQWCevXfjuX+OboNgGwQ8CPon82eQftkeOyQIvha3mDP7q9ZVRWuVGFQb1D/qTKe45yFG
FvWLK86R06P+05fHz0+zWKFlqdniQRvw10oYq/SuKebBx1h9BuezIK+Xqk6X3uwpeS1F+AlWqBGsgo27
pjGiNvllJXMa2ZUXroTkR4PjvJZEDFc/t/jwuPmk3CiZrPGOpVXLRiEzWPNHrJe3TNopWSrJyZw6F1px
I+QG9Z7RH9c45CmhQAZFLIkfNMHwHn5PM+FtIq2Ipcj6DdsnOuvin7XzR138UTt/3MUft/MnXfxJO/+8
i3/ezr/o4l+08y+7+Jft/IHflTDwW8/7nwAAAP//62RCfoUPAAA=
`,
	},

	"/definitions/cinemaz.yml": {
		local:   "definitions/cinemaz.yml",
		size:    1258,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6RTTYvbMBC9+1cMOZQurOx+0RbBbil7LD2Uhj20LEWRJvFQRXKl8Ybssv+9yLYcBTak
0BxM5s2bNzNvbCFEBRCJUYImh1v1UAE4tUUJN0P8owKwym16tUEJ6EQfE0Lud5QVAICAlrmLsmkmgZp9
U1UAWnUTRSvGjQ+EU5x+i+Wt+N76HSx9COh4IWF5e8h+9feERW6IY5HvI+ki/7k35If01ptDn4gq6FbC
zz93E8L34gBeJkL07hKwu0szW78hNxZ3ilsJjeq5bQZ4QNc+bOXwHEJyXc9zN9wqsr/6iGF0cPH4CPWN
d2va1BmFp6fFRO9UjDsfzDExozMRQ/AhtxAn5ho2xxjTlWYg7W9Rsw8SamUxMFxDJx23Qrdkzcs3FwOX
MbKchxrltfa94+TJZFdhyqLh0fj4idzVqxcj4yot8a3HsK+/4D6tEPMKwe+Kk+SRWK0s1sOzmOntBVwD
r7zZHw2awDCegNCaWW56t/aHrYsGppB4nSToCJhLFHOgVZ++Aia2mN+V9P95YagnC8SaLB5H6cxTkUFW
ZON/aRyN1wZcZ22/c9Yrc0488wRp785oRno4sXBh3PtsnFF8nv0u+R479Wzr0u2IaDCc8KoQ/JDbW0Td
/kvFxwuo/gYAAP//qka1Y+oEAAA=
`,
	},

	"/definitions/danishbits.yml": {
		local:   "definitions/danishbits.yml",
		size:    2588,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6yWz27jNhDG73qKgVMEMRpJazt/CeSwibcN0CZN4TRwsRsElDiRCMukStJ23EVep8/Q
+75YIUqipECpe6gOi+Xox0/ffBxZ8X3fA9DcIAFGBddpxI32AARdIoGpLV2WpYyKZEUTS/psUVS4WGji
AQD4kBqTaxKGjUwgVeJ5ADHNKyqmBhOpOFZrgDMCcCPXHHV4PS1LsAeX2cpXdAs/8GxZgRMHgrv2CoN6
gW3uiADcP0Dnqji4f6igYwLwccW47IFuVpov6oeOCfxiUlThDdex4/galcaKGY1ITwfTh2nb1fjs3yA4
+O2W3/8+rBXHBC6lXOjwUyTlooI/+dG3vxNUFXTem4dVu5QLrb/9Zepnj/uefT2Fl/HJUcdk71lc1x6b
YaiNjs/fUe5ryp5Lwzn2/qFP+oRAT/I/b1knhNGEwN1VeEPjtmix/JjnuoaOLHSXSoG+1bSQjHjWwY5J
z0y0p2F0UhHhA2co28S6KDS2TglcSaFlhi2pWy4MCiZhlvOsJj/0JXgrFePFXL+ZolMCcHcFnWuvqLTb
OLPd/kiXzWRYpv3Uc+cvvJtNaiaj25mhhkvRhsdHzuKssTjrGJu4d6Bt7J5GGRrLhe7VG38gPe9nc3vU
f/vtUE/ekYGD2SrShooY3ZROmm7nkXwpUfu/dpvHBObz+RvB+XzuWWApWfOjpZGqOCXw+Y/HqmLWflM8
LAAtxSFg/lhsV0Wo5eacmpRAyAXDlyBPc6/UyzA2UhEY7GlDjX6yOyAoFoNCIZMJFx0FW3EKz1Itif3X
LrnIV8bZXWlU5e/54OtXCK6keOZJUFfh9XVQgTnVeiMV64J11YGolFS1uN9vyGaGWhffC69JtOk02FAl
uEiIkOYgSDljKMoTM6gNcY6quCJZxuW59NthGKkUCqPd47sBwHeKbggMaFwM94XATaTkRuN+KXVRNPvr
CtU2+Am3RasaXl/3c4VPZpvjRa2+b1dlBsX+Zhzqpkwx8XsV/2RXYFRQFcqT4pgxt7X6IG6bjFpijAiT
+nHKM3YwGgJ1DDVG8WhVfLRThc+u/swzg0q3A/erL7nCBF9yrz3dVCWawKDc9HTw5Qv7fnghRT0MhpsM
d/oaDyGIlcyrFg2+mH6jVq66wdBQnun/T7yVApMbkUnK+sTpZ+viYjDlChcGHWz/khk87pDW/M/dgUyG
jlC4lGskEKlDxte1P2p2ixwPITB8ibuC1Iisc+Dv6J3WpjLEOP0vO86G3j8BAAD//7VXlH8cCgAA
`,
	},

	"/definitions/demonoid.yml": {
		local:   "definitions/demonoid.yml",
		size:    1974,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6xUTW/jNhC9+1cMjAK2kciOk+62SyAJ0hTooShatMVeAtegxLFEhOIo5Ched7H/vaBE
fXjhNDnsRcI8zQwf571RkiQTAK8ZBSgsyZJWEwArSxTw8wAo9JnTFWuyR7iRNq9ljgLQJrUPiLaPXkwA
ABIomCuxWu33+6UKFcsSV5MJQCarmJNJxpycxhgDrAXAb/Ss0UfgUgDc1UpTjK8EwN8fY/C9APjjfvWL
LPv8dw0Ug/cC4Hcu0MX4x6/iD02z1Z3VJXYELgT8RPTY9VuvY7y6p1JnPXwlWlqrj1phR279Q4c2z5To
cdJ8KkkNd/QoXVYIeHjaRISfkwE8Dwme7DlgtQnlhnJt2+JKciFgJbOMasvbQlpl0C2romqPQS5ICajI
cwNoW9XcH2x19tiKO/38GZb3ZHc6X9YeXUDhy5dpTKyk93ty6jixQ0eJDrl2NrKK2F91WmoW8R3Bvebi
ScC6CdE5ch2pBDwazJicgB258kGr61mabjNSuA3AbANnwCmpA+zI8kNmpPfXM4dqthGWiyQrtFHz9SI2
DGPwPtiyB+DbnxHuXtIzCvB8MHgOaZDKSdZ0JFWY7jYjy47MtpIWTa/WwEk2IoAu8+aDZHY6rcNasmbT
WnOnDaPzw9haJR3m+KnqaUmXewHTPxseMF+eLb6bBmLRXmNmO23i1hy75KlGd4jK/4qHoLgfSR539iDg
ovNznQ7gnTF9H2k0H0EeUaEScBlj/MTBe2ZAaq26xo72o5XpRrXMWKYGm5Gi5a2lbSUV3EADh3ej4g2w
E5Z4/iCNzu1m0QSNlEZ6jnom9uxq0Wkqd4yusyiAkowFSjWa+ZjHtG18PTO449lm2qd8pdP/aDXS6587
pVDB/FaQhcXtfHl2G4Rr+6FRfbt+0CcosTqyKtyAHHwxWKpwuHsL28YHnp22+QnKHZHuDxZ8+iqpyzGp
jMoSLZ+c7stVL11F0d4akurVbu/e0s3rf1+/zfvF2NanbXJU8CEcHf4vMdEgZsVbCtcXfeVk8l8AAAD/
/4SEpxS2BwAA
`,
	},

	"/definitions/eotforum.yml": {
		local:   "definitions/eotforum.yml",
		size:    4785,
		modtime: 1487335360,
		compressed: `
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
`,
	},

	"/definitions/ethor.yml": {
		local:   "definitions/ethor.yml",
		size:    2927,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/7SV3W7bNhTH7wv0HQ7kYUvaypYdO2kJZENr92PrimZ1lwZIUoAWjyQiNKmStBPX8APs
UfIMu9hFHmivMFCKKCt1umzDbmzp8Hf+54M81J+//xGG4f17AIZbJIA2U9q9SjpFAs/da1uiha33mdLf
GfiZSrbtAIYm1jy3XEkCwVN4oVHGGaSoJWoqwGoan6EOHCuoTGc0RQKJDpNCX3B5Zoh7AgghszY3pNPB
Kl7n/j23FtO8gmJqMVWaY2UA6PUIHAyhBU/zXPCYulxMZ8TnqE3FdDchB0P4cfy2Qvr9jYzSlk4EemqP
wBs152g6z8TsHV1AC15wMS1eNV34pCLPjQ5HHhodjrSX6nnkVU28GsE7nvu8n3hmXDPjBjPYhBwORz6O
q6vzkk7RQAt+wtlF52DoF7sEhkoaJdAv3ih5p++RzsG477FxLzwY73jqcU194LyiPnAevqyjRTV0NFEX
FeWed3ajCtsl8HTGuIIWvJkZ/mlW57J3vdQ55AzXgPLdx3lM4P1h2dbx1aXmGL6/uhRXlzd2qT8ouPGX
XGOjdm5Ta27W3m1izf3q7RB4ptSZ243nYfHkI+0SeGsz1NCCq9/mV5dtMLnSlie+AWvEa84+e9FBbX/D
TVzZnxA4OjqClvsthwlgqtja+BikOs4IHH86rUx2HtbWR44wSj4CzE9LCYPWcpnWg1teE7FSZ9xvlV3k
SMDiha0sgk5QuANQYsX8q5TLa5mpG3vWlOEyn9k613KJQLBcQnuoZMLTdmmD1SooKYvGeoec2ozARKtz
g+08y8uo2g33NfMlUTZFYGyVJmByKluFwzv3UzWg6M7XFG5k/o2m50XemsoUoT3099hqFS+X7dVqv/vt
comS+ULqrSnK/WWGetF+jYtzpZmBNYyyeUwtgSDyJi5jwZAyAkG3liu3JJh4C8ubXuZXg3ooeHy2ZneF
rZ2WqjE5PITinoDvwU4UW7h/TTJqtuhxpjH5uB90GFrKhXFtCU63S5GEo2C1IlPnUijKvGE9yC1KRZTJ
du1BrdV8MnPfLudQLyRcWNRmTbw+rxpzQWNcXwKgOjUEjpsBH0HQqfIsM/CnnFuB/yH16+/Z4m8k6tP1
Q0ztfh3/X5b+yR0mYzWX6abyY2rrkZtOUVrzf2xPwgV+TfjBftByjODG1iUb/nlzxy0j0mZhnHHBtnZ9
+FTTyeYoDYe97buemhQv8k1dC7ZOTtjD7bXxRdZUuSX0Yx9aIMbZnXyeeB9G7R36MbhjeTTPUbKN5UEY
DUgUBdACy6f4WUkElSQG7SYhl1ZOtdk4X0EvinbDqBtGve6ARH0SDSCM9pz6zXthrsRsigktaqq1Ymqw
UUDAp+lxMY4f90+CFxqxaCaBk+A0aFxzX9KvqEhgg0t70HR6EDQu1Fn+jxL03n8FAAD//2p0TkZvCwAA
`,
	},

	"/definitions/eztv.yml": {
		local:   "definitions/eztv.yml",
		size:    968,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5SST4+bMBDF7/4UIy5bpABN/6my1GN76alSuoeuVqsJHrBVY1N7SJpFfPeKhJDtphXa
o3m/eTO8mSzLBEA0TBLokXcCwGFDEj7/2NwKAIuu7rAeVZd1UQCQK70yrpbwffMl+zgyxv2MUgAAZKCZ
2yiLYjTLsS6EACixnfQSmWofDE1vgLWEza04PhqvLt95l0XCUGoJd79WEAmjdyug9n6EJ+nItshaQtL3
poL8W0fhkH+lw94HFYfhBBZ9fy31PdlIw1D0PTk1DMnRLfj9PEMkSyX7ICFh3FrKKx+65kETKgoPWx8U
BeBwNyb26Ub7HYWb+39ByWRYGcsUZv8xr1Pa6FSDXGpxosiqGZoiO1yKmH6zhPU5KMOWLuJlZlbSsc5K
bax69SYFnBhFjMbGl5QAIHMw2248Ex2oOlv5vbMe1aLX2xQwb7B2xKtrZaFLNI/LP/guPc+EvEy/T2fi
aidPttK25NQTAQBDHSUkgLVP5ish9ZfBf3p+SJ+ltvO2a6jCI/tsu8nrs33XLrLrRPwJAAD//+5SQ9PI
AwAA
`,
	},

	"/definitions/filelist.yml": {
		local:   "definitions/filelist.yml",
		size:    1819,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6yUX2vjOBTF3/MpLtmHbqGJ/zRtt4J96LZ0CzMDw3QohdIHxb6xRWTJla4TMqXffZAt
Oc7QPxOYl4DOPfenoyvFk8lkBGAFIYOFkCiFpRGA4hUyuBYSP3eC5KpoeIEMUE0a6xShlpaNAAAmUBLV
LIoCYWp0NBoBZLz2lowTFtoI9GuAdMYA4PtddKFEhV5MEideNLnQwXbilC96JdBGx1dBjQfqf7L5xjeh
4gqhcq0NikL50vGwdHUXWLOhfBPU5Jy9xUmGDbd9QxvpUiurZThOy4Cvl9H/vEIbjGesFUPgdGf5j28J
1cRP6eYK/gKL3GoFNc+WgZYee8PtG4akN9TaUBD98P2yTQT39/ejVqh0vr0ni9xkJYOHp0ev0GqyFY/8
nkeA9aNrl7oQqmuuOZUMolaZ1mXdigttKtb+PvCMhFb/HhBfYm866LYRqm6oD9FYNN2bHD8/w/RSq4Uo
pkGFl5exN9bc2rU2+a4xqANjo6TOlgzGydilNpyE3kldbfrIFiVmpA2DqSVOjZ1zA7lYMUXlJCuFzP9O
D38RZodga64c2w9rAJ8bvbbY83fPGmbbHuATblxwO0iecWIQtwuH2bb1IUkbg4qMXncTFyjz3uf/iZuw
fq2T+Fzi4DBJdxjgfQ8nMmLeuM9GaXDR6wshCY3dwt3Xobu5pwbNxpIRqhhUAbgpLHOxwvMSJHGfdOl7
6VqaL+RIXEj7x9iDk+d6raTm+T7s2W+yrfix1zzOPHehFfVtBiu9Qvf0QmJOe1FPDz8Avnv5breaG4uv
XH1ywuIZi0/iNIqTKI3j0/5JY77D+zjluf/bdX6JmJV7EpLYI34GAAD//wd031MbBwAA
`,
	},

	"/definitions/freshon.yml": {
		local:   "definitions/freshon.yml",
		size:    1347,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6yTQW8TMRCF7/srRttLK3U3irj5AAcEFw4IUfWCUDVZT9YWXs/imbSEqv8d2dndJKjQ
A1wS+fnN8/jzbNM0FYB4JQPbROI4VgARBzLwPq8/xpvbCsCSdMmP6jkaqG9u4TMlTwIY9xzpTV0BBIz9
DnsyQLHZSVZ8/CamAgBowKmOYlar6ZRW71dVBdDhOFk6VOo5px7WAGsDN7dVWQxsj7reN0KYOmfgy/dr
EELheA00fs3mwL2PB+uI6gysitKObiziltNgym9Z+jjudIneCaXD7evHR2jfctz6vp1VeHqqJ+OIIg+c
7LlxVhcjpcRpDm9AKFCnnAy4tek4Kvool/W74qqvJh/AQCIZ5SLASWl9US4EY32ynWjgezIgug90DZuy
pSRqloYLizHx1gcqNCqAhOr5DJbjgRZWp4eKogqMJqprOueDvXx1BTJiPFHWV3VOnV7nNHaT+EGOwefU
59csKD/tKO3bD7TPIGUhmeuP/rktxU2gC+d7F3zvFF6Dbtju838ykfWy7Tg4Qntgu/UU7BIzTdz+iFnp
hxpYz3PmNZy8wfFUbJVToqh3eSru8pwvLlRNfrPLH1SpnzYsKfog/5DmEm3nxnkYKOp/SrP8EAOjfS5N
bVsQlxwBPHvrF3LF/3yW3pKZDdXvA7xJc1uofy9Ha8n+uV6ILKVnGR07OHgAyyBPzkDUuZcqZxPgrwAA
AP//ODev5kMFAAA=
`,
	},

	"/definitions/funfile.yml": {
		local:   "definitions/funfile.yml",
		size:    1844,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/4xUTW/cNhC9G/B/GGiLIkYiKWsHSUvADlyj6aE16gKGL8a24JKzErEUyZAjC9vF/veC
+nYRO3sRqTdvHmceR0rT9PQEIChCBpvabJTGCBheIYMvtfnSAxKD8MqRsoZBcg0FGvRcA3kutuiTyNHc
FDUvkAGatA4tpMw2sLgDSKEkcoHledM0WX9YZn2Rn55EhuBuoApOWFivcAAAPnxgcP+QXxtVISygXYfY
+TmDu5uIOqeV4LHKMCZeMLiupbJ5+1xbu43MuIdfrN2OxPNPrANgAb9G2qjAAO5u8t94hTHWrmPsPYM/
qUSf36ogYAFxQa25QVuPrOXPDG7tk2rzu80Q+sigL2YBt3VQYghcLGNT+V1pDabtEbCAO+uJr/XY+ScG
cP8AC7h/6EwEqKyc2RaQe1EyePy6GiB6Sif0XWQEa94BulUnoW2hTC/gOJUMcuJbbOHMla4/Bqm0koGz
gTpkY33F2mf3royraaqkDui7sUr2e8hurNmoIhtQOBySgel4CI318jlzQOfMrlRI/ohrj6L31o+nphBQ
oyDrGVSbf4Q1hKYvmDAQmw5tO612XYsR9nGQnhux9rYJOLkwiUv1xEoe3kj1tBDaiu0ZXMEa3sL6Lazh
CoLjppPtzX9V93/e/eB507rhuSkQspvx+zgcxH6fHQ6Xyx/3ezRyZs5wya2Jf9Xod9nvuIsWhrmHyggt
kUsGyzGztE1wvGqhDoz1zYZq6Lodxsc48I5LqUxxmZwnK7gCWlu5i6tvXSGZedtcnPWTolDLSY4UaRzf
5vL8sY2tHkuPm78vE4nElQ7Rp2Q1JXAir9Z1/Ie1/CHS/0Z2L2j3opP3nwWnyxeEI3kKbJQm9GEmHEet
G++v0epAXpliHgbgvggsFjXWZ6sKDYXX65s1/VnJY+qTtjHacvkd3Z71iptz0aD+/fYlkWSGylSUSss3
P51NQ4LyuUUvpCzfjzkaUZTHJS3HJMnpiMI+nh13e9w5NPJbF5cAL2xyevJfAAAA//+hsze7NAcAAA==
`,
	},

	"/definitions/hdarea.yml": {
		local:   "definitions/hdarea.yml",
		size:    1932,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/7RUT2/bNhS/61M8RJcGqyS7SdpNQFckUwYDS9AtMZTDXBS0+CwRoUmNfKrjFf3ug0SJ
stoG3g478vdPPz5SjKIoALCCMIWKM4MsAFBsiykssku3lEyVDSsxBVRRY1tEqEebBgAAEVREdZoku90u
dglxoZMgAChY3YsKRlhqI7BfA5zPXqUAyxwgXOZwjy3nqbNDqtI7z8xnKcCt/iTQJosMIPy1kXKRwYRf
5t/n5l95F9mEGmwH8Nk3jmU+IUfPSMym3wmvZHPH9hPW2b5i5ueOuc+6dZjl2YTqI7P8O/zFyLtQCPvl
Moc73DZPXvp6lJ6NAe4wkofrqyi7cWD4gOuIy4mkNz5cX2U33wjeuIz3y8X1HYTid8YnVO99TxWaKT3r
d57potmiImbGcV046lKJLXqw3cNlw4VOcsFRQ3jbWFFAtxhvUd/nvtaGPPhTCtBVSG6FLTz8o0+80dZK
tDbouK3m4521yExRpfDnXx96hD5FI/iyFVitXgLWH1q71KVQzlwzqlJIiD1ih8Z1VbsPIFWap1Br60oK
VTfkP9lYNO5vPPn8GeJftNqIMh5Q+PLlpBfWzNqdNnwqHFAvRGO0GcKj52u5ata2f70H2glILEibFIjH
hE+uMqGl1PfoEnuv7QIDP7qDWZwkpI1BRZ3mnRO8bbv/0aDZx7/hvm1uh+ZG7w4Owtdga4nxkAS01nwP
ZGKllSVRPH7UDaH5uC4750ag5D6lf5T24wYPt5cqqqKiEpK/mJ/Cz8AmgLcwIiPWTfd8Gtx4fCMkobGH
w4v6d9VgiU/1AQHATGlTOFmt3hWM3r5YrfgPp8PJkiCJRzu+Op3Oorse/Tg4MFj3ARyJCWn/hz1zvVNS
M340+vw/R1vx9/EJXAwpnNFxdVfC1kwd7dHN31885JNTfSb99ZAkEYvq3zjenELwTwAAAP//y+1gwIwH
AAA=
`,
	},

	"/definitions/hdme.yml": {
		local:   "definitions/hdme.yml",
		size:    2540,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6RVf0/jRhD9359iZNqKqMRJIPxaCSoOrqXq5UrDgSJRdFp7B3uLvevbXRNyUb575bXX
TqivIB1/EHn2zcyb92btfr/vAWhukEDCMvQABM2QwOXF5L0HkFIRFzRGAij6hS4jXDxq4gEA9CExJtdk
MChTAyw8DyCieX0cUYOxVBzrZ4DdMYFPt4MzwTOsIrAF9skB9glcnQ+GjC6gBeT51/p8fEhgIp846sHl
RXN+e3554QocNIB3aTEty2zBu7RQdFEj9scdJVg4fT+5mbkuow7IhYyKDIWh5UCu1rAD+Kss1G8NoXFX
u1pcS/iIwFnBuBx8kFqnqHV5/hdMCs0jV+OoowY/p4K5ocb7HYjfP5oPrssxgT9Nggqg1XXCddPh2Op+
lUiB/Qq4BRMZ8tQ5s9c1qg3AaHg05A7WJV0Lyx1s95uww90Wtee0ueUMZY0qhQEbcDbs2bWC9b8tuKLR
Y+PTXke7KWbFswMcdAFuZm6s/S59r/lHfm1w6lgc2OW+zqUyDQv71PA8tAhXoUJ8uoVrLHdqU8ejV6GN
Ssf/C13Tc7+5PWvIm8sL2B0dNKDxkMBsNtsUczabefY8k6y9zhqpihICd1/u64h56rfBnRKgpdgBzO/L
9FTGXFTJOTUJgYGhj2ijQZ5UBDI0iWQEcqmNDTxIlRH73z5ykRemYVBoVNX7yl8uITiX4oHHgYvCauXX
wJxqPZeKbQJdtAGiUlK54n3QmGJkpCJgWIBZiIwh85wsGWpdvhq9VqmNBIPP1QgGtSENETt4trATewCK
Gi43RKnPNuvpnIpAZzRNH6QwcArlTyUQTw0q3bKuBFEY43PecKMq1gT8qe22Hfz8S+8mTyVlyPz/JOYp
jfBF5p2/4++A71sja4/XSYdKzjU2xDdt+kHRuVVeUREjBOfNl2G1ipbLYLU6Gf20XKJgrWNuj6xff+Ci
9EmvOcpFlDKkpaMjFwtTWmYMq+eSUbuqjTE0TPFuzplJTvzRcPijfw+nYELJFuWvIgnV24YFoTRGZnch
jR5jJQvBTvzPPKMx6kEYfzZSKRQm+CeP/fte7QOmrGlYf/wW7Xas7wYRJulHCU/Z9m4PaKu1MYqHhf0a
K3xo4i88XrfrS4FqoY3iIvbWL21lW0SNu5rcpPgqm70enHbzsfluOJmVn0L9HeXWxmNyLspVfLXaaPSW
cpp/fX3Og16DUJjJJyQQKseHmrcI9RZzXlzC1hf/rHyVELA38W/RLj2yjVLf6H7kuqeIUfKWjOOe928A
AAD//w1IUlzsCQAA
`,
	},

	"/definitions/hdspace.yml": {
		local:   "definitions/hdspace.yml",
		size:    2419,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/7RVW08jRxN9R+I/lIbvQ6AwvoABb0feaMEPljZRsixBSIhIzXR5prXt7tnuGoxjOb89
mlvPmNiOX/LC5fSpU1WnqnvCMDw8AHCSkEEiXMojzAHNZ8hgMg6/1ohAF1mZkjSaQcBhMgayPPqGNsiP
FddxxmNkgDrMXAFJ/c2x/C+AEBKi1LFuNxFhkaVjbNw9PMiPI57WvIgTxsZKrAGA/iWDX8yrRNe9Udkd
X+TgEdyoLLzji5o16HnWZFxCcAR3OMvevNDVJspkHI4fxp4z3MQpEbg+76We+GEHsd8b9mrieZ/B/UND
qrPeP6zpnZ9vpa2pDQra2ETZDDVxu8hpYxOti11uY61pXResT1rOsEmZ/8vzKa8rDndz27oXPQafMiFN
92fjnELnila+lKBn9WvWgxRoagMzJyMoEOeZFwweHx8B1qzJkXaFF4NtrLXarhj8SgnaNda95VKhbTJe
M/jtFt5pfTVTmnOLnjXcpFUgfiv7m7Zk8Bl+n4zL3QeYGdHadofcRgmDp+/PNUSvYYOe5Qxn9Blg+lxK
2HwElUDKKWHQlVrgWydN0p9SHuOIjLWoqerPocKIjGVA/EVhR0lHHD4CvRixyH/b/Id44krGehREqAlt
8Mwio4lL7U6CuyJjcFrqTaUitE0LYfV6WIzxzc8HgNvYMaiC4aTzw2nwz5BUVe9NO+YpCMMwOIOgF1Q9
KxNLvbPnglE5jJQYwSA1jkpE6jSjpuJMCgbBcgmdW6OnMu5kDm1eEaxWvsZ0/o6UcufmxoqGROjIi+4z
ivYwhHw9mqHOygarie8z1WMekXzFUe/YFM+zG/WOq3d0Mepf/v/iZtDb2Pb/LJ8zCDx3ubRcxwidW/8K
r1bLZWe1+nG5RC1aZtQLWfjxJUO76HzGRe6Ga3tWVsbA37+qwgIpMWvm7e3fvZqOFgrhLyBbrx4q0YST
JIWs2Z6WnF9nhVMKnuEj8KfE4vSPUbDZ1lAgcalc4O9h7dPGBP+i5pqJNIIAnMjKl6z4+FqcNgfv71T7
knzP7XZkpY7bx/VlqRP5ss0s/wi4/8yXrW0IM9fKcLHbsJqVZ9pD1Mk/t42YaUrCKJFKnFydNiuFYt3J
LSFDH6IQo2SvmA8+RnDao6zL08ODvwMAAP//LfdHz3MJAAA=
`,
	},

	"/definitions/hdtorrents.yml": {
		local:   "definitions/hdtorrents.yml",
		size:    1871,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5RUa2vjRhT9rl9xcUs3YlcPP5MObGG7ISyUUroNIRDyYay5lobIM8rMlV3X+L8XyRo9
nARnvxjO0bmvc+c6CAIPwEpCBpkgbQwqsh6A4mtk8O06uO24nKu05CkyQBWUNSPVk2UeAEAAGVFhWRRl
InCJQm3SyPMAEl40uoQTptpIbDDAmAH8qTcSbfR7Xn7nu4aftHxDTDvht+uGm3XcP46bv6JbzNhLbtpy
f1GGxoX/yuD2rieLa9yXTOOhZHpV47aDxXz4fTZj8KUUUrvv4xM8GeL55RDP5if6xYn+isH9/X10fddW
PDIOXfbRokE1XGvRrcIiN0nG4OH5sWFoE3Tkp0pgtfoEWDxW4blOpToGF5wyBlHNhEVW1ORKmzWrf2so
VVFSW6yUgsFov4fwq1YrmYalRVO9OzgcRo2m2J5oCm7tVhvRatAYbVzKACzmmJA2DIgvcwyTZa6TJwsr
regh0bk2nz/8dHMTx3H84TgjoSUX3wyxlTbLpaV6Dg/AcJJ6MKZ3dMvVChXflBZIMEVZkGQyFxdT/2iB
zAmN7To8XpbBFP8tGhKAm9QyGH2vC8FF+NG/+fL3qKrduN8r3h6Xc3lo68+Gb2vPDFcpQvi1PbjDoTm+
3cPj5/0+PBx+2e9Ric5ut+ra8T9wVzlte+vgCckNMhjX2Oht7+UMfV9zqWrrE60IFREBmZ476uPUZxm3
FyReaP0mp8G1rqqRccpqKdw6YzEXbX03GWst7XXUX8vYh9+AD4huC0RGLsv6z9DgquVPdtjf43OJZmfJ
SJX2vrqFuq7cMUnK8WyHUx+WwBuVQOIytz8W9NYkQm9Vrrk4m23un81l5X/nR7ly3gpO59WX/nscr1IV
3Fh8xe/xnMUzFs8B4kkUj6NJHC/aB4pikPGt9xH7sGxEOWKSvStoXAV53v8BAAD///r7SUlPBwAA
`,
	},

	"/definitions/immortalseed.yml": {
		local:   "definitions/immortalseed.yml",
		size:    2354,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/4yVXW/bNhfH7wv0OxwoNzGeyE6c14dDMmRxsQxbt24JigG2EdDiiUSYJlnyyKlb9LsP
1Ltiu1suFOnkd/76nxcxcRy/fQPgJSEDuVoZR1x5RBGimq+QwS+vogJ94qQlaTSDiEOKGh1XQI4nS3RR
YBTXac5TZIA6zn0RknrpWbgDiCEjsp6NRt1XDlc4evsmEAm3NZpwwtQ4iXUA4JTBH5Shg87PAdwqBb/n
y9JjgY0ZPH4c3Wq5wg4WHmtkfMrgwx3AKyVrfSNyzuA2F9KMiuvCmGUg6vuWO2HwUwj0lO4yqYRD7Ud3
3JExuuHH44ofvSs0S/5dT3N8HtyNfuYr9K1m8VgjJxcM3pu1RD+6nzRIGYnvJw122WAPW9hv5gUm+Nyw
Z99hHxrJ0+OqMb2C3+deJjVydr49qIMyUiOXxYgerHHUQYpnD4/rRqniJj2lx48Qw9nVsa2xq4K6n8A2
di/TLFQptQyL2whf7RF+mEAMn8cXZzX5/22wQ/69lk1nzva5eEDujYYPPFn6YKnJuNir3c8I3S9zVkZ0
PgiP3CUZg+mneR2iddxGjwLhjT4CtPNSwnGSphKwnDIGo4UzLx6HNqv66VFhQsYxEHJ9QMbCTbgbWi6E
1CncgLdcM01ZnIQ9PxwPysRnqQhday+uDhKHKX5uhgXAXeoZRH8VVuBwOpuJ2Wx4NP/fICpNKpNK3TdJ
fIlFuPW5QsqMYGCNpzIitc2pfX/u0ZUOoq9fYXhn9LNMh3UUvn2LatJy71+ME32yjrYkoSfWJhXOgl5i
W1uvG5ghF+iqHn5BmyjufVlnNak90/hRmOuS2FndEjfBnK88/5mj2wx/rYLd4qROVC7wSSAXT2ScQ02e
waY9TaqzdsPguL9aT7SxyICeQsPqHQwG2XatxBcKD3w41RcKiwvcAC2M2ITfjmXcHwq5nhYNuI7IGEXS
xonRhJqiebNFqET7ApKkkLXL02vuXq2y3Z0dPRlsFbtLk8Qw13UNY7gB3llbIicXefiHmbn23Nze++7u
fwpj8eSkTrt/rj+C2kxjzqxWxXh2mePTohnX0V1FRfN/dSfMi1aGiz3Vdjp0Ogjl7visv6Pu5Zfds+kp
nw/abUHR79SelMsmRSEm2X/KuWpyBKfdtqJegh5UaxI2M5xpU08bhdez6FkZTgycTDP6YRbNw8H0TwAA
AP//rmWJjjIJAAA=
`,
	},

	"/definitions/iptorrents.yml": {
		local:   "definitions/iptorrents.yml",
		size:    2414,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5RWf2+jRhD9359ixJ2quBdjm/MPWKlXXV1VrdqTUjWqIl1Saw1jWAV2ud1xEtf1d68A
L2CcnHP/JJo37735AYw8GAx6AEYQMhA5Ka1RkukBSJ4hg9+urhsoQhNqkZNQkoHzZgxXWjxwQrjWPLxH
7fQAUi7jDY+RAcrBplClQt4b1gMAGEBClBs2HDal3FBlw17RAxIJGdfUqoFQqXuBJQRA2xwZED7RAUj5
ClMGi4rUAwh5fjAIOWGstMBDDDD3GHxSDwKNBeYWGP718wHzg1MsGD2DzU6xGYMTbOLXvF8tNp10Gpl5
p2bvG+EvSqOIpSX7HbU3Oi0xf6aV+XsG139b0awImuR0ehzP/ePYK8VNAW/SyXf0s46/73Xy01Yzfruz
edBx9o4rd/MTBseTVHHND4KOfsrg4yYSqrWWduyPDvHwD2VMisbuOegQZ6PS+KMUmX1BxwzgamHzAYOr
xfATDy0wsc7l35VS9/ZhTxn8pNR9XWpyiIcLlYnQwv6Ywc3NjY0YtKPpUc5vR8GxblJFZZipqPlCDHId
Jgw+f7mzX9zDoAEvC4JR8hIwvyvkqYqFrMQZUqKio+9VyHxDtXeVYODsduAulFyL2K0w2O+dkkRoyNJz
TgmDoT0Kbp7kRUXNSaiKc2AIGeGTyKlkVEOkGJLSDNxwWfKr61JO0VI6Q3Ke6fOt5o9ll5rLGMFd1Hdk
v9/t3P3+u90OZbTff/mhmOTPDeqt+ztuH5WOTD2KVo+trdqGHOKrFN/Y0wcfgFYq2hb/NZOUDMJEpNGF
fOf1WcLNBUUuLcsb169s1wLTqFlp1drWxu1SFLUMx334ALxmcSItVpvi5Cca1zW+FimhNo1dc4Y1xviU
txIAXMeGgfPP7e2PF7e30bv+W8e+NIJSPNuU124qQuIiNd8memmSUGVZsd+zbtPXuEXqUaaKR2fdJq9x
M+Lf86uZ9W1tTq9apEvLkJoz9NVHafJU0DNP8jM4/zmXMBjD3bepYLUF5xJGtc4gRkfFX2jct2OmiGHy
GkVgFbHmq/P0eb/zDB9UuslwzUturQ65wfamTM6lS0vi8XKtEZdlewyckdMiOd87DJyxhTb5y/7FD5aK
+38AAAD//4OKp+RuCQAA
`,
	},

	"/definitions/kickasstorrent.yml": {
		local:   "definitions/kickasstorrent.yml",
		size:    1299,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6xT3WrbTBC911MMuvoCkfOlP1AWEkgD7UUa2kJaAsYtk92xNGi9q+yMbKch715kS4pb
GtxAb7TM2XNmzs6MiqLIAISVDNRsaxTRmBIFzQACLsjABdv6TORqhD2GssWSDFAoWskAKNjoOJQGvly9
K950HA61mAwAoIBKtRFzdFSjTqq4yjIAi01/bVGpjImpjwEWcdlFcLk5e1CXBq6+DoxW2Bo4ax3HHrqJ
sRYDb7ujh97joktzHoNETz2ITePZonIMYuDTeQ+v12sD19fXfRi1omTgY3dk2daUe3QohMlWBqa3s9Fe
8QgedgSJ4RComXXy/mrDbVArA/kWmjRVk29gDk2rY4FbA/n9PUw+t5TuJhd0t4rJCTw8bMkprna8eLIa
kwHFG09T61HkJHeomM9A05Td1uScybtR1rf9boh3E0mDYcru20luUb/nMzgF0RRDCaeAIx9VE9+03eZU
ieYjPmevlOQxcbcC21XSxIsdGABTKQaOhiayevqTIXUmaFXYir377/gATsHxcvzi8GZL3l8ihw8c6nyY
jCNF9vKv0z71fhdXwUd0z6mH025axSA9yfdWEf6xv1EvDgZPqPvZLw/GhSL3y/yeELwaBJ7IVn+jeD0a
IrGJm+4nfP5cuu38rdnL6NsFzXEjHxMqrdVA/n/eI22zl3ucZ9nPAAAA//90rEmMEwUAAA==
`,
	},

	"/definitions/morethantv.yml": {
		local:   "definitions/morethantv.yml",
		size:    1472,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6RUzW6cPBTd8xRXfJsvUmAy6b+3XVZVN1E3VTS6A3fAqrHp9QWaRtn3KfpwfZIKY5iZ
NhGVmkWEj8+P78FMlmUJgNdCChrHJDVa6RMAiw0peO+Ybmq0+c3HBMCgrTqsSAHZrPMjou1nrxIAgAxq
kdarzWYYhnz2yqXfJAlAgW3kFShUOdYU1wDbMafX5OP6WkGIG/+eKfggNXES1o0rjzJPyEWt4NOX24hI
nx3By5Hgnb0Eam9HuXGVtpO4RakVbAKSt3UbwIPjRoX/OXZS78ansKFt28kS23niqZz0/h7yt84edJXP
KDw8pJHYoveD4/KcOKMLkZgdz+bZ40cLs5P3Y/cLMDZgqBDHCtLjkSEfkK22FRRkhXhKEfKiloOFCHHt
9ipEJACMop2K1K+iIP35/Uc67sRGkzMlM1nxy/nOG5oUXjhO/o7uxon9MjK74YQ8zyC4N/Rf9N6FFQjn
EZhekSZTLtJ4k+6OlZyYlXmB4neFM1DqfmEU6M8qTCdaE+5fqmD75570qYLr+YppMfRE4F5XO20PDvKK
XddOz6is1FlRa1P+f30RhSUJauP/2QcARVjvu/HzrZkOs78brHFYrgX4Fu2Z83bN2etvT4x/4vJymRNl
nf3iAnLRDT2WHOpebgqVxE+UduL3ek43REX9N4o3s+KgDa3Tn8/0inG/Tn918dtb6Z3pGjpg4C7q+Nld
zT8fXbvK3abJrwAAAP//xivm6sAFAAA=
`,
	},

	"/definitions/myspleen.yml": {
		local:   "definitions/myspleen.yml",
		size:    2245,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6xUUW/bNhB+L9D/cFCGIV4rJU7XLSXQAkOAvnTtsCXrHtoUOItniQhFquQpmmf4vw+U
RElO4Szr9mTr+7777nh3ZJqmjx8BeMUkoNr4WhOZgBisSMDbzeWIaDRFgwUJIJM2PkBkciuVKQT8fvU6
Pe9Uytx4Ef4BpFAy116cnLRtm0X3zLri8aOgyLGO0hyZCusURQDg2VLA1Xs4gp9koxkuW1WN1GmkjKqQ
lTWROXs+MBfo2FoD74hb627GyEjbiuRmjPphD4YLMuxQj/T3A/328urZmxE9j+jV+xF7MWDvVD7mPDsV
8AuX5OCo/x2LORvUl4wO/kDn+74AVFbOOuEJXV4K+PD5OkJ8m07o06Dw1jwFqq97C20LZQaDGrkUcMJ4
Qx2c1WU9pCEurRRQW889okzd8JS58eT6TUi2W8gurFmrIoso7HZJVNbofWud3FdGdK50xI0zbAUkJwNI
zlk3Jk3Bk6acrRMg1e1Rbg2TYVGiPy7PRPhEZfxx8nM4DbxGpUkmi0Ufz+RZTFV1R18523qazg2zDL5G
k93QZub7W1ipZAFPevIWdUN9W11g9tt61/srnYdp3mv9j+P6xmHbDcChKQiyi/FW7Xb5dpvtdi+X3263
ZORsHnGPurn92pDbZG9oE6bm52NTJteSMAx4OYKsWIflOB2QULD4sseMK01HFSqTsnWODHt4BbyychN+
hwuxVqTlFN57x689O5mxdWm3gq8AJwkyO7VqwlvWRUdmeF02B9w+5Bq9//QyyZHT5PqgZ+loPRFrpZmc
n3mG1e1vy+fQRs9OmWJOA6ArvAj1RFgSo9L+q88ZaurB0dG2RluU91kGzYOOuVaaDhUnDJdpXiotj58v
xqGrvw4NbaY/H/WFw9UDErxY3Nf42HZHBf1Zzzve9zs5/vhRPlnMVp7kvsOBtMvT/5b3053EmigvH5Z5
+f9mlsgPmMuPi7tbdGt1U9EaO/EUn6Onvc3v3jPP6GavwYxZa1YVfcEl3yV770lT/5ucMfrvAAAA//8X
ac+CxQgAAA==
`,
	},

	"/definitions/ncore.yml": {
		local:   "definitions/ncore.yml",
		size:    2107,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/7RUTW/jNhC961cM3MJ7sRQ0WRQoAS1QZFEUaLcf6GIvQWDQ0lgiTJFaciQnUfXfC0oU
TWft7aW9BObMe2/evFBM0zQBsIKQgSq0wQRA8QYZ/HY/nyRXVccrZIAq7ayrCHWwLAEASKEmai27uZnI
WVEkCUDBW98vOGGljUB/BnjqRbmtO8Wm0wfdC7Q3P2mDolIRZG5HkL/e+27ZOz47754LlH3EvyTww784
cJDrDupogSsCdWxggfz8PlrQovEqHz9d2N+i8QofP51ZP/G+JM7tZXJMrGPel8Q65rl2sNq0d9GyP3al
0K+4TXsX7TojPvxx57tSWyvR2lnkEn9BsIj/q695SCFFy16NCJAJ0+jydMUsclPUDB4+P/oK9empuHEA
q9UGsH10dKkroWZyy6lmcDNVsrZup+Jem4ZNf6ejUG1HYZjCnsFqGCC712ovqqyzaNwXBOO48piWu/Vi
kKsctSkDCI3RZtFML/uYFkVr3ceYnCK3KLEgbRisvqnFjncvqJCAyoeeS1Gp/E0jylLim8d5FKElFpxN
c1qj90LiNCkJAcaJkDYGFdlg5jyFbw0/Mvg7uBoGw1WFkN2HB2AcD6LnktsX0kRbEm1nHx7zYcjGcT0M
qMpxbITB3MX0Z4fmOfsFn11IFsZx3YgdqtwFu56o+TCI/XV9ftge9IsmGgaUFseRS7nVR+UHrW23awRl
T/nbt8vv5/y77/3v/PfDmnhl8zkxo4/R5VrSznb6aetz2XIp4d1Zab46AmUZuP41fD79916pccnbrWgq
eAc8YDiREbvOvdC1wX2o74UkNDa+Cql/uz+7+CwZoaqoC8BNZRlM+S1fhiCJF/0sm9ETOTubs8rtNYeT
3vIcIXEh7X+oHu1f6qOSmpf/j/pX0zXYSl7ghWQfVrwgoVXuV19tIFS839XyJlnxcjl4dxEaNEi3y6qc
riP3KElL6vE2IAw2ukcGYgM7E64tlmfrvJKxtyETiVjUX8NKh02SfwIAAP//HmVTkzsIAAA=
`,
	},

	"/definitions/nethd.yml": {
		local:   "definitions/nethd.yml",
		size:    2425,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6xV72/bNhD9XqD/w4HehriNFNmL04JAMnQZug3Dfm/AADcYGPEsEZZJ9Xiylxn+3wdK
liVnCuIP/RJL7717d7w7KlEUvXwB4A2jBIuc6/Bq1Qol/IT83TfhVaNPyZRsnJUg3sHaIAeJR2BS6RJJ
BFmhbFapDCWsTbS2AUKbOm1sJuHPP95Hb2uVsUsvwxNABDlzKS8u6syxo+zi5YtApapsNalizBwZbAGA
y2Qi4Ue3NuhhtH/ouKmEd5U2LlCVN2nHfCnhl9uLb0PlMILw23GXgYMR/O4WvFHUY2YSfuYcCUbw/Upl
PeZKwtfOLYNZ+G1KB1g53SvWo6I0lzD/eNdCvI469DwovLPngOVdY1G4zNi9Qak4l3DBaok1HJd5uU+D
nDstoXSeG8TYsuIuc+WRmkmK7RbiW2cXJotbFHY70SpL5f3GkT5WtminRCJHB/8IPBaYsiMJC0erUV1g
eIIb8KWy8UaRNTZrAhg9yy5jcyxHhJZ9c6pAkmLjHp39SNT0tM2rzbo+URQSR/gPo9VwA7oIf1imzrIy
1p+J32pfMYbXoOsd30/gmUzPNvkzUpu6baRshhDfHtZ1t0u323i3u558sd2i1b1+t9Ovu/1rhfQQ/4AP
odce/if7WxEqCckjNGxZHy09K+4jxqbFvXPLlaIl6keERqUlTBqI3Ka/r21vmWSu/BnrOOzLuBEsDBa6
U7PhAg9vR9FNWJiDWcMNqDj0ECkqCdcGN/M69K4LVcxk7qtwhppqmf0H4OGpLC0fcgy65YSLjliYgpF8
zy1scnNPPoZReKbDzh7cKPMyVNLCGlmZwn/yk/drTd1qFbZxMImaB+mrazFqZeJ5R+02tnBKDzvG7bbM
nU0Lky6HDffkaR0lLAuV4lA35+J8oQqPY3EOol/8idGEXJGFtuizYNMeMFzgr4y+7my9+fepPZWW8yjN
TaHPZuPuGqA+PtQTIVeHkAIxzU+KeXOIyUjdnxDw9hCgwyV/Vn85Pm08qizR6qH+CnidvEkSMRQVaigV
+cHBiGmSXEXJJEqmk5lMLmUyg+jIqh3S2hXVCheqrr+zSpXHo2JF/b+kUPdYdB/0D+I9IX4QYyFBHJf5
hH6WfN7K49lxwKuATg5YVX6C8qZ/7bNNB3P9FwAA//8i4IYQeQkAAA==
`,
	},

	"/definitions/newretro.yml": {
		local:   "definitions/newretro.yml",
		size:    3647,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/7yW324btxLG7wPkHebQBwfxgSWtJNtpCLiBbTUJkLpx7cQ1YLsAtRztEuKSG5JrRTX8
Wr1pgV7kgfoKBVfaJWUrf4oavVlbsz9+Q858JPfP3/7odDqPHwFY4ZCCwplBZ7SPKFYghbc5wg84g5Mm
zNGmRpROaEWB7MNLNAVTkKFRaJgEZ1g6RUM8K5nKKpYhBY4djnVIqKml/j+ADuTOlbTXUzjr1Hm7nCvb
Veh6jx95JmVlA6fMYaaNwCYA0E/6FN6e9faVKBA24IWQBUIH/G+BgRpQONLXAm3vQFYnbB6hB7IybN6i
w36LvnE5mpjUUs5nWvOgO2zh0dkoQkdnowBtfxrqnITEg09gr0aDWG6YNFzEvJbMWjFFE7LurME+DHa3
A7G7hji/FiPowUhcnzfg7jMK5+fnsAH+Tzt8sE1hv+JCwwYcVVZMfdnlGFUgBkuiVz/HWk8j9tXHX824
SvOAD+8LHokPoYvDBugdHQ9j6HgYRHYa5kxwjKXq3y3XH1I4PoQNODY6M6yoF/+9UNWHgGx7pHfE0jvY
EUsDtLNG5yehuJ7ZAD31JoUN2JcywwKFgg6cohFRsfq7tZFHOq0KVI4Z79EF462ip1Wk903NnpbauJiq
A4F6RqGx8IGQPDLHIKFwoPXUd/4735fw5mkY9OLj7wbeKYsGYSpRKOvi7vYD+YopPodTV00m4X2U/lQr
60SGYQlJ/LYUKGune5kw/wGFQ62slhhDrxehqHL9uk8vWVE7uQV943p3e5qsZ48PF6cNQKF5dL5YZCbN
KVy8v2pC7roToluesFptAZZXCwmpM6GWAiVzOYWeY1Osw90yL5dp0OWaUyi1XXZMqLJyIXNl0SzOX3Jz
A91DrSYi6zZRuL0lDVkya2fa8FWyicakQVcZ5TQF0lsG0Rht2sPYosTUaUPBsbHEbv0UaqwNR0NTrRwT
yj4h+6pAySuVwUu0aY7CoXH/IZvwLbix5nP/1/gHX0iwRQKH1tEw7bo2fkVpuSiMf2WYE3q1fDGyaEoz
Sy6uN2Th+GK6K9mj2Z7UkovZcapc3klzIfmTweZCbyKkQxNK31lefAZLydL2/AFgJrMULkiXbAEhV1/N
b3m+S5YGWbpnZYljo2cWwxLvmOG/hs3q9hqmMoTuYXsT3t6mNzfd29u9/v9ublDxqNuNS2tX/FihmXdf
49x7wsamsLlPvUT7bbhu+XhOgTDOkQdcG0eB+G+ANiZUKjkyvhRYhP2Kon30GWNd2KoomJlfrTSw6Q1K
HnS4nimpGaehyEGaXeQGJz/vkYby9XzutDGo3F7oFwBzzohx5b93/JB2awsn8QvS6JiQtlYWfI9c0ZzZ
J+PNZtDyG2X+eZXQ7ucp+6qp3Xdp7Lz3vrvWGaGy+HXjwJS1t0KqC3+92LXza7erb8D9HRXvnLt7adtH
2BeXYcUv6wv8D1L3fWS8EmgyZIaNH3ypwzpfIzoREh88xb0lDTbDTkK+aoMHyTj4dBElYpr/GymHbUrO
3IPbZGfz6/bSvVMcopP8skoSlvjjHOJNG0b7mZfM2LXjSTLoJv3uIEl2ob9Dk22a7JC7R9u1llWBE1Yv
OKikzOLKfMlEK3eRaqnN3iUxyC/JVbj0LskbJefwrvSSl2ST+OxkZfz/ycqBX5V/K/ty9F8BAAD//7rG
rvY/DgAA
`,
	},

	"/definitions/norbits.yml": {
		local:   "definitions/norbits.yml",
		size:    1742,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5RUXWvjOBR996+4uMu2YWt7k34i6MJu2e7DTocOU/pSMkW2bmIRW/JIN0kzwf99kCOr
CTSk8xJ0r845OvdYSpIkEYCVhAyUNrkkGwEoXiODz6GuuJrO+dRh8kRp15FqZlkEAJBASdSwLPP8VCFl
UQRQ8MZDCk441UairwGGDO71QqL19YjB45NfnzF4uPXrc7fO/uN1QF4w+HsupPblJYN/tJ7Z7N9c65lv
XnlM1v1ubVz3G09SoI66bq3Fmy+L3BQlg+fvY9+hRfLWPHUAq9UpYDN2dMNJ6g254VQyyKQS+Jo2ZRNt
9CosSBsG8VGHjR2r0lOpdlhdJ7Am2tSs++1KqZo5BYtzi2bzgeL1GtJbrSZymvZdaNvYAxtu7VIbsQvs
uwGIxmjTiyfvG9odJS20IlSU61fImSu4VPbk+A5ldTwIhBqtdZcmNLY19vAILbHgv3PiJhNIXFa28xOF
r7QdYG700mIwvBsZ/Gb4sovBcDVFSG/DhWzbmkv1UnB6Ht+s12nb/r5eoxJtuznkxkX3ZY5mlf6PKxec
Dcm5I98uTj8Z8bzCI9LGoKJHVwDlWqyADFNUJkUpK3Gi/hhtJp5IrESQ8S9lxd6JncQWfziAv0DIRcAV
3O5EHacKlwWnl26+O1nVaGIGw72Qx6eYwWjv9oPRU8PrjcrZXtjXRlZVzOB8L+J+buVsFjO42G9FCmtn
Rk4oZnC5F/ZpJfJZ5+fqwGkL99w75HXUv2tJFR6MeeRi5ju5BwonMjKfu7/OTsxv+Lv6Qenn0uDk202c
eVo8fk/fgXp5vVSV5uIX9XvaQX0rfxyO5apPQXA6jL4YhEeCAs3haIZ/9owKsSg/RBkOouhnAAAA//9/
5GKGzgYAAA==
`,
	},

	"/definitions/pretome.yml": {
		local:   "definitions/pretome.yml",
		size:    2086,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/4xUXY7jRBB+X2nvUPIitIPWDllgF1oa0DJCPMCIQRpWkaKAOu6K3Uq7u7e7vNkQ5Z1T
cBEeuAsX4Aqobbd/opmQeZjYX331ddVX1f73r7/TNH36BMBLQgbWIZkKA6B5hQzuHN6b2wYQ6HMnLUmj
GSQcCtTouAJyPN+iSwJHcV3UvEAGqNPaN5DUW8/CE0AKJZH1bDbrDsqk3pjZ0ychnHMbeTknLIyTGAGA
ly8Z3N3A5O8ZvLFWyZyHmnzPfM3gW2O2fsL8bh2gyPk8iM2+5xX6Ead5j5TP5gx+ohLdROZW+hyV4hpN
3VPnXzG4Ne/lSCxQGyRyXjF4UwtppvXf1l7mkfKawf1bOG3x/m2MfzFnsFgsTuKLxaL1D6AyYuSYR+7y
ksHy3SpC9D4d0BeB4Y1+AWhXrYRHIqmLYVztElipe4W9RQaEHygiiq9RMbiLnJhUe3Th6X8zf5kQ+zO5
9zvjxDT9FI2H93CzcaaQumvBcioZzIhvsYEzW9rOLKTSCAbW+K6ijXEVa/6371LbmgY/HVLtNBkGySzp
Cwiiv1kZrsThANmN0RtZZFZqOB57VvRiSoromBkbPJHr0DGz7RKSH8Nvh6JzxvUFp+BRYU7Gha40LXOj
jLtOHIpkxUrun5fzq0gOhngfru6AwDmBlkboiQ3FN2ZX++iyCzeTReYHYpD888efSVy2ZhMnc1o7s/M4
DOlkBB85vmuccVwXCNlN/504HnNOy9X14ZAdjx8fDqjFyKu49I2nP9fo9tkPuA+O+rGlxAvPoE8iBvM+
tGHAlWpfQ5GjaxYtIr5WCF8DrY3Yh1+Xtf106yVRiSGPJKmx2YMOXzax1bJ0uPn1OhFIXCofXElWQwIn
cnJdh892w4+R7uO5f0S7Ex2c/qZ17hHpQB8CG6kInZ/sSLyx74KvnpzUxTgMwF2wtTmlr9FUFWryZ2r8
5Dp5FmkX1CbMTivDxfmuI+uMl2NRL39/eEQkmKYyzUupxPMvr4ZdQDG155GU+ad9jkLMy8uShtsqOF1Q
2KuryybHrUUtHhpaArwwyX8BAAD//3WJUE4mCAAA
`,
	},

	"/definitions/privatehd.yml": {
		local:   "definitions/privatehd.yml",
		size:    1262,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6RT34sTMRB+379i6INYuGz9hUrgTkQfBBEUy73IIelm2h1MkzWZtNTj/nfJ7mabgqXC
9WHpfPPNNzPf7AohKoBAjBI6TzvF2OoKwKotSvg6IJ8+VgBG2U1UG5SAVsSQELK/gqwAAAS0zF2Qi8Uk
UrOrKoBGdSOnUYwb5wnHOP1my1vxvXV7WDrv0fJMwvL2mP3idoRFro9DkY+BmiL/PmrquwJsnT42Cqh8
00r48ftuRHgnjuBVIgRnrwC7u1Ru3IbsUNwpbiUsVOR20cM9unZ+K/tnH5LtIk/dcKvI/IwB/eDi7P4e
6g/OrmlTZxQeHmYjvVMh7J3Xp8SMTkT03vncQpyZq98cQ0h3moC0v8GGnZdQK4Oe4QY6abkVTUtGP30x
77mMgeU01CDfNC5aTp6MdhWmzBY8OB/ekb1+9mRgXKclvkX0h/ozHtIKIa/g3b44SR6J1cpg3T+LmV7O
4QZ45fThZNAE+uEEhEZPcuPLdZi2LvR1ofA8KdAJkCsUs6dVTF8CExvMb0r6fzSz1IV6NECsyeBplI48
FmlkRSY8SuNkvNbjOmu7vTVO6UvimSeocfaCZqA/ZxYufHudfdOKL7NfJdtDp/7ZunQ7IGr0Z7wqBN/k
9gaxaf+n4u0cqr8BAAD//47AQgPuBAAA
`,
	},

	"/definitions/sceneaccess.yml": {
		local:   "definitions/sceneaccess.yml",
		size:    2626,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6SVbW/bNhDH3/tTHOxhWLtKsiUnzQQEQ2CjS7EZ0ZbC1VBnBiOebaKyqJF0bK/Idx8o
UaJkJw7Q5EUs3f93D7wjRcdxOgCSKQxBJpghSRKUsgOQkfWRLSXZckOWGAJmzqawsOyrDDsAAA6slMpl
6HkNJxc3XqcDkJDcYAlRuOSCoXnXfxchwIQ/MJTe7RigZ57H07HzVw35fk3VxM4/H9bA+2M9fmC01ofH
ekSSr7IGBjrCp2lZQ/HrtBL4Rr82+vWBPhgUeim2UvvnDaWdNAgBIBp5v5G1LgyKXy8a1cBZCCOeSZ6i
F90G0DPAbWCj95tIBJaJLHNhmc+M1cxnxiwTWCa+57vgvF9h5tWiv4RPVd1aWcGYeKeQaORNSPIMMeyH
AFcbyrj3B5cyRamnt5Es8T6kJLG9D2puEgWlzXCTyLZqcFZjU0aRQw8MNP04tjO5qChoRTrYLmZyH29v
oAdXeS71o12bkfuU7Cu9P7762/rrDXmjVii8CZMJ9MY3sfUeHIj6v+3KwJ6Y6zH0Ij/y6tf2phz6rcPV
RA/39zA4OIcNtn2Sho2DoqHi2cpnjXNi5OtGd4eVbE36fMRxDNCL47jdpuB9Szs4PWel2P4SDHxjjR9Y
I0mjaR+4QLbMoFc9PPnJCfyXPVqJg/7LDu2aTC+O4E/Tg8jB82ARsSDXnNrvqkQiklUIX/69Mxb14Fjj
Ow1Inr0DzO+0e8qXLCudc6JWIXiFpQyMasVpCAsu1oWBZflG1ak2EkV5YXS/fQN3xLMFW7qVFR4fuwbM
iZRbLmgbrKw1iEJwUQV3YI1S6mun7oZeXIqJ4iKEbq8oc37Pd3OKMikjKJQqrJMWq2HZPd/phZoeNFdK
0vSJZcEPgmxD6JYOl7riPzco9u7vuNf1Snh8/LFszaVf5hV82xhAVaMSrlJzwbeFsmCY0poyF+Lerq7h
R12lxFztcwTyZSVw8c9l99eEqMvuXY0TpQS73+grXBO1fcFShUI2u+aYa13gEnd5QwAgYilD6M5mRfif
ZjP685tqaoqpFE/UVwyZPFVQ4WkEioqwVH5XnMbCEr5eY6ZeHYfybZZyQp+LQ+c0fTGIZP+daoyWXzcP
PYi3rv43m0nXfeO+rYZCiTqVmlCKtN6ISFuZjwstCQOkiMnqtEOFdP4PAAD///re17FCCgAA
`,
	},

	"/definitions/scenetime.yml": {
		local:   "definitions/scenetime.yml",
		size:    2512,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6yVb08jNxDG359032G0rSo4kU0ChIAlQPwpbdW7ljYnrhKhlbGHjcuuvWdPCGmU7155
N2tvUO7ak+4d+5tnxjPP2KTT6bx+BeAUIQMnUCOpAj3SvEAGI4/er5BEJ6wqSRnNIDmDDDVangNZLh7R
Jl6Tc51NeYYMUHemziPUwkilMwYzpaWZuU5/d7BbiZV+dMz/BdCBCVHpWLc7m83S0EoqTNF9/cprBC8b
seCEmbEKGwDQZ/DOPCl03dFlw/YCu7wJcH+4QTnYCI8C/DHAg00HHewHuBfgYS/AK2NRZTpEYo3zfPo7
n4fAbgj8ShO0YbjeJwNxxg/fn1++jYFBEwhNMri+6P7Ai4j2DxlcGO1Mjt0/7s1z4EeRX4+ugyG9Nt4L
uB/xB6UCHkT8y+UodjZsV9lv+C6D9zctT/f3PGi+jqpoaw1r0eHwRfKw1o9KYyke3Kvgy2VUlbx7nWjf
oPKqJ+NqhgzOjXkM1g12K8k7LkJLg5Wke2EKJaJyzyvDWAzOplKZ0FV/Bbo3SqKprzpAYWTrcjvkVkwY
3H68axA9dSLd8Qpn9A5geVeXcEikdBafV/2ghTGPCkOReYkMCJ+DSTm/x9zvp5ZVr9RkSq/KFEgTI9fL
KF1OKfZahxgkiwWkF0Y/qCytGSyXSa0idBQSSk4TBt1inpaTsj7SclJmJXgRru3IUZCxDKR6+mZEnKZO
qic48d/p1du8lj2onNDGzhoPLGb4XDYUgNvMMUjG462t9M3p9ni8nfg2ahcri9c6ubdm5vCvs+ufYkeN
L6VxtNGVby2fVZ5YrjOE9CL8B1suxWKRLpfH/e8WC9QymBTXXln52xTtPP0Z5zNjpYOWTHDiDObxWTsU
DP7mz81l8g237lJjHtm0nqWxC3MZdaQoRxZdinn8dmLx4c/jRCJxlTtvQ3IHJ/BgNLWa8gPOP1+hPt8X
OBWcjpO71laIrLqf+h8mL46Bl2ttr/ajN8mRVTprh5sdCx77M0WBmtwXTPifvUkz07nh8mvW/Py8Fsuc
C9w06237oFMlj5MdSJoOPeu2T48leVmilpsqJl05LYp5SsZa1BSvqfpn80UhyTRNOmKicrl1sB1vIMr1
gT6RMgwpOaKY/K+cw5AjOW1uy5Vcp5jz0qG85IQbd1Bd/5eLfTL5tMAHXtWJaYI7XFtP4l8CnMA9E0YT
V9ptjZMri1jNMU62EwZJL1lLeeNZP7Bp+UUHrrL/DQAA//+kAaCi0AkAAA==
`,
	},

	"/definitions/shareisland.yml": {
		local:   "definitions/shareisland.yml",
		size:    4357,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5xX227bSBJ9D5B/qKUWQXbXkizJtxBwFraVG3adKJajeOEYQotdpmpNdjPdTdkaw380
fzAP8zAfNL8waN6asiXD1ktC1alzum5dpP/87fdms/nyBYAmgz7oKVNIOmKCW6NgMfowXDRy1IGixJAU
PngHEKJAxSIgwyJiAoxiwRUqz/pGTIQpC9EHMk0y1oQikJxE6MO30/fNvcyLxJX27RNAE6bGJH67XQul
JVX48oXFA5aUjgEzGEpFWBoAoNf14YuZooIGjCip7FvbPpydneVW++QYPcfAmwVSx4djOSPU0CgenN6b
Emt/7EMDPja7O1sVut15iG47rkPLg4+YMlIKp9/ZqnyGVqFPMxdyZ+ehAAmEz3jtFHr3ItztbtbKsbkI
djb36qjTP4zSEzaHBhxGKdinPukrd8a9PA/7J/Xy7VZof5QlMeq7M3r3EhwtcHd8GBxBAw6SJKKA2VGr
FWfPou1NnkU2OFrh9SbzOmaBbR8LSBipp8t9t7Yy38FUCmweCK4kcXt68cSWcXZ9OEg5SaueagoqoLtZ
IO0RcZTZzHGUDu+U+PGgZ9lJz2E+HEmhZYTQgA8sro1cr4Lag+EWNGAoxRwGQ5eEc/hOBA34TuSEdxx6
NpE30ICziXQz1d3NKpAdmRXVIXuO+bk/hAZ8JmFQcJdRb8uHQymvLPMdp0BGzN2FrsO+ptIQJyZcXNs9
B/+XJgoVOW5v24fTkU0VFeHpyLEyIB+6DIPTGXx009XbcZf6hBI4nOc77FO1w/KK7Tq3PukgYhSjcvCe
gwfdAQg011JdQb6JAGLJa7tHI1PB1IfznxelycyazrphPbQUG4DJRS4RyZBEIZAwM/Whzf7PbtqZvZVM
i/sQo5lK7kMitcktl1LFfvZv/ptEkhoXSsafyJtxjPEEVb7EvdtbaB1JcUlhK9W5Fe7uvAekhGl9LRVf
pJTWOoUFxTsg4z5UUpgH4INnVIqFAyolVRVrEzRGGBipfOA0a2VoDhrUpvIr6vPvZH8qY3yVEN/vuNLX
FBr5mWMSl3I8YSovtbK3d7HUS6ScEMXhuSET4b53YqneBfwLdMJELle0dUFvaSsSm7tUCoXRVYESssXt
dSvD3xW7zsqtmAgRWkfVq+3uLiB+fuHf3rbu7l7d3qLgtQZc4dz2RBet+pqimrf+UxjrncrjHZt5YifB
tv4eRLaNWcJeOd5KXtenu6yMYZMIG2VO4+znOIiY1hTAWzATyef2f+ULaV6fE9/3ljuPp8i4d/GPYqQJ
I+7Oy0Jx7/X6+bxV6I2zEX4L1cIovgjmS3l2uEqHMcUsXKACMGMUTVL7ETRVeOmAS4oMqtpHRja1+aX6
aSuujSIR1mEApkLtQ0DVsglkHNsaPDenlYFxeS0iyfhSQeZPmX5dG+J+4Q2n+TlV3R85QdMvK1tgMTcb
yBcLtOiaw9V2QAymj7iXeAmHik1WOgcyTiI0yO9XZSajNMZLlnk6csA0LjTSczX64b0/effubz+8C88H
b9Nb6VbV8jiNDCUR2Q232douma3t1dxjmVMCZqRCKKUe5f/TWjuVLU3WTvBb8iDuXnFqb3XMS1jdgtV9
FqtT5dh5Ro1yoSdFuoL5hGhXMB+J+F5XODNP2lfORWEsZ+gD28heLBt2Q23ARD1t8yhMIhbgsq1zXtQe
OXgb4HkXz+bDZL4uk5k1mV/CkCz1VHI2X4P/CVXG/x9qg2o9jSgVyP/41cocS7GeRsyUKUVOU9TrqRyj
CmRU6nxHLtZV+kByVsUzTdWaMiP7R34h817ReiJDNmFGWokhM+maTerLGAUFLJNJ1+wSsCjCclL/CgAA
//8ZpoSEBREAAA==
`,
	},

	"/definitions/skytorrents.yml": {
		local:   "definitions/skytorrents.yml",
		size:    1198,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5xSza7TPBDdX+m+w8irr59uGsqfkKV7l2wQC0TppiqSG08aq44d7ElDiPJkLHgkXgE5
aX4kiohY2T5zzvF4jn9+/xFF0f0dgFeEHPy5JuscGvIBNCJHDh/PNcxRLcypFCfkgCYqOwhNYqUyJw6f
tm+jNx1LmbPnYQcQQUZUeB7HVVWtZ5eslYnv7wIpEcXATgThyTqFAwCw3XHY7obTe3sJxeva6wFyK2cK
j8IlGYf9l8MA0SWa0IfA8NY8ABaH3uJa7OmFoIwD67FYaB03jUph/aFEV6/fYV1ZJ33bomwa1B7bVoSd
kW0bb+Km+Z3JemNnq3mXGhOyjgMjcdQIT0BHK+uwuqsgVajlJCFFGsfT3IIkN5RFSaa0/G+zgicQPFXO
Uw/tO+VhUgoip45lyL0rDRWJJJT2/3LHTfPMYTo+WH1b0Pzz1UBJlcYFnbxYTRNFiW6B5NUo0YhJtkjz
ejXNyCdOFaSsWTKniZIIPx8AgMpPfTKPbIdOpQolCCMhF+6Mkh043MTnFux/xoGNkBS0YMYvp7fYymgr
5B80IeN9yPDzI4tDHuz2H5rHPFherC5zTEXnNakIvxIH9mxsuSz+zt6wXwEAAP//m2FCKa4EAAA=
`,
	},

	"/definitions/speedcd.yml": {
		local:   "definitions/speedcd.yml",
		size:    1961,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/4xVTW/jNhC961cMvJcYiKRYtpOGhxQbexcpFtt1N4FboO6BFicWEYXUkqMkbuD/XogS
Kdvdr4tBvvdmODNvbMdxHAFYScjAVogiFxGA4o/+njhAoM2NrEhqxWBwt4RbNBItcLXVCn8dRAAlV5ua
b5ABqri2DSLVg2URAEAMBVHF0tTnTKMIIOdVx+eccKOblO0dYJoxuFumN/Pm8uZumV7Hn/nWkyNH3gZy
vpzHnzsy2+Mc+a6SVgu0HT+5PMp8M/dpp475RAUax3yQIkSNjrgFzx88OT1r36y0IUe6k2fHLftWyUfs
oBGDj/pJom0T+jcyD/vIY2AyDoGh6snFkSj7JYiuy7qf2qTHxyH4LGDzpQezKYPFLP3t9pOv45LBTCur
S0z/lNLHTnt0cTv22h71yLjX/bXWL+PzM5/h3D20KLTCeH8S2cQRZ+LA8sUs/chzL7lgcK118CA7Z/C2
FlL73JPDe3bZ3dOlFKgjBz82e+FXjp5ii9zkBYO/v5yCRW61OgWs/mnEhpPUrbTiVDBI10Y/W0yqonKo
xRJz0obBWostXIGQT0l+g1zAFRBfl5hYss2548k0H4IpKuK8kKU4GQ3bsD0oG7rk97IkNKHUuPuKGtzg
S9WBANxsLIOTFazE6+g0262SlXjNdisYNh2UeiPVQQfEH/Cdg0MXj0iFFgwqbckBUlU1hYdri6Z9evD6
CslMq3u5STwKu92gE1bc2mdtxKHQo0GIxmjTd9XP0A2J5VoRl8qeDN7LEuF3TfBe10oM2qkQWmLhQddS
U0oaNe12Xn7HscPWvPeu3j9qNNvkA26bam0ot4nv9b7WwRvSxqCiu8ZlSNb6ZaYVoSLvfOc5mUHnJpai
3ztJJbJg4t4MRFLeU7sScAUcrmDdybofzO3Xw76/Ug7i/98xt0FERq7r5u+gMHgf8KP929/BL82oLBmp
NnusX8acU4cK/axKzcUPSx4f1zf6UX2C088M0FZcJY22O38tqzPDOyz//UbaveKmw7APKA5G9I2Acx9Q
IubFz0RcDKPovwAAAP//hb9Kd6kHAAA=
`,
	},

	"/definitions/t411.yml": {
		local:   "definitions/t411.yml",
		size:    2364,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6xVTW/jNhC9+1cM1KLYANGHLVuOCWyK3RTtoQi6RdMiwGIPtDiWiMiklqQSeA3/94LU
l6XIzR7qgwy/efM482ZE+74/A9DcIAGznM9nAILukcCDVApFAxVUZBXNkMBO+Ts1A0CRSsZFRuDvh1/9
G8vh4kmTGQCAD7kxpSZhaPODgs9mACktm3BKDWZScWx+t594syLwoWJcDuBkEU/Cy0UDh3+YHNUguIxa
qdA9t1I+jQgJgY9SPukRvG7g8E7ueTqO3kwnbSbhedRq3dOMfuMCh4TFMiHw6S78je5HkXgTE3jd1SJe
uYR7mo5wJzRyreZ+yqVA/7VWEt9M5GxeYcvVisDDP+EHwfc4UlhfjMQuMnXq0kV+kWm1R2GoOoziGxsf
YXMC9/KZj0xa1qeMuKuLJycu8lcplRnqJHMCj4+PU2uULC6FkrgOjXaylpo5dC9Zv+IaqUpzAp+/fmkQ
8+z34LUlaCmuAcsvNr2QGRd1cklNTiCsNCodOjys9dHkkhEopa474qKsTHdirQDe8QjBnRQ7ngVWwr7e
cDp5DaukWr9IxYbEFu2IqJRUrbIPGgtMjVQEGH/+YY9a0ww13ELpKAa1Id0BrnpT3yc6rHsOO1taIRoU
MpOVsc0rargcND8bkhl/Dlx/H6mCW9AlFUSY3E9zXrB36yuLGSVF5vJ2vDCodF9+fcUpLAua9qtLVaYJ
fPauvWvwAs+NoZnQWSneq1a8CfN/VPTFWaqoyBCCu+7OO510tU2peX88BqfTT8cjCtaPo90IN4w/K1SH
4Hc82FHobhZKvpytVWuJodsC4RbMVrKD/VZN71iwjt5cvYfu4j3LZ2cOzq2DtCVRYxTfVvY/Ile4a+GR
refWfrWVa6N4M4GhxbUB7WvATYG9xoWCFq6gQYX93PrynFgTYGgoL/T/JX3WOZMvopCUvSkdD2y85ON/
GTne0fM97RdR7GT4M2fv7eL2aFtmHWqvHc2/vW13ctXtF7JBaecJQVW2dw1iml/m2Upa76j5zmmzwj6G
8PdYZk8oqdJTpnmLKEr8aO5HC5ivSLQk0Qre+dGaRNFV+w5mim7fXpv11b8BAAD//zkDgSM8CQAA
`,
	},

	"/definitions/thepiratebay.yml": {
		local:   "definitions/thepiratebay.yml",
		size:    1499,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/7yTX0/bMBTF3/sprsIL1ZqkBbZRS9O0qZqQNrR/hRfEJBPfxhbG9uyblhL1u09JU7f8
Ezytb/fc3z09vnbSNO0BBEXIgCQ65TnhFV/2AAy/QQZTifCjVeFzK2tuyoqXyABNWoUeAJrCCmVKBmfT
L+lxwyhzHVgPACAFSeQCy/Nd+8z6Mn+uHxrALR71nbe3y8x5m/d6AAV33V8UnLC0XmFXAxwMRwxO7Vxh
iMrBRskn55OoHjL4VAll83Ml0ALswWkVVAHzptwOH8Xh7yTRN1hTQqGV21JvGUzPY/WuqTq++e3BCTdC
ohYReR9tT7aRjtu5HWEcqcMojsf3E/Xaxo0V2yXQPA3IfSEZXPwdQEAerBkAussG7lot6zhJBkldqxlk
Pyv0y+wrLhfWi7BarcG8rh+38mE+HufDukYdcLXKPRZoqK7RiNUqaa29XcRAATUWZD2DZG/t+gtDpQno
yoolkGeSh30S2Rw9TWV/7TBTqEX06O56ual3XUkwQzItpNJif9QHzjQPtK4jzom8uqqa9y49zqI+U5rQ
h61v8/TWn0BwWtGODsB9GRhcQJInA0hHcLlZuSKNT0XLBNI3Za67lkDiSodXkM8FFnZhtOXixT0c9IFf
tLn+fEgm3VRy+YJ9UHdPnuOh9cwaahJPMBSvWaXHEm/dE7tMfqs7hP3szcf+INmckdN/DnHmmu2geBAk
IIp7Xs9kOex3hEYs5GsmjvoP7nNudXWDM96ycZrwlhgkw02gyr3IjpJ/AQAA///Cjl7V2wUAAA==
`,
	},

	"/definitions/theshinning.yml": {
		local:   "definitions/theshinning.yml",
		size:    3743,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6xW627iSBb+31K/Q8lZrXa3A9hAbpbSLQLewHJtTBJWmUxU4BNTiqlyVxW5TMiTzY95
pHmFURm7yoZ0T0aaPySc7/vO5atTxr//+lupVPr4ASFBJLhILkAsCKWEhipI8RJcNFkA8nPRAMSck1gS
Rl1kNdA58CWmKAROgeMISY7n98AtxY0wDVc4BBcFUApAhYDOWUBo6KJHQgP2KEpO9aCakAm9F676D6ES
WkgZC7dSyfVUZjz8+EER5jjOmHMsIWScQBZAaA/9l0RLyL5Wj13UZw8ERMVvoT3UHo7Hw7FGDzTaulSw
31q3LlsZfHJSFPuttT8ajidaXt/Bn6qHdQ1Xd+DpZUdnrx66aDqdpsDTdDPdZoY2CReoBXeEEuW1lpiG
zqLVGD+jPXTWuxg3/p8xHMNoJ/O21o59bMc6g6PxWorXdEeOfbSjrne11N4Bj6oms2M7O/jY619MDWE3
QdHOo8yQdmpIJqwWhX6/0eut2628Y/2VIPdaUXNRYxUQhvZQs90YT/wMqTkpUukxISIQQm1Mr9HUZ36U
EZLPGWP3ydp4Y3/U8XqeTmRnvP6opuqPaqbfQ1Peb/RHPW+cb3UMEWABXPP1jqp2x7WrtjmRwxw0HHS9
HvqfN9DTOEc5+GnSX08fSMEVHzjRd8FxXDS53Oxi67K19r1xxxvowRNQH0wBq1UVqKbsTL1t0LFTcNRo
drewY1PR39UdpMLNeSaoBrOKF76W5cZiVEgS6slOjl00aqI91BiNtDknSV+VFpuvlkAl5uq2tIbdC8Nw
XHTG2L1yzzsbDrvG2HrO2G5nMGx3zA45By4aygVwdZbe1wvPNytcyyDV5XDgTzrnnhHa9WTDK50lDkGA
RHvq67rTb5x7hfliApGZ7tBFTUYFi6AyaPloD503+p6/HrTMJFVlQOUcL5OON/jILHXdZBj5dcMwCQ4M
44oQzbjqdDSlZijTGXvSnOnZsPDwuuIgZJT+YiTrs9mCmHGZPhcaV9qTbczvN5pdfQEOC/Bgqp0+KgCj
0WUG1AuAXyhW3caKxQ4K8GTQUFNtwCULcj8yAjCfL1x0/e0mC8mHkonuK4ZgdB9BfLNJEbGQ0DRBjOXC
RZUkVI4X6SN0CXLBAhfdMb7cRAiNV9JUXQngm59l6+UFlZuM3pGwnEXR66uVMWMsxCPjQZGZRQ0TOGdc
5y8hARHMJeMuCshDWcgAOL8l9I7dPnKctilBSNcUSiaZcfYowIyCcpnIMryWREZwao2xJMy6QZ8Q2ZjC
VaBoynaq9yVKrf9hpi03/8HxY+IPxzQEVG7qd4nX1/nLS/n19dT558sL0CBnbHbEia1fV8Cfy114VqaK
vP9ioUqnVEeHCZ1HAeCgEGQ8AD57dpGFgwACk4Nx6SJLvXGlMTWPu2uwxLMIyktMKPqM5IwFz+ovd+eM
Skyo+JfViCJAE8Y5UCmsf6NPSHLFCdRHok4+CZ0lzRTSbOrdEYgCUzxgjzRiONCBfD/4esHh7udTK2OV
hIjUMXyRmxZOrRujw1JyMlupN1Al09dJHfSb6dVmJmhuJ7cSJXCGpO+Izz/u1SzLlzl+V4fKlEgCF7nE
6g5tbug3tRtC8txTMM3GQ+GqprSXIDGJxDtmRZ8R/tO2BPnl+77NOLvlJFzIW6zOfva+WTjEEZ7DW3Nc
W2VrH1l5u96p21e6shGGHM++b0Kuccg3LgCCYtu7uhl7cm4xzssigPniHbpbXNQFWP7YXKVRJPQZiRjT
v8Hfn1a2jW1lFnrbZVUuxly8qbfsatl2ylXbPlSvLXbdtQ+s7Vv8wKLVEu5wMorJMscCCv2qibIn8JBG
zxexdaNKWHmS9R+r8HhbxX+lhFH/EQAA//92nFktnw4AAA==
`,
	},

	"/definitions/torrentbytes.yml": {
		local:   "definitions/torrentbytes.yml",
		size:    2621,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5xWfU/jNhj/v5/iEZ0mul0b+gIHltB0tOM4bd0YRVDpiiY3edpYpHZmP6V0Xb/75CRO
nB0cp/FP69/L84qTttvtBoARhAxIaY2S5ltC0wCQfIUMbnPwogATLpdrvkQGKNvrDBHy0bAGAEAbYqKU
BcFms+n40ToSKWg0AEKeFtqQEy6VFlicAQZdBnB7B9C8vYNrHj4WeL+f4cFklFOTkSNOc+KqIK5G4Jhe
zRKM7kY3InV5jhjAWD0JNNAsPr10g4FH/6bkJESJwfRJRL4gq/MltntW2vP05fdCcFzxVx5/5fje0csB
Rnej9o1L0qs0F8n6hm9L2eU6SSDHXL3Heb1VuRfeMI5f6tYX2Ol/WEdCZcfmeG1EGPjzOnWC4FdlTILG
ONVlwsNCdcLAqcbXfdfqcWm9ExEqaEJutPtyqxy8oslOxhVx8t8en3snA5/8okdfYMd5PQw+8pWVQPZZ
69EKhkoaleCr/PUwGPPwRdo2/2X8oaPfMxc8uJ70oVnwEzen/lElmM7Vs1PY726UXox7IcBJ7oVwihNW
1ZCtMlcoilHXNbYN+JrGjeJVTZfBdDrNZmA8jTcTe02m0ymUf01rqO72GXuBri5JQWfoPc4dPMjzflrx
JRqkwvZp/HHys7s62aacvse8Q5/ZZ8YHKVbooNMMulQaxVJ6oLtsNaLfZfC7nUIwFiZsZOhKRdUTziDX
Yczg818PBUJP7Qp8ZwVGyXeA6YO1J2opZG5OOcUMggzppHF+O1dIsYoYLJReZYCQ6ZrKdGuDOn+IH+x2
0BkquRDLjkNhvz8ohCk3ZqN0VBc6tBSi1kq74G1XEvFHrJeVl2aMfU00qg0aTDAkpRlQ1CF8powjNMTK
OrKIQs7VcxatUc7MH8Jcq43BMl29afhO8w2Dg9x3bvv5Y4162/kFt7Yb4/UNTafe7TSXS4TOsHwt7fcL
kRDqP0NOnx/Od7vOfv/9bocy2u/z4Ib01+LbKqvNu94PYlolMFfRFiLx1BRSot5onqaoM2DFhVxoux97
CpUklGS/Q4iSUAPxeYJAWQTSTCo6ZAuhDbXDWCRRK0++EJhEZfriZbut1uEvg0mKc/NhtwX8XQ3ptYCX
Lk6kxXxtfyzEGhclnk/K+NtuF78fNC7xOfUIAK6XhsFBvsbZzC5yNvsp5HR+OJtFP7b+Oez80HI7IkEJ
vlm3rbJ2LAwREheJ+d/+11qO1EYmikffFPiNWEb8/XaD78uOOL2tPm6V/3cY1RbziuHMGRLEMP4WR/eo
1fg3AAD//09KDN49CgAA
`,
	},

	"/definitions/torrentday.yml": {
		local:   "definitions/torrentday.yml",
		size:    2975,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5SW/0/buBvHf+9f8SidPqIaadqkFGbpsxOjx02666l3IDZpcJObPCQWrp3FDrTr9X8/
xYmbZLSF+QfAj19+P99sHNd1OwCKaSSgZZah0BFddQAEXSCB69I0MSZORZzTGAmgcHNVWJh4UKQDAOBC
onVKPE9Hbp5GVGM/lAuvU6ij1kzEW7CUDqV8YGhMAHqVFgHgUlcGTufICVyUUAcgpGklEFKNscwYVnOA
dwRgKh8ZKqhHtzJVjH9CKoN3NfmR8UZng7QCh8Mt+PE5+IHnmSlGMU62fvej7mXOecUHNb8jhsnNxP3b
Rjs8FO10NrKcf4j7Uwr3VxFzphKbXHCIn9HwwdZrNDpEXk28pT+2UYzODpVs6Y9PrPuDBfj8yCYVeEYA
rm+gPbpwfWPzHhG4vmmI2PVmJwPfQB+fQ60uBsN9Ws2OBMFBilmno/E+birnjNsDPzQJ7Eiw1QR/r1i7
AaemXLu4BhTsrVmjQ/5epUZ3/HdG6VywBTYgM7dlIACzC+83umhcyi6YuS3BWXG/hZIcvdlVYJHZVdA4
BDUwq4GZVRjUCp8Ys8Anxhr/GCzweS6XwXhQAMWfbjAeVNSYAJznEZPtnKe5YqH1dEpKxJvOgh8QzyzY
0mzBP6RSHJVqgJecWkE/sNylzJDFoin4/N6OhmRviK0TMxxb3RsWoWyDxmQVfQuan3MpH0wPjZMPUtaK
PoHZBbRHF87T9LvN5Wwv0T7NA1Iqt7mms2BgTtZEhvkChabZyiCNuQ0/KHx6Uxq2C0LDjiEWMqrfB4U0
CxMCX77d2ffm0a2NxwWgpDgGTO+K7VzGTJSbF6gTGbVeKybSXG+1ywUCznoN/Qsp7lncL22w2TgG0qi0
xVOqEwLeApWiMap+mqSFx4xqJkumIpiIcGmWyww4hlpmBJwuDTV7ZHo1YY+gUiqI0IkbJoxHR37PKR9c
k1pDzvG0syP4Nxl9MqFnVMQI/Yvt07rZrNf9zeZ/6zWKaLP59v8ivb9yzFb933H1JLNIbfPL5FOj1DZQ
Teccu9U3xXUxgfeg5zJaFb+zRtTird8jCVVHOurrr+bl7xm5e4Y8qitdBrey85azqKE37MF7oFuKap2x
eV584SQZ3m/t94xrzFQtV3+bZBjjMu00zynNYkXA+ef29pej29vobe+NY88S0xxfDMpvBhWhpoyrn9u0
L5NQLoq78bLa6DVqkXwSXNLoRbXgNWqKfX+5NCc965vqVxWyr7+Gun5rDrZSpZzpHZ38As6/zjG4Q7j7
uV0wX4FzDIPtPoUYtZzvCXxs0+SIYfKaHae9/wIAAP//0xltjp8LAAA=
`,
	},

	"/definitions/torrentheaven.yml": {
		local:   "definitions/torrentheaven.yml",
		size:    4024,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/7yX7U7cOBfHv1fqPRyFRxT6MO8ML5F4qmGGMuiBMtthaVdAK098JrFI7NR2gCnitvbL
rrQfekF7CysnEyfphi67qvYLk/j8zouP/47N77/81mg0nj8DUEyjC1pIiVwHSG6QgxnnJEIXzrLxcTpu
hikqT7JYM8FdcAZwiDIiHHzkKEkIWhLvGqVj0JBwPyE+ukCxQTEdYvxaueYJoAGB1rFyW61K8ma0mH+i
vMn4XLSePzOsR+LcySMafSEZ5gMAHRcAJsPWIYlQwQocDk4Opq3JMLf3jH0ouBIhWvNUcKWZjznU33Ut
1JpMN4s4IVlMNTHzzdmtdoWdFOx0YpleiXnHmGXeMWaZ7RLzfibueltty73fP30Pva22hXfcdJaTQHBs
nOoAZZF3NIAWjAmnixzf7pZivxlNLftmNLUMuAAn4oahao3OR7ACJ6fnRwdT85IzJmvOTEvI1BK97cej
wDTGz4yEObvZKdjXQiLzecEfcY2Sp50uPDrtjpt7jEvBt7vtuGC6tUynvVOGerXQ4Hw4HhXQpoX2w+Qt
WRTgfphIsijILUv2SuF6RaxNF2CQUCZgBQY/jo5OWyeJYtfW3s/trfTvTIhrS46//OwFsy+/egFK67D1
bQepYoahVXSvwI+FUiEqZenXIfHsmnQtl6sqg6Yi4TTdz8rukh3LnjOKX83sxgxZtmPYs/NMNeYXJUMO
hXA6uxkwrgLjEZQXt9t+jKos7+ZubbaSkvudjBgJL4mQayIXGToS10mprH73W1yhlX4v5SrphiKKQ9QI
k0rXNpe1xULqJW8e7dbeyuzFjjg7r98NXVPbvhDXZilPjqbD1kH6Zu0mU76KqT3/zhVI/2vkTHAVJxJ5
oe0d15Qw4CzCnEpfbJDsY2SWfzI5PhoOzo5O35S/ud3duq9VBT4RM1ZSa7s24tef6e167Jjx5M4ynSz3
CfH+lNPIPuMiQUtniEIivcCFi09X+ZC+aRSjG4ZQgm8AxldZCIVaM+4Xp1l2YHpCXDNbr17E5nDFO7vY
IZlh6MJwiaWHovAZX4aJUAeCVsMwHie6qDUzueDc30NzKPic+c1sDB4enIzSqLR1iIkOXGgxTvGuGQfx
K6XlO5ydkzDBPeJ5IuF6NRsbeEZwex4a9WXVSaPBZax/HilrcoieFtIFFRO+chCi2V9vTfy8pWm/63PV
9qKU395hqrasDncZO7f9R5LbtIOScB+hObTXiocHyqS6v28+POx1Vu/vkVPb1bxCpSXj/nIFfkhQLpr/
x8WtkFRBGRZS26o+EkqR5qbZwgVaFYk96an4mAtPJV5Q7EvNogqHmrBQuaCZDjHXtRS3ZVnnHddkFmJz
WUv6Av8DPRN0YX6ly4VecwOi1jStUAESinJ9PYs4ZxjSInya2L6V05GLQOL8w55TL5VlhqpUlvNxroqA
RGvJZom5oAoeiUShuCkORFNPqFGqUgnFRpTo411ctgAQ6SsXnA8SdSI5zBIW0jMzzcvLtRdrzZev1l9s
XF6q/9olXF43F99tktnKrlIm9x6ZqIn5tBl+MsrLtFg3Tcpk8cmIzFZT/9JaladAxS0PBaHfL/UyoFLh
KqNP6aJin+tVqmnTE2ES8a7LddDwAhbStc31Yvcgrfa+xrVTct22riGiubr9hW857Y71pUTXlzsXXKd7
lKwXZomRuEEXyNM0o+KQ6Tq1XIBzyZ0NaMNV/W6KQ+JhraczxkSjswHOmaBk4fz9AIeozG3HhPgpfayE
8SWZPdbJ2u7nErkxXcY5SeHC3yMKK81xWORfKOntXTo6wAhVi0XEN//PHPfbzZj7l86V44LTbvadJ/nN
JWIqgbJz1fWlGevYsST+RsHm/rCk/wgAAP//211qFrgPAAA=
`,
	},

	"/definitions/torrentleech.yml": {
		local:   "definitions/torrentleech.yml",
		size:    2314,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5yWbW/iRhDH3/tTjJxTdaixHSCPK/WkC1Ef1F4vLdG1UoiqxR7sVewd3+46hCK+e+XF
awPHFXR+gbW/edz5LwtBEHgAWhhkYEgplCZHjDMPQPICGTys4W8NzLlMK54iA5RBpWsi5LNmHgBAAJkx
pWZRNJ/Pw810Iak08jyAmJeNc8wNpqQENmuAawbwgV4E6mh8B81zAiNeNA43+x0exo29f8b22ccxokTl
nPobTo4N9rBhy35u2fkeduEYQFcTbulVo9GN0+CGwcOn6I7iqkBpuFo4w6U1tFUHV/Uatp8v0g0HNqrr
4YrB/Sj6iRfoXPrXDEYkNeUY/T2lV4dvtvHw8sxV7neW+/HQ0cEmvXd0I/lfQriuzjr6+10ryiWD91Ui
KPokEiTn22+oW5/bHb2XokCHLvaPYsSVIZJuoxesng49t6Me2ln8Mv7owLkFH3jswIUF9xlJDD6arD0a
w3XoWcIXnkUFJd3x1MhVnDF4/PzUEPMSdPC0dtAkTwHLpzo8p1TIdXDJTcYgqjSqiMcxVdJE1hxZ84xU
weynXQpZVqYtWwetv4r+cgnhiORMpKGjsFr5jWPJtZ6TSrYdHW0dUSlSLnkAGnOMDSkG/olt6UdSxW1O
8TOElQi04QYDG+PqABSodX0HeJ0ux6eB0t+IU1jQCzLQJZenoI0imVqzQW1YuzM7vuZC0dFU0Vyjt1s5
LLCYoppy9Q/PTS2B4kbQlgTebrPrmFvebG8mcoNKd/NZj15hiq9l2zdXqWbg/2nTh28nk+T7ySSsXz2/
rtucio3C/m730XIpZhD+UaFahL/iotZIr1ZCJvgafa5pVKu4a18uUSZOyu2D8kbxudVecZkihKP2gl2t
mst28fj0w3IZrlbfbeWpG2JfjNPwaY4nTdt2Ae/ATClZ1G/VjAvzpI11VbqDsZEtYdJkQZyJPHnb78E7
4N04jVFiWtW/QZnCWct3xPgfQTZEmUwiK0jvjTtnRpgcD/Y06EFoPdvGEjRc5PobIr+2pZiK+vo/nPL8
mAElNJc58eRgtmHvYC4t/j08ooueq8zNMQM9Rkhd5sLs0fERfCAJ/ikEfXhqDycmW3m+UvrKlbZ/Po6J
uO79FwAA//8/6UIoCgkAAA==
`,
	},

	"/definitions/torrentsectorcrew.yml": {
		local:   "definitions/torrentsectorcrew.yml",
		size:    4276,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6xW224btxZ9D5B/2Bgd5MSxdb96DnwAx0IuaB27lu0IsF2DGm7PEB6RE5KyrRr+rb70
oQ/9oP5CQUpDzqiKqxR50Yy41t577Qs5/PO336vV6ssXAIppDEELKZFrhZEWMpJ4byBOphjC6QKCkcXg
YAlSVJFkmWaChxDsw3uUU8IhRo6SpKAliW5RBoaaEh7PSIwhUKxStEuM36rQvAFUIdE6U2G9rlW0tKtx
1PWXLwwhIlnOjIjGWEiG+QJABfazTOX/et0Qjg/qx4ngWN3nVApGDWXxlrMG7QLr49EIKsCORjnabPQt
3KBkDhX4kfHZQ451BhY6JBFU4HD/wNvsGgAqMBJcaRZjjrRaRWefGV8ktZQ+o0w4ZidcLNTt70SIW6gA
eSvErfJxchJUYMhUJP7rsPZgnf2HP36VKmOYuhJ5B4czxW5dUVoeODVNcI5bu7njc0bR4PapipkMxe3M
iWy2Qzg9rw9FNJsi10Sa1NvDHO/31sAfPDxYAx8XBfW7axijoY/fWYevNKbfX0M641rMogRpMbd3LJ06
s24nhENxx1DVPwyhAs3GoJE50GHtYSnlbrds1W95o2azmYN2VJkP1mk7s7fp7MSO0OLF9abhGMNz43l4
PjxxngvoOyGRxRwq8JHrmo+9W4h9glM/6s1Ws6z5zLeo3StDD61ex9k5j/WRwcZ3bFis5gglQ+4FtGwb
bNqr6fd2LWZDFOZjYbE+oZ4dntWB6S0mzuopzEkjJ5+eV0eJuC9N9Oq8NJv9EI50gtJuvTuUCv0WGYRg
NypUAEs7trNQO8qE1MareXoBbe/ydKaFZCR1lrse+0zSNCMZSpdQJ4TxeGzKOx6XVJutXtDcXX9etfsh
HAiuRIr1T0NzAH5iXCOn7jTqWMP3ZGon49hZdgfe8njUMdho7OP1vOjVAnYb3nI8EQ9G/Nsjq35BmQpa
ONgVEhklIVx8ucqX9F3Vr+4YhhJ8BzC7WrhQqDXjsf+uLD5gGXPjpueZ+djhg+tBSiaYhnCcc3KjmUJp
3v7R8qxEdDGJUvdC0rL56moe3C3bT6OIGV+mkBGdhOb7SRmPa1myPDSmqBNBQ7gRcrpYMW/F/4xnM+2L
mWcTQvD4CLUDwW9YXMtX4ekpyJm5xDIzXy0xGV8hMe5xlFJIF78KClN7fwghqNgMryfi4dpcIpYGGpUO
vQqT93TuUwYwT0k0E6Xa1IskH0TTGid37FqLLIwE14Rx9ToYIuMIJ9ZJsLWsHEs1SlWQuiiUxCwlkes/
AJGxCuGi5ASCHQiCq41ta9/I3zH8WuDm2w5/Mf2JFPcKfQVW+v4fSe5tkyThMULtwF2gnp6ix8fa09Ne
89XjI3JaaGy+x2xvf5qhnNd+wLlpvyr2n/EopUjMpDTdopAU5WQeQkAoRep9CqlDe2dcrBjZha3uZiNp
wXYG2zCRsA2aTFKs2d9FnvB/0BNB5+YpL5Sep7h3GSTI4kSH0OlmD/+7DK52QEvf9MvgMzKN8jLYMmbU
/Pw7v0E+MJhSL14znWLoW+eTIRcWu7pIJN78vBdQ1ISlyjTLjwAA0Vqyyczewg0/R5aX3flXfBunb/aC
iOi9r3gzDA+sznlx/r6YJistGY+LcD6HEXGHXiSm5raknhf1SgtD3EQYFfc8FYQ+49CUbsm6Vip9pn5F
xzcsxfUyNQ25TqpRwlL6ur2VU2JJJhsY7G5tVtO/7Wko7OvqlKTlo0CxX9ZPUSl473sEXzmENrbzh1Gu
EWlZwFdU9zdWHeNDti548Prykm5vueMkRYySjWIPvntsSvQGneoW4kqcijs0h/VmWkyEjEi1tg9Bo1Vr
NGutRqPX7IaNTtjoBqub6U6ksyneEKvLO4mIwlLA4EZwXTwqL454Oj/LjI+ry2ArMOGC5y2qnQYMl3Gd
Ta1XtnoTlL4Ss+ybVC6t/woAAP//FDaxcbQQAAA=
`,
	},

	"/definitions/torrentsyndikat.yml": {
		local:   "definitions/torrentsyndikat.yml",
		size:    3133,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6yWa08bORfH31fqdzgKjx61EsmQhHKxxFZACpW2abMNpawQWznjw4wVx57aHmgWsZ99
5Zn4EhYqtFpekMw5v/M/F3vsdLvdly8ADLdIwCqtUVqzlIzPqXUOSRdI4Kx1dKeJh6HJNa8sV5JA5xBO
US+ohAIlairAaprPUXccKqgsalogAYZdhs6EMleMy4LAl7OT7l5DcTk3xH0D6EJpbWVIlq1q6vqiekoX
2csXDstp5fmcWiyU5ugNAAMCMDmGDTisKgMZfOWSqVvj3f0hWXN/4LL+4Z3bTWw2pnkExjT37p0H0mNu
gq8/cLrZKV2ggQ2YVhwFQgaTY0/sEYBjJY0SmE2mkzVoeuFMntxNyK+cp+RXzj01HJBAXczUjxS7OPp0
EZrqR+7jaJpiH7m0KJmC0TSMb0BgrG44mmw4gg044WLh0OEo5CXgiSNRf6bLhGoNYSb9oPXJlqgT8PO7
8ZdYYcz5Ps056O9sVR7aj2nXoP7WXoQGW48r7Q4SJlY1Ok+h0XnosR91pikyjVN4urdDyRcYuN3AnSiN
vJAJ+U4Wgpu8DGl3CJydt6Wfnbt8qDnK7H0s7E1DTB8SsbDt7QcaE5rPTfYlagyGjxIJsPsgSQvEHIO9
BhipvF6gtFQvPTlS83qN3P85GZMOt9qkldI2NOceQmMt0EzXA+ujbqcXx9wgD2c82CZwWDOusg/KGIHG
va6NATI4+XAY3tfBGw+OJ8OEGU+GIWNA/B7wUPMclHZWWAI0Iw2ruud1mv8zpeYJSo+Uiuxw6NlzzjCV
bJ6j5i6BJhA2moMKMsA1oTftYTcplcSur39FjtWMizjYfQIPgCMuWC2L9kAGWCiWHMEGqc5LApffr7zJ
3nSjddMRRslNwOqqlTBoLZdFvAfaCyhXas5DHXZZudsKf4Q9IegMhTveWqy5TlTB5UpmgbZUbF2Gy6q2
sdbWRaBzdwe9YyWvedFrbXB/32kpi8aGgIrakkBW51WvKqu3NHc34QHDa1oL25agqeVqFbDCZ1rdGnQR
rdmgwNwqTYBfWm4FHnRO0dCFbWI7V/AXmIpKP5xmck8IvmXqoCUebfB/mt427WkqC4Tecbgy7+/zu7ve
/f1B//93dyhZ6DeuYDOV32rUy96vuLxVmhn4B8YlgaYHb+YyFwwpI9D3Jo3iW7t8W37PuAaSLePnYelM
YG91/X9rnuAXsDPFlu5TE6nsK1JS88qyXq5EiZS9ft3qXHMULIo2RYWnNAlth351WWq8/uOgw9BSLowb
Z+cqBlBrNZ/VzW+ktMHVz47lE9or0WSJcmoPnhB2cHRcc2FRm0Q4vgzf3ToYq7l77yD5o7owxBUVN/XC
nbbmP+o9LZGpWykUZT9v3VPPFDX8z8fXyTIibdnNSy7Yq53XcbcgW5/SEyF7IUQg5uWzYvZDDKP28bLc
qxm26Gp+kdO4UDdIgG8Cfd7KaqwEzfGxVb3svMfaYmcTOmeK0WU6zWcKnKKxqKWT+L35+u9k4EY1Gp2r
vwMAAP//AjfJAD0MAAA=
`,
	},

	"/definitions/torrentz2.yml": {
		local:   "definitions/torrentz2.yml",
		size:    1893,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5yU247cRBCG7/0UJXORXRGvQxYk1BJCYTlJsFKQBjRSwkW5u8Zutk90l+ewUZ6MO54M
2WN7xuCwo9xY6r++v/vvctlFUWQASTMJYB8jOX58mQE4tCRgdaYYdHWLNQkgV7QpA+BDIAGhrYyWGQA5
6ZV2tYBfV98XX3YW7R6SyAAACmiYQxJlOZ1yQ22ZZQASwwBJZKp91DSsAfKtVuSBt7mA1W8zMRdw77ea
0hy1nQaNmsrlj9+OBKrWcC5gvV6PElXeP4D0VsuUC/jG+4dU3vXLGTLVvutWcze37KNGMzErko3TEs2c
s1jjo3Y0cfeDMB2FrdIebJu0BBtucwGvOqW8f3073+oIDsGOTP88D3e+mfEpGUppon8ehAkOwWiJrL3L
Bby+K18oPIzFGi0d1R/QntJ2Muwrv88F3HmXvKFyXfn9rN1gtdnMe37UG9TxsFRA17fyP3pl/O4PX6V5
rQ8R5EK8/v7jhbNetV6dZisRRtkIePPn74PC2+IkPu+A5N1zoDAC/XAV58bu4yFm7eok4M2w7su9JSA3
YlBe9Yp2oeUpA2wE5O/ewc0vLcXDzU902PmoErx/n/dE9LuzvIYk+yggb9gaqLw6wCe7iAFuIqXWcAJl
RIPpCq/zwbTRhimezoNi+LTRKYssm+xIkVETxJoNnRyncxUDDrIiRm3SExQAMkddtd3vpYm0Gd1+54xH
9ZH2hWudLhappn04KwBg7N5OXl69fbv7dGrNuSlECuTUksti7YjF13v+qo1OVKwbMW4x/LAOy/eYxEjW
b7ueX5Kfo7ZLOf7+ayn4h2gYYYW8/C4VpIBOOG4K2Wijrl5eXxIPwwca9Qyw9s+WImpLWPtxivXjhXlu
r6fBJzVL83+uz0eXIZLNxbYvrv81mltvWksb7A3TFkx7FpC/GLvbhifZz/LsnwAAAP//0pXSYGUHAAA=
`,
	},

	"/definitions/totheglory.yml": {
		local:   "definitions/totheglory.yml",
		size:    4102,
		modtime: 1487335360,
		compressed: `
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
C/xX8Zeq/xsAAP//m5APZAYQAAA=
`,
	},

	"/definitions/transmithenet.yml": {
		local:   "definitions/transmithenet.yml",
		size:    2405,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6RV224bNxB9L9B/GDBFYDXetXVBWhBQjKJBUKB9KWrkRVKD0XK0S3iX3JIjKYrhL+tD
P6m/UHDvchXLTV8s88yZ4ZnhIffvP/+KoujrrwC8ZpLADo0vNGdkiANssCAJtx0cN7ginzhdsrZGgvgB
bt+H3OSOnAjhHE26xZQkkIm2voK0ufMy/AcQQcZcenl1xY8qh3iCZUtMkCm1TlMLAIwl3L5vF5OwuPrp
bbueVuvf3taFAAqrBqme0CWZhMUfqxbiXdSjl4HhrbkEKld1idym2jQFSuRMwlUFxWVWNlsQZ1ZJKK3n
GtlYV8jqb73Wptxyr2LrydVjFff3EP9ozUancYvCw4NomSV6v7dOHTNbdMi8Iypzm6YUuOMOrtWD+CX8
Nig5Z10nJgJPOSVsnQRfoon36Iw2aR1n8ix7NVX7bJ0jw76eQAg6ZG1lm/GRJYgoaKhsVQ/3aIDHFU5M
qE5qSoXOf92SO8Q/0yH07YeNW6fIfVgfJAjWBT3C9xgCwatdAJPGs6h2aBJSovWKs/uhhHYqjOucXjSa
P1QreAO8tuoQfl3chJqz15SrvgxrzqlbDcviQiFj5F2y6sPI7PR6G25iG+2DG50zOT+oFo6vdpKjlD6W
wwgAutRLEBfxq5vRxXIZL5ZLtVzuV/fTy9nD6OabbiTDu3xKqtK7mDH1bay5lIdnkQPd05FoXGSONnMx
9MFNfS7zNXqdvGRMc+15nimxkiAm4svSfZ0+PU6XiTWM2vgLMZtOrksx+vcWA85k/PosZzybnedcf3+W
8915ObO2yue7mr4+S5nMBpQXPUd8K6oHBF6AtwWFK+XB7nbknFbkAQ1g7gjVAQrkJCPVuaFzhy2KcCxf
5Ppwsp0t7d7kFtVnCgXq76dd0Ga+1Gouzm+00TmdlssK3gRHDyb3riKL0f+/l8ulejXqTsjrT08+FJVI
/YmeeC1aSv+GkToWNuxMGs6iHD1HSaZzdTHpWsqJkuz5ieMuUSGf7qH6sgQznRRfPZJtJHW4fu7O09Fj
r+xsvi1og1XGE0+QLtIF5jwX7xxR1W/1UlwfXZj2MnQf7vK/bNFn/xMAAP//vvUwt2UJAAA=
`,
	},

	"/definitions/tspate.yml": {
		local:   "definitions/tspate.yml",
		size:    4194,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5SXb0/jRhPA359032FkHj167tE5sfMfS1CF5Cj0jl5KcoCEaLXYE2eFvevbXRMo4pP1
RT9Sv0K1dry2g6HpGxzP/ObPzuyOl7/++NO27ffvACRV6IGSCVGo3xmJ0YOFeQ9Q+oIminLmgTWGY4HM
X0GIgqEgEShB/DsUlmYjwsKUhOjBUthLkYkou5Oe/gVgw0qpRHrt9nq9buUhWzG237/Tep8kBegThSEX
FAsBQN/xYHHRHjMaI+xB9vRgwVNVED3Hg9lE65JEenBJWcDX0mhdrW2fEd8QZ8Q32k6mna04Q3vMAsFp
YLjNu2G7Ffb069xwpwkP2qeJlhu2V8tpnCqBZUr9mvILZelDoesMPTji/E7CHnyys1+auK+Yd0YNyDgN
KDfE/oZoT3hM/RqYSwqy6zT4OiMsJCXSaUCOglLvNuhnAqU01Rh0GxCcfTsqANeDM35PUbbnU9iDYxrF
0oPpxfScJmZRbzBw8XW+OD43KRl0elFnbcP0DHNSQYYdxwTsNxKuMyqRgUGOovScPJbY5v2hM+gBKv3s
F0ZDY9SdbhvY3ampidsY/6hak/3mhZ6bWK7zWtXmE1OK/qAxUnxS5tJpzuW8kovbbWQuPx1VoP7oDahe
/f23yFoXRo1rnJAY2rCYm141r+HbKujdmY3T1weESR5h++qWP8Ae/ITpgwdXR1+vDDQooUtKDXN5emqQ
YYnM5jODzOYzg4xK5Ofp3CBTk3B3v3E25djWbOr2qgG7lYBdU9B+ZRDm2sog7ObD7UcSoyytJ+a0DDYj
pn02097PUkm/p6iHxcxE6A0L6AuXMkIpa+Txl3Hpb7RBa8Tl2dgA+01AfY66vezTkDV9jtlXA06miwuj
72t9Rbm42BoU7qCZMGm4w5fAbDz5bPSjLIWT7RTqG7Q/fEktLk6mW9l0nKZsTKj9V0JVT82goST1Ezhw
XrppOH4D91WstrSO+zJgNmeh7EOnqPICI7ynMrtRlL8N1mvA5gn6ivhR5fPXzSMmXKgXMBeqBPu1m0MV
zGQl2MnAKffTGJki4nELNxpa2X1914OrqyvY03+zG4nM7zMAMQ8qNxiJRPgrD66/3xQidW+X0o+akJx9
BExuchcRDynbOEiIWnlAfJ+nTNmZppWsNg2IUa144EHC5eY+RFmSqjJ4KlHkNzvr6QlaE86WNGwVUnh+
tgoyIVKuuQjqZCGtkgJjjG9ReGA9oqyIVSqY4h5Y7Y0QheDCXAIlRugrLjwI6P0eW/JbIuAQkr1Yhp80
mYMKpfLKrPTqFRcCmZKbquXr14ggivJaoRrRvA3V6DQP78AhpBEcQkQ9pla2v6JR8D/3AxzCkjPlRUSq
XJjH23Rtp4BbnfiPIOustoKwEKE1MXfd52f/6an1/Hzg/vfpCVlQKXWxS7KW/JKieGx9xkfdEFntCGV+
FCDRvXOtYhcKvq5uwmL1itxG2FLZ47cVkoAyhroP6pYHj/opWsoWfF2YkqXSvXbz9yXFKCj9BnzNIk4C
I6iGItcrgctfD6yC0sX5gQYH1k2JE6UEvU31PyOaNkeEqgjf9moqH6AiNJL/4D3zWGg2/2k87hYh8+wT
tUviukSRQiErrvXuz0/hd91DqQRlYVUNQEQoPZ2WyZDHeubItzMsKJ3hDslJ+ntzUVVQ2f/DDwUSCnLb
nEHNYPSh3GcY1Nf+ism+MYkQ/dVONq5jjAKimhciYxLpA01364fAJCI+NvXi2prqIGB9BKta2h2NI3zD
UqefECEbbS237zk9z+mD07Ed1+44zsDaPnD3PEpjXJJs1aUXn0isLZTG4bUU/oFFYxKibC8FYiukS+vG
A8uxqqj1f2szQPL3NPk3gQrrvwMAAP//eM/ou2IQAAA=
`,
	},

	"/definitions/uhdbits.yml": {
		local:   "definitions/uhdbits.yml",
		size:    3265,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/8xW327bthe+D5B3OD/+uiFuYzlx62YlkBZd2mHAMGwD1t5YnkGLxxIRidTII/+Z4du9
wB5nF3uXvcBeYaAlWbLhDkYcDLtJxHO+73wnHw/J/PX7H91u9/wMwClCDkUiJ4qcD2iRIYcPX7/7sgpI
dJFVOSmjObC3MFNIHuQQYtRoRQpkRXSPlnl4KnRciBg5oO4Wmwqp0veO+y+ALiREueO9XqUZGBv3zs98
NhJ5DYsEYWyswjoAcM3hWzNT6OpAn8PbQipTr59z+PFjvXixlxxw+P6uXrzk8B0laEtZgMzIlpBDYaOE
w/DnUR2iWbeJXnqEM/oSMB+VJVITK10VyAUlHHqbUJAneRmdGpvxzc9AFJSM/VeZUTovqFEvHNpyD9hq
BcGd0VMVB3UU1mtWI3Ph3NxYuYuso23kPWKemjhGj72uwmitsdtdcZhiRMZyaPqD1xDMhdVKx3WlDJ3z
m1uvoUVkB5mVGqEj3nS+sYhMfn1VWuQzVpAyuya2ELtaLhc6sK9elcxqb/ao1qIm17D3jH5ixXxjnRU6
RgjutiO3Xk9VSmjHkaDhahWs16Pb689XK9SyZWop6shW/v9QoF0G3+DSu+/a9pOI3ZiWObbcBzBWoh1P
lhwYqQz34nPhE/7sbRMiqs6gkDOhI5T7rRSTTFGlUaasmbfnunaPxCTF/1cOjTcreA00MXLpf9ugSlWj
qzCVTRlp5jo1Qh6cATEkRSnesncVio0amCCyalL4CyexOG0Spd2uVdGPZHkKLOapiLCdAhA2dhyGLCR2
CQwakY36wc6kmgWxNUU+VnpqGoTFzMywnKjLDcpv1yP0Bj34yiL6BtsmHM/+qN4/nPzBoX04+/0iSgun
Zvi/h9f489ff+oPPTuIPrk7j35yi31/ASfr9xYn731/scuuDNzNpkeFUbOa6YUfC4c6UMkfW6DiY+Pd3
nCGPjCahtLsIme8sZB3GgV2xY0mDq5oSDI4m9Qdb0s3xrJuG1d9jPWU712iRP5Yn/UWl2f9Hwepfk+Un
7j9/s/10y9rPz5vmOXmE6zDGRX5oZlgYDi/CUD570wnDUdOuyTLfyOF265t+7Kv/qze1RBIq/c+1dZzX
F8HTTv1+hvLZk+0Tq345/PqQ5JqSbpSoVF4MOjUktmJy2IIdwsstYapSPILwvNM8+ih3nfkE5WZLSRGj
5CjOF1uOFHTEH/6ic9yWiTxHLQ8OOYjYsPOzvwMAAP//hzm/3MEMAAA=
`,
	},

	"/definitions/worldofp2p.yml": {
		local:   "definitions/worldofp2p.yml",
		size:    3767,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/5yWbW/bNhCAvxfofyCUYUi7yrYsv0VAVjRJuw6LES/J4gBtBtDiWSIikSpJxc6C/PeB
tERKqdJ2+xJFd8+9n2j6vv/yBUKSKojQhouM8HUxLLSM4RwitNSys/ViuNAyAjIWtFCUswh571ACDATO
kBI4vgXhaSbDLClxAhEC5pdSi4DFnFCWROivyw/+zFCU3cpI/4eQj1KlChn1+y6DHgP18oXWx7iowRgr
SLigUAsQOojQ5VX/HaM5oD1knrUqGEdocdwfEHyvVUUh/VPKyq3VT4x+juNaPccxZYrL1CJTgyxSzsA/
UykIi/IVzVyoCD0NtaSM8I2sidFBhN6VhHKtN89L3TOrHweVvm/+rji/rckjzm8tNoiQfpdoD73XkHUw
DCN0zJnkGfTrTH/DOUh/TrdAaiwcWswCi2YiQ1OxUTjg2BbarUbntHAeBs9lsrgIamj2PDO0jgIHLS7C
BhJaZNhERg1kVCPTJrFoEAtb1MghS0otsqTUjq+BXK/4NpwMLKbfUTgZWHbcZtvgGbNbMwo7wa0B7cAG
EZrzOwqy//EE7VX/+8FgNijsVkw6mWEwccxwZJmwwYQn1klogYsGMGrGGXYi4+nEhRl3pjIdOi9BYJGj
rDw3n0yFHWWlwPfWV3dZR+fn1MWbdqZ0jHMQ2EZ0PTy5alAnVyfn1tHMMvVCVtRZRijY7Q6D57jWZzSe
duZ+Drk7gMLwOVcXBRfKZnZgseX7o5NThy1h1T85te46O3F9R+2Ig7A+ZE65lBlI/Q3PS0lj/0OGY7uX
NTVfhBaYL+xHF8zcWbZTtkoP6qOuf0UJOMi8uT5OzcltmnN59WSfRy1da4/HY6O7qHTN9QwPWmbNnQtn
LVV7z8K2y+ZWhJOvVG73RoFW7jStDoyGxmr5/sg347q8ejKqcTud7XBij6txhJoRm+MbTSJ0fX2N9tD2
eut/dPKpkfd/z3ECElQF7F4tNNtBJtGKaCW9c97XASv1LvZOn3PS+NWVgEWcRujTl5tapO58J32jCcnZ
GwTFzc5FxhPKKgcFVmmE+grfghH3irRqag4q5SRCBZfV/q+5yCPzd/dOWVEql0kpQexuKt7DA+odc7am
Sa+WosdHryYLLOWGC9Ima2mT3KWKvFP9rKQgBBc2qo8kZBArLiKkSE8qksuk+t1SIFXkgppKdT5xoct8
i2N9ezoksMZlVl1xqsa1mrMSfCPBdeZJ3T8JvDGVCMwSQL1jezV6fIwfHnqPj4fBzw8PwEijsHpApgF/
liDue3/AvS5fNuunLM4IYN2p4Imt6YuiKgOv3gydZ2MxbFvwKoPerooNJSoNBgP0K1IrTu71U0Qplvv4
Uypg/fehR/iGZRwT0yPFhQCmDr2bV7XjNc0UCBdIz2A3d8xIjlVc3dnWFDLiMJNqw8jl5yKDwjSTJjAl
h57d6Oq2ef9tczentzFWDXOEsFKCrkp9tdawU3xdjCvnix6LVIKypKlGCItERjopmx/Pc2BK/s/yns2v
HsV3/HYP7Lve1zSD7pQViZhK/TilGdkf28lL+k/3BFv8zPEApN3cZ0yCgbXJAOL0x4wCa0Sw+oHEppZP
BF59Z1iS6U2Gzml9a2kEJLAtuvbF2//8mfzyyns63DuelTmssUnAGcZYQisC7lG25lHMmcKUyX3vgwDw
XkXIG3hNjubJJyni14defy0AegVLvJuvKO+11zpUyuK/5OKs/w0AAP//GdUzt7cOAAA=
`,
	},

	"/definitions/xthor.yml": {
		local:   "definitions/xthor.yml",
		size:    4105,
		modtime: 1487335360,
		compressed: `
H4sIAAAAAAAA/6SXb0/byBPHn+dVjPz76QS9Js4f/q5KUSGl1bW0HEEpEsdJG3ti72HvurvrQEC8IF4H
b+zkTex4jcHcnR+AdjyfnZnvzFqbdrvdAlBMI4EbHQrZAuA0RgLny1VEeZDSAAlMZXtqLIxfKdICAGhD
qHWiiOsauDO5dVstAI8mSwePagyEZLhc58//4IhFsbJsPQLHYsZQuQdRekrnlie4sLQenVpUv4E6xTi9
sYhBQQyGTzICFz5yT/gIg6FFbRTU51qq193pJnDT39qwi+p2a8AabNPCNhuibffrgm0X1KiWGg2fMlsN
zPmM+RawUwDD8bAGGI5t4XYL5crWlf+PjwfDrzax2UicshNb5j6Bs7E7FF4aI9dUzitJ5S+YRHsYei/H
GqJSjCv4wFn8+GAP7GBggn7XIcpauRP0NPUiO+JuoV8ZXFHj76Ozo9NW5biMHh+yU1SZrY0sAbCfwrc4
MjY0qGEs6IR6VzA+sqlNU2p5GCvU52H1ZPZMdtZYVZizcZUZLPqYqY31zEsN6XV7NXi1tgVbGYNlromQ
uj6ueWUz3SYlP8ZMKSY4nI3tEdipH5xqDxaTYFW43UR+HtZyW03c2fiZyfvKZrIyeP0NAgdCXCmoc80+
uiKm3GZ2t5aMe0wDesu4OXElpjDb2G5DqAPKfVTgm7F4fHhyRPp52EMRM09V8YW1wgxqYlqZ8oCqJzqJ
gHkMo4pUmwRODt1PNC5/YUre4MJvmN7AyaHNbRM4FFyJCOujZEMS0bnSVDPBbXarmT2fiHMb2ikg9wdj
tdAPxmxmd8V8G45qmW/Dkc30MkGg8tjMhySJmHoiSd9IeUy9RvKYejY6MOhJKDi2i0Ngo6OYSp1kHlln
Y+HndxaFVHohgYufl2atZ+2V6W32Wgn+FjC5zMBIBIwvwITqkICr6RUaaycJE/MiRh0Kn0Ai1OKTMhUy
JuavWTKepLq4MaUK5eJO5tzdQedQ8CkLOrkV7u+dpWNClboW0rcdc2vJUaWTmGmz43luk6hTybUg4LgL
m0alSbG1KSUL6iWmjlYhTKnWiRTXCos67Tr+L+m1SU1SHiB0Dou74f29d3fXub/f6/1yd4fcL2VqQjBO
wNFMR5jbGfciH6lPoLu0BFKkyWqZd8go8XuKct75gvNMB1UIkSVLCv8IPS0kAZ/NOlRq5kUI70HTSYSE
67DthSzy13q99cw6Ef48+y9JSNUavQglTv/cc3xxzSNB/UyBfS2kRK73nMv1RZMZRn4R0JSzuhSv4q82
Q01ZpMxezN9zLpfOyzv1/EV41Yl9j+oVDEC1lmySZjf+zLewT1mkUVr39Pbyt8DPTD+lJeOBdayoDBTJ
8skzE3F2yVL/qqznMss1fXnTeuEbtlbstrYF2i91fGv9NRJJTCLqYY08F84nJpy34ASTUkKvI48XZPzP
yS8L8mpFKkTfSv6ZenfyeiNEL3wNsZsTPtXNevazA0TrOiN4LFKFYobyv0r+R9rt0m6mwDPKBXiT1JDO
h79Eqh8fIEIC79z0/do+ebe/vtb5dX/9nTt571SGciaiNMYpNXUW+3lUYTljFgcXSnpv9hx3KhGLYQ3Y
1Lkk4HSdkrPzxiHg9HJTmrw+TgUNJJ00d2/b9OPvAAAA//+CZOBvCRAAAA==
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
