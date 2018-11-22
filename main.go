package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%v [dirs]", os.Args[0])
	}

	var args = os.Args[1:]
	for _, arg := range args {
		err := readTree(arg)
		if err != nil {
			log.Printf("readTree %s: %v\n", arg, err)
		}
	}
}

func readTree(basePath string) error {
	// exists?
	file, err := os.Stat(basePath)
	if err != nil {
		return fmt.Errorf("could not stat %s: %v", basePath, err)
	}

	// ignore hidden files
	if file.Name()[0] == '.' {
		return nil
	}

	// isFile? handle it
	if !file.IsDir() {
		if err := handleFile(basePath); err != nil {
			return err
		}
	} else {
		// isDir! handle it
		fis, err := ioutil.ReadDir(basePath)
		if err != nil {
			return fmt.Errorf("could not read dir %s: %v", basePath, err)
		}

		// create list of files under the first dir
		var fiList []string
		for _, fi := range fis {
			if fi.Name()[0] != '.' {
				fiList = append(fiList, fi.Name())
			}
		}

		// call recursively those files
		for _, file := range fiList {
			err := readTree(filepath.Join(basePath, file))
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func handleFile(filePath string) error {
	extension := filepath.Ext(filePath)

	switch extension {
	case ".yaml":
		if err := scanYaml(filePath); err != nil {
			return err
		}
	case ".yml":
		if err := scanYaml(filePath); err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

func scanYaml(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, isPrefix, err := reader.ReadLine()

	result := make(map[string][]string)
	for err == nil && !isPrefix {
		s := string(line)
		if strings.Contains(s, "import_tasks") || strings.Contains(s, "include_tasks") || strings.Contains(s, "import_role") {
			result[filePath] = append(result[filePath], s)
		}
		line, isPrefix, err = reader.ReadLine()
	}

	for f := range result {
		if f != "" {
			fmt.Println(f)
			for i := range result[f] {
				if i == len(result[f])-1 {
					fmt.Printf("└── %s\n", result[f][i])
				} else {
					fmt.Printf("├── %s\n", result[f][i])
				}
			}
			fmt.Printf("\n")
		}
	}

	if isPrefix {
		fmt.Println("buffer error")
		return err
	}

	if err != io.EOF {
		return err
	}

	return nil
}
