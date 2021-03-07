package lib

import (
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Read for read
func Read() []byte {
	dat, err := ioutil.ReadFile("lib/file/t.txt")
	check(err)
	return dat
}

//Write files
func Write(thing string) {
	dat := []byte(thing)
	err := ioutil.WriteFile("lib/file/a.txt", dat, 0644)
	check(err)
}
