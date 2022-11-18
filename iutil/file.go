package iutil

import (
	"github.com/kardianos/osext"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//读取指定目录下的文件
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	// PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if len(suffix) == 0 || strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, strings.TrimRight(strings.TrimRight(dirPth, "\\"), "/")+string(os.PathSeparator)+fi.Name())
		}
	}

	return files, nil
}

func FileExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}

	return true, nil
}

//AbsolutePath 绝对路径
func AbsolutePath(datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(datadir, filename)
}

func GetFileInfo(path string) (os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileinfo, erri := file.Stat()
	if erri != nil {
		return nil, erri
	}
	return fileinfo, nil
}

//获取目录下的文件夹
func ListFolder(dirPth string) (folders []string, err error) {
	folders = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			folders = append(folders, fi.Name())
		}
	}

	return folders, nil
}

// WalkDir 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil {
			return err
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

//获取当前目录
func GetCurrentDirectory() string {
	// dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// return strings.Replace(dir, "\\", "/", -1)
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		folderPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			folderPath = filepath.Dir(os.Args[0])
		}
	}

	// 读取链接
	linkFolderPath, err := filepath.EvalSymlinks(folderPath)
	if err != nil {
		return folderPath
	}
	return linkFolderPath
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}

	defer dst.Close()
	return io.Copy(dst, src)
}
