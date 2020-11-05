# 「 gargantua 」

The fast website crawler

You can use「 gargantua 」to quickly and easily

- **warm-up** your frontend caches
- perform small **load-tests** against your publicly available pages
- **measure** response times
- **detect** broken links

from your command line on Linux, macOS and Windows.

![Animation: gargantua v0.1.0 crawling a website](files/gargantua-in-action-crawling-a-website.gif)

> Note: Press `Q` to stop the current crawling process.

## Usage

Crawl **www.sitemaps.org** with 5 concurrent workers:

```bash
gargantua crawl --url https://www.sitemaps.org/sitemap.xml --workers 5
```

see also: [A short introduction video of gargantua on YouTube](https://www.youtube.com/watch?v=TSCMvUvc0qo)

### Customize the user-agent

You can specify a customized user agent using the `--user-agent` argument:

```bash
gargantua crawl --url https://www.sitemaps.org/sitemap.xml --workers 5 --user-agent "gargantua bot / iPhone"
```

### Log all requests

You can specify a log file with the `--log` argument:

```bash
gargantua crawl --url https://www.sitemaps.org/sitemap.xml --workers 5 --log "gargantua.log"
```

```
Date and time       #worker   Status Code     Bytes   Response Time   URL                                                          Parent URL
2020/11/05 09:23:14 #001:     200             4403    148.759000ms    https://www.sitemaps.org                                     https://www.sitemaps.org/ko/faq.html
2020/11/05 09:23:14 #002:     200             4403    290.536000ms    http://www.sitemaps.org/                                     https://www.sitemaps.org/ko/faq.html
2020/11/05 09:23:14 #003:     200            45077    283.243000ms    https://www.sitemaps.org/protocol.html                       https://www.sitemaps.org/ko/faq.html
2020/11/05 09:23:14 #004:     404             1245    155.376000ms    https://www.sitemaps.org/protocol.htm                        https://www.sitemaps.org/ko/faq.html
2020/11/05 09:23:14 #005:     200             4403    155.577000ms    https://www.sitemaps.org/index.html                          https://www.sitemaps.org/ko/faq.html
2020/11/05 09:23:14 #001:     200             2591    286.451000ms    http://www.sitemaps.org/schemas/sitemap/0.9/siteindex.xsd    https://www.sitemaps.org/ko/faq.html
2020/11/05 09:23:14 #003:     200            10839    143.738000ms    https://www.sitemaps.org/terms.html                          https://www.sitemaps.org/ko/faq.html
2020/11/05 09:23:14 #005:     200            15681    141.580000ms    https://www.sitemaps.org/faq.html                            https://www.sitemaps.org/ko/protocol.html
2020/11/05 09:23:14 #002:     404             1245    286.175000ms    http://www.sitemaps.org/protocol.htm                         https://www.sitemaps.org/ko/faq.html
```

[gargantua.log](files/gargantua.log)


## Download

You can download binaries for Linux, macOS and Windows from [github.com »andreaskoch » gargantua » releases](https://github.com/andreaskoch/gargantua/releases):

```bash
wget https://github.com/andreaskoch/gargantua/releases/download/v0.3.0-alpha/gargantua_linux_amd64
```

## Docker Image

There is also a docker image that you can use to download or run the latest version of gargantua:

[andreaskoch/gargantua](https://hub.docker.com/r/andreaskoch/gargantua/)

```bash
docker run --rm andreaskoch/gargantua:latest \
       crawl \
       --verbose \
       --url https://www.sitemaps.org/sitemap.xml \
       --workers 5
```

**Note**: You will need the `--verbose` flag in order to prevent the command-line UI from loading. Otherwise gargantua will fail.

## Roadmap

- Increase the number of workers at runtime
- Silent mode (only show statistics at the end)
- CSV mode (print CSV output to stdout)
- Web-UI
- Save downloaded data to disk

## License

「 gargantua 」is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
