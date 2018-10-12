package fileChecker

import (
	"log"
	"os"
	"io"
	"time"
)
const FILESMAX = 5
const CHECKPERIOD = 5 // secs
// Сущность, постоянно проверяющая рабочий каталог
type FileChecker struct {
	fname string;
	ticker *time.Ticker 
}
// Функция добавляет один слайс в другой. Может, слишком велосипедно, но по-другому у меня не получилось
func appendWithFiles(oldFiles,newFiles []string) []string {
	for _, file := range newFiles {
		oldFiles = append(oldFiles,file);
	}
	return oldFiles;
}
func fileOpen(fname *string) (*os.File, error){
	workDir, err := os.Open(*fname);
	if (err != nil) {
		if !os.IsNotExist(err) {
		 	errMsg("Some issue with work directory creating",err);
		 	return nil,err;
		 }
		 workDir, err = os.Create(*fname);
		 if (err != nil) {
		 	errMsg("Some issue with work directory creating",err);
		 	return nil,err;
		 }
	}	
	return workDir,nil
}
// Функция для получения текущего списка файлов рабочего каталога  

func getFiles(fname *string) ([]string, error) {
	var oldFiles,curFiles []string;
	var err error;
	workDir, err := fileOpen(fname)
	if err != nil {
		return nil,err
	}
	defer workDir.Close()
	oldFiles = make([]string,0,FILESMAX);
	curFiles, err = workDir.Readdirnames(FILESMAX)
	if err != nil && err != io.EOF {
		log.Println("getFiles error during readdirnames: "+err.Error())
		return nil,err 
	}
	len1 := len(curFiles)
	for err == nil {
		oldFiles = appendWithFiles(oldFiles,curFiles)
		curFiles, err = workDir.Readdirnames(FILESMAX)		
	}
	len2 := len(curFiles)	
	if err != io.EOF {
		return nil,err
	}
	return oldFiles,nil
}

// Основная функция данного пакета. Получает текущий списко файлов, сравнивает его с предыдущим и выдает наружу
// свежий или предыдущий, если они не изменились

func checkThisDir(fname *string,oldFiles []string) ([]string, error) {
	var err error 
	log.Println("checkThisDir work dir name: "+*fname)
	curFiles, err := getFiles(fname)
	if err != nil && err != io.EOF {
		return nil,err
	}
	oldLen := len(oldFiles)
	curLen := len(curFiles)
	if oldLen != curLen {
		return curFiles,nil
	}
	for i,d:= range oldFiles {
		if d != curFiles[i] {
			return curFiles, nil
		}
	}
	return oldFiles,nil
}
// Вспомогательная фунекция - пишет лшибки в лог
func errMsg(msg string, err error) {
	log.Println(msg+err.Error());	
}
// Основная функция - регулярно проверяет рабочий каталог
func (fc *FileChecker) Process() error {

	var err error
	
	curFiles, err := getFiles(&fc.fname)
	if err != nil {
		return nil
	}
	fc.ticker = time.NewTicker(CHECKPERIOD*time.Second)
	for _ = range fc.ticker.C {
		oldFiles := curFiles;
		curFiles, err = checkThisDir(&fc.fname, oldFiles)
		if err != nil {
			errMsg("Error during file checking: ",err)
		}
	}
	return nil;
}