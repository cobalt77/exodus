package main

import (
	"os"
	"path"
	"strings"
	"syscall"
)

func main() {
	// Determine the location of this launcher executable.
	selfExe, err := os.Readlink("/proc/self/exe")
	if err != nil {
		panic(err)
	}
	currentDir := path.Dir(selfExe)

	// Prefix each segment with the current working directory so it's an absolute path.
	libraryPaths := strings.Split("{{library_path}}", ":")
	for i, libraryPath := range libraryPaths {
		libraryPaths[i] = path.Join(currentDir, libraryPath)
	}
	newLibraryPath := strings.Join(libraryPaths, ":")
	
	// Construct absolute paths to the linker and the executable that we're trying to launch.
	fullLinkerPath := path.Join(currentDir, "{{linker_dirname}}", "{{linker_basename}}")
	fullExePath := path.Join(currentDir, "{{executable}}")

	// Construct all of the arguments for the linker.
	defaultArgs := [...]string{"{{linker_basename}}", "--library-path", newLibraryPath}
	argsLen := len(defaultArgs) + len(os.Args)
	if {{full_linker}} {
		argsLen += 2
	}
	linkerArgs := make([]string, 0, argsLen) // Preallocate the slice with the number of args
	linkerArgs = append(linkerArgs, defaultArgs[:]...)
	// We can't use `--inhibit-rpath` or `--inhibit-cache` with the musl linker.
	if {{full_linker}} {
		linkerArgs = append(linkerArgs, "--inhibit-rpath", "--inhibit-cache")
	}
	linkerArgs = append(linkerArgs, fullExePath) 
	linkerArgs = append(linkerArgs, os.Args[1:]...)

	// Execute the linker.
	syscall.Exec(fullLinkerPath, linkerArgs, os.Environ())
}