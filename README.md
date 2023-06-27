# sconf
extremely Simple CONFig file reader

## why?
beacuse I need it

## format
file format is toml like, ex:

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

## "api"
 import "github.com/maxb-odessa/sconf"

 err := sconf.Read("/path/to/conf.txt")

 str, err := sconf.Str("scope", "str key")
 str2 := sconf.Str("scope", "str key", "default value")

 intv, err := sconf.Int32("another scope", "number key")
 intv2 := sconf.Int32Def("another scope", "number key", 123)

 intv, err := sconf.Int64("another scope", "number key")
 intv2 := sconf.Int64Def("another scope", "number key", 123123567890)
 
 floatv, err := sconf.Float32("another scope 2", "another float key")
 floatv2 := sconf.Float32Def("another scope 2", "another float key", 123.45)

 floatv, err := sconf.Float64("another scope 2", "another float key")
 floatv2 := sconf.Float64Def("another scope 2", "another float key", 0.1234e+10)

 boolv, err := sconf.Bool("scope N", "yesno")
 boolv2 := sconf.BoolDef("scope N", "yesno", false)

 err = sconf.Read("/path/to/another/conf.txt") // will override already read values or add new

 err = sconf.Dump("/path/to/conf.dump")
