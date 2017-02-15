# 「 gargantua 」

The fast website crawler

You can use「 gargantua 」to quickly and easily

- **warm-up** your frontend caches
- perform small **load-tests** against your publicly available pages
- **measure** response times
- **detect** broken links

from your command line on Linux, MacOS and Windows.

![Animation: gargantua v0.1.0 crawling a website](files/gargantua-in-action-crawling-a-website.gif)

> Note: Press `Q` to stop the current crawling process.

## Usage

Crawl **www.sitemaps.org** with 5 concurrent workers:

```bash
gargantua crawl --url https://www.sitemaps.org/sitemap.xml --workers 5
```

## Roadmap

- Increase the number of workers at runtime
- Personalized user agent string
- Silent mode (only show statistics at the end)
- CSV mode (print CSV output to stdout)
- Dockerfile
- Web-UI

## License

「 gargantua 」is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
