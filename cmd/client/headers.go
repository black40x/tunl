package client

import (
	"strings"
	"tunl-cli/cmd/options"
)

func ArrToHeaders(a []string, sep string) options.Headers {
	h := make(options.Headers)

	for _, s := range a {
		kv := strings.Split(s, sep)
		if len(kv) >= 2 {
			h[kv[0]] = strings.Join(kv[1:], "=")
		}
	}

	return h
}
