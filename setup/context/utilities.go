package context

import (
	"bufio"
	"fmt"
	"github.com/kardianos/osext"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type ByString []string

func (ip ByString) Len() int {
	return len(ip)
}

func (ip ByString) Swap(i, j int) {
	ip[i], ip[j] = ip[j], ip[i]
}

func (ip ByString) Less(i, j int) bool {
	return ip[i] < ip[j]
}

func ReadFileLines(path string) ([]string, error) {
	lines := []string{}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func CopyFile(src, dest string) error {
	r, err := os.Open(src)
	defer r.Close()
	if err != nil {
		return err
	}

	w, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	defer w.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	return nil
}

func CopyDir(source, dest string) error {
	dir, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !dir.IsDir() {
		return fmt.Errorf("Source '%s' is not a directory", source)
	}

	if err := os.MkdirAll(dest, dir.Mode()); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	for _, file := range files {
		sfp := source + "/" + file.Name()
		dfp := dest + "/" + file.Name()
		if file.IsDir() {
			if err := CopyDir(sfp, dfp); err != nil {
				return err
			}
		} else {
			if err = CopyFile(sfp, dfp); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadTFVars(varsFile string) (map[string]string, error) {
	tfvars := map[string]string{}

	if _, err := os.Stat(varsFile); err == nil {
		lines, err := ReadFileLines(varsFile)
		if err != nil {
			return nil, fmt.Errorf("Failed to get tfvars: %s", err.Error())
		}

		for _, line := range lines {
			line = strings.Replace(line, "=", "", 1)
			split := strings.Split(line, "\"")
			key := strings.Replace(split[0], " ", "", -1)
			val := split[1]
			tfvars[key] = val
		}
	}

	return tfvars, nil
}

func GetExecutionDir() (string, error) {
	// check if were executed via "go run main.go" or from a binary
	// todo: use better method
	if strings.Contains(os.Args[0], "main") {
		return os.Getwd()
	}

	return osext.ExecutableFolder()
}
