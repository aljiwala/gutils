// FIXME:
// - Doc is not proper.
// - Test cases are remaining.

// Package filepathx provides extra utilities for file and filepath related operations.
package filepathx

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aljiwala/gutils/strx"
)

var curpath = SelfDir()

type dirNames struct {
	names []string
	err   error
}

// Dir gets directory name of the filepath.
func Dir(file string) string {
	return path.Dir(file)
}

// Ext should return extension of the given file.
func Ext(file string) string {
	return path.Ext(file)
}

// Rename should rename the file.
func Rename(file string, to string) error {
	return os.Rename(file, to)
}

// Remove should remove the file.
func Remove(file string) error {
	return os.Remove(file)
}

// Basename gets base name of provided filepath.
func Basename(file string) string {
	return path.Base(file)
}

// SelfPath gets compiled executable file absolute path.
func SelfPath() (path string) {
	path, _ = filepath.Abs(os.Args[0])
	return
}

// SelfDir gets compiled executable file directory.
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// SelfChdir switch the working path to my own path.
func SelfChdir() {
	if err := os.Chdir(curpath); err != nil {
		log.Fatal(err)
	}
}

// IsDir should return true and nil as error if provided path is of directory.
func IsDir(path string) (bool, error) {
	// Stat returns a FileInfo describing the named file.
	// If there is an error, it will be of type *PathError.
	fileInfo, err := os.Stat(path)
	if err != nil {
		// no such file or dir
		return false, errors.New("No such file or directory")
	}

	// Return true if it's directory.
	if fileInfo.IsDir() {
		return true, nil
	}
	return false, nil
}

// GetWDPath gets the work directory path.
func GetWDPath() (wd string) {
	wd = os.Getenv("GOPATH")
	if wd == "" {
		panic("GOPATH is not setted in env.")
	}
	return
}

