package main

import (
	"fmt"
	"math"
	"strings"
)

const prettyPythonScriptName = "texttable.py"

// PrettyTableFormatter prints input in tabular format
type PrettyTableFormatter struct {
}

// Entry a key value pair
type Entry struct {
	key   string
	value string
}

func (prettyTableFormatter *PrettyTableFormatter) format(toFormat string, opts *Opts) string {
	var decodedJSON = make([][]string, 0)
	decodedJSON = decode(toFormat)
	normalizedTable := prettyTableFormatter.normalizeTable(decodedJSON)
	keys := ""
	seperator := ""
	values := ""
	for _, entry := range normalizedTable {
		maxLength := int(math.Max(float64(len(entry.key)), float64(len(entry.value))))
		keys = fmt.Sprintf("%s| %s%s ", keys, entry.key, strings.Repeat(" ", maxLength-len(entry.key)))
		values = fmt.Sprintf("%s| %s%s ", values, entry.value, strings.Repeat(" ", maxLength-len(entry.value)))
		seperator = fmt.Sprintf("%s--%s-", seperator, strings.Repeat("-", maxLength))
	}
	keys = fmt.Sprintf("%s|", keys)
	values = fmt.Sprintf("%s|", values)
	seperator = fmt.Sprintf("%s-", seperator)
	return fmt.Sprintf("%s\n%s\n%s\n", keys, seperator, values)

}

func (prettyTableFormatter *PrettyTableFormatter) normalizeTable(toNormalize [][]string) []Entry {
	var toReturn = make([]Entry, 0)
	keys := toNormalize[0][:]
	values := toNormalize[1:][:]
	skip := false
	for _, entry := range values {
		for i, value := range entry {
			if skip {
				skip = false
				continue
			}
			if value == "" {
				continue
			}
			if strings.Contains(keys[i], "Key") {
				toReturn = append(toReturn, Entry{value, entry[i+1]})
				skip = true
				continue
			}
			if !strings.Contains(keys[i], "Value") {
				toReturn = append(toReturn, Entry{keys[i], value})
			}
		}
	}
	return toReturn
}
