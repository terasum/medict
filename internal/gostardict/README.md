GO STARDICT
=======

To download and install this package run:

`go get github.com/dyatlov/gostardict/stardict`

Source docs: http://godoc.org/github.com/dyatlov/gostardict/stardict

Disclaimer
---
The code is currently undocumented, and is certainly **not idiomatic Go**. Pull requests are welcome!

Samples
---
Sample code can be found in [`gostardict/samples`](https://github.com/dyatlov/gostardict/tree/master/samples).

Project Overview
---
The project was started as an attempt to read stardict dictionaries in language learning webservice and grew into a tool supporting several dictionary formats.

Current limitations:

  * Whole dictionary and index are fully loaded into memory for fast random access
  * DictZip format is not supported, it is processed as a simple GZip format (means that no random blocks access is supported as in DictZip)
  * Syn files are not processed
  * There's no recovering from errors (means that dictionaries should be well formed)

Not tested but should be working in theory (I didn't find dictionaries with those properties in place):

  * 64bit offsets
  * multi typed dictionary fields


Dictionary Functions
---
`func NewDictionary(path string, name string) (*Dictionary, error)` - returns `gostardict.Dictionary` struct. `path` - path to dictionary files, `name` - name of dictionary to parse

`func (d Dictionary) GetBookName() string` - returns dictionary' book name (from ifo file)

`func (d Dictionary) GetWordCount() uint64` - returns dictionary' word count (from ifo file)

Additional links
---
  * Stardict format: https://code.google.com/p/babiloo/wiki/StarDict_format
  * Dictionaries: http://abloz.com/huzheng/stardict-dic/

Work Opportunities
---
If you need a Golang developer feel free to contact with me :)
