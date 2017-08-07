package id3v2

import (
	"bytes"
	"errors"
    "fmt"
    "log"
    "os"
    "path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

// The Verbose variable can be used to set the level of verbosity of the
// package. Possible values are 0 (no output), 1 (output discovered tags)
// and 2 (output ID3v2 headers information).
var Verbose	uint

// A processed tag, as stored in the map of all processed tags.
type ProcessedTag struct {
	Name string		// understandable name
	Value string	// tag value as a string
}

// Information for a MP3 file
type MP3Info struct {
	AllTags map[string]*ProcessedTag	// the map of all processed tags
	BitRate int							// bitrate (from the first sample)
}

// ID3v2 tag
type mp3Tag struct {
	tag	string		// 4-char tag
	size uint		// payload size
	flags uint16	// abc00000ijk00000 (see http://id3.org/id3v2.3.0)
	extra struct {
		uncSize	uint	// uncompressed payload size if flgCompressed
		encType byte	// encryption type if flgEncrypted
		groupID byte	// group ID if flgGroupId
	}
	payload []byte
}

// Function for the processing of a tag.
// This takes a mp3Tag pointer as input and returns a couple of strings,
// for a label and its value.
type tagF func(frame *mp3Tag) (string, string)

const (
	// mp3Tag flags
	flgTap        = 0x8000	// a: Tag alter preservation
	flgFap        = 0x4000	// b: File alter preservation
	flgRO         = 0x2000	// c: Read Only
	flgCompressed = 0x0080	// i: Compression. A 4-byte uncompressed length follows the tag header.
	flgEncrypted  = 0x0040	// j: Encrypted. A 1-byte encyption type follows the tag header.
	flgGroupId    = 0x0020	// k: Grouping Identity. A 1-byte group ID follows the tag header.
)

var (
	// mp3Tag tags processing functions
	tagmap = map[string]tagF{
		"AENC":doAENC,	// Audio encryption
		"APIC":doAPIC,	// Attached picture
		"COMM":doCOMM,	// Comments
		"COMR":doCOMR,	// Commercial frame
		"ENCR":doENCR,	// Encryption method registration
		"EQUA":doEQUA,	// Equalization
		"ETCO":doETCO,	// Event timing codes
		"GEOB":doGEOB,	// General encapsulated object
		"GRID":doGRID,	// Group identification registration
		"IPLS":doIPLS,	// Involved people list
		"LINK":doLINK,	// Linked information
		"MCDI":doMCDI,	// Music CD identifier
		"MLLT":doMLLT,	// MPEG location lookup table
		"OWNE":doOWNE,	// Ownership frame
		"PRIV":doPRIV,	// Private frame
		"PCNT":doPCNT,	// Play counter
		"POPM":doPOPM,	// Popularimeter
		"POSS":doPOSS,	// Position synchronisation frame
		"RBUF":doRBUF,	// Recommended buffer size
		"RVAD":doRVAD,	// Relative volume adjustment
		"RVRB":doRVRB,	// Reverb
		"SYLT":doSYLT,	// Synchronized lyric/text
		"SYTC":doSYTC,	// Synchronized tempo codes
		"TALB":doTALB,	// Album/Movie/Show title
		"TBPM":doTBPM,	// BPM (beats per minute)
		"TCOM":doTCOM,	// Composer
		"TCON":doTCON,	// Content type
		"TCOP":doTCOP,	// Copyright message
		"TDAT":doTDAT,	// Date
		"TDLY":doTDLY,	// Playlist delay
		"TENC":doTENC,	// Encoded by
		"TEXT":doTEXT,	// Lyricist/Text writer
		"TFLT":doTFLT,	// File type
		"TIME":doTIME,	// Time
		"TIT1":doTIT1,	// Content group description
		"TIT2":doTIT2,	// Title/songname/content description
		"TIT3":doTIT3,	// Subtitle/Description refinement
		"TKEY":doTKEY,	// Initial key
		"TLAN":doTLAN,	// Language(s)
		"TLEN":doTLEN,	// Length
		"TMED":doTMED,	// Media type
		"TOAL":doTOAL,	// Original album/movie/show title
		"TOFN":doTOFN,	// Original filename
		"TOLY":doTOLY,	// Original lyricist(s)/text writer(s)
		"TOPE":doTOPE,	// Original artist(s)/performer(s)
		"TORY":doTORY,	// Original release year
		"TOWN":doTOWN,	// File owner/licensee
		"TPE1":doTPE1,	// Lead performer(s)/Soloist(s)
		"TPE2":doTPE2,	// Band/orchestra/accompaniment
		"TPE3":doTPE3,	// Conductor/performer refinement
		"TPE4":doTPE4,	// Interpreted, remixed, or otherwise modified by
		"TPOS":doTPOS,	// Part of a set
		"TPUB":doTPUB,	// Publisher
		"TRCK":doTRCK,	// Track number/Position in set
		"TRDA":doTRDA,	// Recording dates
		"TRSN":doTRSN,	// Internet radio station name
		"TRSO":doTRSO,	// Internet radio station owner
		"TSIZ":doTSIZ,	// Size
		"TSRC":doTSRC,	// ISRC (international standard recording code)
		"TSSE":doTSSE,	// Software/Hardware and settings used for encoding
		"TYER":doTYER,	// Year
		"TXXX":doTXXX,	// User defined text information frame
		"UFID":doUFID,	// Unique file identifier
		"USER":doUSER,	// Terms of use
		"USLT":doUSLT,	// Unsychronized lyric/text transcription
		"WCOM":doWCOM,	// Commercial information
		"WCOP":doWCOP,	// Copyright/Legal information
		"WOAF":doWOAF,	// Official audio file webpage
		"WOAR":doWOAR,	// Official artist/performer webpage
		"WOAS":doWOAS,	// Official audio source webpage
		"WORS":doWORS,	// Official internet radio station homepage
		"WPAY":doWPAY,	// Payment
		"WPUB":doWPUB,	// Publishers official webpage
		"WXXX":doWXXX,	// User defined URL link frame
	}
	
/*
Example of data frame first word:
FF FB D0 00 = AAAAAAAA AAABBCCD EEEEFFGH IIJJKLMM
              11111111 11111011 11010000 00000000
AAAAAAAAAAA = sync
BB = 11		MPEG V1 (V1)		00: V2.5, 01: reserved, 10: V2, 11: V1
CC = 10		Layer II (L2)		00: reserved, 01: L3, 10: L2, 11: L1
D = 1		Not protected
EEEE = 1101	V1,L2 => 320 Kbps
FF = 00		MPEG 1 => 44100 Hz
G = 0		no padding
H = 0
II = 00		Stereo
JJ = 00
K = 0		Not copyrighted
L = 0		Copy of original media
MM = 00		no emphasis

// Bitrate depending on MPEG Version (BB), Layer (CC) and bitrate index (EEEE)
EEEE	V1,L1	V1,L2	V1,L3	V2,L1	V2, L2 & L3
0000	free	free	free	free	free
0001	32		32		32		32		8
0010	64		48		40		48		16
0011	96		56		48		56		24
0100	128		64		56		64		32
0101	160		80		64		80		40
0110	192		96		80		96		48
0111	224		112		96		112		56
1000	256		128		112		128		64
1001	288		160		128		144		80
1010	320		192		160		160		96
1011	352		224		192		176		112
1100	384		256		224		192		128
1101	416		320		256		224		144
1110	448		384		320		256		160
1111	bad		bad		bad		bad		bad
NOTES: All values are in kbps
V1 - MPEG Version 1
V2 - MPEG Version 2 and Version 2.5
L1 - Layer I
L2 - Layer II
L3 - Layer III
*/
	bitRateTable = [2][3][16]int{
		{         // MPEG 2 & 2.5
			{0,  8, 16, 24, 32, 40, 48, 56, 64, 80, 96,112,128,144,160,0}, // Layer III
			{0,  8, 16, 24, 32, 40, 48, 56, 64, 80, 96,112,128,144,160,0}, // Layer II
			{0, 32, 48, 56, 64, 80, 96,112,128,144,160,176,192,224,256,0},  // Layer I
		},                                                                    
		{       // MPEG 1                                                     
			{0, 32, 40, 48, 56, 64, 80, 96,112,128,160,192,224,256,320,0}, // Layer III
			{0, 32, 48, 56, 64, 80, 96,112,128,160,192,224,256,320,384,0}, // Layer II
			{0, 32, 64, 96,128,160,192,224,256,288,320,352,384,416,448,0},  // Layer I
		},
	}
)

func decodeISO8859(src []byte) string {
	dst := make([]rune, len(src))
	for i, b := range src {
		dst[i] = rune(b)
	}
	return string(dst)
}

func decodeUtf16With(buf []byte, enc encoding.Encoding) string {
	decoder := enc.NewDecoder()
	dst := make([]byte, len(buf)*2)
	nDst, _, err := decoder.Transform(dst, buf, true)
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(dst[:nDst])
}

// Return the text of a Txxx tag.
func textFrame(payload []byte) string {
	et := payload[0]	// encoding byte
	pl := payload[1:]	// actual payload
	switch et {
	case 0x00:	// ISO 8859-1
		return decodeISO8859(pl)
	case 0x01:	// UTF-16, starting with a BOM
		return decodeUtf16With(pl, unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM))
	case 0x02:	// UTF-16BE string without BOM???
		return decodeUtf16With(pl, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM))
	case 0x03:	// UTF-8 string
		return string(pl)
	default:
		return fmt.Sprintf("UNKNOWN ENCODING (0x%02x)", et)
	}
}

