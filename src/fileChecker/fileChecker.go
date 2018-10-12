package fileChecker

import (
	"log"
	"os"
	"io"
	"time"
)
// Сущность, постоянно проверяющая рабочий каталог
type FileChecker struct {
	ticker *time.Ticker 
}
// Функция добавляет один слайс в другой. Может, слишком велосипедно, но по-другому у меня не получилось
func appendWithFiles(oldFiles,newFiles []string) []string {
	for _, file := range newFiles {
		oldFiles = append(oldFiles,file);
	}
	return oldFiles;
}
// Функция для получения текущего списка файлов рабочего каталога  
func getFiles(workDir* os.File) ([]string, error) {
	var oldFiles,curFiles []string;
	var err error;
	oldFiles = make([]string,0,10);
	
	for curFiles, err = workDir.Readdirnames(10);err == nil;oldFiles = appendWithFiles(oldFiles,curFiles) {
		
	}
	if err != io.EOF {
		return nil,err
	}
	return oldFiles,nil
}

// Основная функция данного пакета. Получает текущий списко файлов, сравнивает его с предыдущим и выдает наружу
// свежий или предыдущий, если они не изменились
func checkThisDir(workDir *os.File,oldFiles []string) ([]string, error) {
	var err error 
	curFiles, err := getFiles(workDir)
	if err != nil && err != io.EOF {
		return nil,err
	}
	old_len := len(oldFiles)
	cur_len := len(curFiles)
	if old_len != cur_len {
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
func (fc *FileChecker) Process(fname string) error {
	var workDir *os.File
	var err error
	workDir, err = os.Open(fname);
	if (err != nil) {
		if !os.IsNotExist(err) {
		 	errMsg("Some issue with work directory creating",err);
		 	return err;
		 }
		 workDir, err = os.Create(fname);
		 if (err != nil) {
		 	errMsg("Some issue with work directory creating",err);
		 	return err;
		 }
	}	
	curFiles, err := getFiles(workDir)
	if err != nil {
		return nil
	}
	fc.ticker = time.NewTicker(5*time.Second)
	for _ = range fc.ticker.C {
		oldFiles := curFiles;
		curFiles, err = checkThisDir(workDir, oldFiles)
		if err != nil {
			errMsg("Error during file checking: ",err)
		}
	}
	return nil;
}