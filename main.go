package main

import (
	"debug/buildinfo"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rymdport/easypgo"
	"golang.org/x/sync/errgroup"
)

func main() {
	stop := easypgo.Generate()
	defer stop()

	run(parseFlags())
}

func parseFlags() (set flags) {
	flag.BoolVar(&set.verbose, "verbose", false, "enable verbose output")
	flag.BoolVar(&set.verbose, "v", false, "")
	flag.BoolVar(&set.release, "release", false, "install in release mode (-ldflags=\"-s -w\")")
	flag.BoolVar(&set.release, "r", false, "")
	flag.BoolVar(&set.noUpdate, "no-update", false, "no update, only rebuild")
	flag.BoolVar(&set.noUpdate, "n", false, "")
	flag.Parse()
	return set
}

type flags struct {
	verbose  bool
	release  bool
	noUpdate bool
}

func run(set flags) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to load user's home directory: %v", err)
	}

	gobin := filepath.Join(home, "go", "bin")
	updater := updater{options: set}
	updater.updateBinariesAt(gobin)
}

type updater struct {
	options flags
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

	ldflags := ""
	if u.options.release {
		ldflags = "-s -w"
	}

	version := "@latest"
	if u.options.noUpdate {
		version = "@" + info.Main.Version
	}

	cmd := exec.Command("go", "install", "-ldflags", ldflags, info.Path+version) //#nosec Variables are safe.
	if u.options.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
