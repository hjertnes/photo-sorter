package main

import (
	"fmt"
	"github.com/hjertnes/photo-sorter/dateparser"
	"github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


func help(){
	fmt.Println("photo-sorter is a simple utility that renames files based on dates")
	fmt.Println("It first tries to read exif data, then file system date or to \"now\" if neither succeeds")
	fmt.Println()
	fmt.Println("Usage: ")
	fmt.Println("  photo-sorter [input-dir] [output-dir]: dry run runs through but doesn't write anything")
	fmt.Println("  photo-sorter [input-dir] [output-dir] --flat: moves into a flat structure like outdir/YYYY-MM-dd HH:MM:SS orignal-filename")
	fmt.Println("  photo-sorter [input-dir] [output-dir] --nested into a nested structure like outdir/YYYY/MM/YYYY-MM-dd HH:MM:SS orignal-filename")
	fmt.Println()
}


func main(){
	mode := "dry"

	if len(os.Args) < 3{
		help()
		os.Exit(0)
	}

	if len(os.Args) == 4{
		if os.Args[3] == "--flat"{
			mode = "flat"
		}
	}

	if len(os.Args) == 4{
		if os.Args[3] == "--nested"{
			mode = "nested"
		}
	}

	in := os.Args[1]
	out := os.Args[2]

	if !utils.FileExist(in){
		fmt.Println("Input dir doesn't exist")
	}

	if !utils.FileExist(out){
		err := os.MkdirAll(out, 0700)
		if err != nil{
			fmt.Println("Failed to create not existing output dir")
			os.Exit(0)
		}
	}

	err := filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if info.IsDir(){
			return nil
		}

		if strings.HasPrefix(info.Name(), "."){
			return nil
		}

		d := dateparser.GetDate(path)

		outputPath := ""

		if mode == "flat"{
			outputPath = fmt.Sprintf("%s/%s %s", out, d.Format("2006-01-02 15:04:05"), info.Name())
		}

		if mode == "nested"{
			folderPath := fmt.Sprintf("%s/%s", out, d.Format("2006/01"))
			if !utils.FileExist(folderPath){
				err = os.MkdirAll(folderPath, 0700)
				if err != nil {
					fmt.Printf("Faield to create output folder")
					return nil
				}
			}

			outputPath = fmt.Sprintf("%s/%s %s", folderPath, d.Format("2006-01-02 15:04:05"), info.Name())
		}

		if utils.FileExist(outputPath){
			fmt.Printf("File %s exist\n", outputPath)
			return nil
		}

		f, err := ioutil.ReadFile(path)
		if err != nil{
			fmt.Printf("Failed to read file %s\n", path)
			return nil
		}

		if mode != "dry" {
			err = ioutil.WriteFile(outputPath, f, 0600)
			if err != nil{
				fmt.Printf("Failed to write file %s\n", outputPath)
				return nil
			}
		}

		return nil
	})

	if err != nil{
		fmt.Println("Error")
		fmt.Println(eris.ToString(err, true))
		os.Exit(1)
	}
}