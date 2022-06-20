package v1

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_Disk(t *testing.T) {

}

func GetAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			GetAllFile(pathname + fi.Name() + "\\")
		} else {
			fmt.Println(fi.Name())
		}
	}
	return err
}

func Test_ArrFiles(t *testing.T) {
	path := "D:\\go\\src\\fsfc_store\\fs\\testDir\\"

	rd, _ := ioutil.ReadDir(path)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("Dir: %v\n", fi.Name())
		} else {
			fmt.Println(fi.Name())
		}
	}
}

func Test_de(t *testing.T) {
	dir := fmt.Sprintf(".\\%s", string(time.Now().UTC().String()))
	fmt.Printf(dir)
	os.MkdirAll(dir, 0777)
}
