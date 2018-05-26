package gexto

import (
	"testing"
	"os/exec"
	"io/ioutil"
	"log"
	"os"

	"github.com/stretchr/testify/require"
	"crypto/rand"
)

func createAndMountFs() (string, string, error) {
	f, err := ioutil.TempFile("", "gextotest")
	if err != nil {
		log.Fatalln(err)
	}
	blank := make([]byte, 1024*1024)
	for i := 0; i < 1024; i++ {
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
	if true {
		defer os.Remove(devPath)
	} else {
		log.Println(devPath)
	}

	text := []byte("hello world")
	err := ioutil.WriteFile(mntPath + "/smallfile", text, 777)
	require.Nil(t, err)

	largefile, _ := ioutil.TempFile("", "gexto")
	len := 987654321
	for len > 0 {
		dataLen := 512*1024
		if dataLen > len {
			dataLen = len
		}
		data := make([]byte, dataLen)
		n, _ := rand.Read(data)
		m, _ := largefile.Write(data[:n])
		len -= m
	}
	err = largefile.Close()
	require.Nil(t, err)
	defer os.Remove(largefile.Name())
	err = exec.Command("cp", largefile.Name(), mntPath + "/largefile").Run()
	require.Nil(t, err)
	unmountFs(devPath)

	fs, err := NewFileSystem(devPath)
	require.Nil(t, err)

	{
		file, err := fs.Open("/smallfile")
		require.Nil(t, err)
		out, err := ioutil.ReadAll(file)
		require.Nil(t, err)
		require.Equal(t, text, out)
	}

	{
		file, err := fs.Open("/largefile")
		require.Nil(t, err)
		comparefile, err := os.Open(largefile.Name())
		for err == nil {
			a := make([]byte, 1024*1024)
			b := make([]byte, 1024*1024)
			var na int
			na, err = file.Read(a)
			nb, err2 := comparefile.Read(b)
			require.Equal(t, na, nb)
			log.Printf("Read %d (%d)", na, nb)
			require.Equal(t, a[:na], b[:nb])
			require.Equal(t, na, nb)
			require.Equal(t, err, err2)
		}
	}

}
