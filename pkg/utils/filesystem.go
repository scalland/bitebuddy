package utils

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-resty/resty/v2"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func (u *Utils) FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (u *Utils) IsDirectory(path string) bool {
	finfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	if finfo.Mode().IsDir() {
		return true
	}

	return false
}

func (u *Utils) OverwriteFile(path string, data string) (int, error) {
	var bytesWritten = 0
	f, err := os.Create(path)
	if err != nil {
		return bytesWritten, err
	}

	bytesWritten, err = f.WriteString(data)

	err = f.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("Error closing file connection %s", err.Error()))
	}

	return bytesWritten, err
}

func (u *Utils) FileOverwrite(path string, data []byte) (int64, error) {
	var bytesWritten = 0
	f, err := os.Create(path)
	if err != nil {
		return int64(bytesWritten), err
	}

	bytesWritten, err = f.Write(data)

	err = f.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("Error closing file connection %s", err.Error()))
	}

	return int64(bytesWritten), err
}

func (u *Utils) AppendToFile(path string, data string) (int, error) {
	var bytesWritten = 0
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return bytesWritten, err
	}

	bytesWritten, err = f.WriteString(data)

	err = f.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("Error closing file connection %s", err.Error()))
	}

	return bytesWritten, err
}

func (u *Utils) Watch(fileName string) {
	// Create a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		watcherErr := watcher.Close()
		if watcherErr != nil {
			slog.Error(fmt.Sprintf("Error closing watcher: %s", watcherErr.Error()))
		}
	}()

	// Set the file to watch
	filename := fileName

	// Add the file to the watcher
	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Create a channel to receive events
	events := watcher.Events

	// Start a loop to watch for events
	for {
		select {
		case event := <-events:
			// Check if the event was a modification
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Reload the file
				file, err := os.Open(filename)
				if err != nil {
					log.Fatal(err)
				}

				// Print the contents of the file
				slog.Info(fmt.Sprintf("Reloading file: %s", filename))
				fileCloseErr := file.Close()
				if fileCloseErr != nil {
					slog.Error(fmt.Sprintf("utils.Watch: error closing file handle for %s. %s", filename, fileCloseErr.Error()))
				}
			}
		case err := <-watcher.Errors:
			slog.Error(fmt.Sprintf("error: %s", err.Error()))
		}
	}
}

func (u *Utils) GetFileContentType(pathToFileName string) (string, error) {
	// Code from: https://www.tutorialspoint.com/how-to-detect-the-content-type-of-a-file-in-golang

	// Open the file whose type you
	// want to check
	file, err := os.Open(pathToFileName)

	if err != nil {
		return "application/octet-stream", err
	}

	// to sniff the content type only the first
	// 512 bytes are used.

	buf := make([]byte, 512)

	bytesRead, err := file.Read(buf)

	if err != nil {
		return "application/octet-stream", err
	}

	slog.Info(fmt.Sprintf("Read %d bytes from %s. Closing the file now", bytesRead, pathToFileName))

	fileCloseErr := file.Close()
	if fileCloseErr != nil {
		slog.Error(fmt.Sprintf("Utils.GetFileContentType()::Error closing handle for file %s. %s", pathToFileName, fileCloseErr.Error()))
	}

	// the function that actually does the trick
	contentType := http.DetectContentType(buf)

	return contentType, nil
}

func (u *Utils) CreateFileWithDataAtURL(sourceFileURL string, targetFilePath string) (*resty.Response, error) {
	slog.Info(fmt.Sprintf("Downloading file: %s", sourceFileURL))
	client := resty.New()
	//client = client.SetLogger(slog.Log)
	return client.R().SetOutput(targetFilePath).Get(sourceFileURL)
}

