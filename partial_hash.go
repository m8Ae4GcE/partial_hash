package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// Bar ...
type Bar struct {
	percent int64  // progress percentage
	cur     int64  // current progress
	total   int64  // total value for progress
	rate    string // the actual progress bar to be printed
	graph   string // the fill value for progress bar
}

// Define some variables to receive command-line values
var filename string
var path_to_file string
var buffer_size uint64
var export_format string

const YYYYMMDDhhmmss = "2006-01-02_15-04-05"

var export_filename string
var version = "2.0"

// Array to store the usage order
var UsageOrder []string

func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "#"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 1 {
		bar.rate += bar.graph // initial progress position
	}
}

func (bar *Bar) getPercent() int64 {
	return int64((float32(bar.cur) / float32(bar.total)) * 50)
}

func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last {
		var i int64 = 0
		for ; i < bar.percent-last; i++ {
			bar.rate += bar.graph
		}
		fmt.Printf("\r[%-50s]%3d%% %8d/%d", bar.rate, bar.percent*2, bar.cur, bar.total)
	}
}

func (bar *Bar) Finish() {
	fmt.Println()
}

func main() {

	// -h
	flag.Usage = parameters
	// Setup flags and parse command line parameters
	set_parameters()
	flag.Parse()

	now := time.Now().UTC()

	//<YYYY-MM-DD_hh-mm-ss>.partial_hash_<version>.buffer_<buffer_size_in_bytes>.csv
	export_filename = (now.Format(YYYYMMDDhhmmss)) + ".partial_hash_" + version + ".buffer_" + fmt.Sprint(buffer_size) + "." + export_format

	// Conditions
	if (path_to_file != "" && filename != "") || (len(os.Args) < 2) || (path_to_file == "" && filename == "") {
		parameters()
		return
	} else if filename != "" {
		if export_format == "stdout" {
			stdout_header()
		} else {
			fmt.Print("Output format ", export_format, " does not exist or not possible", "\n\n")
			parameters()
			return
		}
		calc_hash_file(buffer_size, filename, false)

	} else if path_to_file != "" {
		if export_format == "stdout" {
			stdout_header()
		} else if export_format == "csv" {
			create_output_file(export_filename)
			stdout_header()
		} else {
			fmt.Print("Output format ", export_format, " does not exist", "\n\n")
			parameters()
			return
		}
		recusive(buffer_size, path_to_file)
	}
}

func stdout_header() {
	fmt.Println("Version : ", version)
	fmt.Println("Size of the buffer (in bytes) : ", buffer_size)
	fmt.Println("")
}

func create_output_file(export_filename string) {

	// create the file
	f, err := os.Create(export_filename)
	if err != nil {
		fmt.Println(err)
	}
	// close the file with defer
	defer f.Close()
	f.WriteString("filename;hash;is_partial_hash\n")

}

func parameters() {
	fmt.Println("Version : ", version)
	fmt.Println("Size of the buffer (in bytes) : ", buffer_size)
	fmt.Println("Directory : ", path_to_file)
	fmt.Println("Filename : ", filename)
	fmt.Println("Output format : ", export_format)
	fmt.Println("")

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
	fmt.Printf("Usage %s [[-f <file>] or [-d <recursive_folder>]] [-b <buffer_size>] [--output csv] :\n", os.Args[0])
	for s := range UsageOrder {
		fmt.Println(usageMap[UsageOrder[s]])
	}
}

func set_parameters() {
	// Set usage order for display
	UsageOrder = []string{"file", "f", "directory", "d", "buffer", "b", "output"}

	flag.StringVar(&filename, "file", filename, "File to be hashed.")
	flag.StringVar(&filename, "f", filename, "")

	flag.StringVar(&path_to_file, "directory", path_to_file, "Folder where files will be recursively hashed.")
	flag.StringVar(&path_to_file, "d", path_to_file, "")

	flag.Uint64Var(&buffer_size, "buffer", 100000000, "Size of the buffer in bytes.")
	flag.Uint64Var(&buffer_size, "b", 100000000, "")

	flag.StringVar(&export_format, "output", "stdout", "Output format.")
}

func recusive(buffer_size uint64, path_to_file string) {
	fsys := os.DirFS(path_to_file)

	var bar Bar

	var count_file int64 = 0
	var in_queue_file int64 = 0

	fs.WalkDir(fsys, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if info.IsDir() == false && path != export_filename {
			count_file++
		}
		return nil
	})
	if export_format == "csv" {
		bar.NewOption(0, count_file)
	}

	fs.WalkDir(fsys, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if info.IsDir() == false && path != export_filename {
			if runtime.GOOS == "windows" {
				path = strings.ReplaceAll(path, "/", "\\")
			}
			uniq_file := path_to_file + path
			in_queue_file++
			calc_hash_file(buffer_size, uniq_file, true)
			if export_format == "csv" {
				bar.Play(int64(in_queue_file))
			}
		}

		return nil
	})
	if export_format == "csv" {
		bar.Finish()
	}

}

func calc_hash_file(buffer_size uint64, uniq_file string, is_recursive bool) {
	fi, err := os.Stat(uniq_file)

	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	// get the size
	size := fi.Size()

	if size < 2*int64(buffer_size) {
		// open file
		file, err := os.Open(uniq_file)
		if err != nil {
			log.Printf("Err: %v\n", err)
			return
		}
		defer file.Close()

		buf := make([]byte, 30*1024)
		md5 := md5.New()
		for {
			n, err := file.Read(buf)
			if n > 0 {
				_, err := md5.Write(buf[:n])
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

		hash := md5.Sum(nil)

		if export_format == "csv" && is_recursive == true {
			output_csv := fmt.Sprintf("%s;%x;%t\n", uniq_file, hash, false)
			f, err := os.OpenFile(export_filename, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}

			defer f.Close()

			if _, err = f.WriteString(output_csv); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("%s    %x\n", uniq_file, hash)
		}

	} else {
		// open file
		file, err := os.Open(uniq_file)
		if err != nil {
			log.Printf("Err: %v\n", err)
			return
		}
		defer file.Close()

		s := []byte("")

		buffer := make([]byte, buffer_size)
		first_bytes, err := file.Read(buffer)
		byte_first_bytes := (buffer[:first_bytes])
		s = append(s, byte_first_bytes...)

		_, err = file.Seek(int64(-buffer_size), io.SeekEnd)

		last_bytes, err := file.Read(buffer)
		byte_last_bytes := (buffer[:last_bytes])
		s = append(s, byte_last_bytes...)

		//calc hash
		hash := md5.Sum(s)

		if export_format == "csv" && is_recursive == true {
			output_csv := fmt.Sprintf("%s;%x;%t\n", uniq_file, hash, true)
			f, err := os.OpenFile(export_filename, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}

			defer f.Close()

			if _, err = f.WriteString(output_csv); err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("%s    %x\n", uniq_file, hash)
		}
	}

}