func doAENC(frame *mp3Tag) (string, string) {	// Audio encryption
	return "", ""
}
func doAPIC(frame *mp3Tag) (string, string) {	// Attached picture
// Text encoding   $xx
// MIME type       <text string> $00
// Picture type    $xx
// Description     <text string according to encoding> $00 (00)
// Picture data    <binary data>
	et := frame.payload[0]		// encoding byte
	i := bytes.IndexByte(frame.payload[1:], 0x00)
//fmt.Println("i = ", i)
	if i == -1 {
		return "Error", "Cannot find MIME type termination in APIC frame"
	}
	mimeType := textFrame(frame.payload[:1 + i])
	pictureType := frame.payload[1 + i + 1]
//fmt.Printf("mimeType = %s, pictureType = %d\n", mimeType, pictureType)
	j := bytes.IndexByte(frame.payload[1 + i + 1 + 1:], 0x00)
//fmt.Println("j = ", j)
	if j == -1 {
		return "Error", "Cannot find Description termination in APIC frame"
	}
	var db = make([]byte, 1 + j)
	db[0] = et
	copy(db[1:], frame.payload[1 + i + 1 + 1:1 + i + 1 + 1 + j])
	description := textFrame(db)
	bytes := len(frame.payload) - (1 + i + 1 + 1 + j + 1)
	return "Picture", fmt.Sprintf("%s, 0x%02x, \"%s\", %d (0x%x) bytes", mimeType, pictureType, description, bytes, bytes)
}
func doCOMM(frame *mp3Tag) (string, string) {	// Comments
	return "", ""
}
func doCOMR(frame *mp3Tag) (string, string) {	// Commercial frame
	return "", ""
}
func doENCR(frame *mp3Tag) (string, string) {	// Encryption method registration
	return "", ""
}
func doEQUA(frame *mp3Tag) (string, string) {	// Equalization
	return "", ""
}
func doETCO(frame *mp3Tag) (string, string) {	// Event timing codes
	return "", ""
}
func doGEOB(frame *mp3Tag) (string, string) {	// General encapsulated object
	return "", ""
}
func doGRID(frame *mp3Tag) (string, string) {	// Group identification registration
	return "", ""
}
func doIPLS(frame *mp3Tag) (string, string) {	// Involved people list
	return "", ""
}
func doLINK(frame *mp3Tag) (string, string) {	// Linked information
	return "", ""
}
func doMCDI(frame *mp3Tag) (string, string) {	// Music CD identifier
	return "", ""
}
func doMLLT(frame *mp3Tag) (string, string) {	// MPEG location lookup table
	return "", ""
}
func doOWNE(frame *mp3Tag) (string, string) {	// Ownership frame
	return "", ""
}
func doPRIV(frame *mp3Tag) (string, string) {	// Private frame
	i := bytes.IndexByte(frame.payload, 0x00)	// <text string> $00 <private data> expected
	if i == -1 {
		return "Error", "Missing 0x00 in PRIV frame"
	}	
	return string(frame.payload[0:i]), fmt.Sprintf("%d bytes", len(frame.payload) - i - 1)
}
func doPCNT(frame *mp3Tag) (string, string) {	// Play counter
	return "", ""
}
func doPOPM(frame *mp3Tag) (string, string) {	// Popularimeter
	return "", ""
}
func doPOSS(frame *mp3Tag) (string, string) {	// Position synchronisation frame
	return "", ""
}
func doRBUF(frame *mp3Tag) (string, string) {	// Recommended buffer size
	return "", ""
}
func doRVAD(frame *mp3Tag) (string, string) {	// Relative volume adjustment
	return "", ""
}
func doRVRB(frame *mp3Tag) (string, string) {	// Reverb
	return "", ""
}
func doSYLT(frame *mp3Tag) (string, string) {	// Synchronized lyric/text
	return "", ""
}
func doSYTC(frame *mp3Tag) (string, string) {	// Synchronized tempo codes
	return "", ""
}
func doTALB(frame *mp3Tag) (string, string) {	// Album/Movie/Show title
	return "Album", textFrame(frame.payload)
}
func doTBPM(frame *mp3Tag) (string, string) {	// BPM (beats per minute)
	return "BPM", textFrame(frame.payload)
}
func doTCOM(frame *mp3Tag) (string, string) {	// Composer
	return "Composer", textFrame(frame.payload)
}
func doTCON(frame *mp3Tag) (string, string) {	// Content type
	return "Content type", textFrame(frame.payload)
}
func doTCOP(frame *mp3Tag) (string, string) {	// Copyright message
	return "Copyright", textFrame(frame.payload)
}
func doTDAT(frame *mp3Tag) (string, string) {	// Date
	return "Playlist delay", textFrame(frame.payload)
}
func doTDLY(frame *mp3Tag) (string, string) {	// Playlist delay
	return "", textFrame(frame.payload)
}
func doTENC(frame *mp3Tag) (string, string) {	// Encoded by
	return "Encoded by", textFrame(frame.payload)
}
func doTEXT(frame *mp3Tag) (string, string) {	// Lyricist/Text writer
	return "Lyrics by", textFrame(frame.payload)
}
func doTFLT(frame *mp3Tag) (string, string) {	// File type
	return "", textFrame(frame.payload)
}
func doTIME(frame *mp3Tag) (string, string) {	// Time
	return "Time", textFrame(frame.payload)
}
func doTIT1(frame *mp3Tag) (string, string) {	// Content group description
	return "Content group", textFrame(frame.payload)
}
func doTIT2(frame *mp3Tag) (string, string) {	// Title/songname/content description
	return "Title", textFrame(frame.payload)
}
func doTIT3(frame *mp3Tag) (string, string) {	// Subtitle/Description refinement
	return "Also", textFrame(frame.payload)
}
func doTKEY(frame *mp3Tag) (string, string) {	// Initial key
	return "Initial key", textFrame(frame.payload)
}
func doTLAN(frame *mp3Tag) (string, string) {	// Language(s)
	return "Language(s)", textFrame(frame.payload)
}
func doTLEN(frame *mp3Tag) (string, string) {	// Length
	return "Length", textFrame(frame.payload)
}
func doTMED(frame *mp3Tag) (string, string) {	// Media type
	return "Media type", textFrame(frame.payload)
}
func doTOAL(frame *mp3Tag) (string, string) {	// Original album/movie/show title
	return "Original album", textFrame(frame.payload)
}
func doTOFN(frame *mp3Tag) (string, string) {	// Original filename
	return "Original filename", textFrame(frame.payload)
}
func doTOLY(frame *mp3Tag) (string, string) {	// Original lyricist(s)/text writer(s)
	return "Original lyricist(s)", textFrame(frame.payload)
}
func doTOPE(frame *mp3Tag) (string, string) {	// Original artist(s)/performer(s)
	return "Original artist(s)", textFrame(frame.payload)
}
func doTORY(frame *mp3Tag) (string, string) {	// Original release year
	return "Original release year", textFrame(frame.payload)
}
func doTOWN(frame *mp3Tag) (string, string) {	// File owner/licensee
	return "Owner", textFrame(frame.payload)
}
func doTPE1(frame *mp3Tag) (string, string) {	// Lead performer(s)/Soloist(s)
	return "Artist(s)", textFrame(frame.payload)
}
func doTPE2(frame *mp3Tag) (string, string) {	// Band/orchestra/accompaniment
	return "Band", textFrame(frame.payload)
}
func doTPE3(frame *mp3Tag) (string, string) {	// Conductor/performer refinement
	return "Also", textFrame(frame.payload)
}
func doTPE4(frame *mp3Tag) (string, string) {	// Interpreted, remixed, or otherwise modified by
	return "Modified by", textFrame(frame.payload)
}
func doTPOS(frame *mp3Tag) (string, string) {	// Part of a set
	return "Part", textFrame(frame.payload)
}
func doTPUB(frame *mp3Tag) (string, string) {	// Publisher
	return "Publisher", textFrame(frame.payload)
}
func doTRCK(frame *mp3Tag) (string, string) {	// Track number/Position in set
	return "Track", textFrame(frame.payload)
}
func doTRDA(frame *mp3Tag) (string, string) {	// Recording dates
	return "Recorded on", textFrame(frame.payload)
}
func doTRSN(frame *mp3Tag) (string, string) {	// Internet radio station name
	return "Radio", textFrame(frame.payload)
}
func doTRSO(frame *mp3Tag) (string, string) {	// Internet radio station owner
	return "Radio owner", textFrame(frame.payload)
}
func doTSIZ(frame *mp3Tag) (string, string) {	// Size
	return "Size", textFrame(frame.payload)
}
func doTSRC(frame *mp3Tag) (string, string) {	// ISRC (international standard recording code)
	return "ISRC", textFrame(frame.payload)
}
func doTSSE(frame *mp3Tag) (string, string) {	// Software/Hardware and settings used for encoding
	return "Encoding settings", textFrame(frame.payload)
}
func doTYER(frame *mp3Tag) (string, string) {	// Year
	return "Year", textFrame(frame.payload)
}
func doTXXX(frame *mp3Tag) (string, string) {	// User defined text information frame
	return "User defined", textFrame(frame.payload)
}
func doUFID(frame *mp3Tag) (string, string) {	// Unique file identifier
	return "", ""
}
func doUSER(frame *mp3Tag) (string, string) {	// Terms of use
	return "", ""
}
func doUSLT(frame *mp3Tag) (string, string) {	// Unsychronized lyric/text transcription
	return "", ""
}
func doWCOM(frame *mp3Tag) (string, string) {	// Commercial information
	return "", ""
}
func doWCOP(frame *mp3Tag) (string, string) {	// Copyright/Legal information
	return "", ""
}
func doWOAF(frame *mp3Tag) (string, string) {	// Official audio file webpage
	return "", ""
}
func doWOAR(frame *mp3Tag) (string, string) {	// Official artist/performer webpage
	return "", ""
}
func doWOAS(frame *mp3Tag) (string, string) {	// Official audio source webpage
	return "", ""
}
func doWORS(frame *mp3Tag) (string, string) {	// Official internet radio station homepage
	return "", ""
}
func doWPAY(frame *mp3Tag) (string, string) {	// Payment
	return "", ""
}
func doWPUB(frame *mp3Tag) (string, string) {	// Publishers official webpage
	return "", ""
}
func doWXXX(frame *mp3Tag) (string, string) {	// User defined URL link frame
	return "", ""
}
// Read an exact count of bytes
func readExact(f *os.File, b []byte) error {
	n, err := f.Read(b)
//fmt.Printf("Read 0x%x bytes\n", n)
	if err == nil {
		if n != len(b) {
			err = errors.New("Error reading MP3 file")
		}
	}
	return err
}

