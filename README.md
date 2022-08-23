[![forthebadge](https://forthebadge.com/images/badges/0-percent-optimized.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/built-by-crips.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/powered-by-black-magic.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/uses-badges.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)

# partial_hash
Partial sha256 of file (depending if file > 100MB or not) because when you are very hurried and that is necessary, the script saves your time 

# Run

```
USAGE:
    partial_hash [-f <file>] or [-d <recursive_folder>]

OPTIONS (only one of them):
    -f, --file <file>                  File to be hashed.
    -d, --directory <recursive_path>   Folder where files will be recursively hashed.
```

For example:

```
./partial_hash -f hello.txt 
hello.txt;2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824;1/1

./partial_hash -d /folder/subfolder/
/folder/subfolder/file1;426175749bf6a302057e80109bb9a90fd7c56bda22c73c5a34bbc85197c25c2d;1/2
/folder/subfolder/hello.txt;2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824;2/2
```

For Windows, same options but open *partial_hash.exe* with a command interpreter (tested with powershell.exe and cmd.exe).

# How it works

If the file is less than 100MB, so the hash is a real sha256.  
If the file is greater than 100MB, so this take te first 100MB and the last 100MB of the file and hash this (200MB) with sha256.

# FAQ

## Why?

Why not...

## There are only 3 comments in your code?

"Comments are like the H of Hawaii"

## Why not a real sha256 of the file? The integrity of this file is questionable with your tool and there is a possibility of hash collision in my lab...
1. The time :
```
$ time ./partial_hash -f 11G_urandom_file 
11G_urandom_file;f1a5ec97f391a2d72e269f8cb1c91a516bc5246855b0dcf9073578df463891b6;1/1

real    0m7.674s
user    0m6.285s
sys     0m3.945s

$ time sha256sum 11G_urandom_file 
53f52c5d14cd467246a4d5bcaabd48ee032af8e679470654e2a0d014ab548e5e  11G_urandom_file

real    6m7.376s
user    4m20.007s
sys     0m45.039s
```
2. A hash of 200MB is not really questionable if you don't lose your eye index.

## How I can compare the hash afterwards ?
 
Uh, good question... I think, with your eyes and a ruler because the output is in lexical order. (Good news)

# Integrity (SHA256, the real one because the files are smaller than 100MB)
adcd38f1e80d83fc077fc9638f8e412996710bd8ff5bf100e41970b3bc696d26  partial_hash_amd64
f07a7e3e26dc5317de0bdf61b26f4f663a2736e4556c7ae66005f6c380231bad  partial_hash_amd64.exe
1d70ad3094cb2a3c592bcc1ea87772c41b09ae85a9a0426632b6a2fe038cc705  partial_hash_arm64
12b052db0aa13d46b97252b015b01cad60edbb2177401be4387ae4e719abb783  partial_hash.go