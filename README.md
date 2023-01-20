[![percent](https://forthebadge.com/images/badges/0-percent-optimized.svg)](https://forthebadge.com) [![love](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com) [![go](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![black-magic](https://forthebadge.com/images/badges/powered-by-black-magic.svg)](https://forthebadge.com) [![ask](https://forthebadge.com/images/badges/you-didnt-ask-for-this.svg)](https://forthebadge.com)  [![release](https://img.shields.io/badge/Release-v2.0-green?style=for-the-badge)](https://forthebadge.com)

# partial_hash
Partial md5 of file (depending on the chosen buffer) because when you are in a hurry and that is necessary, this tool saves you time.

# Run

Use the right executable depending on your operating system and processor architecture. See the release for more information.

> **Note** Next, replace the name of the executable with your own.

```
USAGE:
    partial_hash [[-f <file>] or [-d <recursive_folder>]] [-b <buffer_size>] [--output csv]

OPTIONS (only -f or -d):
    -f, --file <file>                  File to be hashed.
    -d, --directory <recursive_path>   Folder where files will be recursively hashed.
    -b, --buffer <buffer_size>         Size of the buffer in bytes.
    --output                           Output format. (Default is stdout)
```

For example:

```
./partial_hash -f hello.txt -b 200
Version :  2.0
Size of the buffer (in bytes) :  200

hello.txt    5dd4a8e994400c28b21f8f42f54f443f

./partial_hash -d /folder/subfolder/
Version :  2.0
Size of the buffer (in bytes) :  100000000

/folder/subfolder/hello.txt       5dd4a8e994400c28b21f8f42f54f443f
/folder/subfolder/file1       2f282b84e7e608d5852449ed940bfc51
```

For Windows, same options but open *partial_hash.exe* with a command interpreter (tested with Powershell and cmd.exe).

> **Warning** The format of file paths on Windows is "C:\\..." if C is the root of the drive to be hashed. 

# Output format (--output)

No argument products :  \<filename>   \<hash>

"--output" is not possible in single file mode.

"--output csv" products : \<filename>;\<hash>;\<boolean if is a partial hash or not>


# How it works

A buffer (defined by the argument "-b" and by default 100MB) is used. 
If the file is smaller than the buffer size, then it will be calculated in md5.
If the file is larger than the buffer size, then the tool will take the first X and last X bytes (where X is the buffer) to calculate the md5. Without any modification of the buffer size, this creates a 200MB hash (the first 100MB and the last 100MB of the file).

# FAQ

## Why?

Why not...

## Why not a real md5 of the file? The integrity of this file is questionable with your tool and there is a possibility of hash collision in my lab...
1. The time :
```
$ time ./partial_hash -f 11G_urandom_file 
11G_urandom_file;f1a5ec97f391a2d72e269f8cb1c91a516bc5246855b0dcf9073578df463891b6;1/1

real    0m7.674s
user    0m6.285s
sys     0m3.945s

$ time md5sum 11G_urandom_file 
53f52c5d14cd467246a4d5bcaabd48ee032af8e679470654e2a0d014ab548e5e  11G_urandom_file

real    6m7.376s
user    4m20.007s
sys     0m45.039s
```
2. A hash of 200MB is not really questionable if you don't lose your eye index.

## How I can compare the hash afterwards ?
 
Uh, good question... I think, with your eyes and a ruler because the output is in lexical order. (Good news, I'm working on this)

---