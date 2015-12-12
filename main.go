package main

import (
	"flag"
	"fmt"
	"github.com/exherb/Dashing"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadUrlsToDir(urls []string, dirpath string) error {
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		os.Mkdir(dirpath, os.ModePerm)
	}
	for _, url := range urls {
		filename := filepath.Base(url)
		toFilepath := filepath.Join(dirpath, filename)
		if _, err := os.Stat(toFilepath); !os.IsNotExist(err) {
			continue
		}
		fmt.Printf("\tdownloading %s ...\n", url)

		toFile, err := os.Create(toFilepath)
		if err != nil {
			return err
		}
		defer toFile.Close()
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if _, err := io.Copy(toFile, response.Body); err != nil {
			return err
		}
	}
	return nil
}

func restoreAssets(dirpath string, resource string) {
	resourcePath := filepath.Join(dirpath, resource)
	if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
		dashing.RestoreAssets(dirpath, resource)
	}
}

func prepareDependencies(dirpath string) {
	restoreAssets(dirpath, "dashboards")
	restoreAssets(dirpath, "jobs")
	restoreAssets(dirpath, "src")
	restoreAssets(dirpath, "public")
	restoreAssets(dirpath, "config.json")
	restoreAssets(dirpath, "package.json")
	restoreAssets(dirpath, "webpack.config.js")
	restoreAssets(dirpath, "server.go")
}

func build(dirpath string) error {
	prepareDependencies(dirpath)
	return nil
}

func main() {
	flag.Parse()
	action := flag.Arg(0)
	dirpath := flag.Arg(1)
	dirpath, _ = filepath.Abs(dirpath)

	fmt.Printf("%s ...\n", action)

	var err error
	switch action {
	case "new":
		err = build(dirpath)
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("done.\n")
	}
}
