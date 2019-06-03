package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var err error
	if err != nil {
		log.Fatalln(err)
	}
	var DEST, SRC string
	args := os.Args
	if len(args) <= 3 && len(args) >= 2 {
		SRC, err = filepath.Abs(args[1])
		if err != nil {
			log.Fatalln(err)
		}
		if len(args) >= 3 {
			DEST, err = filepath.Abs(args[2])
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			DEST = SRC
		}
		if !strings.HasSuffix(DEST, ".zip") {
			DEST = DEST + ".zip"
		}

	} else {
		fmt.Println("Use zip src dist")
		os.Exit(1)
	}
	log.Println("src ", SRC)
	srcInfo, err := os.Stat(SRC)
	log.Println(srcInfo)
	if err != nil {
		log.Fatalln(err)
	}
	var paths []string
	if srcInfo.IsDir() {
		err := filepath.Walk(SRC, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else {
				if !info.IsDir() {
					if &paths == nil {
						paths = []string{path}
					} else {
						paths = append(paths, path)
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		paths = []string{SRC}
	}
	var dest *os.File

	dest, err = os.OpenFile(DEST, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new zip archive.
	w := zip.NewWriter(dest)
	for _, file := range paths {

		name := filepath.Clean(strings.Replace(file, SRC, "", 1))
		name = strings.TrimPrefix(name, "\\")
		name = strings.TrimPrefix(name, "/")
		if name == "." {
			_, name = filepath.Split(file)

		}
		log.Println(name, file)
		f, err := w.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		c, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(c)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println(args[0])
	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

}
