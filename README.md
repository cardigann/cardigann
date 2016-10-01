Cardigann [![Build Status](https://travis-ci.org/cardigann/cardigann.svg?branch=master)](https://travis-ci.org/cardigann/cardigann) [![Go Report Card](https://goreportcard.com/badge/github.com/cardigann/cardigann)](https://goreportcard.com/report/github.com/cardigann/cardigann)
=========

A server for adding extra indexers to Sonarr, SickRage and CouchPotato via [Torznab](https://github.com/Sonarr/Sonarr/wiki/Implementing-a-Torznab-indexer) and [TorrentPotato](https://github.com/CouchPotato/CouchPotatoServer/wiki/Couchpotato-torrent-provider) proxies. Behind the scenese Cardigann logs in and runs searches and then transforms the results into a compatible format. 

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

You can set a password requirement by either passing the `--passphrase` flag to the server command, or by setting `global.password` in the [Configuration](#Configuration).

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

## Supported Indexers

Cardigann simply provides a format for describing how to log into and scrape the search results of various forums and sites. It is not endorsed by the various sites, nor is it intended for piracy. You are using Cardigann at your own risk.

* AlphaRatio
* BIT-HDTV
* IPTorrents
* Freshon
* Demonoid
* HD-Torrents
* BeyondHD 
* FileList
* MoreThanTV
* NCore
* ThePirateBay (TPB)
* EZTV
* Torrentleech
* TorrentDay
* Speed.CD

I'm happy to add new trackers, please either open a new issue, or a pull request with whatever details you have for the tracker.

## Credits

Inspired by Jackett, or at least born of frustration with it always crashing and requiring a mono runtime.
