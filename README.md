# 「 gargantua 」

The fast website crawler

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
