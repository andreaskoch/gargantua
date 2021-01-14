# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [v0.5.0-alpha] - 2020-11-16

### Changed
- Ignore invalid SSL certificates

### Added
- Log response headers

## [v0.4.1-alpha] - 2020-11-16

### Fixed
- Don't visit URLs twice (introduced with v0.4.0-alpha)

## [v0.4.0-alpha] - 2020-11-05

Logging

### Added
- Add "Save downloaded data to disk" to the roadmap (feature request #1)
- Log results to a log file
- Add Github actions

### Changed
- Add support for non-sitemap URLs
- Capture the parent URL
- Switch from go 1.14 to 1.15

## [v0.2.0-alpha] - 2017-02-07

### Fixed
- Shut down the UI 10 seconds after the crawler finished

### Changed
- Remove the TERABYTE constant from the byte formatter so cross-compilation for ARM works
- Remove the `reset` command hint from the troubleshooting section of the README
- Start the timer only after the links in the XML sitemap(s) have been read
- Stop counting the seconds after the crawler finished

### Removed
- Remove the debug console from the command-line ui

## [v0.1.0-alpha] - 2017-02-07

The Prototype

### Added
- The gargantua prototype
