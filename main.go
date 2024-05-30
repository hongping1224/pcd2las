package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/hongping1224/lidario"

	"github.com/hongping1224/pcd2las/lidarpal"
)

var numOFCPU int

func main() {
	workerCount := flag.Int("cpu", runtime.NumCPU(), "set Cpu usage")
	dir := flag.String("dir", "./", "input Folder")
	headersample := flag.String("header", "./header/headersample.las", "input Folder")
	flag.Parse()
	runtime.GOMAXPROCS(*workerCount)

	//check directory exist
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatal(err)
		return
	}
	//find all las file
	xyzs := findFile(*dir, ".pcd")

	convert(xyzs, *headersample)
}

func convert(xyzs []string, header string) error {
	for _, xyz := range xyzs {
		headerExample, err := lidario.NewLasFile(header, "rh")
		if err != nil {
			return err
		}
		laspath := strings.Replace(xyz, ".pcd", ".las", -1)
		fmt.Println(laspath)
		// open las file
		las, err := lidario.InitializeUsingFile(laspath, headerExample)
		las.Header.PointFormatID = 0
		headerExample.Close()
		if err != nil {
			return err
		}
		writechan := make(chan lidario.LasPointer, numOFCPU*8)
		writer := lidarpal.NewWriter(writechan)
		writer.Serve(las)
		var wg sync.WaitGroup
		file, err := os.Open(xyz)
		if err != nil {
			return err
		}
		//setup reader for each file
		scanner := bufio.NewScanner(file)
		wg.Add(1)
		reader := lidarpal.NewReader(scanner, &wg)
		reader.Serve(writechan)
		wg.Wait()
		writer.Close()
		file.Close()
	}
	return nil
}
func findFile(root string, match string) (file []string) {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if strings.HasSuffix(info.Name(), match) {
			file = append(file, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Total shp file : ", len(file))
	return file
}