// DoesExist checks whether a file or directory exists. It returns false if not.
func DoesExist(path string) bool {
	// Stat returns a FileInfo describing the named file.
	// If there is an error, it will be of type *PathError.
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// DoesNotExist should check if file/folder doesn't exist on provided path.
func DoesNotExist(path string) bool {
	// Stat returns a FileInfo describing the named file.
	// If there is an error, it will be of type *PathError.
	_, err := os.Stat(path)

	// If the file doesn't exists, we will get an error.
	// Thus, we can use this to check:
	if err != nil {
		// IsNotExist returns a boolean indicating whether the error is known to
		// report that a file or directory does not exist.
		// It's satisfied by ErrNotExist as well as some syscall errors.
		if os.IsNotExist(err) {
			// Not exist.
			return true
		}
	}

	// Exist.
	return false
}

// BasePath should return base file path.
func BasePath(path string) string {
	n := strings.LastIndexByte(path, '.')
	if n > 0 {
		return path[:n]
	}
	return path
}

// RelPath returns a relative path that is lexically equivalent to targpath.
func RelPath(targpath string) string {
	basepath, _ := filepath.Abs("./")
	rel, _ := filepath.Rel(basepath, targpath)
	return strings.Replace(rel, `\`, `/`, -1)
}

// // RealPath gets absolute filepath, based on built executable file.
// func RealPath(file string) (string, error) {
// 	if path.IsAbs(file) {
// 		return file, nil
// 	}
// 	wd, err := os.Getwd()
// 	return path.Join(wd, file), err
// }

// FileMTime should get and return modified time of the file.
func FileMTime(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.ModTime().Unix(), nil
}

// FileSize should get return size as bytes.
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exist.
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// DirsUnder should return list of dirs under given dirPath.
func DirsUnder(dirPath string) (container []string, err error) {
	var fs []os.FileInfo
	if !DoesExist(dirPath) {
		return
	}

	fs, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	sz := len(fs)
	if sz == 0 {
		return
	}

	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				container = append(container, name)
			}
		}
	}

	return
}

// FilesUnder should return list of files under given dirPath.
func FilesUnder(dirPath string) (container []string, err error) {
	var fs []os.FileInfo
	if !DoesExist(dirPath) {
		return
	}

	fs, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	sz := len(fs)
	if sz == 0 {
		return
	}

	for i := 0; i < sz; i++ {
		if !fs[i].IsDir() {
			container = append(container, fs[i].Name())
		}
	}

	return
}

// SearchFile searches a file in paths.
// Often used in search config file in `/etc`.
func SearchFile(filename string, paths ...string) (fullPath string, err error) {
	for _, path := range paths {
		if fullPath = filepath.Join(path, filename); DoesExist(fullPath) {
			return
		}
	}
	err = errors.New(fullPath + " not found in paths")
	return
}

// IsDirExists should return true if provided path is directory.
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}

// IsBinaryExist searches for an executable binary named file in the directories.
func IsBinaryExist(binary string) bool {
	if _, err := exec.LookPath(binary); err != nil {
		return false
	}
	return true
}

// CreateDir should create directory (if it doesn't exist) based on provided path.
func CreateDir(dirStr string) error {
	if _, err := os.Stat(dirStr); os.IsNotExist(err) {
		err := os.MkdirAll(dirStr, 0777)
		if err != nil {
			log.Printf("Failed to create directory `%s`: %v", dirStr, err)
			return err
		}
	}
	return nil
}

// FilenameReplace should replace illegal filename with similar characters.
func FilenameReplace(filename string) (rfn string) {
	// Replace “” with "".
	if strings.Count(filename, `"`) > 0 {
		var i = 1
	label:
		for k, v := range []byte(filename) {
			if string(v) != `"` {
				continue
			}
			if i%2 == 1 {
				filename = string(filename[:k]) + `“` + string(filename[k+1:])
			} else {
				filename = string(filename[:k]) + `”` + string(filename[k+1:])
			}
			i++
			goto label
		}
	}

	replace := strings.Replace
	rfn = replace(filename, `:`, `：`, -1)
	rfn = replace(rfn, `*`, `ж`, -1)
	rfn = replace(rfn, `<`, `＜`, -1)
	rfn = replace(rfn, `>`, `＞`, -1)
	rfn = replace(rfn, `?`, `？`, -1)
	rfn = replace(rfn, `/`, `／`, -1)
	rfn = replace(rfn, `|`, `∣`, -1)
	rfn = replace(rfn, `\`, `╲`, -1)
	return
}

// ExcelSheetNameReplace should replace the illegal characters in the excel
// worksheet name with the underscore.
func ExcelSheetNameReplace(filename string) (rfn string) {
	us := StrUnderscore
	replace := strings.Replace
	rfn = replace(filename, `:`, us, -1)
	rfn = replace(rfn, `：`, ``, -1)
	rfn = replace(rfn, `*`, us, -1)
	rfn = replace(rfn, `?`, us, -1)
	rfn = replace(rfn, `？`, us, -1)
	rfn = replace(rfn, `/`, us, -1)
	rfn = replace(rfn, `／`, us, -1)
	rfn = replace(rfn, `\`, us, -1)
	rfn = replace(rfn, `╲`, us, -1)
	rfn = replace(rfn, `]`, us, -1)
	rfn = replace(rfn, `[`, us, -1)
	return
}

func grepFile(pattern, filename string) (lines []string, err error) {
	var isLongLine bool
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}

	prefix := ""
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)

	for {
		byteLine, isPrefix, er := reader.ReadLine()

		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}

		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}

		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}

	return
}

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(pattern string, filename string) (lines []string, err error) {
	return grepFile(pattern, filename)
}

// walk recursively descends path, calling w.
func walk(path string, info os.FileInfo, walkFn filepath.WalkFunc, followSymlinks bool) error {
	stat := os.Lstat
	if followSymlinks {
		stat = os.Stat
	}
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	c, err := readDirNames(path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for names := range c {
		if names.err != nil {
			return walkFn(path, info, names.err)
		}
		for _, name := range names.names {
			filename := filepath.Join(path, name)
			fileInfo, err := stat(filename)
			if err != nil {
				if err = walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
					return err
				}
			} else {
				err = walk(filename, fileInfo, walkFn, followSymlinks)
				if err != nil {
					if !fileInfo.IsDir() || err != filepath.SkipDir {
						return err
					}
				}
			}
		}
	}
	return nil
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked UNORDERED,
// which makes the output undeterministic!
// Walk does not follow symbolic links.
func Walk(root string, walkFn filepath.WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(root, info, walkFn, false)
}

// WalkWithSymlinks walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked UNORDERED,
// which makes the output undeterministic!
// WalkWithSymlinks does follow symbolic links!
func WalkWithSymlinks(root string, walkFn filepath.WalkFunc) error {
	info, err := os.Stat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(root, info, walkFn, true)
}

// WalkDirs traverses the directory, return to the relative path.
// You can specify the suffix.
func WalkDirs(targpath string, suffixes ...string) (dirlist []string) {
	if !filepath.IsAbs(targpath) {
		targpath, _ = filepath.Abs(targpath)
	}

	err := filepath.Walk(targpath, func(retpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !f.IsDir() {
			return nil
		}

		if len(suffixes) == 0 {
			dirlist = append(dirlist, RelPath(retpath))
			return nil
		}

		_retpath := RelPath(retpath)
		for _, suffix := range suffixes {
			if strings.HasSuffix(_retpath, suffix) {
				dirlist = append(dirlist, _retpath)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("utils.WalkRelDirs: %v\n", err)
		return
	}

	return
}

// Converters, Readers & Writers -----------------------------------------------

// ReadFileToBytes reads data type '[]byte' from file by given path. It returns
// error when fail to finish operation.
func ReadFileToBytes(filePath string) (b []byte, err error) {
	b, err = ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	return
}

// ReadFileToString reads data type 'string' from file by given path. It returns
// error when fail to finish operation.
func ReadFileToString(filePath string, escapeNL bool) (string, error) {
	b, err := ReadFileToBytes(filePath)
	if err != nil {
		return "", err
	}

	str := string(b)
	if escapeNL {
		return strx.TrimRightSpace(str), nil
	}
	return str, nil
}

// Helpers ---------------------------------------------------------------------

// readDirNames reads the directory named by dirname and returns
// a channel for future results.
func readDirNames(dirname string) (<-chan dirNames, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	c := make(chan dirNames)

	go func() {
		defer f.Close()
		defer close(c)

		for {
			names, err := f.Readdirnames(1024)
			if err != nil {
				if err == io.EOF {
					if len(names) > 0 {
						c <- dirNames{names: names}
					}
					return
				}
				c <- dirNames{err: err}
				return
			}
			c <- dirNames{names: names}
		}
	}()

	return c, nil
}

// -----------------------------------------------------------------------------
