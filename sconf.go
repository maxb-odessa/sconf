/* simple config file reader
   the format is toml like
     # comment
     ; comment too
     [scope]
     key = val
     [another scope]
	 another key = another value
     ..
*/

package sconf

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// config var key-val pair type
type kvPair map[string]string

// parsed and stored config data itself
var kvScopes map[string]kvPair

// parsed line
type parsedLine struct {
	scope string
	key   string
	value string
}

// read and parse config file
func Read(path string) error {

	// init config data storage
	kvScopes = make(map[string]kvPair)

	// open config file
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)

	// read config file
	for lineNo := 1; scanner.Scan(); lineNo++ {

		line := prepareLine(scanner.Text())
		if line == "" {
			continue
		}

		pl, err := parseLine(line)
		if err != nil {
			return fmt.Errorf("line %d: %s (near '%.16s ...')", lineNo, err, line)
		}

		// a line is a scope line, nothing yet to save
		if pl == nil {
			continue
		}

		kvSet(pl.scope, pl.key, pl.value)
	}

	return nil
}

// prepare config file line for parsing: trim and unescape it, discard comments
func prepareLine(line string) string {
	l := strings.TrimSpace(line)
	if len(l) < 1 || l[0] == '#' || l[0] == ';' {
		return ""
	}
	return l
}

// parse prepared config line
var currScope string

func parseLine(line string) (*parsedLine, error) {
	lineLen := len(line)

	// no config line can be shorter than 3 chars (i.e. A=B or [G])
	if lineLen < 3 {
		return nil, fmt.Errorf("too short line")
	}

	// got a scope defining line: remember scope name and return
	if line[0] == '[' && line[lineLen-1] == ']' {
		currScope = strings.TrimSpace(line[1 : lineLen-1])
		if currScope == "" {
			return nil, fmt.Errorf("ikvalid scope name")
		}
		return nil, nil
	}

	// no scope defined yet
	if currScope == "" {
		return nil, fmt.Errorf("expression without scope")
	}

	// try to get name=value pair here
	tokens := strings.SplitN(line, "=", 2)
	if len(tokens) != 2 {
		return nil, fmt.Errorf("can not parse")
	}

	key := strings.TrimSpace(tokens[0])
	if len(key) < 1 {
		return nil, fmt.Errorf("param name missed")
	}

	value := strings.TrimSpace(tokens[1])
	if len(value) < 1 {
		return nil, fmt.Errorf("param value missed")
	}

	value = strings.ReplaceAll(value, `\n`, "\n")
	value = strings.ReplaceAll(value, `\r`, "\r")
	value = strings.ReplaceAll(value, `\t`, "\t")
	value = strings.ReplaceAll(value, `\\`, `\`)
	value = strings.ReplaceAll(value, `\'`, `'`)
	value = strings.ReplaceAll(value, `\"`, `"`)

	return &parsedLine{
		scope: currScope,
		key:   key,
		value: value,
	}, nil
}

// set (overriding) name-value pair in aprropriate scope
func kvSet(scope string, key string, value string) {

	if _, ok := kvScopes[scope]; !ok {
		kvScopes[scope] = make(kvPair)
	}

	kvScopes[scope][key] = value
}

// get array of configured scopes
func Scopes() []string {
	var scopes []string
	for sc, _ := range kvScopes {
		scopes = append(scopes, sc)
	}
	return scopes
}

// get string value from specified scope
func ValAsStr(scope string, key string) (string, error) {
	if kvp, ok := kvScopes[scope]; ok {
		if val, ok := kvp[key]; ok {
			return val, nil
		}
	}
	return "", fmt.Errorf("key '%s' not found in scope '%s'", key, scope)
}

// get int32 value from specified scope
func ValAsInt32(scope string, key string) (int32, error) {
	if val, err := ValAsStr(scope, key); err != nil {
		return 0, err
	} else if i, err := strconv.ParseInt(val, 10, 0); err != nil {
		return 0, err
	} else {
		return int32(i), nil
	}
}

// get float32 value from specified scope
func ValAsFloat32(scope string, key string) (float32, error) {
	if val, err := ValAsStr(scope, key); err != nil {
		return 0, err
	} else if i, err := strconv.ParseFloat(val, 32); err != nil {
		return 0, err
	} else {
		return float32(i), nil
	}
}
