# stravaganza

![CI Status](https://github.com/jackal-xmpp/stravaganza/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jackal-xmpp/stravaganza?style=flat-square)](https://goreportcard.com/report/github.com/jackal-xmpp/stravaganza)
[![Coverage](https://codecov.io/gh/jackal-xmpp/stravaganza/branch/master/graph/badge.svg)](https://codecov.io/gh/jackal-xmpp/stravaganza)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/jackal-xmpp/stravaganza)
[![Releases](https://img.shields.io/github/release/jackal-xmpp/stravaganza/all.svg?style=flat-square)](https://github.com/jackal-xmpp/stravaganza/releases)
[![LICENSE](https://img.shields.io/github/license/jackal-xmpp/stravaganza.svg?style=flat-square)](https://github.com/jackal-xmpp/stravaganza/blob/master/LICENSE)

### Installation
```bash
go get -u github.com/jackal-xmpp/stravaganza/v2
```

### Example
```go
package main

import (
	"fmt"
	"os"

	"github.com/jackal-xmpp/stravaganza/v2"
)

func main() {
	iq, err := stravaganza.NewBuilder("iq").
		WithValidateJIDs(true).		
		WithAttribute("id", "zid615d9").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", "get").
		WithChild(
			stravaganza.NewBuilder("ping").
				WithAttribute("xmlns", "urn:xmpp:ping").
				Build(),
		).
		BuildIQ()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		return
	}
	_ = iq.ToXML(os.Stdout, true)
}
```

Expected output:
```
<iq id='zid615d9' from='ortuman@jackal.im/yard' to='noelia@jackal.im/balcony' type='get'><ping xmlns='urn:xmpp:ping'/></iq>
```

### Contributing
- Fork it
- Create your feature branch (`git checkout -b my-new-feature`)
- Commit your changes (`git commit -am 'Add some feature'`)
- Push to the branch (`git push origin my-new-feature`)
- Create new Pull Request

### License

[Apache License 2.0](https://github.com/jackal-xmpp/stravaganza/blob/master/LICENSE)
