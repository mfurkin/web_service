package fileChecker

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
    "testing"
)
func testErr(err error, msg string, t *testing.T) {
	if (err != nil) {
		log.Println(msg+": "+err.Error())
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
	len1  := len(files)
	log.Printf("afterTest pt 1 files len=%d",len1)
	for err == nil && len(files) > 0 {
		for _ , d := range files {
			fname2 := d.Name()
			log.Println("Removing file "+fname2)
			filepath := path.Join(*fname,fname2)
			err = os.Remove(filepath)
			testErr(err,"Error during afterTest file deleting",t)
		}
		files, err = workDir.Readdir(10)
		len1  := len(files)
		log.Printf("afterTest pt2 files len=%d",len1)		
	}
	workDir.Close()
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
func checkFileNamesTest (resFiles,files []string, t *testing.T, msg string ) {
	if !checkFileNames(resFiles,files) {
		log.Println(msg)
		t.Fail()
	}
}
func checkFileNames(resFiles,files []string) bool {
	res := false
	lenRes := len(resFiles)
	len1 := len(files)
	log.Printf("checkFileNames lenRes=%d len1=%d\n",lenRes,len1)
//	if len(resFiles) == len(files) {
	if lenRes == len1 {
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
/*
func checkFilePath(fname string, dir* string) bool {
	path := path.Join(*dir,fname)
	_,err := os.Stat(path)
	res := err == nil
	log.Printf("%s = %t\n", path,res)	
	return res;
}

func checkFilePaths(files []string, dir *string) bool {
	var res bool 
	log.Printf("checkFilePaths len=%d\n",len(files))
	for _,d := range files {
		res = checkFilePath(d,dir)
		if !res {
			break
		}
	}
	return res;
}
*/
func TestGetFilesEmptyDir(t *testing.T) {
	fname := "/home/mikhailf/web_work/empty"
	resFiles, err := getFiles(&fname)
	if err != nil || len(resFiles) != 0 {
		t.FailNow()
	}
}
func testGetFilesWithFiles(allFiles []string,t *testing.T,fname *string) {
	var files = allFiles[0:]
	resFiles, err := getFiles(fname)
	if err != nil || !checkFileNames(resFiles,files) {
		t.FailNow()
	}	
}
func TestGetFiles1File(t *testing.T) {
	fname := "/home/mikhailf/web_work/file1"
	allFiles  := [1]string{"big.txt"}
	testGetFilesWithFiles(allFiles[0:],t,&fname)
}

func TestGetFiles3File(t *testing.T) {
	fname := "/home/mikhailf/web_work/file3"
	allFiles := [3]string{"big.txt","blocks.txt","blocks_0.txt"}
	testGetFilesWithFiles(allFiles[0:],t,&fname)
}

func TestGetFiles6File(t *testing.T) {
	fname := "/home/mikhailf/web_work/file6"
	allFiles := [6]string{"big.txt","blocks.txt","blocks_0.txt","extracted_from_370000 samples.txt","just.txt","sinusoids.txt"}
	testGetFilesWithFiles(allFiles[0:],t,&fname)
}
func copyFile(file string, fromDir *string, ToDir *string) error {
	srcFile := path.Join(*fromDir,file)
	dstFile := path.Join(*ToDir,file)
	cmd := exec.Command("cp", srcFile, dstFile)
	stderr, err := cmd.StderrPipe()	
	if err != nil {
		log.Println("error during stderr getting: "+err.Error())
	}
	err = cmd.Start()
	if err != nil {
		log.Println("error during cmd.Start: " + err.Error())
	}
	buf, _ := ioutil.ReadAll(stderr)
	log.Printf("%s\n",buf)
	if err := cmd.Wait();err != nil {
		log.Println("error during wait: "+err.Error())
	} else {
		log.Println("File "+srcFile+" copied successfully to "+dstFile)
	}
	return err
//	return cmd.Run()
}
func copyFiles(files []string, fromDir *string, ToDir *string) error {
	var err error
	for _,d := range files {
		err = copyFile(d,fromDir,ToDir)
		if err != nil {
			log.Println("copyFiles error: "+err.Error())
			break
		}
	}
	return err;
}
func removeFile(dir *string, fname string) error {
	path := path.Join(*dir,fname)
	return os.Remove(path)
}
func TestCheckThisDir1(t *testing.T) {
	fname := "/home/mikhailf/web_work/workdir1"
	fromDir := "/home/mikhailf/gotest_files"
	allFiles := [6]string{"big.txt","blocks.txt","blocks_0.txt","extracted_from_370000 samples.txt","just.txt","sinusoids.txt"}
	resFiles,err := getFiles(&fname)
	testErr(err,"error in checkThisDir pretest",t)
	err = copyFiles(allFiles[0:1],&fromDir,&fname)
	testErr(err,"error in file copying",t)
	res2Files, err := checkThisDir(&fname,resFiles)
	testErr(err,"error in checkThisDir test 1",t)
	checkFileNamesTest(res2Files,allFiles[0:1],t,"incorrect values checkThisDir 1")
	err = copyFiles(allFiles[1:2],&fromDir,&fname)
	testErr(err,"error in file copying 2",t)	
	res3Files, err := checkThisDir(&fname,res2Files)
	testErr(err,"error in checkThisDir test 2",t)
	checkFileNamesTest(res3Files,allFiles[0:2],t,"incorrect values checkThisDir 2")
	err = removeFile(&fname,allFiles[0])
	testErr(err,"error during file deleting",t)
	err = copyFiles(allFiles[2:3],&fromDir,&fname)
	testErr(err,"error in file copying 3",t)	
	res4Files, err := checkThisDir(&fname,res3Files)
	testErr(err,"error in checkThisDIr test 3",t)
	checkFileNamesTest(res4Files,allFiles[1:3],t,"incorrect values checkThisDir 3")
	defer afterTest(&fname,t)
}