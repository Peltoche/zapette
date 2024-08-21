package sysinfos

import "time"

type Infos struct {
	Hostname string
	Uptime   time.Duration
}
