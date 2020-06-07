package main

import (
	"fmt"
	"os"
	"path/filepath"
	//"transcoder/builder/cmd"
)

func main() {
	buildPath := filepath.Join("build","server")
	fmt.Printf("Create folder %s\n",buildPath)
	if err:=os.MkdirAll(buildPath,os.ModePerm);err!=nil {
		panic(err)
	}
	//cmd.BuildServer([]string{"windows-amd64","linux-amd64","darwin-amd64"})
	//cmd.Execute()
}