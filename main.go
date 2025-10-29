package main

import (
	"debug/buildinfo"
	"flag"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func main() {
	run(parseFlags())
}

func parseFlags() (verbose bool, release bool) {
	flag.BoolVar(&verbose, "verbose", false, "enable verbose output")
	flag.BoolVar(&release, "release", false, "install in release mode (-ldflags=\"-s -w\")")
	flag.Parse()
	return verbose, release
}

func run(verbose, release bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to load user's home directory: %v", err)
	}

	gobin := filepath.Join(home, "go", "bin")
	updateBinariesAt(gobin, verbose, release)
}

func updateBinariesAt(path string, verbose, release bool) {
	err := os.Chdir(path)
	if err != nil {
		log.Fatalf("Failed to change directory to %s: %v", path, err)
	}

	binaries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	wg := errgroup.Group{}
	for _, binary := range binaries {
		wg.Go(func() error {
			if binary.IsDir() {
				return nil
			}

			return installLatestVersionOf(binary, verbose, release)
		})
	}

	if err := wg.Wait(); err != nil {
		log.Fatalf("Failed to install latest version: %v", err)
	}
}

func installLatestVersionOf(binary fs.DirEntry, verbose, release bool) error {
	name := binary.Name()
	info, err := buildinfo.ReadFile(name)
	if err != nil {
		log.Printf("Failed to read build info for %s: %v", name, err)
		return err
	}

	var ldflags string
	if release {
		ldflags = "-s -w"
	}

	cmd := exec.Command("go", "install", "-ldflags", ldflags, info.Path+"@latest") //#nosec Variables are safe.
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
