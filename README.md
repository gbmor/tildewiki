# Tildewiki [![Go Report Card](https://goreportcard.com/badge/github.com/gbmor/tildewiki)](https://goreportcard.com/report/github.com/gbmor/tildewiki) [![GolangCI](https://img.shields.io/badge/golangci-check-blue.svg)](https://golangci.com/r/github.com/gbmor/tildewiki) [![Travis CI](https://api.travis-ci.org/gbmor/tildewiki.svg?branch=master)](https://travis-ci.org/gbmor/tildewiki)
A wiki engine designed around the needs of the [tildeverse](https://tildeverse.org)

## [v0.5.2](https://github.com/gbmor/tildewiki/releases/tag/v0.5.2)
A ton of refactoring has gone into 0.5. Here are some noteworthy changes:
* Various performance improvements
* Index page is now being cached
* Refresh interval for the index page is configurable
* Logging can be output to `stdout` (default), to a file, or to /dev/null for some peace and quiet.

I've been stress testing it with [tsenart/vegeta](https://github.com/tsenart/vegeta).
With all the internal changes in v0.5, it can handle a *lot* of requests before performance
starts to degrade.

### [Development Branch](https://github.com/gbmor/tildewiki/tree/dev)
Contains all the new changes going into the next version

### Currently powering the [tilde.institute](https://tilde.institute) wiki: 
* [https://wiki.tilde.institute](https://wiki.tilde.institute) 
* [gtmetrix report](https://gtmetrix.com/reports/wiki.tilde.institute/F1tzxEch)

## Features
* Speed is a priority
* Mobile-friendly pages
* Markdown!
* Uses [kognise/water.css](https://github.com/kognise/water.css) dark theme by
default (and includes as an example, a simple but nice local CSS file)
* YAML configuration
* Automatically reloads config file when a change is detected.
* Generates list of pages, then places at an anchor comment in the index page
* Caches pages to memory and only re-renders when the file changes
* Very configurable. For example:
  * URL path for viewing pages
  * Directory for page data
  * File to use for index page
  * Logging output (file, stdout, null) and file location
* Runs as a multithreaded service, rather than via CGI
* Easily use Nginx to proxy requests to it. This allows you to use your
existing SSL certificates.

### Notes
* Builds with `Go 1.11` and `Go 1.12`. Not tested with any other version.
* Tested on Linux (Ubuntu 18.04LTS, Debian 9) and OpenBSD 6.4
* If you have access to other environments and can test, please let me know.
It will be much appreciated.

For [tildeverse](https://tildeverse.org) projects, we tend to use a PR
workflow. For example, wiki pages are submitted to the repo via pull
request. That's what I'm initially designing this around. I will likely
add authentication and in-place page editing last, after everything else
is done, including unit tests.

Uses a patched copy of [russross/blackfriday](https://github.com/russross/blackfriday)
([gopkg](https://gopkg.in/russross/blackfriday.v2)) as the markdown
parser. The patch allows injection of various `<meta.../>` tags into
the document header during the `markdown->html` translation.

* The patched `v2` repository lives at:
[gbmor-forks/blackfriday.v2-patched](https://github.com/gbmor-forks/blackfriday.v2-patched)
* The patched `master` repo lives at:
[gbmor-forks/blackfriday](https://github.com/gbmor-forks/blackfriday).
* The PR can be found here: [allow writing of user-specified
&lt;meta.../&gt;...](https://github.com/russross/blackfriday/pull/541)

