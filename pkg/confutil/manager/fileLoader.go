package manager

import (
	"io/ioutil"
	"path"
	"strings"
)

type fileLoader struct {
	dir string
}

func NewFileLoader(dir string) *fileLoader {
	l := new(fileLoader)
	l.dir = strings.TrimRight(dir, "/")
	return l
}

func readDir(filePath string, configSlice []string) ([]string, error) {
	fileList := make([]string, 0)
	var err error
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		return fileList, err
	}
	for _, v := range files {
		if v.IsDir() {
			fileListTmp, err := readDir(filePath+"/"+v.Name(), configSlice)
			if err != nil {
				return fileList, err
			}
			fileList = mergeSlice(fileList, fileListTmp)
			continue
		}
		if contains(strings.Replace(path.Ext(v.Name()), ".", "", -1), configSlice) {
			fileList = append(fileList, filePath+"/"+v.Name())
		}
	}
	return fileList, err
}

func mergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

func (l *fileLoader) Read(configSlice []string) (map[string][]byte, error) {
	var err error
	data := make(map[string][]byte)
	files, err := readDir(l.dir, configSlice)
	if err != nil {
		return data, err
	}
	for _, v := range files {
		data[strings.Replace(v, l.dir+"/", "", -1)], err = ioutil.ReadFile(v)
		if err != nil {
			return data, err
		}

	}
	return data, err
}
