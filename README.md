Cardigann
=========

**Note that this is still in active development, very little is working**

Provides [Torznab](https://github.com/Sonarr/Sonarr/wiki/Implementing-a-Torznab-indexer) and [TorrentPotato](https://github.com/CouchPotato/CouchPotatoServer/wiki/Couchpotato-torrent-provider) interfaces for [private torrent trackers.](http://lifehacker.com/5897095/whats-a-private-bittorrent-tracker-and-why-should-i-use-one).

Cardigann can be used to add any supported private tracker to your software of choice (e.g Sonarr, SickRage, CouchPotato). This is done by proxying requests to the individual trackers and scraping the responses and converting them to the correct format. Unfortunately this means that as the individual trackers sites update their markup, they need to be addressed with updates to Cardigann.

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

Once the server is running, visit http://localhost:3000 and configure via the web interface.

## TODO

 * [ ] HTTP API for Torznab
 * [ ] Web UI (Basic ReactJS UI)
 * [ ] Proxying downloads and rewriting feed links
 * [ ] DRY up indexers, is surf + goquery the best way to do this?

## Planned Trackers

 * [ ] BIT-HDTV
 * [ ] IPTorrents
 * [ ] TV Torrents
 * [ ] Demonoid

Open to suggestions on these, please file an issue.

## Credits

Inspired by Jackett, or at least born of frustration with it always crashing and requiring a mono runtime.