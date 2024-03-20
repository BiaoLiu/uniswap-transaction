package file

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// GetSize get the file size
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)
	return len(content), err
}

// GetExt get the file ext
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckNotExist check if the file exists
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// CheckPermission check if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// MkDir create a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Open a file according to a specific mode
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// MustOpen maximize trying to open the file
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd e: %v", err)
	}
	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, e: %v", src, err)
	}
	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}
	return f, nil
}

func ReadLines(path string) ([]string, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var lines []string
	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lineCount++
	}
	return lines, lineCount, scanner.Err()
}

func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func CheckDir(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return true
}

func CheckCreateDir(dir string) error {
	//if not exist, create it
	_, err := os.Stat(dir)
	if err != nil {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func GetDirFileList(dir string, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	nameList := []string{}
	for _, f := range files {
		fileName := f.Name()
		tmp := strings.TrimSuffix(f.Name(), suffix)
		if tmp != f.Name() {
			nameList = append(nameList, fileName)
		}
	}
	return nameList, nil
}

// GetFilesBySuffix .
func GetFilesBySuffix(dirPth, suffix string) (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			subFiles, err := GetFilesBySuffix(dirPth+"/"+fi.Name(), suffix)
			if err == nil {
				files = append(files, subFiles...)
			}
		} else {
			ok := strings.HasSuffix(fi.Name(), suffix)
			if ok {
				files = append(files, dirPth+"/"+fi.Name())
			}
		}
	}
	return files, nil
}

func MoveFile(dstDir, srcDir, fileName string) error {
	oldPath, newPath := srcDir+"/"+fileName, dstDir+"/"+fileName
	//return os.Rename(oldPath, newPath)
	//fmt.Println(fmt.Sprintf("MoveFile oldPath(%v) newPath(%v)", oldPath, newPath))
	//cmd := exec.Command("cp", oldPath, newPath)
	//_, err := cmd.Output()
	//if err != nil {
	//	return err
	//}
	//os.Remove(oldPath)
	//return nil
	if err := Copy(oldPath, newPath); err != nil {
		return err
	}
	os.Remove(oldPath)
	return nil
}
