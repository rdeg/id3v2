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

Even with a zero-verbosity (the default for option `-v`), some ouput showing the progression of the copying process is made.
This output goes to stdout and can easily be redirected to a file. When a big-enough music library has to be processed, it is
recommended to first run `mp3copy` without giving a destination directory and to redirect the output to a file, to make sure
everything will go as expected. In case of trouble, it can be useful to raise the level of verbosity.

# Example execution

Under Windows, given a D:\src directory with the following files:

```
01-01- Here with Me.mp3
01-04- Awake.mp3
01-09- American Daydream.mp3
01-10- Too Much Hope.mp3
AUIM.mp3
B004SBHBPO_(disc_1)_01_-_Everglades_(feat__Arnaud_Rebotini).mp3
B004SBHC6M_(disc_1)_10_-_Cold_Nights.mp3
B004SBJE40_(disc_1)_05_-_L'animale_(feat__La_Food)-1.mp3
B004SBJE5E_(disc_1)_07_-_Eraser.mp3
B004SBJE5Y_(disc_1)_08_-_Endless_Disco.mp3
B004SBLM18_(disc_1)_03_-_Life_In_Mono.mp3
B004SBLM9K_(disc_1)_09_-_Bad_Obsession.mp3
B0088E2DNM_(disc_1)_01_-_One_Day___Reckoning_Song_(Wankelmut.mp3
BEAK.mp3
ENLK.mp3
GIAH.mp3
IEHZ.mp3
MXGE.mp3
PPTD.mp3
QWJE.mp3
SVZO.mp3
XJCY.mp3
```

Running the `mp3copy d:\src d:\dst` will build the following tree:

```
D:\DST
+---Air
|   \---Moon Safari
|           02 Sexy Boy.mp3
|
+---Alain Souchon
|   \---Nickel
|           07 Les cadors (live).mp3
|           20 On Avance (Live).mp3
|
+---Asaf Avidan & the Mojos
|   \---One Day - Reckoning Song (Wankelmut Remix)
|           One Day - Reckoning Song (Wankelmut Remix) (Radio Edit).mp3
|
+---Dido
|   \---No Angel
|           01 Here with Me.mp3
|
+---Electric Guest
|   \---Mondo
|           02 This Head I Hold.mp3
|           03 Under the Gun.mp3
|           04 Awake.mp3
|           08 Troubleman.mp3
|           09 American Daydream.mp3
|
+---H-burns
|   \---Night Moves
|           10 Too Much Hope.mp3
|
\---Rafale
    \---Obsessions
            01 Everglades (feat. Arnaud Rebotini).mp3
            02 Beyond Bad.mp3
            03 Life In Mono.mp3
            04 Never Ever.mp3
            05 L'animale (feat. La Food).mp3
            05 L'animale.mp3
            06 Marine Aircrash.mp3
            07 Eraser.mp3
            08 Endless Disco.mp3
            09 Bad Obsession.mp3
            10 Cold Nights.mp3
```

And the following output (verbosity = 0):

