package main

import (
	"flag"
    "fmt"
	"io"
    "log"
    "os"
    "path/filepath"
	"runtime"
	
	"github.com/rdeg/id3v2"
)

var (
	srcDir 		string		// source directory
	dstDir		string		// destination directory (may be empty)
	moveFiles	bool		// delete the source after copying
	
	mp3Files	uint		// count of MP3 files
	wmaFiles	uint		// count of WMA files (for inforamtion)
)

// -------------------------------------------------------------------------
// Process a MP3 file.
//
func doMP3(fname string, osfi os.FileInfo) error {
	// Retrieve all the tags in the MP3 file header.
	mi, err := id3v2.ProcessAllTags(fname)
	if err != nil {
		return err
	}

	// Make an MP3 file name using the tags previously harvested.
	// Result should be something like:
	//	artist/album/{disk-}track title.mp3
	dname := mi.MakeFileName()
	if dname == ".mp3" {	// no tags in this file!
		dname = osfi.Name()	// preserve source base name (will end in the root of dstDir)
	}
	br := mi.BitRate
	mtime := osfi.ModTime()
	fsize := osfi.Size()
	
	// If dsdDir is empty, just display a few file information and return.
	fmt.Printf("<<< %s\n", fname)
	fmt.Printf(">>> %s\n", dname)
	fmt.Printf(">>> %d kpbs - %s - %d bytes\n\n", br, mtime, fsize)
	if dstDir == "" {
		return nil
	}
	
	// Create the destination directory, if needed.
	dir := filepath.Join(dstDir, filepath.Dir(dname))
	err = os.MkdirAll(dir, 777)
	if err != nil {
		panic(fmt.Sprintf("%s: directory creation failed: %s", dir, err.Error()))
	}

	// Create the destination file.
	dname = filepath.Join(dstDir, dname)
	df, err := os.Create(dname)
	if err != nil || df == nil {
		panic(fmt.Sprintf("%s: creation failed: %s", dname, err.Error()))
	}
	
	// Open the source file.
	sf, err := os.Open(fname)
	if err != nil {
		panic(fmt.Sprintf("%s: open failed: %s", fname, err.Error()))
	}

	// Copy the data to its new file.
	_, err = io.Copy(df, sf)
	if err != nil {
		panic(fmt.Sprintf("%s -> %s: copy failed: %s", fname, dname, err.Error()))
	}
	
	// Close the file handles.
	sf.Close()
	df.Close()

	// Restore initial file mode for non-Windows OSes.
	if runtime.GOOS != "windows" {
		err = df.Chmod(osfi.Mode())
		if err != nil {
			panic(fmt.Sprintf("%s: failed to set file mode: %s", dname, err.Error()))
		}
	}
	
	// Set the access and modification times of the new file to the
	// modification time of the original file. Please note that the
	// creation time if not modified.
	err = os.Chtimes(dname, mtime, mtime)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to set file time: %s", dname, err.Error()))
	}

	// Remove the source file if this is what the user wants.
	if moveFiles {
		err = os.Remove(fname)
		if err != nil {
			panic(fmt.Sprintf("%s: failed to remove file: %s", fname, err.Error()))
		}
	}
	
	return nil
}

// Walk the source file tree.
func doWalk(fname string, info os.FileInfo, err error) error {
    if err != nil {
        log.Print(err)
        return nil
    }
	
	if id3v2.Verbose >= 1 {
		fmt.Println(fname)		// display the whole source filename
	}
	ext := filepath.Ext(fname)		// file extention
	switch ext {
	case ".mp3":
		mp3Files++
//		fmt.Println(base)		// display the base name
		return doMP3(fname, info)
	case ".wma":
		wmaFiles++
//		fmt.Println(base)		// display the base name
	}
    return nil
}

// Program entry.
func main() {
	// Setup options and args.
	flag.BoolVar(&moveFiles, "m", false, "move files (delete after copying)")
	flag.UintVar(&id3v2.Verbose, "v", 0, "verbosity level (0 = none, 1 = tags, 2 = headers)")
	flag.Parse()

	if flag.NArg() < 1 || flag.NArg() > 2 {
		fmt.Printf("Usage: %s [options] source_path [destination_directory]\n options:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	
	srcDir = flag.Arg(0)
	if flag.NArg() == 2 {
		dstDir = flag.Arg(1)
	}
fmt.Printf("move = %t, src = %s, dst = %s\n", moveFiles, srcDir, dstDir)
	
	// Walk the source directory
    err := filepath.Walk(srcDir, doWalk)
    if err != nil {
        log.Fatal(err)
    }
	fmt.Printf("\n%d MP3 files (%d WMA files)\n", mp3Files, wmaFiles)
}