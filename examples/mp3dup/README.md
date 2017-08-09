# mp3dup

This utility uses the id3v3 package to find duplicates in an MP3 library.

# Usage

```
Usage: mp3dup [options] source_path
 options:
  -v uint
        verbosity level (0 = none, 1 = tags, 2 = headers)
```

`source_path` is the path where can be found the MP3 files to process. Only files with a '.mp3' extension will be taken into account. Every MP3 file in this tree (or flat) directory will be processed.

Even with a zero-verbosity (the default for option `-v`), some ouput showing the progression of the copying process is made. In case of trouble, it can be useful to raise the level of verbosity.

The output should be redirected to a file for review of suggested duplicates.

# Behavior

The detection of duplicates is based on the names generated automatically from the tags of the ID3v2 header of the files: when two files produce the same name, it can be considered as a duplicate.

This is where the intelligence stops and we will have to examine the list of duplicates to choose which file to eliminate in each pair. This choice can be oriented by bitrate or file size, or even by its dating, but in case of doubt, only a hearing of each file will allow a truly informed choice.

# Example output

```
Maria Montell\Bossa For My Baby\13 And So The Story Goes (DiDaDi) (Album Version).mp3
1: D:\Vrac\B0026UPEM6_(disc_1)_13_-_And_So_The_Story_Goes_(DiDaDi)_(Album_Version).mp3
   48 kpbs - 2015-01-10 18:26:20 +0100 CET - 17685702 bytes
2: D:\Vrac\01-13- And So The Story Goes (DiDaDi) (Album Version).mp3
   48 kpbs - 2016-05-16 22:03:30.9389242 +0200 CEST - 17685702 bytes

Gorillaz\Clint Eastwood [Explicit]\01 Clint Eastwood.mp3
1: D:\Vrac\B002Q2ED9O_(disc_1)_01_-_Clint_Eastwood.mp3
   128 kpbs - 2015-01-10 18:26:02 +0100 CET - 11160916 bytes
2: D:\Vrac\01-01- Clint Eastwood.mp3
   128 kpbs - 2016-05-16 21:53:56.3975566 +0200 CEST - 11160916 bytes

(!!! Chk Chik Chick)\Strange Weather, Isn't It-\07 Jump Back.mp3
1: D:\Vrac\B003XRMTOW_(disc_1)_07_-_Jump_Back.mp3
   128 kpbs - 2015-03-11 00:08:45 +0100 CET - 6229259 bytes
2: D:\Vrac\01-07- Jump Back.mp3
   128 kpbs - 2016-05-16 21:48:57.1265767 +0200 CEST - 6229259 bytes

Poni Hoax\A State Of War\06 Life In A New Motion.mp3
1: D:\Vrac\B00BVYYNGQ_(disc_1)_06_-_Life_In_A_New_Motion.mp3
   48 kpbs - 2015-01-10 18:26:04 +0100 CET - 7492711 bytes
2: D:\Vrac\01-06- Life In A New Motion.mp3
   48 kpbs - 2016-05-16 22:06:10.2229845 +0200 CEST - 7492711 bytes

(!!! Chk Chik Chick)\Thr!!!er\08 Careful.mp3
1: D:\Vrac\B00CE6C7NG_(disc_1)_08_-_Careful.mp3
   48 kpbs - 2015-02-07 11:50:46 +0100 CET - 12614599 bytes
2: D:\Vrac\01-08- Careful.mp3
   48 kpbs - 2016-05-16 21:49:32.0007973 +0200 CEST - 12614599 bytes

Balthazar\Leipzig\Leipzig.mp3
1: D:\Vrac\B00HQ18ENQ_(disc_1)_01_-_Leipzig.mp3
   48 kpbs - 2015-04-02 00:32:17 +0200 CEST - 5770409 bytes
2: D:\Vrac\01-01- Leipzig.mp3
   48 kpbs - 2016-05-16 21:45:37.9371106 +0200 CEST - 5770409 bytes

The Avener\The Wanderings Of The Avener\1-01 Panama.mp3
1: D:\Vrac\B00QV1JLLG_(disc_1)_01_-_Panama.mp3
   56 kpbs - 2015-03-10 23:48:50 +0100 CET - 8887899 bytes
2: D:\Vrac\01-01- Panama.mp3
   56 kpbs - 2016-05-16 21:44:37.7625238 +0200 CEST - 8887899 bytes

The Avener\The Wanderings Of The Avener\1-04 Castle In The Snow.mp3
1: D:\Vrac\B00QV1JVUW_(disc_1)_04_-_Castle_In_The_Snow.mp3
   56 kpbs - 2015-03-10 23:47:54 +0100 CET - 6855071 bytes
2: D:\Vrac\01-04- Castle In The Snow.mp3
   56 kpbs - 2016-05-16 21:44:35.8872501 +0200 CEST - 6855071 bytes

The Avener\The Wanderings Of The Avener\1-08 Hate Street Dialogue [feat. Rodriguez].mp3
1: D:\Vrac\B00QV1K86I_(disc_1)_08_-_Hate_Street_Dialogue_[feat__Rodriguez].mp3
   56 kpbs - 2015-03-10 23:48:23 +0100 CET - 9197447 bytes
2: D:\Vrac\01-08- Hate Street Dialogue [feat Rodriguez].mp3
   56 kpbs - 2016-05-16 21:45:16.3945266 +0200 CEST - 9197447 bytes

Joni Mitchell\Both Sides Now\08 Sometimes I'm Happy.mp3
1: D:\Vrac\Joni Mitchell - Sometimes I'm Happy_1.mp3
   256 kpbs - 2017-04-18 15:42:32.1010864 +0200 CEST - 7736232 bytes
2: D:\Vrac\Joni Mitchell - Sometimes I'm Happy.mp3
   320 kpbs - 2017-04-18 13:46:48.83423 +0200 CEST - 9535112 bytes

   
346 MP3 files
```
