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

You can set a password requirement by either passing the `--passphrase` flag to the server command, or by setting a password in the global config file.

## Configuration

Config is stored in a `config.json` file that is read from the current directory or your user level config dir (e.g `$HOME/.config` in linux). Tracker credentials will be stored here, but you can also set the following keys:

<table>
<thead><tr><th>Key</th><th>Values</th></tr></thead>
<tbody>
<tr><td>apikey</td><td>A 16 character hex value</td></tr>
</tbody>
</table>

## Installation

Cardigann is distributed as a binary and a collection of tracker definition files. These are available from the [releases page](https://github.com/cardigann/cardigann/releases) for macOS, Linux and Windows. The following example shows how to run the daemon interactively under Linux:

```bash
curl https://github.com/cardigann/cardigann/releases/download/(VERSION)/cardigann-linux-amd64 
chmod +x cardigann-linux-amd64 
curl https://github.com/cardigann/cardigann/releases/download/(VERSION)/defs.zip
unzip defs.zip
./cardigann-linux-amd64 server
```

At this point you can visit the web interface on http://localhost:5060.

If you want to run this service non-interactively, you can install it as a service (supports windows services, macOS launchd, linux upstart, systemv and systemd):

```bash
./cardigann-linux-amd64 service install
./cardigann-linux-amd64 service start
```

Install your definitions in `/etc/xdg/cardigann/definitions` for them to be found.

## Supported Indexers

Cardigann simply provides a format for describing how to log into and scrape the search results of various forums and sites. It is not endorsed by the various sites, nor is it intended for piracy. You are using Cardigann at your own risk.

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

I'm happy to add new trackers, please either open a new issue, or a pull request with whatever details you have for the tracker.

## Credits

Inspired by Jackett, or at least born of frustration with it always crashing and requiring a mono runtime.
