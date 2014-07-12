package resourcelocator

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var embeddedFiles = make(map[string][]byte)

func SetEmbeddedFiles(newMap map[string][]byte) {
	embeddedFiles = newMap
}

func Locate(path string) []byte {
	path = path[1:]
	ret := embeddedFiles[path]
	if ret != nil {
		return ret
	}
	return locateExternal(path)
}

func locateExternal(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b := new(bytes.Buffer)
	b.ReadFrom(file)
	return b.Bytes()
}

func CreateEmbeddedLocator(fileName, packageName, resourceDir string) {
	info, err := os.Stat(resourceDir)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	appendHeader(file, packageName)
	defer func(f *os.File) {
		appendFooter(f)
		f.Close()
	}(file)
	processFile(file, info.Name(), info)
	fmt.Println("Done")
}

func appendHeader(writer io.Writer, packageName string) {
	fmt.Fprintf(writer, "package %s\n\n", packageName)
	fmt.Fprintf(writer, "import \"github.com/baardsen/resourcelocator\"\n\n")
	fmt.Fprintln(writer, "func init() {")
	fmt.Fprintln(writer, "\tembeddedFiles := make(map[string] []byte)")
}

func appendFooter(writer io.Writer) {
	fmt.Fprintln(writer, "\tresourcelocator.SetEmbeddedFiles(embeddedFiles)")
	fmt.Fprintln(writer, "}")
}

func processFile(writer io.Writer, path string, info os.FileInfo) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %+v\n", path, err)
		return
	}
	defer file.Close()

	if info.IsDir() {
		files, _ := file.Readdir(0)
		for _, f := range files {
			processFile(writer, path+"/"+f.Name(), f)
		}
		return
	}
	slice := make([]byte, info.Size())
	_, err = file.Read(slice)
	if err != nil {
		fmt.Printf("Error reading %s: %+v\n", path, err)
		return
	}
	fmt.Fprintf(writer, "\tembeddedFiles[\"%s\"] = []byte{", path)
	for _, b := range slice {
		fmt.Fprintf(writer, "0x%x, ", b)
	}
	fmt.Fprintln(writer, "}")
}
