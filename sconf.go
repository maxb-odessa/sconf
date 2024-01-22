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

// strict parsing mode flag
var strictMode bool

// config var key-val pair type
type kvPairT map[string]string

// parsed and stored config data itself
var kvScopes map[string]kvPairT

// parsed line
type parsedLineT struct {
	scope string
	key   string
	value string
}

var configFilesRead []string

func init() {
	Clear()
}

func Clear() {
	kvScopes = make(map[string]kvPairT)
	configFilesRead = make([]string, 0)
}

// ToggleStrictMode toggles strict parsing mode, returns prev mode
func ToggleStrictMode() bool {
	strictMode = !strictMode
	return !strictMode
}

// we won't read very large files, say >16MB
const maxConfigSize = 16 * 1024 * 1024

var configSizeLimit int64

// SetReadLimit for config file size we may read
func SetReadLimit(limit int64) error {
	if limit > 0 {
		if limit < maxConfigSize {
			configSizeLimit = limit
			return nil
		} else {
			return fmt.Errorf("new limit exceeds upper boundary of %d bytes", maxConfigSize)
		}
	}

	return fmt.Errorf("invalid value")
}

// ReadLimit returns config size limit
func ReadLimit() int64 {
	if configSizeLimit == 0 {
		configSizeLimit = maxConfigSize
	}

	return configSizeLimit
}

// Read reads and parses config file
// return error or nil
func Read(path string) error {

	// open config file
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	confLimit := ReadLimit()
	fileInfo, err := fp.Stat()
	if fileInfo.Size() > confLimit {
		return fmt.Errorf("file size exceeds limit of %d bytes", confLimit)
	}

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

		if err := kvSet(pl.scope, pl.key, pl.value); err != nil {
			return err
		}
	}

	// preserve config files we've read for Dump() call
	configFilesRead = append(configFilesRead, path)

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

func parseLine(line string) (*parsedLineT, error) {
	lineLen := len(line)

	// no config line can be shorter than 3 chars (i.e. A=B or [G])
	if lineLen < 3 {
		return nil, fmt.Errorf("too short line")
	}

	// got a scope defining line: remember scope name and return
	if line[0] == '[' && line[lineLen-1] == ']' {

		currScope = strings.TrimSpace(line[1 : lineLen-1])

		if strictMode {
			if _, ok := kvScopes[currScope]; ok {
				return nil, fmt.Errorf("scope '%s' already defined (strict mode ON)", currScope)
			}
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

	return &parsedLineT{
		scope: currScope,
		key:   key,
		value: value,
	}, nil
}

// set (overriding) name-value pair in aprropriate scope
func kvSet(scope string, key string, value string) error {

	if _, ok := kvScopes[scope]; !ok {
		kvScopes[scope] = make(kvPairT)
	}

	if strictMode {
		if _, ok := kvScopes[scope][key]; ok {
			return fmt.Errorf("key '%s' in scope '%s' already exists (strict mode ON)", key, scope)
		}
	}

	kvScopes[scope][key] = value

	return nil
}

// Scopes returns an array of configured scopes
func Scopes() []string {
	var scopes []string
	for sc, _ := range kvScopes {
		scopes = append(scopes, sc)
	}
	return scopes
}

func getVal(scope string, key string) (string, error) {
	if kvp, ok := kvScopes[scope]; ok {
		if val, ok := kvp[key]; ok {
			return val, nil
		}
	} else {
		return "", fmt.Errorf("scope '%s' is not found", scope)
	}

	return "", fmt.Errorf("key '%s' is not found in scope '%s'", key, scope)
}

// Str returns configured value as a string from within specified scope
// return configured value as a string or an error if (either scope or key) not found
// if default value specified it will be returned instead of rising an error
func Str(scope string, key string, def ...string) (string, error) {
	val, err := getVal(scope, key)

	if err != nil && len(def) > 0 {
		return def[0], nil
	}

	return val, err
}

// Int gets intXX value from specified scope
func Int(scope string, key string, def ...int64) (int64, error) {
	val, err := getVal(scope, key)

	if err != nil {
		if len(def) > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}

	return strconv.ParseInt(val, 0, 64)
}

// Float gets floatXX value from specified scope
func Float(scope string, key string, def ...float64) (float64, error) {
	val, err := getVal(scope, key)

	if err != nil {
		if len(def) > 0 {
			return def[0], nil
		} else {
			return 0, err
		}
	}

	return strconv.ParseFloat(val, 64)
}

// Bool gets boolean value from specified scope
func Bool(scope string, key string, def ...bool) (bool, error) {
	val, err := getVal(scope, key)

	if err != nil {
		if len(def) > 0 {
			return def[0], nil
		} else {
			return false, err
		}
	}

	switch strings.ToLower(val) {
	case "0", "no", "f", "false", "none", "never", "negative":
		return false, nil
	default:
		return true, nil
	}
}

// Dump current config values into specified file
// useful to create "override" configs
func Dump(fname string) error {

	if len(configFilesRead) == 0 {
		return fmt.Errorf("no config file(s) read yet")
	}

	confData := "# Generated from: " + strings.Join(configFilesRead, ", ") + "\n"

	for scope, kv := range kvScopes {
		confData += "[" + scope + "]\n"
		for key, val := range kv {
			confData += "  " + key + " = " + val + "\n"
		}
	}

	return os.WriteFile(fname, []byte(confData), 0644)
}
