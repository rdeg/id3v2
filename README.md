# id3v2

`go get github.com/rdeg/id3v2`

Package id3v2 can be used to retrieve most of the ID3v2 tags of an MP3 file.

It currently offers the two functions ProcessAllTags and MakeFileName, as
well as a Verbosity indicator that can be used to tune the level of verbosity
of the package (and, hopefully, of the application that uses it too).

See https://godoc.org/github.com/rdeg/id3v2 documentation for details.

## Motivation

I wrote this package (and the mp3copy utility) after the frustration that followed downloading my music library from Google Play Music on my PC. I found myself with over 5000 pieces of music in one folder, each song being named in a way that did not really fit me. Indeed, I tend to prefer full albums to playlists and I favor the usual classification of songs into artist / album / songs.

So if **ProcessAllTags** does the greatest job of harvesting ID3v2 tags from an MP3 file, **MakeFileName** is the tool that allows you to place each song in its 'ideal' tree. Of course, if what is ideal for me is not what you want, feel free to change the code accordingly ;-)

See the **mp3copy** example for more information.

## Tags processing

The tags laying in the ID3v2 header of an MP3 file are processed in the **ProcessAllTags** function.

Many tags are taken into account but only a few one are actually processed. Namely, they are currently the APIC, PRIV and Txxx tags, but this may change in the future.

Tags of interest for later creation of a path and filename (ie, TALB, TIT2, TPE1, TPE2, TPOS, and TRCK) belong to this list. Please note that APIC tag processing only retrieves meta-information from the image that can be embedded in an MP3 file and not the actual bits of that image.

## Bitrate

Bitrate determination is minimalistic and only relies on the information word of the first data frame of the file. It is thus not guaranteed to reflect the actual bitrate of a file, specially in case of variable bitrate.
