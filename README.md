# gonixutils
redo of coreutils in go with multithreading, multi-byte encoding (utf8) and cleaned up syntax in mind

## filesystem
### ls
list file system entries.

```bash
ls ~/go/src/github.com/ProhtMeyhet/gonixutils/filesystem/ls/
api.go         constants.go   decorators.go  input.go       library.go     ls  tools.go

ls -l ~/go/src/github.com/ProhtMeyhet/gonixutils/filesystem/ls/
-rw-r--r-- pm  100  1235 Aug 19 14:32:23 api.go
-rw-r--r-- pm  100  590  Aug 15 19:25:04 constants.go
-rw-r--r-- pm  100  2224 Aug 15 19:39:51 decorators.go
-rw-r--r-- pm  100  813  Aug 17 13:01:51 input.go
-rw-r--r-- pm  100  5808 Aug 19 14:26:50 library.go
drwxr-xr-x pm  100  30   Aug 19 19:30:55 ls
-rw-r--r-- pm  100  685  Aug 17 14:24:35 tools.go
```

### mk
create one file system entry or recursivly directorys (even with one file or link). mk implements:
 * mkdir
 * ln
 * mktemp
 * file creation
 
 ```bash
 # by default create a directory
 mk myDirectory
 mk --file myFile
 mk --link myFile myLink
 mk --symbolic myFile mySymbolikLink
 
 # recursive
 mk -r myDirectory/anotherDirectory
 mk -rf myDirectory/anotherDirectory/moreDirectory/file
 # link needs implementation for recursive
 ```

### rm
remove one file system entry or, with -r --recursive, recursivly remove all file system entries under a path.
it removes empty directorys without -r --recursive or error. rm implements
 * rm
 * rmdir

## text

### cat
print one file or concate two or more files and print them.

### head
print first parts of file, by default 10 lines. can also print first bytes and first runes.

### hashsum
print the hashsum of FILES... or compare computed hashsums in a file to files. hashsum implements:
 * cksum
 * md5sum
 * sha1sum
 * sha256sum
 * sha512sum

## miscellaneous

### true & false
exit with 0 or 1 respectivly.

### sleep
sleep for N seconds or, with -u --until, until a given clocktime or date.
