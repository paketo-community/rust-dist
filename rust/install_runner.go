package rust

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/pexec"
)

//go:generate mockery -name Executable -case=underscore

// Executable allows for mocking the pexec.Executable
type Executable interface {
	Execute(execution pexec.Execution) error
}

// InstallerRunner can run the rust installer shell script
type InstallerRunner struct {
	exec   Executable
	stdout io.Writer
	stderr io.Writer
}

// NewInstallRunner creates a new InstallRunner with a given Executable
func NewInstallRunner(exec Executable, stdout io.Writer, stderr io.Writer) InstallerRunner {
	return InstallerRunner{
		exec:   exec,
		stdout: stdout,
		stderr: stderr,
	}
}

// Install runs the install script
func (r InstallerRunner) Install(downloadDir string, destDir string, version string) error {
	extractedFiles := filepath.Join(downloadDir, fmt.Sprintf("rust-%s-x86_64-unknown-linux-gnu", version))

	path := fmt.Sprintf("%s%c%s", os.Getenv("PATH"), os.PathListSeparator, extractedFiles)
	os.Setenv("PATH", path)

	return r.exec.Execute(pexec.Execution{
		Stdout: r.stdout,
		Stderr: r.stderr,
		Dir:    extractedFiles,
		Args: []string{
			fmt.Sprintf("--prefix=%s", destDir),
			"--without=rust-docs,rls-preview,clippy-preview,miri-preview,rustfmt-preview,llvm-tools-preview,rust-analysis-x86_64-unknown-linux-gnu",
			"--disable-ldconfig",
		},
	})
}
