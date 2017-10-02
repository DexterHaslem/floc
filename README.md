## FLOC

This is a simple tool to output physical lines of code by extension time for
all files in a directory, recursively.

### Usage
FLOC has 3 command line parameters, which need to be in the go format of `-name=value`:
    - `-dir` - sets starting directory to report.
    - `-filter` - comma separated list of file extensions to report, eg `-filter=go,ts`. **Optional**
    - `-ignore` - comma separated list of paths or filenames to ignore, eg `-ignore=node_modules,.git`.  **Optional**

### Output
The output (stdout) is JSON and looks something like this
```json

{
  "byDir": {
    ".": {
      "go": 53
    },
    "foo\\bar": {
      "ts": 27
    },
    "foo\\baz": {
          "go": 111
        },
  },
  "total": {
    "go": 164,
    "ts": 27
  }
}

```