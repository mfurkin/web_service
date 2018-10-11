package fileChecker

import (
	"log"
	"os"
	"io"
	"time"
)

type FileChecker struct {
	ticker *time.Ticker 
}

func appendWithFiles(oldFiles,newFiles []os.FileInfo) []os.FileInfo {
	for _, file := range newFiles {
		oldFiles = append(oldFiles,file);
	}
	return oldFiles;
}

func getFiles(workDir* os.File) ([]os.FileInfo, error) {
	var oldFiles,curFiles []os.FileInfo;
	var err error;
	oldFiles = make([]os.FileInfo,0,10);
	
	for curFiles, err = workDir.Readdir(10);err == nil;oldFiles = appendWithFiles(oldFiles,curFiles) {
		
	}
	if err != io.EOF {
		return nil,err
	}
	return oldFiles,nil
}
func checkThisDir(workDir *os.File) { 
	
}

func errMsg(err error) {
	log.Println("Some issue with work directory creating"+err.Error());	
}

func (fc *FileChecker) Start(fname string) error {
	var workDir *os.File
	var err error
	workDir, err = os.Open(fname);
	if (err != nil) {
		if !os.IsNotExist(err) {
		 	errMsg(err);
		 	return err;
		 }
		 workDir, err = os.Create(fname);
		 if (err != nil) {
		 	errMsg(err);
		 	return err;
		 }
	}	
//	oldFiles, err = getFiles(workDir)
	if err != nil {
		return nil
	}
	fc.ticker = time.NewTicker(5*time.Second)
	for _ = range fc.ticker.C {
		checkThisDir(workDir)
	}
	return nil;
}