// ProcessAllTags processes all the tags of an MP3 file and saves related
// information in an MP3Info structure returned to the caller.
func ProcessAllTags(fname string) (*MP3Info, error) {
	var	mi MP3Info
	mi.AllTags = make(map[string]*ProcessedTag)

	// Retrieve the header, if any.
	sf, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer sf.Close()

	// Make sure we have a valid ID3v2 header.
	b := make([]byte, 10)	// ID3v2 header size
	err = readExact(sf, b)
	if err != nil {
		return nil, err
	}
	if b[0] != 'I' || b[1] != 'D' || b[2] != '3' {
		fmt.Printf("%#v\n", b)
		panic(fmt.Sprintf("NO ID3v2 HEADER IN %q", fname))
	}
	// 49 44 33 yy yy xx zz zz zz zz	// yy yy = version, xx = flags, zz zz zz zz = size
	if b[3] == 0xff || b[4] == 0xff || b[5] & 0x1f != 0 || b[6] & 0x80 != 0 || b[7] & 0x80 != 0 || b[8] & 0x80 != 0 || b[9] & 0x80 != 0 {
		return nil, errors.New("Invalid ID3v2 header")
	}
	hdrsz := (uint(b[6]) << 21) + (uint(b[7]) << 14) + (uint(b[8]) << 7) + (uint(b[9]) << 0)
	if Verbose >= 2 {
		fmt.Printf("ID3v2.%d.%d header", b[3], b[4])
		aux := ""
		if b[5] & 0x80 != 0 {
			aux += "unsync"
		}
		if b[5] & 0x40 != 0 {
			aux += "exthdr"
		}
		if b[5] & 0x20 != 0 {
			aux += "eXprmt"
		}
		if aux != "" {
			fmt.Printf(" (%s)", aux)
		}
		fmt.Printf(", 0x%x bytes\n", hdrsz)
	}
	
	// Read all the ID3v2 frames in a single buffer.
	hb := make([]byte, hdrsz)
	err = readExact(sf, hb)
	if err != nil {
		return nil, err
	}
		
	// Read the information word of the first data frame.
	err = readExact(sf, b[:4])
	if err != nil {
		return nil, err
	}
	var info uint32 = uint32(b[0]) << 24 | uint32(b[1]) << 16  | uint32(b[2]) << 8  | uint32(b[3])
	if info & 0xffe00000 != 0xffe00000 {	// invalid synch pattern?!?
//		err = errors.New(fmt.Sprintf("Invalid sync pattern (info = 0x%04X)", info))
//		return nil, err
		fmt.Printf("Invalid sync pattern (info = 0x%04X)\n", info)
//		mi.BitRate = 0
	} else {	// looks like a valid sync pattern
		// Compute the bitrate from the indexes in the info word.
		vi  := (info & 0x00180000) >> 19	// MPEG version (00: V2.5, 01: reserved, 10: V2, 11: V1)
		li  := (info & 0x00060000) >> 17	// Layer (00: reserved, 01: L3, 10: L2, 11: L1)
		bri := (info & 0x0000f000) >> 12	// bitrate index
//fmt.Printf("info = 0x%04X (vi = %d, li = %d, bri = %d)\n", info, vi, li, bri)
		mi.BitRate = bitRateTable[vi & 1][li - 1][bri]
	}
	if Verbose >= 1 {
		fmt.Printf("Bitrate = %d kbps\n", mi.BitRate)
	}

	var ihb uint = 0		// start here

	// Check for an extended header.
	if b[5] & 0x40 != 0	{ // an extended header follows
		var xsiz, xflg, xcrc, xpad uint
		xsiz = (uint(hb[0]) << 24) + (uint(hb[1]) << 16) + (uint(hb[2]) << 8) + (uint(hb[3]) << 0)
		xflg = uint((uint16(hb[4]) << 8) + (uint16(hb[5]) << 0))
		if Verbose >= 2 {
			xpad = (uint(hb[6]) << 24) + (uint(hb[7]) << 16) + (uint(hb[8]) << 8) + (uint(hb[9]) << 0)
		}
		ihb += 10	// skip the extended header size
		if xflg & 0x7fff != 0 {
			return nil, errors.New(fmt.Sprintf("Invalid flags for extended tag (0x%x)", xflg))
		}
		if xflg & 0x8000 != 0 {	// CRC data present
			if Verbose >= 2 {
				xcrc = (uint(hb[10]) << 24) + (uint(hb[11]) << 16) + (uint(hb[12]) << 8) + (uint(hb[13]) << 0)
			}
			ihb += 4	// skip the CRC
		}
		ihb += xsiz	// seek the first tag
		if Verbose >= 2 {
			fmt.Printf(" * Extended header: length = 0x%08x, flags = 0x%x, padding = 0x%08x", xsiz, xflg, xpad)
			if xflg & 0x8000 != 0 {	// CRC data present
				fmt.Printf(", CRC = 0x%x", xcrc)
			}
			fmt.Println()
		}
	}
		
	// Process the tags until the big header is consumed.
	for ;ihb < hdrsz; {
//		ihb0 := ihb	// save beginning of the tag
		
		// Slice the tag header. Reuse b.
		b = hb[ihb:ihb + 10]
		ihb += 10
		
		var t mp3Tag
		t.tag = string(b[:4])
		t.size = (uint(b[4]) << 24) + (uint(b[5]) << 16) + (uint(b[6]) << 8) + (uint(b[7]) << 0)
		t.flags = (uint16(b[8]) << 8) + (uint16(b[9]) << 0)
		if t.flags & 0x1f1f != 0 {
			return nil, errors.New(fmt.Sprintf("Invalid flags for tag %s (0x%x)", t.tag, t.flags))
		}
		if t.tag == "" || t.size == 0 {
			break
		}
		
		// Read the extra bytes following the tag header, if any.
		if (t.flags & flgCompressed) != 0 {
			b = hb[ihb:ihb + 4]
			ihb += 4
			t.extra.uncSize = (uint(b[0]) << 24) + (uint(b[1]) << 16) + (uint(b[2]) << 8) + (uint(b[3]) << 0)
		}
		if (t.flags & flgEncrypted) != 0 {
			t.extra.encType = hb[ihb]
			ihb += 1
		}
		if (t.flags & flgGroupId) != 0 {
			t.extra.groupID = hb[ihb]
			ihb += 1
		}
		
		// Slice the payload.
		t.payload = hb[ihb:ihb + t.size]
		ihb += t.size

		if Verbose >= 2 {
			fmt.Printf("  %s: length = 0x%07x, flags = 0x%x\n", t.tag, t.size, t.flags)
		}

		// Extract tag's value.
		tfn, ok := tagmap[t.tag]
		if !ok {
			fmt.Printf("%s: UNEXPECTED TAG!\n", t.tag)
			continue
		}
		lbl, val := tfn(&t)
		if Verbose >= 1 {
			fmt.Printf("%s: %s\n", lbl, val)
		}
		mi.AllTags[t.tag] = &ProcessedTag{Name:lbl, Value:val}
	}
	return &mi, nil
}

