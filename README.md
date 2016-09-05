Cardigann [![Build Status](https://travis-ci.org/cardigann/cardigann.svg?branch=master)](https://travis-ci.org/cardigann/cardigann) [![Go Report Card](https://goreportcard.com/badge/github.com/cardigann/cardigann)](https://goreportcard.com/report/github.com/cardigann/cardigann)
=========

Provides [Torznab](https://github.com/Sonarr/Sonarr/wiki/Implementing-a-Torznab-indexer) and [TorrentPotato](https://github.com/CouchPotato/CouchPotatoServer/wiki/Couchpotato-torrent-provider) interfaces for [private torrent trackers](http://lifehacker.com/5897095/whats-a-private-bittorrent-tracker-and-why-should-i-use-one).

Cardigann can be used to add any supported private tracker to your software of choice (e.g Sonarr, SickRage, CouchPotato). This is done by proxying requests to the individual trackers and scraping the responses and converting them to the correct format. The rules for scraping sites is expressed in a custom YAML format to make updating it easy without having to write code. 

Cardigann is implemented in golang, which means it's simply a single binary to execute/install, no runtime. Windows, Linux and OSX will be supported, although it should run on any platform that golang supports.

## Usage

Cardigann provides a cli tool for debugging and development:

```bash
cardigann query bithdtv t=tv-search "q=mr robot" ep=1 season=2
```

Or you can run the proxy server:

```
cardigann server
```

Once the server is running, visit http://localhost:5060 and configure via the web interface.

You can set a password requirement by either passing the `--passphrase` flag to the server command, or by setting a password in the global config file.

## Configuration

Config is stored in a `config.json` file that is read from the current directory or your user level config dir (e.g `$HOME/.config` in linux). Tracker credentials will be stored here, but you can also set the following keys:

<table>
<thead><tr><th>Key</th><th>Values</th></tr></thead>
<tbody>
<tr><td>apikey</td><td>A 16 character hex value</td></tr>
</tbody>
</table>

## Supported Trackers

* BIT-HDTV
* IPTorrents
* Freshon
* Demonoid
* HD-Torrents

## Planned Trackers

I'd love assistance adding these trackers, either via invites or pull-requests. 

* Abnormal
* AlphaRatio
* AnimeBytes
* Avistaz
* bB
* BeyondHD
* BitMeTV
* BitSoup
* BlueTigers
* BTN
* DanishBits
* Demonoid
* EuTorrents
* FileList
* Fuzer
* HD-Space
* Hebits
* Hounddawgs
* ILoveTorrents
* Immortalseed
* PassThePopcorn
* MoreThanTV
* MyAnonamouse
* NCore
* NextGen
* Pretome
* PrivateHD
* RevolutionTT
* SceneAccess
* SceneFZ
* SceneTime
* Shazbat
* SpeedCD
* TehConnection
* TorrentBytes
* TorrentDay
* TorrentLeech
* TorrentShack
* TransmitheNet
* TV Chaos UK
* World-In-HD
* XSpeeds
* Xthor

## Credits

Inspired by Jackett, or at least born of frustration with it always crashing and requiring a mono runtime.