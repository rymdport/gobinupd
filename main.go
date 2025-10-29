package main

import (
	"debug/buildinfo"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

func main() {
	run()
}

func run() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to load user's home directory: %v", err)
	}

	gobin := filepath.Join(home, "go", "bin")
	updateBinariesAt(gobin)
}

func updateBinariesAt(path string) {
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

			return installLatestVersionOf(binary)
		})
	}

	if err := wg.Wait(); err != nil {
		log.Fatalf("Failed to install latest version: %v", err)
	}
}

func installLatestVersionOf(binary fs.DirEntry) error {
	name := binary.Name()
	info, err := buildinfo.ReadFile(name)
	if err != nil {
		log.Printf("Failed to read build info for %s: %v", name, err)
		return err
	}

	cmd := exec.Command("go", "install", info.Path+"@latest")
	return cmd.Run()
}
