package jagex

import (
    "io/ioutil"
    "path/filepath"
    "testing"
)

func loadBytes(t *testing.T, name string) []byte {
    path := filepath.Join("testdata", name)
    bs, err := ioutil.ReadFile(path)
    if err != nil {
        t.Fatal(err)
    }
    return bs
}