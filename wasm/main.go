package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"runtime"
	"syscall/js"
)

func main() {
	c := make(chan bool)
	js.Global().Set("listZip", js.FuncOf(listZip))
	<-c
}

func listZip(this js.Value, args []js.Value) interface{} {
	buf := make([]byte, args[0].Get("byteLength").Int())
	js.CopyBytesToGo(buf, args[0])

	fmt.Println("bytes is loaded into Go, len:", len(buf))

	r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return err
	}

	runtime.GC()

	r.RegisterDecompressor(zip.Deflate, nil)

	fmt.Println(len(r.File), "files in zip")
	for _, file := range r.File {
		if file.FileHeader.FileInfo().IsDir() {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		b, err := io.ReadAll(rc)
		if err != nil {
			return err
		}

		args[1].Call("postMessage", map[string]interface{}{
			"action": "exec",
			"sql":    "INSERT INTO files VALUES ($name, $size);",
			"params": map[string]interface{}{
				"$name": file.Name,
				"$size": len(b),
			},
		})
	}

	return nil
}
