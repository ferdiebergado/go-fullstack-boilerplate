//go:build dev

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
)

const (
	cssEntry  = "../web/app/css/styles.css"
	cssOutDir = "../web/static/css"
	tsEntry   = "../web/app/ts/**/**/*.ts"
	tsOutdir  = "../web/static/js"
)

var ErrAssetInvalid = errors.New("invalid asset type")

func getCSSBuildOptions() api.BuildOptions {
	return api.BuildOptions{
		EntryPoints: []string{cssEntry},
		Outdir:      cssOutDir,
	}
}

func getTSBuildOptions() api.BuildOptions {
	return api.BuildOptions{
		EntryPoints: []string{tsEntry},
		Outdir:      tsOutdir,
		Target:      api.ES2015,
		Format:      api.FormatIIFE,
	}
}

func addBasicOpts(opts api.BuildOptions) api.BuildOptions {
	opts.Bundle = true
	opts.Engines = []api.Engine{
		{Name: api.EngineChrome, Version: "58"},
		{Name: api.EngineFirefox, Version: "57"},
		{Name: api.EngineSafari, Version: "11"},
		{Name: api.EngineEdge, Version: "16"},
	}
	opts.Sourcemap = api.SourceMapLinked
	opts.Write = true

	return opts
}

func enableProdOpts(opts api.BuildOptions) api.BuildOptions {
	opts.MinifyWhitespace = true
	opts.MinifyIdentifiers = true
	opts.MinifySyntax = true
	opts.Sourcemap = api.SourceMapNone

	return opts
}

func build(opts api.BuildOptions) {
	fmt.Printf("Bundling %s to %s...\n", opts.EntryPoints[0], opts.Outdir)
	result := api.Build(opts)
	if len(result.Errors) > 0 {
		log.Printf("Rebuild failed with errors: %v", result.Errors)
	} else {
		log.Println("Rebuild succeeded")
	}
}

func watch(opts api.BuildOptions) {
	// Set up a file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Watch the source directory
	sourceDir := parseSourceDir(opts.EntryPoints[0])

	err = watcher.Add(sourceDir)
	if err != nil {
		log.Printf("watch add: %s %v", sourceDir, err)
		panic("Failed to watch directory")
	}

	log.Printf("Watching for changes on %s...\n", sourceDir)

	// Debounce map to track recent events
	lastEvent := time.Now()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Debounce: only rebuild if enough time has passed
			if time.Since(lastEvent) < 100*time.Millisecond {
				continue
			}
			lastEvent = time.Now()

			// Rebuild on file changes
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				fmt.Printf("File changed: %s\n", event.Name)
				build(opts)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func parseSourceDir(entry string) string {
	const unwantedSuffix = "/**/**/*.ts"

	var sourceDir string
	if strings.HasSuffix(entry, unwantedSuffix) {
		sourceDir = strings.TrimSuffix(entry, unwantedSuffix)
	} else {
		sourceDir = filepath.Dir(entry)
	}

	return sourceDir
}

func bailout(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v", err)
	os.Exit(1)
}

func main() {
	isProd := false
	isWatched := false

	flag.BoolFunc("prod", "Bundle for production", func(_ string) error {
		isProd = true
		return nil
	})

	flag.BoolFunc("watch", "Watch for file changes", func(_ string) error {
		isWatched = true
		return nil
	})

	flag.Parse()

	cssOpts := addBasicOpts(getCSSBuildOptions())
	tsOpts := addBasicOpts(getTSBuildOptions())

	if isWatched {
		asset := os.Args[len(os.Args)-1]

		if asset != "css" && asset != "ts" {
			bailout(fmt.Errorf("%w %s", ErrAssetInvalid, asset))
		}

		switch asset {
		case "css":
			watch(cssOpts)
		case "ts":
			watch(tsOpts)
		}
	} else {
		buildOpts := []api.BuildOptions{
			cssOpts, tsOpts,
		}

		for _, opts := range buildOpts {
			if isProd {
				opts = enableProdOpts(opts)
			}

			build(opts)
		}
	}

	fmt.Println("Done.")
}
