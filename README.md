# coffee

![alt text](https://travis-ci.com/Hadyn/coffee.svg?branch=master "Build Status")

A Runescape file editing library for Go

## Installing

```bash
git clone https://github.com/hadyn/coffee
cd coffee
go test ./... && go install
```

## Examples

```bash
# Read the map manifest from the cache, write to file.
coffee cache read --index main_file_cache.idx255 \
                  --blocks main_file_cache.dat2  \
                  255 5 > manifest.maps.c.dat

# Read the map manifest from the cache, decompress it, write to file.
coffee cache read --index main_file_cache.idx255  --blocks main_file_cache.dat2 255 5 | \ 
coffee archive decompress > manifest.maps.d.dat


# Read the map manifest from the cache, decompress it, decode it, write to JSON file.
coffee cache read --index main_file_cache.idx255  --blocks main_file_cache.dat2 255 5 | \ 
coffee archive decompress | \
coffee index decode > manifest.maps.json
```