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
 ...

## "api"
 import "github.com/maxb-odessa/sconf"
 err := sconf.Read(configFile)
 strval, err := sconf.ValAsStr("scope1", "str key")
 intval, err := sconf.ValAsInt32("another scope 2", number")
 floatval, err := sconf.ValAsFloat32("another scope 2", "another key")
 strval2, err := sconf.ValAsStr("scope1", "another key")
