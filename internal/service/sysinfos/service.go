package sysinfos

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Peltoche/zapette/internal/tools"
	"github.com/Peltoche/zapette/internal/tools/clock"
	"github.com/spf13/afero"
)

const (
	uptimeFilePath = "/proc/uptime"
	hostnamePath   = "/etc/hostname"
)

type service struct {
	fs       afero.Fs
	clock    clock.Clock
	uptime   time.Duration
	hostname string
}

func newService(fs afero.Fs, tools tools.Tools) *service {
	return &service{
		fs:    fs,
		clock: tools.Clock(),
	}
}

func (s *service) fetch(ctx context.Context) error {
	uptime, err := s.fetchUptime(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch the uptime: %w", err)
	}

	hostname, err := s.fetchHostname(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch the hostname: %w", err)
	}

	s.uptime = uptime
	s.hostname = hostname

	return nil
}

func (s *service) fetchHostname(_ context.Context) (string, error) {
	rawFile, err := afero.ReadFile(s.fs, hostnamePath)
	if err != nil {
		return "", fmt.Errorf("failed to read %q: %w", hostnamePath, err)
	}

	return strings.TrimSpace(string(rawFile)), nil
}

func (s *service) fetchUptime(_ context.Context) (time.Duration, error) {
	rawFile, err := afero.ReadFile(s.fs, uptimeFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read %q: %w", uptimeFilePath, err)
	}

	uptime, err := time.ParseDuration(strings.Split(string(rawFile), " ")[0] + "s")
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q: %w", uptime, err)
	}

	return uptime, nil
}

func (s *service) GetInfos(ctx context.Context) (*Infos, error) {
	return nil, nil
}
