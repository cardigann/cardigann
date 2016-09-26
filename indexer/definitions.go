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
		size:    1875,
		modtime: 1474593344,
		compressed: `
H4sIAAAJbogA/6xUbU/bMBD+3l9hwTRRjSTjHSIxicG0SVO1DiZWiVaVm1wTi9TObIfCqvz3XZzYSVkR
+zA+lNxzz709Z9vzvB4himkICc3ylEqqmUCI0wVCFxV03UAZ5UlBE4SBe4WqEMbvVYgfhHgk1TpXYRC0
afwoCnrojWjesCKqIRGSQWMTshcS8uM2uLlq7P3a/mLtA+sn9m8b7avbq2uWN5TDjZQhje5d1iOb9S+K
K3SMlIF4wNbaZk5azPFOu7wmlbHXCp51I5+zXK699yEZjUaEdNoyLAQtBQW6KGImglsWg6gp+GndqNfw
MviM21IWQskuBVcig2A0E48WPmzh4c2BRY9a9CdjFj2uslrjxJQY0MgCZwYYpoKD902nILGnizxXAzFj
GdhF1rNZC8f4KMS9Cj7N8J9F9+1w5rfrOWg81sbuTalgwFTUM+hCxO05UkBllIbk7tekQfSD14K7FUEJ
vksgn1ThmUgYr4NzqpETGMTP0/pYzYVchObXmIznhXbFCgWyviBbqxXxUcA5S3yLkrLcaog5VWopZLxO
tKgjgpRChi5mQzs4LShVXb6ePSoKMog0xhF/SSVnPAm50Dt+yuIYeN/wNCj9LC/jeCZM3p5TrSsDZpTA
tXKl10cnbyRdmnEkPgeAM7krXZZzlmmQU7zld5Pz1covy7erFfC4LOtCSsvzSofvBcgn/ys8VSooJ4MU
y84+7XSazjLYbtqaGoto6TdAvSwGWexCm0fmaZNUOvbRraaRyAh1fqq1ZLOiegRTCXOH1/O0s5PqnasX
LyGBx7zjwCwyUahMK8J4fDcex+/G48n5TvXRt8dCM51t3CS2N2PJlPG5IH4iRZHX3x8IDblOvShlWbyz
329CY9CUZeo/ZHpJg1gseSZo/FoJlVO+lnnvtcyK/X5Bgk6WYzcp7vRV9lGf+JotYFNlI7k7WxCvrfWF
fKe2egYQpf8Scdbv9f4EAAD//9B7ZyJTBwAA
`,
	},

	"/definitions/beyondhd.yml": {
		local:   "definitions/beyondhd.yml",
		size:    2333,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/5SVXW/6NhTG7/kUFrsZGiFQwlvu1tKJ6V+kijJWqaomkxyIVcdObac07frd57w4gUED
/7vkye95znGO41iW1UBIEgUuWkPCmR/4WmA41MJ1JsymWqCYbWO81SIwK5apQtiLdPUFQhYKlIqka9sm
ovOefNjfPQvBbqRFQSnCtmVGXtPj/IVAJiGkkkhLCt5VIVC8BuqimxzSooejIsDDCrZcECjuEeqPXDTn
b1qxr2m8wInWfskVZCMtWQInBTvqlWx/mkt7bH9aYOP+KewPgefwsOSLChxVtWfHeb3uuBvZpGAnTh3r
/DCR4zpsdNWNCnBQmze/dwquV9vjAsL43SR2T5F/MgWCYarh6g1kSzMdDy627bXvTC52LW7nfz0WtmHv
jG02fTwwjIdnDA83zsL0pF/AclVxObtcoYeA7442k+PU07fM4z5IQ4/r6dl0uTItT+rRe+y9lLHDenZ/
vs6gnv379tqa3pnPSgf/HvuE23dcSgpSZtsmlsTTqNHMVCYGnt/3y2QDa83ssIHhVsQH/j/u8HsZDOvY
/Z10VQeamWVwmA7EHBwSsPACFz29PpuD6M2qxHYKSM7aCKLn1E75lrDcHIIKuH9wjBEWxarMzh+4qPn5
iTo3nG3ItpNr6OurmUEKpDJ4hJWuaYf6jerTV3aiIMrPzqyZRoWsBd9JyJ4fFzW9Z0V/QLLjwpdlvdRZ
kRQ8xYU+dfGaQkdfCmCquFn/44PChCIlXKYCywsI9X9lv121MvuGAPWrleYHcmLuD8L9PX+vhXDJYKUE
Wcfp/ygQsCn1DaH6w5RVWPXLeI1BJFLb2HbvqY4S+ueStmGGSBSFs904VTf5YuXPWL5bgM93jHLsn83q
n8+S5OP8MsatkhAQ8rf0/y5MM3oylyxKRpi5FEuVS0eBpF1l1k5IRpSoE7N5Qs1/m23UQ88nTHrp4QlP
s9Ms9yr4BwW/219d8y4ogBdcZOm1Go3GfwEAAP//QB5cFR0JAAA=
`,
	},

	"/definitions/bithdtv.yml": {
		local:   "definitions/bithdtv.yml",
		size:    1471,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/5xTT0/cPhC951OM+B0gEkl2+RVoI/UA5UBV9dIihIQ4eOPZxCKxgz3ZNEV89zqOk822
oK162c28eW/+O4qiAMAIwhRWggpOG2tLVln78vNNdH11c2sBjibToiahZArXqkJQa7gWeQFXuBZS9A5w
zJLJvGG5laOMGtMjQj6a1H4ARFAQ1SZNkrZtY5sv6hPGmaqSwBIyVntixghzpQV6G2CZgk2QXEhRoYdO
LPRVbSwpuSybb6zz+LuBeqWypkJJTI+O08HxvVaaPHRmoYuGC+Xt8ymmB96PhORWcBxpy0U6tOsMW9vd
3d1onfQu+A8MMmOnUrPs0QTOWSm+bci6dVakcP/04BHaRFvw2OuPAeuHXl6qXMhBXDOynMQhcV3UDlwr
XaXu15lC1g1NyRqDeljqwfMzxJ+UXIs8HlF4eTnwxJoZ0yrNd4kjOhFRa6XTSePKIfaIuyXZjtGY/hi8
2XddYkZWC8RWJcYciYkSYsIfNJE0VmpjazXUlXgMK+cgNPRbxqpzqYJpmLPp2BzaLt9MxewOZBy06/IL
dn13ZjYHe4ApLJyhVTuT7ZR/3wpOxcfD89PF4QPQSvEOSKdS0VH/R0WUFaLkR8swHHYksORTMH/k3avT
4TP5SQhs4jAiLVZN/2ALjesJX4uSUJttsP65DTt/alB3xspkPvPaUDo3aV/GeICCytd3Na/m/zeqcWrv
GNZq/jnWrDOuWlkqxvfGWoZQ741mxM/9HZ6Hf5ziSo/l2KXtDXAWTveCfGcpbwg+jIISMSv+RrFchEEQ
/AoAAP//71AKDb8FAAA=
`,
	},

	"/definitions/demonoid.yml": {
		local:   "definitions/demonoid.yml",
		size:    1816,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/6xUwW7bOBC9+ysIYwHbSGTFySa7EZAGbgr0UBQt0CIXwzUocSwRoUiFHEVxg/x7RzIl
W62D5NCL6fc4M3x8M2IQBAPGnESImIDcaCMFEZrnRHzYEQJcYmWB0uger7hOS55SMOigdDUj9Z2L6A9j
AcsQiygMq6qattWnRRUOaDvhhQ9LOEJqrASPGZtFjH02D8R44pSIeSmk8fiM8PdbD/4l8PUm/Eia2/jz
hvLggsAXzMB6/P9v+LIpFs61zKEVcBKx98bctfVmM4/DG5PLpKNJRyMrvJUCWnGz/1q2+Y0pb9Bs5Ubs
7uiA2ySL2OJ+6Rl8CHbkcR3gjD5mUCzrdGVSqbfJBUeKCXmSmFLjKuNaKLDTIiu2xwBmRkSsMA4bQuqi
xO5gLZO7bX+HT09semP0WqbT0oGtWfb8PPSBBXeuMlb0A1t2L9ACllZ7VZ77Vsa5xMivnqwkZvcRmzUQ
rDV254aCBAmztbH5QoqrURyvEjJsVROjJTtiGBuxoX2Ni0SRiquRBTFaRhqzIMmkEuPZxJejSXD1VHr4
1+vXd87NA3nocKPgmMV1i3z39nu0lsoPZb8J9yXYjTf2E2xqQ92eo/6ToIiT1qAy3pFzpbo6XEnsUQ5A
ADXt1GN4xLq1aseUUrSFrancnz2YJshjBeSPRqAJ02ZVcMHesYau18YqWm2kDY4XJCLVy0kDGr/IP/Sm
BfrobNIax9ckpp0AelboShlwAdYd6tVwW/hqpGCNo+WwCyFbsZdTPzbbobaQwmOxt0GH2tRRsR9zQcaw
8XVkNJtcj6dH15N/tjXXEpToynVGH5CEojcP5AHvojiilXFZv6WZhfVb1DZz4ChNpwckt0LaB0KiOjjT
PVGn+6ISk+fUwYPuvpz10lWEqbQyXLxa7fwt1Zz8+fptLib7Y314THoJl/XR9UfsAxVAkr0lcXbSZQ4G
vwIAAP//sIT1nhgHAAA=
`,
	},

	"/definitions/eztv.yml": {
		local:   "definitions/eztv.yml",
		size:    751,
		modtime: 1474349693,
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
		size:    1613,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/6xT0WrbMBR991eI7KEL1HGcdOsq2EPX0BW2wVhHCZQ+KPaNLaJIriQnZKX/vitbchxY
twX6Yqxzzzk6utKN4zgixHALlCy5AMGNRUCyNQLXCHxtAcFkUbMCQZBxbRzC5cpQ/CEkJqW1FU2S4DDS
KomwlLHKUzJmoVCag18TMjmj+P15l1xKvgYPpqkDL+ucq0B755BvaoPSZDoL6LiHfhL1D7YLFVcIlWul
gRfSl6b90uwueJ314ZuAphf0JZ+0L7jtBE2kKyWNEuE4jQf5fpV8xn6aQDynDRgCTw6WH7wkVFPfpZsZ
eUMMMKMkqVi2Cm6TqSfcvkBIO0KltA2gb75fNonIfD6PGmCt8v09oaXOSkruHx88YjfxHjz1e54SqB6c
XKiCy1ZcMYucpEFGVVk14FLpNW2+9yyzXMmPJ5atoCOdtNtwWdW2C1Eb0O2bHDw9kRF2ecmLUUDJ8/PA
EytmzFbp/JAY0B6xlkJlK6SlA5faH6gXe6HV1kAX+zBPOH+zyRfYOXPTc8fnTsm4WTibvUxAZpWmZIRf
DdJite0KB5F3PD8tu7D+k9KyhQAqbRlnJRf523RITMUkYZ2GWav5onajXWpYdjhOqQVt9uZugtvuPtag
dwZlsuhV0UoXhrpY4Qlwi5sfkW7yt3SNmy/kYBkX5tW8eyfP1RYvneXHeE//09vwX0f143yIIyBtp9Cw
Vhtwry6ExSdwjOH7A0MDkB9c8b8NLoaeLACy8khxOh5GUfQ7AAD//++sw7lNBgAA
`,
	},

	"/definitions/freshon.yml": {
		local:   "definitions/freshon.yml",
		size:    1262,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/6xTMW/cPAzd8ysIZ0mA2IdbNXzfULRLh6JokKUoAtmiLaGypIp0UjfIf6+ks32X4toM
7XBn8/HxiXyi67q+ACDDKKCPSNq7FDs5pvhdjj+427uEKKQumsDGOwHV7R18wmiQQLrZO/y/ShQr3TDJ
IRWiqyfKiHFfSaQXgBo0cyCx2y2nNPywu0ipToaF0knGwWfVQwywF5AOL8Ho1RHnh5pQxk4L+PztBtI7
eXcDGL5ksvWDcQdqkJw4u4I0QYcC9j6OovyX0Lgw8SY9EcbD9NXTEzRvvOvN0KwoPD9XCzFIokcf1Uvi
im5EjNHHVbxOrVrsOCGg96LzjqVxdFW9LazqeuGlcZEoW7kBcFJaXZaBIFQn6Yijf0htE88Wb6AtKUZi
sTVcvAjR98ZicSO7tRh5alcb/SPh5tdLg1bjy9QfJ4xz8x7nPDNtQ+f6I39tm2Vr8VKbQdv0Y/gPuPVq
zs8onOerpvNWo1QHG3qDVm0yy3LMR0cYv7OA/boShu2JXcdTZZMeER3f5wu8zyu5sSRzNO2Ud7/ULwmF
6Vos/YWajtivjftxTPx/pKb8o7NeqnNqrJpicdFJX6ZwrOtOG6uu9tev6JL5cda9TTMTLn7dtTaubaXb
+WO5VArV7+sJUWE869GxgwMHJFCQbmFaxE6/VrmSQP4MAAD//ym0NnnuBAAA
`,
	},

	"/definitions/hdtorrents.yml": {
		local:   "definitions/hdtorrents.yml",
		size:    1720,
		modtime: 1474169009,
		compressed: `
H4sIAAAJbogA/5RUXWvbMBR9z68Q3VgbVsd2EqedoYNtpRTGGGylFEofFOvGFlUkV5IbspD/PlmRZTtr
SZeHkHPuuR/n5tpBEAwQUlRDigqihZTAtTIUx0tDXV8GNy3HMM8rnBseeFBZhvJHlZofCAWo0LpUaRgW
JGgKjYTMw4GJZ7h0ugxryIWk4DBCcYrQD/FsmPArq37htePHnnfEpBVeXzpu2nK/Gy55QTczwn+4ied+
6gJkk/4pRTe3HVlkcVcyifqSybnFfoJZ0o9PTfcvFaGiicd7eNzHyVkfT5M9/WxPb/rf3d2Fl7e+445p
0FkXzRyycClI+1cowDIrUnT/9OAY/Ry05GktUIKfIigf6nQmcsp3ySXWRhNaZlQWpSUXQi5T+20h5WWl
fbOKkhQdbTZo9E3wBc1HlQJZ3x3abo+cplztaUqs1EpI4jUgpZDt/Awyc3wp0njOYJTNmcgelZmA6/tM
MCEvjt9dXUXmc7xzqEHpJttZWFFVMKq0dTHwS+m49Ofd+Owbey/xyk4tzQMDZnR/8tutO//1/cPFZjPa
bj9sNsBJa7hZtvX8Hda1V9VZCM40fTZPYGyxFCv1mvclptzaz4x5M63WSMuU6yLICsrICf84GbpUCUtR
FzXxAqsTTUa1f6x28QUFRnybxkCDe41Jp348RJ8R7hE+BWst6byybx0JC88vKNMgVVu7fq/sXkVPFci1
Mmk870RNKZmr1E/VXC3VDA5OOBmiOcJORUBjytT/Jb3mhIgVZwKTg9WS4cFaiv45bOW82S0xizioPhu+
ZeN1qRJLBS/sO07SaJpGCULROIzicBxFM3+HQHoVX7uPyGzSiRhAVrwpKa6TBoO/AQAA///M5E7quAYA
AA==
`,
	},

	"/definitions/iptorrents.yml": {
		local:   "definitions/iptorrents.yml",
		size:    2394,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/5RVbY/jNBD+vr/Cyp7Qrtg0ba8vaSRARwGB4KRFrNBKtwty42ljrWPnbLdLKf3vTNI4
b+1de1+azjPzPDOeSca+718RYriFiPDMKq1BWoOQpClCv9w/1BADE2sM4kpGxLsekHvNN9QCedA0fgHt
YYygcrWmK6SC9Nc5S3D5YiL8Q4hPEmszEwVBnaoXqzS4QndMszIsRs2V0hxKm5DpMCLv1QYRB0wdEPzx
Q4mFs2Ns1j+BTY4xhI6wUVjF/eyw8ahTyGR4LPa2Jv6kNPCVdMFhhz3sH6eYnihl+jYiD3860iQ3aud4
3LanYdseFuQ6wXDU8Xf4k45+OOz4x41iwmZl01lHedjO3PVjIe2THOwqfjbr8DHg3Zpx1WhL0w77pR38
powRYFyfZ53ASb8Qfid5CiU0wNz3c+fHxPfz4D2NHTByysXvQqkXN2ys6Xs0q1Sj0g7mKuWxg0PUf3x8
dBYma1jjli9sWrM2b3SwCjNVrP5CDFAdJxH58PG5ROzGr8G7PMAoeUcge87pQq24PJAzajEmKJBelmQF
uFQ6jYrfwuQyW9sq2dqAPuwHb7cjvbmSS77qOZTs914ZmFFjXpVm7UCHVoGgtdJO3MdSBcS4HyKyUGwb
xUpayqW58eY0s3FCyQY0X3JcE7iJyJJyAcy7LenYF5x7voEqgHQUGw4NqdrgORb6RPZk0Mj9Y1HjpWm8
66KdJPNOJDN2K+COLAqXBWOjql3FJNJtMYaraqqNMXmB9U5M5I2mr0WTNe5fwE5XK3S/3+16+/1Xux1I
tt9//CYfxO9r0Nver7DNx2CqOWj12nih3FksXQi4dvuafEts3sT8qSOp7E3+sIkfJ1ywm8HtoUNLDoJV
YuVG39bNasizFh11aRVFrdV8sc4vp0TDssKXXFjQptl7v7yxNKzgn6zhQBW9Mticv56evrt5emJf375x
U7HcCjhb1LBZFAN8IYT5MtKnToJ3X5o39aza+BI1pl6lUJSdVRtdomb4v+dbM3HfA8MJX9KTnv07tvXW
/ewoTSa4PTHJD8T7z7sj/oA8fxmLLLYEif2KZwBYK/knCg/dMQVAnFzCmN3+HwAA//+qvYfQWgkAAA==
`,
	},

	"/definitions/morethantv.yml": {
		local:   "definitions/morethantv.yml",
		size:    1255,
		modtime: 1474349693,
		compressed: `
H4sIAAAJbogA/6RTwY7TMBC95ytG5UKljUuLQJArR4S4rLigVeUm08TCsc14kmhZ7b/jOE7aQldFYg8r
z5v33nhe3DzPMwCvGAtoLSE30nAfICPbAH0J0H2AxP23gGlp6k7WAUeTd35ElPnhi3AAyKFhdr7YbIZh
ELOX4H6ThX4pXeKVkrG2pDDVANtxTh+AVO8KiOPGv7cFfOUGKYt1a6uTzKOksing+8+HhHCfn8C7keCt
uQN0D6Nc21qZSewkB84mIsI1LoJHS20R/wvZcbMfT7GhjOt4Gdt5pCmc1dMTiE/WHFUtZhSen1eJ6KT3
g6XqkjijCxGJLBWL5srFwt7o/Zh7KsfdNZYcdLA6XRbEIMkoU0OJhpEmf0bPf9izdds30T5bYswu+kTB
wS83uExgUnimtNlnfBw38stKZIcz8nxTlgeNr5L3PlbAJBIwfQKFulqk6aU8XlubKxHafl9aDZXqF0Yp
/VlMEOKJtDa+r1UB27973Ad8Nz8hxfpqzmHgQdV7ZY4WRE22c9NZFoabvGyUrl7v1klYIUul/X/7AEhm
Uodu/Hk2hMfZ3w5GW1ndGuCdNBfO21vOXv16Yf0zl/fLnuEL3WS/W4Ng1eK1yTHu5aVghfRCaGd+H+bp
GrFs/kXxcf07AAD//01qxqHnBAAA
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

	"/": {
		isDir: true,
		local: "/",
	},

	"/definitions": {
		isDir: true,
		local: "/definitions",
	},
}
