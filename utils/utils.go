package utils

import (
	"archive/zip"
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Home, _ = os.UserHomeDir()
var SdkmanPath = filepath.Join(Home, ".gosdkman")

func SetEnv(variable string, value string) error {
	cmd := exec.Command("setx", variable, value)
	_, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadFile(url string) (filename string, err error) {
	start := time.Now()

	request, _ := http.NewRequest("GET", url, nil)
	filename = path.Base(request.URL.Path)

	if !Exists(SdkmanPath) {
		os.Mkdir(SdkmanPath, os.FileMode(0755))
	}

	file := filepath.Join(SdkmanPath, filename)

	// Create the file
	out, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	headResp, err := http.Head(url)

	if err != nil {
		panic(err)
	}
	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

	if err != nil {
		panic(err)
	}

	done := make(chan int64)

	go PrintDownloadPercent(done, filename, file, int64(size))

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	done <- n

	elapsed := time.Since(start)
	log.Printf("\nDownload completed in %s", elapsed)

	return filename, nil
}

func PrintDownloadPercent(done chan int64, filename string, path string, total int64) {
	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			var percent = float64(size) / float64(total) * 100

			//fmt.Printf("%.0f", percent)
			//fmt.Println("%")

			fmt.Printf("\r%s", strings.Repeat(" ", 100))
			fmt.Printf("\rDownloading %s [Total Size: %s]... %s", filename, humanize.Bytes(uint64(total)), fmt.Sprintf("%.2f%%", percent) /**/)
		}

		if stop {
			fmt.Printf("\r%s", strings.Repeat(" ", 100))
			fmt.Printf("\rDownloading %s [Total Size: %s]... 100%%\n", filename, humanize.Bytes(uint64(total)))
			break
		}

		time.Sleep(time.Second)
	}
}

/*
exists returns whether the given file or directory exists
*/
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
