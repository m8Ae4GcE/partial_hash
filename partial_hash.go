//v1.2

package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"
)

// Define some variables to receive command-line values
var filename string
var path_to_file string
var buffer_size uint64

// Array to store the usage order
var UsageOrder []string

func main() {

	// -h
	flag.Usage = Parameters
	// Setup flags and parse command line parameters
	setParameters()
	flag.Parse()

	// Conditions
	if (path_to_file) != "" && (filename) != "" || len(os.Args) < 2 {
		Version()
		Parameters()
		return
	} else if filename != "" {
		Version()
		fmt.Println("Size of the buffer (in bytes): ", buffer_size)
		calc_hash_file(buffer_size, filename, 1, 1)

	} else if path_to_file != "" {
		Version()
		fmt.Println("Size of the buffer (in bytes): ", buffer_size)
		recusive(buffer_size, path_to_file)
	}

}

func Version() {
	fmt.Println("Version : 1.2")
}

func Parameters() {

	if len(UsageOrder) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Attempting to use Parameters, but UsageOrder is not set\n")
		flag.VisitAll(func(f *flag.Flag) {
			// append f.Name to UsageOrder
			UsageOrder = append(UsageOrder, f.Name)
		})
	}

	usageMap := make(map[string]string, 0)

	// Loop through all defined flags, set or unset
	flag.VisitAll(func(f *flag.Flag) {
		if len(f.Name) > 1 && f.DefValue != "" {
			// Longname usage with default value
			usageMap[f.Name] = fmt.Sprintf("\n  --%s \t%s (Default is %s)", f.Name, f.Usage, f.DefValue)
		} else if len(f.Name) > 1 {
			// Longname usage without default value
			usageMap[f.Name] = fmt.Sprintf("\n  --%s \t%s", f.Name, f.Usage)
		} else {
			// Shorthand usage
			usageMap[f.Name] = fmt.Sprintf("\t  -%s  \t%s", f.Name, f.Usage)
		}
	})
	fmt.Printf("Usage %s [-f <file>] or [-d <recursive_folder>] [-b <buffer_size>]:\n", os.Args[0])
	for s := range UsageOrder {
		fmt.Println(usageMap[UsageOrder[s]])
	}
}

func setParameters() {
	// Set usage order for display
	UsageOrder = []string{"file", "f", "directory", "d", "buffer", "b"}

	flag.StringVar(&filename, "file", filename, "File to be hashed.")
	flag.StringVar(&filename, "f", filename, "")

	flag.StringVar(&path_to_file, "directory", path_to_file, "Folder where files will be recursively hashed.")
	flag.StringVar(&path_to_file, "d", path_to_file, "")

	flag.Uint64Var(&buffer_size, "buffer", 100000000, "Size of the buffer in bytes.")
	flag.Uint64Var(&buffer_size, "b", 100000000, "")

}

func recusive(BYTES_TO_READ uint64, path_to_file string) {
	fsys := os.DirFS(path_to_file)

	count_file := 0
	in_queue_file := 0

	fs.WalkDir(fsys, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if info.IsDir() == false {
			count_file++
		}
		return nil
	})

	fs.WalkDir(fsys, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if info.IsDir() == false {
			if runtime.GOOS == "windows" {
				path = strings.ReplaceAll(path, "/", "\\")
			}
			uniq_file := path_to_file + path
			in_queue_file++
			calc_hash_file(BYTES_TO_READ, uniq_file, in_queue_file, count_file)
		}
		return nil
	})

}

func calc_hash_file(BYTES_TO_READ uint64, uniq_file string, in_queue_file int, count_file int) {
	fi, err := os.Stat(uniq_file)

	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	// get the size
	size := fi.Size()

	if size < 2*int64(BYTES_TO_READ) {
		// open file
		file, err := os.Open(uniq_file)
		if err != nil {
			log.Printf("Err: %v\n", err)
			return
		}
		defer file.Close()

		buf := make([]byte, 30*1024)
		sha256 := sha256.New()
		for {
			n, err := file.Read(buf)
			if n > 0 {
				_, err := sha256.Write(buf[:n])
				if err != nil {
					log.Fatal(err)
				}
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("Read %d bytes: %v", n, err)
				break
			}
		}

		hash := sha256.Sum(nil)
		fmt.Printf("%s;%x;%d/%d\n", uniq_file, hash, in_queue_file, count_file)

	} else {
		// open file
		file, err := os.Open(uniq_file)
		if err != nil {
			log.Printf("Err: %v\n", err)
			return
		}
		defer file.Close()

		s := []byte("")

		buffer := make([]byte, BYTES_TO_READ)
		first_bytes, err := file.Read(buffer)
		byte_first_bytes := (buffer[:first_bytes])
		s = append(s, byte_first_bytes...)

		_, err = file.Seek(int64(-BYTES_TO_READ), io.SeekEnd)

		last_bytes, err := file.Read(buffer)
		byte_last_bytes := (buffer[:last_bytes])
		s = append(s, byte_last_bytes...)

		//calc hash
		hash := sha256.Sum256(s)
		fmt.Printf("%s;%x;%d/%d\n", uniq_file, hash, in_queue_file, count_file)

	}

}
