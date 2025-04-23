# gotail

Installation:

```bash
go install github.com/corani/gotail@latest
```

Usage:

```bash
$ gotail -n 10 -- tree /
│           │   │   │   │   │   │   │       ├── v0.1.13.lock
│           │   │   │   │   │   │   │       ├── v0.1.13.mod
│           │   │   │   │   │   │   │       ├── v0.1.13.zip
│           │   │   │   │   │   │   │       ├── v0.1.13.ziphash
│           │   │   │   │   │   │   │       ├── v0.1.14.info
│           │   │   │   │   │   │   │       ├── v0.1.14.lock
│           │   │   │   │   │   │   │       ├── v0.1.14.mod
│           │   │   │   │   │   │   │       ├── v0.1.14.zip
│           │   │   │   │   │   │   │       ├── v0.1.14.ziphash

```

This executes the command (after the `--`) in the background and follows its
output, showing a "viewport" with the last 10 lines.
