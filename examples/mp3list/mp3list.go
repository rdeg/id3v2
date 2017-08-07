package main

import (
	"flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
	
	"github.com/rdeg/id3v2"
)

var (
	srcDir 		string		// source directory
	mp3Files	uint		// count of MP3 files
)

// Process a MP3 file.
func doMP3(fname string, osfi os.FileInfo) error {
	// Retrieve all the tags of the MP3 file header in a MP3Info.
	mi, err := id3v2.ProcessAllTags(fname)	// mi is a pointer to a MP3Info
	if err != nil {
		return err
	}

	// Make an MP3 file name using the tags previously harvested.
	full := mi.MakeFileName()
	if full == ".mp3" {	// no tags in this file!
//		full = fname	// preserve source path and name
		full = osfi.Name()	// preserve source base name
	}
	fmt.Printf(">>> %s\n", full)
	fmt.Printf("<<< %s\n", fname)
	fmt.Printf(">>> %d kpbs - %s - %d bytes\n\n", mi.BitRate, osfi.ModTime(), osfi.Size())
		
	return nil
}

// Walk the source file tree.
func doWalk(fname string, info os.FileInfo, err error) error {
    if err != nil {
        log.Print(err)
        return nil
    }
//fmt.Printf("fname = %s, info.Name() = %s\n", fname, info.Name())
	
	if id3v2.Verbose >= 1 {
		fmt.Println(fname)		// display the whole source filename
	}
	if filepath.Ext(fname) == ".mp3" {
		mp3Files++
		return doMP3(fname, info)
	}
    return nil
}

func main() {
	// Setup options and args.
	flag.UintVar(&id3v2.Verbose, "v", 0, "verbosity level (0 = none, 1 = tags, 2 = headers)")
	flag.Parse()

	if flag.NArg() < 1 || flag.NArg() > 2 {
		fmt.Printf("Usage: %s [options] source_path\n options:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
//fmt.Printf("src = %s\n", flag.Arg(0))
	
	// Walk the source directory
    err := filepath.Walk(flag.Arg(0), doWalk)
    if err != nil {
        log.Fatal(err)
    }
	fmt.Printf("\n%d MP3 files\n", mp3Files)
}