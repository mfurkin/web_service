package fileChecker

import (
	"log"
	"os"
    "testing"
)
func testErr(err error, msg string, t *testing.T) {
	if (err != nil) {
		log.Println(msg+err.Error())
		t.Fail()
	}	
}

func beforeTest(fname *string, t *testing.T) (*os.File, error) {
	var workDir, err = os.Open(*fname)
	testErr(err,"Error during work dir opening",t);
	return workDir, err
}

func afterTest(fname *string, t *testing.T) {
	var workDir, err = os.Open(*fname)
	testErr(err,"Error during afterTest",t)
	files, err := workDir.Readdir(10)	
	for err == nil && len(files) > 0 {
		for _ , d := range files {
			err = os.Remove(d.Name())
			testErr(err,"Error during afterTest file deleting",t)
		}
		files, err = workDir.Readdir(10)
	}
}



func fileFound(files []string, file string) bool {
	res := false
	for _, d := range files {
		if d == file {
			res = true
			break;
		}
	}
	return res;
}

func checkFileNames(resFiles,files []string) bool {
	res := false
	if len(resFiles) == len(files) {
		res = true
		for _,d := range files {
			if !fileFound(resFiles,d)  {
				res = false
				break
			}
		}
	}
	return res;
}

func TestEmptyDir(t *testing.T) {
	fname := "/home/mikhailf/web_work/empty"
	workDir, _ := beforeTest(&fname,t)
	resFiles, err := getFiles(workDir)
	if err != nil || len(resFiles) != 0 {
		t.Fail()
	}
}
func testWithFiles(allFiles []string,t *testing.T,fname *string) {
	var files = allFiles[0:]
	workDir, _ := beforeTest(fname,t)
	resFiles, err := getFiles(workDir)
	if err != nil || !checkFileNames(resFiles,files) {
		t.Fail()
	}	
}
func Test1File(t *testing.T) {
	fname := "/home/mikhailf/web_work/file1"
	allFiles  := [1]string{"big.txt"}
	testWithFiles(allFiles[0:],t,&fname)
}

func Test3File(t *testing.T) {
	fname := "/home/mikhailf/web_work/file3"
	allFiles := [3]string{"big.txt","blocks.txt","blocks_0.txt"}
	testWithFiles(allFiles[0:],t,&fname)
}

func Test8File(t *testing.T) {
	t.Fail()
	// TODO implement this test
}

