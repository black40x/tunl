package options

import (
	"github.com/black40x/tunl-core/tunl"
	"time"
)

type BasicAuth struct {
	Login string
	Pass  string
}

type Headers = map[string]string

type Options struct {
	ServerAddr      string
	LocalAddr       *tunl.Address
	HttpTimeout     time.Duration
	BasicAuth       *BasicAuth
	Monitor         bool
	MonitorAddr     *tunl.Address
	RequestHeaders  Headers
	ResponseHeaders Headers
	ServerPassword  string
}
