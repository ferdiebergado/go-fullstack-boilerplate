package app

import (
	"context"
	"database/sql"
	"runtime"
)

type service struct {
	repo Repo
}

type Service interface {
	DBStats(context.Context) (*DBHealth, error)
	MemStats() *RAMHealth
	Ping(context.Context) error
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

type DBHealth struct {
	Status string       `json:"status"`
	Stats  *sql.DBStats `json:"stats,omitempty"`
}

type RAMHealth struct {
	Status string       `json:"status"`
	Stats  *MemoryStats `json:"stats,omitempty"`
}

// MemoryStats holds the memory usage information
type MemoryStats struct {
	Alloc      float64 `json:"alloc"`       // bytes allocated and not yet freed
	TotalAlloc float64 `json:"total_alloc"` // bytes allocated (even if freed)
	Sys        float64 `json:"sys"`         // bytes obtained from the OS
	NumGC      uint32  `json:"num_gc"`      // number of garbage collections
}

const bytesMB = 1048576

// ConvertBytesToMB converts bytes to megabytes using binary definition
func ConvertBytesToMB(bytes uint64) float64 {
	return float64(bytes) / bytesMB // 1 MB = 1,048,576 Bytes
}

func (s *service) DBStats(ctx context.Context) (*DBHealth, error) {
	if err := s.Ping(ctx); err != nil {
		return &DBHealth{
			Status: "down",
		}, err
	}

	stats := s.repo.Stats()

	return &DBHealth{
		Status: "up",
		Stats:  &stats,
	}, nil
}

func (s *service) MemStats() *RAMHealth {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := &MemoryStats{
		Alloc:      ConvertBytesToMB(memStats.Alloc),
		TotalAlloc: ConvertBytesToMB(memStats.TotalAlloc),
		Sys:        ConvertBytesToMB(memStats.Sys),
		NumGC:      memStats.NumGC,
	}

	return &RAMHealth{
		Status: "up",
		Stats:  stats,
	}
}

func (s *service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
