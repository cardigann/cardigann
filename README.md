Cardigann [![Build Status](https://travis-ci.org/cardigann/cardigann.svg?branch=master)](https://travis-ci.org/cardigann/cardigann) [![Go Report Card](https://goreportcard.com/badge/github.com/cardigann/cardigann)](https://goreportcard.com/report/github.com/cardigann/cardigann) [![Gitter](https://badges.gitter.im/cardigann/cardigann.svg)](https://gitter.im/cardigann/cardigann?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=body_badge)
=========

A server for adding extra indexers to Sonarr, SickRage and CouchPotato via [Torznab](https://github.com/Sonarr/Sonarr/wiki/Implementing-a-Torznab-indexer) and [TorrentPotato](https://github.com/CouchPotato/CouchPotatoServer/wiki/Couchpotato-torrent-provider) proxies. Behind the scenes Cardigann logs in and runs searches and then transforms the results into a compatible format. 

Cardigann is implemented in golang, which means it's simply a single binary to execute/install, no runtime. Windows, Linux and OSX will be supported, although it should run on any platform that golang supports.

## Usage

Cardigann provides a cli tool for debugging and development:

```bash
cardigann query bithdtv t=tv-search "q=my show name" ep=1 season=2
```

Or you can run the proxy server:

```
cardigann server
```

Once the server is running, visit http://localhost:5060 and configure via the web interface.

You can set a password requirement by either passing the `--passphrase` flag to the server command, or by setting `global.passphrase` in the [Configuration](#Configuration).

## Installation

Cardigann is distributed on equinox.io in a variety of formats for macOS, Linux and Windows. 

https://dl.equinox.io/cardigann/cardigann/stable

Follow the instructions on the above to install the cardigan binary, and then you can run the following to run the server in the foreground:

```bash
cardigann server
```

At this point you can visit the web interface on http://localhost:5060.

If you want to run this service non-interactively, you can install it as a service (supports windows services, macOS launchd, linux upstart, systemv and systemd):

```bash
cardigann service install
cardigann service start
```

## Updating

Cardigann has an experimental upgrade-in-place feature using equinox.io:

```
cardigann update --channel=stable
```

If you like to live dangerously, you can update to the edge channel:

```
cardigann update --channel=edge
```

## Configuration

Configuration is stored in a `config.json` file. It's searched for in a few different locations, in order of priority:

#### All Platforms
  * `$CWD/config.json`
  * `$CONFIG_DIR/config.json`

#### Linux/BSD:
  * `$HOME/.config/cardigann/config.json`
  * `/etc/xdg/cardigan/config.json`
  
#### OSX:
  * `$HOME/Library/Preferences/cardigann/config.json`
  * `/Library/Preferences/cardigann/config.json`
  * `/Library/Application Support/cardigann/config.json`

#### Windows
  * `%APPDATA%\cardigann\config.json`
  * `%LOCALAPPDATA%\cardigann\config.json`

This configuration file will contain your tracker credentials in plain-text, so it's important to keep it secure. 

## Definitions

Definitions are yaml files (see [definitions](definitions/) for their source) that define how to login and search on an indexer. You can either use the included definitions or write your own. Definitions are loaded from the following directories:

#### All Platforms
  * `$CWD/definitions/`
  * `$CONFIG_DIR/definitions/`

#### Linux/BSD:
  * `$HOME/.config/cardigann/definitions/`
  * `/etc/xdg/cardigan/definitions/`
  
#### OSX:
  * `$HOME/Library/Preferences/cardigann/definitions/`
  * `/Library/Preferences/cardigann/definitions/`
  * `/Library/Application Support/cardigann/definitions/`

#### Windows
  * `%APPDATA%\cardigann\definitions\`
  * `%LOCALAPPDATA%\cardigann\definitions\`

## Using with a Proxy

Currently either a SOCKS5 proxy like Privoxy or Tor can be used:

```
SOCKS_PROXY=localhost:1080 cardigann server
```

Or, an http proxy works too:

```
HTTP_PROXY=localhost:8080 cardigann server
```

If you are running as a service, you will need to manually edit the service files to set the environment.

## Supported Indexers

Cardigann simply provides a format for describing how to log into and scrape the search results of various forums and sites. It is not endorsed by the various sites, nor is it intended for piracy. You are using Cardigann at your own risk.

* Abnormal
* AlphaRatio
* AvistaZ
* BIT-HDTV
* BeyondHD
* BitMeTV
* CinemaZ
* DanishBits
* Demonoid
* EZTV
* FileList
* Freshon
* FunFile
* HD-Torrents
* HDArea
* HDME
* Immortalseed
* IPTorrents
* MoreThanTV
* NCore
* Norbits
* PreToMe
* PrivateHD
* Speed.CD
* Sceneaccess
* SceneTime
* Shareisland
* The New Retro
* ThePirateBay (TPB)
* Torrent Sector Crew
* Torrent-Syndikat
* Torrentbytes
* TorrentDay
* TorrentHeaven
* Torrentleech
* Transmithe.Net
* Tspate
* UHDBits
* WorldOfP2P

I'm happy to add new trackers, please either open a new issue, or a pull request with whatever details you have for the tracker.

## Reporting bugs

Cardigann is new software, and relies on scraping indexers, so is inherently prone to breaking. We try and reply as quickly as possible, but please make sure before you report a bug that you've update to the latest version.

If the issue persists, [file a bug][bug_report_template].

## Requests

* Start an issue on GitHub following one of these templates:
  * [Feature request][feature_request_template]
  * [Indexer/Tracker request][indexer_request_template]

## Questions? Wanna chat?

* If none of the templates above is appropriate, [open an issue](https://github.com/cardigann/cardigann/issues/new)
* Join us on [Gitter](https://gitter.im/cardigann/cardigann)

## Credits

Inspired by Jackett, or at least born of frustration with it always crashing and requiring a mono runtime.

[bug_report_template]: https://github.com/cardigann/cardigann/issues/new?title=Bug%20report%3A%20%5Bsummarise%20the%20issue%20here%5D&body=%23%23%23%20Issue%20experienced%0A%0A%23%23%23%20Steps%20to%20reproduce%0A%0A%23%23%23%20Cardigann%20version%20(in%20the%20footer%2C%20or%20via%20%60cardigann%20--version%60)%0A%0A%0A
[feature_request_template]: https://github.com/cardigann/cardigann/issues/new?title=Feature%20request%3A%20summarize%20the%20feature%20here&body=%0A%23%23%23%20Description%20of%20feature%2Fenhancement%0A%0A%23%23%23%20Justification%0A%0A%23%23%23%20Example%20use%20case
[indexer_request_template]: https://github.com/cardigann/cardigann/issues/new?title=Indexer%20Request%3A%20indexer%20title%20here&body=%23%23%23%20Indexer%20Name%0A%0A%23%23%23%20Is%20the%20indexer%20private%20or%20public%3F%20If%20private%2C%20can%20you%20provide%20an%20invite%3F%0A%0A%23%23%23%20Indexer%20URL%0A%0A%23%23%23%20Indexer%20language

