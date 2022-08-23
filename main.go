package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"
)

func main() {
	b, err := ioutil.ReadFile("target.md")
	if err != nil {
		panic(err)
	}
	structInfos := StructInfos{
		StructInfos: Parse(string(b)),
	}
	b, err = ioutil.ReadFile("template")
	if err != nil {
		panic(err)
	}
	t := template.Must(template.New("struct").Parse(string(b)))
	w := bytes.NewBuffer(nil)
	err = t.Execute(w, structInfos)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("out.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(w.Bytes())
	if err != nil {
		panic(err)
	}
}
