# mp3copy

This utility uses the id3v3 package to sort and organize an MP3 library that is probably badly organized
in a nice tree respecting the following organization:

```
<artist>/<album>/[<disk>-]<track> <title>.mp3
```

# Usage

```
Usage: mp3copy [options] source_path [destination_directory]
 options:
  -m    move files (delete after copying)
  -v uint
        verbosity level (0 = none, 1 = tags, 2 = headers)
```

`source_path` is the path where can be found the MP3 files to process. Only files with a '.mp3' extension will be taken into account.
Every MP3 file in this tree (or flat) directory will be processed.

If a `destination_directory` is given, source files will be copied (or moved if the -m option is given) to this directory according
to the hierarchy tree given by their tags. New <artist> and <album> subdirectories will be created as needed.
