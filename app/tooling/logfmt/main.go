package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	service            string
	mandatoryLogFields = [6]string{"service", "ts", "level", "traceid", "caller", "msg"}
)

func init() {
	flag.StringVar(&service, "service", "", "filter which service to see")
}

func main() {
	flag.Parse()

	b := strings.Builder{}

	scanner := bufio.NewScanner(os.Stdin)

	// Scan standard input for log data per line.
	for scanner.Scan() {
		s := scanner.Text()

		// Convert the JSON to a map for processing.
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}

			continue
		}

		// If a service filter was provided, check.
		if service != "" && m["service"] != service {
			continue
		}

		// // I like always having a traceid present in the logs.
		// traceID := "00000000-0000-0000-0000-000000000000"
		// if v, ok := m["traceid"]; ok {
		// 	traceID = fmt.Sprintf("%v", v)
		// }

		// Build out the know portions of the log in the order
		// I want them in.
		b.Reset()

		for _, field := range mandatoryLogFields {
			if _, ok := m[field]; ok {
				b.WriteString(fmt.Sprintf("%s: ", m[field]))
			}
		}

		// Add the rest of the keys ignoring the ones we already
		// added for the log.
		for k, v := range m {
			if isFieldInMandatoryLogFields(k) {
				continue
			}

			// It's nice to see the key [value] in this format
			// especially since map ordering is random.
			b.WriteString(fmt.Sprintf("%s[%v]: ", k, v))
		}

		// Write the new log format, removing the last :
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func isFieldInMandatoryLogFields(field string) bool {
	for _, f := range mandatoryLogFields {
		if f == field {
			return true
		}
	}

	return false
}
