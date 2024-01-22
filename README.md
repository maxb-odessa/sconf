# sconf
extremely Simple CONFig file reader

## why?
beacuse I need it

## format
file format is toml like, ex:

```
 # comment
 ; comment too
 [scope1]
   str key = some line here
   another key = "Some\nText\tHere"

 [another scope 2]
   number = -666
   another key = 123.456

 [scope1]
   forgot to add this key to scope1 = now we did
 ...
```
## "api"

```
 import "github.com/maxb-odessa/sconf"

 // set max config file limit (optional)
 err := sconf.SetReadLimit(1024)

 // enable strict parsing mode: no scopes or keys duplication allowed
 oldStrictMode := sconf.ToggleStrictMode()

 // read config file
 err := sconf.Read("/path/to/conf.txt")

 // get config values as a string
 str, err := sconf.Str("scope", "str key")
 // the same but with the default if a scope or key is not present
 str2 := sconf.Str("scope", "str key", "default string value")

 // this will return int64
 intv, err := sconf.Int("another scope", "number key")
 intv2 := sconf.Int("another scope", "number key", 123)

 // this will return float64
 floatv, err := sconf.Float("another scope 2", "another float key")
 floatv2 := sconf.Float("another scope 2", "another float key", 123.45)

 // ditto for boolean
 boolv, err := sconf.Bool("scope N", "yesno")
 boolv2 := sconf.BoolDef("scope N", "yesno", false)

 // read another config file overriding existing and/or adding new scopes/keys
 err = sconf.Read("/path/to/another/conf.txt")

 // dump current config values into file, useful to create "overrides"
 err = sconf.Dump("/path/to/conf.dump")

 // clear configured data
 sconf.Clear()
```
