package app

import (
	"database/sql"
	"runtime"

	"github.com/ferdiebergado/go-fullstack-boilerplate/internal/pkg/config"
)

type service struct {
	cfg *config.Config
}

type Service interface {
	CPUStats() *CPUHealth
	MemStats() *RAMHealth
}

func NewService(cfg *config.Config) Service {
	return &service{
		cfg: cfg,
	}
}

const bytesMB = 1048576

// ConvertBytesToMB converts bytes to megabytes using binary definition
func ConvertBytesToMB(bytes uint64) float64 {
	return float64(bytes) / bytesMB // 1 MB = 1,048,576 Bytes
}

type Health struct {
	CPU *CPUHealth `json:"cpu,omitempty"`
	RAM *RAMHealth `json:"ram,omitempty"`
}

type ComponentHealth struct {
	DB  *DBHealth `json:"db,omitempty"`
	App *Health   `json:"app,omitempty"`
}

type DBStats struct {
	Driver string      `json:"driver"`
	DB     string      `json:"db"`
	Host   string      `json:"host"`
	Port   string      `json:"port"`
	Stats  sql.DBStats `json:"stats"`
}

type DBHealth struct {
	Status string   `json:"status"`
	Stats  *DBStats `json:"stats,omitempty"`
}

type CPUHealth struct {
	Status string
	Stats  map[string]int
}

func (s *service) CPUStats() *CPUHealth {
	numCPUs := runtime.NumCPU()

	return &CPUHealth{
		Status: "up",
		Stats: map[string]int{
			"num_cpus": numCPUs,
		},
	}
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
