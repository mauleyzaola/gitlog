# Gitlog

[![Build Status](https://travis-ci.org/mauleyzaola/gitlog.svg?branch=master)](https://travis-ci.org/mauleyzaola/gitlog)

Simple parsing of git repositories and get instant metrics in different formats.

The idea behind this project is make easy data transformation of the contents of git files, into more workable information.

So far, we are only focusing to output plain text without UI. This can change in the future.

Data output can come out in different formats such as JSON, DB Engines, XML and so forth.

## Installation
```
go get github.com/mauleyzaola/gitlog
```

## Examples
From within the same repo, just type `gitlog` it will take default value `-directory=./.git`
```
gitlog
```
Result is a JSON array of objects. Each one is a commit (merges are excluded)
```
[{"hash":"052453b347706cef9437eb79e703c0dc625e7bef","author":{"name":"mauleyzaola","email":"mauricio.leyzaola@gmail.com"},"date":"2018-09-14T00:44:32-05:00","comment":"#15 - consider full names for authors","added"...
```

You can point to another directory as well, just pass the path to the `.git/` directory. It can be relative to your current path, or absolute, either will work. For instance, these would achieve the same result, considering you are at `$GOPATH/src/github.com`
```
gitlog -directory $GOPATH/src/github.com/golang/protobuf/.git
```
```
gitlog -directory ../github.com/golang/protobuf/.git
```
```
gitlog -directory golang/protobuf/.git
```
```
gitlog -directory ./golang/protobuf/.git
```

The result goes to stdout, so it can be used to input another program. For instance `jq` to pretty format the result.
```bash
gitlog | jq .
             [
               {
                 "hash": "052453b347706cef9437eb79e703c0dc625e7bef",
                 "author": {
                   "name": "mauleyzaola",
                   "email": "mauricio.leyzaola@gmail.com"
                 },
                 "date": "2018-09-14T00:44:32-05:00",
                 "comment": "#15 - consider full names for authors",
                 "added": 5,
                 "removed": 5
               },
               ...

```

## Parameters
```bash
Usage of gitlog:
  -alsologtostderr
    	log to standard error as well as files
  -directory string
    	the path to the the .git directory (default "./.git")
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -output string
    	the type of output to have: [commits] (default "commits")
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```