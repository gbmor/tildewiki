# Tildewiki [![gpl-v3](https://img.shields.io/badge/license-GPLv3-brightgreen.svg "GPL v3")](https://github.com/gbmor/tildewiki/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/gbmor/tildewiki)](https://goreportcard.com/report/github.com/gbmor/tildewiki) [![GolangCI](https://img.shields.io/badge/golang%20ci-success-blue.svg)](https://golangci.com/r/github.com/gbmor/tildewiki)

A wiki engine designed for the needs of the [tildeverse](https://tildeverse.org)

## Features

* Mobile-friendly pages
* Markdown rendering
* Uses [kognise/water.css](https://github.com/kognise/water.css) dark theme by default
* Watches (YAML) config file for changes and automatically reloads
* Dynamically generates index of pages and places at anchor-point in front/index page
* Basically everything is configurable: URL path for viewing pages, directory for page data, file to use for index page, etc.
* Runs as a multithreaded service, rather than via CGI
* Caches pages to memory and only re-renders when the modification time changes
* Easily use Nginx to proxy requests to it. This allows you to use your existing SSL certificates.
* Speed is a priority

## Notes

Uses a patched copy of [russross/blackfriday](https://github.com/russross/blackfriday) ([gopkg](https://gopkg.in/russross/blackfriday.v2)) as the markdown parser. The patch allows injection of various `<meta.../>` tags into the document header during the `markdown->html` translation.

* The patched `v2` repository lives at: [gbmor-forks/blackfriday.v2-patched](https://github.com/gbmor-forks/blackfriday.v2-patched)
* The patched `master` repo lives at: [gbmor-forks/blackfriday](https://github.com/gbmor-forks/blackfriday). 
* The PR can be found here: [allow writing of user-specified &lt;meta.../&gt;...](https://github.com/russross/blackfriday/pull/541)