// Parse an int value that can be terminated by a non-digit character.
var leadingInt = regexp.MustCompile(`^[-+]?\d+`)
func parseLeadingInt(s string) (int64, error) {
	if s == "1/1" {
		return 0, nil
	}
	s = leadingInt.FindString(s)
	if s == "" { // add this if you don't want error on "xx" etc
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

// Make sure a string that is a part of a full filename is valid.
// We make this for Windows, assuming Linux will be OK is Windows is.
func purify(s string, nofdots bool) string {
	// Remove leading and trailing spaces from given string.
	if nofdots {
		s = strings.Trim(s, " .")
	} else {
		s = strings.TrimSpace(s)
	}
	
	// Remove trailing dots if needed.

	// Replace any Windows invalid chars by '-'.
	b := []byte(s)
	for {
		i := strings.IndexAny(s, "<>:\"/\\|?*")
		if i == -1 {
			return s
		}
		b[i] = '-'
		s = string(b)
	}
}

// tagVal returns the string value of a given tag.
func tagVal(pm *map[string]*ProcessedTag, tag string) (val string) {
	pt, ok := (*pm)[tag]
	if ok {
		val = pt.Value
	}
	return
}

// MakeFileName makes an MP3 filename using the tags previously retrieved
// and saved in an MP3Info structure by ProcessAllTags.
// The resulting string should respect the following pattern:
//	artist/album/{disk-}track title.mp3
func (mi *MP3Info) MakeFileName() string {
	// Pick the values of the tags we need.
	artists := tagVal(&mi.AllTags, "TPE1")
	band    := tagVal(&mi.AllTags, "TPE2")
	album   := tagVal(&mi.AllTags, "TALB")
	disk    := tagVal(&mi.AllTags, "TPOS")
	track   := tagVal(&mi.AllTags, "TRCK")
	title   := tagVal(&mi.AllTags, "TIT2")
	if Verbose >= 1 {
		fmt.Printf("artists:%s, band:%s, album:%s, disk:%s, track:%s, title:%s\n", artists, band, album, disk, track, title)
	}

	// Create an MP3 file name according to its tags.
	var sd, st string
//	var v int64
	artist := purify(band, true)	// prefer the band/orchestra information
	if artist == "" {	// no band info
		artist = purify(artists, true)	// should always exist
	}
//	dir := filepath.Join(artist, purify(album, true))
	v, err := parseLeadingInt(disk)
	if err == nil && v != 0 {	// valid disk number
		sd = fmt.Sprintf("%d-", v)
	}
	v, err = parseLeadingInt(track)
	if err == nil && v != 0 {	// valid track number
		st = fmt.Sprintf("%02d ", v)
	}
//fmt.Printf("track = '%s', v = %d, st = '%s'\n", track, v, st)
	name := sd + st + purify(title, false) + ".mp3"
	return filepath.Join(artist, purify(album, true), name)
}
