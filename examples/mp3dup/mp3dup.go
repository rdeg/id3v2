// This utility uses the id3v3 package to find duplicates in an MP3 library.
package main

import (
	"flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
	"time"
	
	"github.com/rdeg/id3v2"
)

// This structure holds the information that may be needed to identify
// duplicate files.
type tfi struct {	// all File Info
	ffname	string		// final file name, build from the ID3V2 tags
	fname	string		// source file name
	ftime	time.Time	// source file modification time
	fsize	int64		// file size
	bitrate	int			// known bitrate, from the first data frame header
}

var (
	srcDir 		string		// source directory
	mp3Files	uint		// count of MP3 files
	
	afi			[]tfi		// information for all the files
)

// Process a MP3 file.
func doMP3(fname string, osfi os.FileInfo) error {
	// Retrieve all the tags of the MP3 file header in a MP3Info.
	mi, err := id3v2.ProcessAllTags(fname)	// mi is a pointer to a MP3Info
	if err != nil {
		return err
	}

	// Create a tfi element for this file.
	full := mi.MakeFileName()	// make an MP3 file name using the tags previously harvested
	if full == ".mp3" {	// no tags in this file!
//		full = fname	// preserve source path and name
		full = osfi.Name()	// preserve source base name
	}
	var fi = tfi{ffname:full, fname:fname, ftime:osfi.ModTime(), fsize:osfi.Size(), bitrate:mi.BitRate}
	
	// Search for a database entry with the same name (i.e. with the same tags).
	for _, v := range afi {
		if v.ffname == fi.ffname {
			fmt.Printf("%s\n", fi.ffname)
			fmt.Printf("1: %s\n", fi.fname)
			fmt.Printf("   %d kpbs - %s - %d bytes\n", fi.bitrate, fi.ftime, fi.fsize)
			fmt.Printf("2: %s\n", v.fname)
			fmt.Printf("   %d kpbs - %s - %d bytes\n\n", v.bitrate, v.ftime, v.fsize)
		}
	}
	
	// Append the new element to the database.
	afi = append(afi, fi)
	
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

	// Create the initial slice for file info.
	afi = make([]tfi, 0, 1024)
	
	// Walk the source directory
    err := filepath.Walk(flag.Arg(0), doWalk)
    if err != nil {
        log.Fatal(err)
    }
	fmt.Printf("\n%d MP3 files\n", mp3Files)
//fmt.Println(afi)	// debug
}