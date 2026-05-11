package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/go/packages"
)

const mainTemplate = `package main

import (
	"fmt"
	target "{{.ImportPath}}"
)

func main() {
	result := target.Setup()
	fmt.Println(result)
}
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: meh <package-folder>")
		os.Exit(1)
	}

	folderPath := os.Args[1]

	absFolder, err := filepath.Abs(folderPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving path: %v\n", err)
		os.Exit(1)
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedModule,
		Dir:  absFolder,
	}, ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading package: %v\n", err)
		os.Exit(1)
	}
	if len(pkgs) == 0 || pkgs[0].PkgPath == "" {
		fmt.Fprintln(os.Stderr, "could not determine import path")
		os.Exit(1)
	}

	importPath := pkgs[0].PkgPath
	moduleDir := pkgs[0].Module.Dir

	outDir := filepath.Join(moduleDir, ".meh")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating .meh directory: %v\n", err)
		os.Exit(1)
	}

	outFile := filepath.Join(outDir, "main.go")
	f, err := os.Create(outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating main.go: %v\n", err)
		os.Exit(1)
	}

	tmpl := template.Must(template.New("main").Parse(mainTemplate))
	if err := tmpl.Execute(f, map[string]string{
		"ImportPath": importPath,
	}); err != nil {
		f.Close()
		fmt.Fprintf(os.Stderr, "error writing main.go: %v\n", err)
		os.Exit(1)
	}
	f.Close()

	cmd := exec.Command("go", "run", outFile)
	cmd.Dir = moduleDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running generated file: %v\n", err)
		os.Exit(1)
	}
}
