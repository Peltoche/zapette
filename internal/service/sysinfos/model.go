package sysinfos

import "time"

type Infos struct {
	hostname string
	uptime   time.Duration
}

func (i *Infos) Hostname() string {
	return i.hostname
}

func (i *Infos) Uptime() time.Duration {
	return i.uptime
}