```
move = false, src = d:\src, dst = d:\dst
<<< d:\src\01-01- Here with Me.mp3
>>> Dido\No Angel\01 Here with Me.mp3
>>> 56 kpbs - 2016-05-16 21:52:01.2165384 +0200 CEST - 8454426 bytes

<<< d:\src\01-04- Awake.mp3
>>> Electric Guest\Mondo\04 Awake.mp3
>>> 128 kpbs - 2016-05-16 21:52:28.58797 +0200 CEST - 9903130 bytes

<<< d:\src\01-09- American Daydream.mp3
>>> Electric Guest\Mondo\09 American Daydream.mp3
>>> 128 kpbs - 2016-05-16 21:52:18.549844 +0200 CEST - 4733224 bytes

<<< d:\src\01-10- Too Much Hope.mp3
>>> H-burns\Night Moves\10 Too Much Hope.mp3
>>> 56 kpbs - 2016-05-16 21:53:57.9286435 +0200 CEST - 6468970 bytes

<<< d:\src\AUIM.mp3
>>> Alain Souchon\Nickel\07 Les cadors (live).mp3
>>> 256 kpbs - 2017-06-18 12:55:07.9159946 +0200 CEST - 8293221 bytes

<<< d:\src\B004SBHBPO_(disc_1)_01_-_Everglades_(feat__Arnaud_Rebotini).mp3
>>> Rafale\Obsessions\01 Everglades (feat. Arnaud Rebotini).mp3
>>> 128 kpbs - 2015-06-13 15:24:05 +0200 CEST - 12957671 bytes

<<< d:\src\B004SBHC6M_(disc_1)_10_-_Cold_Nights.mp3
>>> Rafale\Obsessions\10 Cold Nights.mp3
>>> 128 kpbs - 2015-06-13 15:24:09 +0200 CEST - 12659426 bytes

<<< d:\src\B004SBJE40_(disc_1)_05_-_L'animale_(feat__La_Food)-1.mp3
>>> Rafale\Obsessions\05 L'animale (feat. La Food).mp3
>>> 128 kpbs - 2015-06-13 15:24:20 +0200 CEST - 8393417 bytes

<<< d:\src\B004SBJE5E_(disc_1)_07_-_Eraser.mp3
>>> Rafale\Obsessions\07 Eraser.mp3
>>> 128 kpbs - 2015-06-13 15:23:25 +0200 CEST - 11996262 bytes

<<< d:\src\B004SBJE5Y_(disc_1)_08_-_Endless_Disco.mp3
>>> Rafale\Obsessions\08 Endless Disco.mp3
>>> 128 kpbs - 2015-06-13 15:23:40 +0200 CEST - 9381255 bytes

<<< d:\src\B004SBLM18_(disc_1)_03_-_Life_In_Mono.mp3
>>> Rafale\Obsessions\03 Life In Mono.mp3
>>> 128 kpbs - 2015-02-05 22:38:09 +0100 CET - 7921506 bytes

<<< d:\src\B004SBLM9K_(disc_1)_09_-_Bad_Obsession.mp3
>>> Rafale\Obsessions\09 Bad Obsession.mp3
>>> 128 kpbs - 2015-06-13 15:24:38 +0200 CEST - 10514242 bytes

<<< d:\src\B0088E2DNM_(disc_1)_01_-_One_Day___Reckoning_Song_(Wankelmut.mp3
>>> Asaf Avidan & the Mojos\One Day - Reckoning Song (Wankelmut Remix)\One Day - Reckoning Song (Wankelmut Remix) (Radio Edit).mp3
>>> 128 kpbs - 2013-12-31 13:20:10 +0100 CET - 6968360 bytes

<<< d:\src\BEAK.mp3
>>> Electric Guest\Mondo\08 Troubleman.mp3
>>> 256 kpbs - 2017-06-18 13:01:41.1277311 +0200 CEST - 16993428 bytes

<<< d:\src\ENLK.mp3
>>> Rafale\Obsessions\05 L'animale.mp3
>>> 256 kpbs - 2017-06-18 13:02:46.5050003 +0200 CEST - 8281516 bytes

<<< d:\src\GIAH.mp3
>>> Rafale\Obsessions\06 Marine Aircrash.mp3
>>> 256 kpbs - 2017-06-18 13:03:19.0606224 +0200 CEST - 6169162 bytes

<<< d:\src\IEHZ.mp3
>>> Alain Souchon\Nickel\20 On Avance (Live).mp3
>>> 256 kpbs - 2017-06-18 12:56:28.6081776 +0200 CEST - 8826537 bytes

<<< d:\src\MXGE.mp3
>>> Rafale\Obsessions\02 Beyond Bad.mp3
>>> 256 kpbs - 2017-06-18 13:04:36.9586571 +0200 CEST - 7219902 bytes

<<< d:\src\PPTD.mp3
>>> Rafale\Obsessions\04 Never Ever.mp3
>>> 256 kpbs - 2017-06-18 13:05:16.4989713 +0200 CEST - 10230880 bytes

<<< d:\src\QWJE.mp3
>>> Electric Guest\Mondo\02 This Head I Hold.mp3
>>> 256 kpbs - 2017-06-18 13:05:36.8443189 +0200 CEST - 5703655 bytes

<<< d:\src\SVZO.mp3
>>> Electric Guest\Mondo\03 Under the Gun.mp3
>>> 256 kpbs - 2017-06-18 13:06:05.8927963 +0200 CEST - 7186440 bytes

<<< d:\src\XJCY.mp3
>>> Air\Moon Safari\02 Sexy Boy.mp3
>>> 256 kpbs - 2017-06-18 13:07:19.2059191 +0200 CEST - 9664361 bytes


22 MP3 files (0 WMA files)
```
