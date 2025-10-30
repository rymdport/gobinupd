package main

import (
	"debug/buildinfo"
	"flag"
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
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&release, "release", false, "install in release mode (-ldflags=\"-s -w\")")
	flag.BoolVar(&release, "r", false, "")
	flag.Parse()
	return verbose, release
}

func run(verbose, release bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to load user's home directory: %v", err)
	}

	gobin := filepath.Join(home, "go", "bin")
	updater := updater{verbose: verbose, release: release}
	updater.updateBinariesAt(gobin)
}

type updater struct {
	verbose bool
	release bool
}

func (u *updater) updateBinariesAt(path string) {
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

			return u.installLatestVersionOf(binary.Name())
		})
	}

	if err := wg.Wait(); err != nil {
		log.Fatalf("Failed to install latest version: %v", err)
	}
}

func (u *updater) installLatestVersionOf(name string) error {
	info, err := buildinfo.ReadFile(name)
	if err != nil {
		log.Printf("Failed to read build info for %s: %v", name, err)
		return err
	}

	var ldflags string
	if u.release {
		ldflags = "-s -w"
	}

	cmd := exec.Command("go", "install", "-ldflags", ldflags, info.Path+"@latest") //#nosec Variables are safe.
	if u.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
