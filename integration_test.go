package gexto

import (
	"testing"
	"os/exec"
	"io/ioutil"
	"log"
	"os"

	"github.com/stretchr/testify/assert"
)

func createAndMountFs() (string, string, error) {
	f, err := ioutil.TempFile("", "gextotest")
	if err != nil {
		log.Fatalln(err)
	}
	blank := make([]byte, 1024*1024)
	for i := 0; i < 10; i++ {
		f.Write(blank)
	}
	f.Close()

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

	return f.Name(), td, nil
}

func unmountFs(fname string) error {
	exec.Command("sudo", "umount", fname).Run()
	return nil
}

func TestIntegrationRead(t *testing.T) {
	devPath, mntPath, _ := createAndMountFs()
	if false {
		defer os.Remove(devPath)
	} else {
		log.Println(devPath)
	}

	text := []byte("hello world")
	err := ioutil.WriteFile(mntPath + "/testfile", text, 777)
	unmountFs(devPath)
	if err != nil {
		log.Fatalln(err)
	}

	fs, err := NewFileSystem(devPath)
	assert.Nil(t, err)
	file, err := fs.Open("/testfile")
	assert.Nil(t, err)
	out, err := ioutil.ReadAll(file)
	assert.Nil(t, err)
	assert.Equal(t, text, out)
}
