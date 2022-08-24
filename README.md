[![forthebadge](https://forthebadge.com/images/badges/0-percent-optimized.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/built-by-crips.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/powered-by-black-magic.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/uses-badges.svg)](https://forthebadge.com) [![forthebadge](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)

# partial_hash
Partial sha256 of file (depending on the chosen buffer) because when you are very hurried and that is necessary, the script saves your time.

# Run

Use the right executable depending on your operating system and processor architecture. See the release for more information.

> **Notes** Next, replace the name of the executable with your own.

```
USAGE:
    partial_hash [-f <file>] or [-d <recursive_folder>] [-b <buffer_size>]

OPTIONS (only one of them):
    -f, --file <file>                  File to be hashed.
    -d, --directory <recursive_path>   Folder where files will be recursively hashed.
    -b, --buffer <buffer_size>         Size of the buffer in bytes.
```

For example:

```
./partial_hash -f hello.txt -b 200
Size of the buffer (in bytes):  200
hello.txt;887a8f98c79fd7c3b036448d56fa2b6d8b91f4461125621587beeef63b1e4f29;1/1

./partial_hash -d /folder/subfolder/
Size of the buffer (in bytes):  100000000
/folder/subfolder/file1;426175749bf6a302057e80109bb9a90fd7c56bda22c73c5a34bbc85197c25c2d;1/2
/folder/subfolder/hello.txt;887a8f98c79fd7c3b036448d56fa2b6d8b91f4461125621587beeef63b1e4f29;2/2
```

For Windows, same options but open *partial_hash.exe* with a command interpreter (tested with powershell.exe and cmd.exe).

> **Warning** The format of file paths on Windows is "C:\\..." if C is the root of the drive to be hashed. 

# How it works

A buffer (defined by the argument "-b" and by default 100MB) is used. 
If the file is smaller than the buffer size, then it will be calculated in sha256.
If the file is larger than the buffer size, then the tool will take the first X and last X bytes (where X is the buffer) to calculate the sha256. Without any modification of the buffer size, this creates a 200MB hash (the first 100MB and the last 100MB of the file).

# FAQ

## Why?

Why not...

## There are only 3 comments in your code?

"Comments are like the H of Hawaii" and there are 15 comments...

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
 
Uh, good question... I think, with your eyes and a ruler because the output is in lexical order. (Good news, I work on this)

---