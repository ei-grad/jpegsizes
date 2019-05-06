jpegsizes
=========

Calculates the sizes of many JPEG images fast.

Install
-------

go get github.com/ei-grad/jpegsizes

Usage
-----

Just run the `jpegsizes` command to output tab-separated `<image> <width> <height>` tuples for `.jpg` files in current directory.

```
  -pattern string
    	glob pattern to match image filenames (default "*.jpg")

  -workers int
    	number of images to process in parallel (default 4)
```
