# id3v2

`go get github.com/rdeg/id3v2`

Package id3v2 can be used to retrieve most of the ID3v2 tags of an MP3 file.

It currently offers the two functions ProcessAllTags and MakeFileName, as
well as a Verbosity indicator that can be used to tune the level of verbosity
of the package (and, hopefully, of the application that uses it too).

See https://godoc.org/github.com/rdeg/id3v2 documentation for details.

## Motivation

I wrote this package (and the mp3copy utility) after the frustration I felt after having donwloaded my music library from Google Play Music on my PC. I found myself with over 5000 pieces of music in a single folder, each song being named in a way that did not really suit me. Indeed, I tend to prefer full albums to playlists and I favor the usual classification of songs into artist / album / songs.

So if ProcessAllTags does the greatest job of harvesting ID3v2 tags from an MP3 file, MakeFileName is the tool that allows you to place each song in its 'ideal' tree. Of course, if what is ideal for me is not what you want, feel free to change the code accordingly ;-)

See the mp3copy example for more information.