// CreateTempFileWithDataAtURL - Download the file at URL specified by `sourceFileURL` into the temporary directory which is
// provided by your operating system. The downloaded file will have a name beginning with the value passed in `fileNamePattern`
// Returns the path of the file that has been downloaded or error in case of one.
func (u *Utils) CreateTempFileWithDataAtURL(sourceFileURL, fileNamePattern string) (string, error) {
	slog.Info(fmt.Sprintf("Downloading file: %s", sourceFileURL))
	client := resty.New()
	//client = client.SetLogger(u.Log)
	filePtr, filePtrErr := u.TempFile(fileNamePattern)
	if filePtrErr != nil {
		return "", filePtrErr
	}
	fileNamePath := filePtr.Name()
	_, resErr := client.R().SetOutput(fileNamePath).Get(sourceFileURL)
	filePtrCloseErr := filePtr.Close()

	if filePtrCloseErr != nil {
		slog.Error(fmt.Sprintf("utils.CreateTempFileWithDataAtURL: unable to close file pointer: %s", filePtrCloseErr.Error()))
	}

	if resErr != nil {
		return fileNamePath, resErr
	}

	return fileNamePath, nil
}

// TempFile - Creates a temporary file in the directory returned by `os.TempDir()` having a name starting with `fileNamePattern`
// and returns a pointer to `os.File` struct along with error if any. See the definition of `os.CreateTemp()` for more details
// on the return value as this function uses the same to create the file.
func (u *Utils) TempFile(fileNamePattern string) (*os.File, error) {
	tempDir := os.TempDir()
	return os.CreateTemp(tempDir, fileNamePattern)
}

func (u *Utils) TempFileAtDir(tempDir, fileNamePattern string) (*os.File, error) {
	return os.CreateTemp(tempDir, fileNamePattern)
}

// CreateTempFileWithData - Create a temporary file with data  specified by `fileData` into the temporary directory which is
// provided by your operating system. The created file will have a name beginning with the value passed in `fileNamePattern`
// Returns the path of the file that has been created or error in case of one.
func (u *Utils) CreateTempFileWithData(fileData []byte, fileNamePattern string) (string, error) {
	filePtr, filePtrErr := u.TempFile(fileNamePattern)
	if filePtrErr != nil {
		return "", filePtrErr
	}
	fileNamePath := filePtr.Name()
	bytesWritten, writeErr := u.FileOverwrite(fileNamePath, fileData)
	slog.Info(fmt.Sprintf("utils.CreateTempFileWithData: wrote %d bytes to %s", bytesWritten, fileNamePath))
	if writeErr != nil {
		return fileNamePath, writeErr
	}
	filePtrCloseErr := filePtr.Close()

	if filePtrCloseErr != nil {
		slog.Error(fmt.Sprintf("utils.CreateTempFileWithDataAtURL: unable to close file pointer: %s", filePtrCloseErr.Error()))
	}

	return fileNamePath, nil
}

func (u *Utils) GetAbsolutePath(path string) (string, error) {
	absPath := ""
	if strings.HasPrefix(path, "~") {
		usr, usrErr := user.Current()
		if usrErr != nil {
			return absPath, usrErr
		}
		absPath += usr.HomeDir

		if strings.HasPrefix(path, "~/") {
			absPath = filepath.Join(absPath, path[2:])
		}

		return absPath, nil
	}
	return filepath.Abs(path)
}

func (u *Utils) ReadFileData(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from file: %s", err.Error())
	}
	return data, nil
}

// ReadFileDataLineByLine reads a file line by line
func (u *Utils) ReadFileDataLineByLine(path string) ([]string, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		// return error
		return nil, err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("utils.ReadFileDataLineByLine: unable to do a defered close of file pointer"))
		}
	}()

	// Create a new scanner
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		// scan lines
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// FileCopyPermissions - reads the permissions of the filesystem entry srcPath and applies it on to destPath
func (u *Utils) FileCopyPermissions(srcPath, destPath string) error {
	// Get the source file information
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %s", err)
	}

	// Get the source file permissions
	srcPerm := srcInfo.Mode().Perm()

	// Change the permissions of the destination file
	err = os.Chmod(destPath, srcPerm)
	if err != nil {
		return fmt.Errorf("failed to copy permissions: %s", err)
	}

	return nil
}

// FilePermissions - reads the permissions of the filesystem entry `path` and returns the same or error in case of one
func (u *Utils) FilePermissions(path string) (os.FileMode, error) {
	// Get the file information
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get source file info: %s", err)
	}

	// return the file permissions
	return fileInfo.Mode().Perm(), nil
}
