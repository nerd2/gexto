package gexto_test

import (
	"testing"
	"os/exec"
	"io/ioutil"
	"log"
	"os"

	"github.com/stretchr/testify/require"
	"crypto/rand"
	"github.com/nerd2/gexto"
)

type TestFs struct {
	devFile string
	mntPath string
	t *testing.T
}

func NewTestFs(t *testing.T, sizeMb int, fsType string) *TestFs {
	f, err := ioutil.TempFile("", "gextotest")
	require.Nil(t, err)
	blank := make([]byte, 1024*1024)
	for i := 0; i < sizeMb; i++ {
		_, err = f.Write(blank)
		require.Nil(t, err)
	}
	err = f.Close()
	require.Nil(t, err)

	td, err := ioutil.TempDir("", "gextotest")
	require.Nil(t, err)

	err = exec.Command("mkfs." + fsType, f.Name()).Run()
	require.Nil(t, err)

	err = exec.Command("sudo", "mount", f.Name(), td).Run()
	require.Nil(t, err)

	err = exec.Command("sudo", "chmod", "-R", "777", td).Run()
	require.Nil(t, err)

	return &TestFs{f.Name(), td, t}
}

func (tfs *TestFs) Unmount() {
	if tfs.mntPath != "" {
		exec.Command("sudo", "umount", tfs.mntPath).Run()
		exec.Command("sudo", "rm", "-rf", tfs.mntPath).Run()
		tfs.mntPath = ""
	}
}

func (tfs *TestFs) Close() {
	tfs.Unmount()
	if true {
		os.Remove(tfs.devFile)
	} else {
		log.Println(tfs.devFile)
	}
}

func (tfs *TestFs) WriteSmallFile(path string, file string, b []byte) {
	err := os.MkdirAll(tfs.mntPath + path, 0777)
	require.Nil(tfs.t, err)
	err = ioutil.WriteFile(tfs.mntPath + path + "/" + file, b, 0777)
	require.Nil(tfs.t, err)
}

func (tfs *TestFs) WriteLargeFile(path string, file string, size int) *os.File {
	largefile, _ := ioutil.TempFile("", "gexto")
	for size > 0 {
		dataLen := 512*1024
		if dataLen > size {
			dataLen = size
		}
		data := make([]byte, dataLen)
		n, err := rand.Read(data)
		require.Nil(tfs.t, err)
		m, err := largefile.Write(data[:n])
		require.Nil(tfs.t, err)
		size -= m
	}
	err := largefile.Close()
	require.Nil(tfs.t, err)
	err = os.MkdirAll(tfs.mntPath + path, 0777)
	require.Nil(tfs.t, err)
	err = exec.Command("cp", largefile.Name(), tfs.mntPath + path + file).Run()
	require.Nil(tfs.t, err)
	return largefile
}

func doTestRead(t *testing.T, fsType string) {
	tfs := NewTestFs(t, 1100, fsType)
	defer func(){tfs.Close()}()

	text := []byte("hello world")
	tfs.WriteSmallFile("/", "smallfile", text)
	tfs.WriteSmallFile("/dir1", "smallfile", text)
	largefile := tfs.WriteLargeFile("/", "largefile", 987654321)
	defer os.Remove(largefile.Name())
	tfs.Unmount()

	fs, err := gexto.NewFileSystem(tfs.devFile)
	require.Nil(t, err)

	{
		file, err := fs.Open("/smallfile")
		require.Nil(t, err)
		out, err := ioutil.ReadAll(file)
		require.Nil(t, err)
		require.Equal(t, text, out)
	}

	{
		file, err := fs.Open("/dir1/smallfile")
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

func TestIntegrationRead(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	doTestRead(t, "ext2")
	doTestRead(t, "ext4")
}
