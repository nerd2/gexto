package main

import (
	"testing"
	"os/exec"
	"io/ioutil"
	"log"
	"os"
)

func TestIntegrationRead(t *testing.T) {
	f, err := ioutil.TempFile("", "gextotest")
	if err != nil {
		log.Fatalln(err)
	}
	blank := make([]byte, 1024*1024)
	for i := 0; i < 10; i++ {
		f.Write(blank)
	}
	f.Close()
	defer os.Remove(f.Name())

	td, err := ioutil.TempDir("", "gextotest")
	if err != nil {
		log.Fatalln(err)
	}

	err = exec.Command("mkfs.ext2", f.Name()).Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = exec.Command("sudo", "mount", f.Name(), td).Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = exec.Command("sudo", "chmod", "777", td).Run()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {exec.Command("sudo", "umount", f.Name()).Run()}()

	err = ioutil.WriteFile(td + "/testfile", []byte("hello world"), 777)
	if err != nil {
		log.Fatalln(err)
	}

	exec.Command("sudo", "umount", f.Name()).Run()
	
	ReadFile(f.Name())
}
