[![Build Status](https://api.travis-ci.com/MrAlias/go-logger.svg)](https://travis-ci.com/MrAlias/go-logger)

# go-logger

A simple and opinionated logging package for Go.

The core goal of this library is to provide severity based logging functionality in a concurrency safe manner.

Outside of this core goal, this library tries to remaining as close to the standard [`log`](https://golang.org/pkg/log/) package as possible.

## Usage

The most common way to use this library is by logging with the default logger.

```go
package main

import (
	"net/http"
	"sync"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	logger "github.com/MrAlias/go-logger"
	severity "github.com/MrAlias/go-logger/severity"
)

var logLevel = kingpin.Flag("log-level", "set minimum log level").Default("INFO").Enum("DEBUG", "INFO", "ERROR")

func main() {
	kingpin.Parse()

	switch *logLevel {
	case "DEBUG":
		logger.SetSeverity(severity.Debug)
	case "INFO":
		logger.SetSeverity(severity.Info)
	case "ERROR":
		logger.SetSeverity(severity.Error)
	default:
		panic("invalid log-level")
	}

	var wg sync.WaitGroup
	var feeds = []string{
		"https://xkcd.com/rss.xml",
		"https://blog.golang.org/feed.atom",
		"not a valid url",
	}
	for _, feed := range feeds {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			logger.Debug("pulling feed: ", f)
			if _, err := http.Get(f); err != nil {
				logger.Errorf("failed to pull %q: %v", f, err)
			} else {
				logger.Infof("pulled %q successfully", f)
			}
		}(feed)
	}
	wg.Wait()
}
```

An example of running the above code might look like this.

```
$ go run code-from-above.go --log-level=DEBUG
DEBUG: pulling feed: not a valid url
DEBUG: pulling feed: https://xkcd.com/rss.xml
ERROR: failed to pull "not a valid url": Get not%20a%20valid%20url: unsupported protocol scheme ""
DEBUG: pulling feed: https://blog.golang.org/feed.atom
INFO : pulled "https://blog.golang.org/feed.atom" successfully
INFO : pulled "https://xkcd.com/rss.xml" successfully
```
