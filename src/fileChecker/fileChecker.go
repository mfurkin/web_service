package fileChecker

import (
	"log"
	"os"
	"io"
	"net/http"
	"path"
)
const FILESMAX = 5
const MAXBUF=1024
// const CHECKPERIOD = 5 // secs
// Сущность, постоянно проверяющая рабочий каталог
type FileChecker struct {
	fname string;
//	ticker *time.Ticker 
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
	for err == nil {
		oldFiles = appendWithFiles(oldFiles,curFiles)
		curFiles, err = workDir.Readdirnames(FILESMAX)		
	}
	if err != io.EOF {
		return nil,err
	}
	return oldFiles,nil
}
// Функция удаляет файл по запросу
func (fc *FileChecker) RemoveFile(fname *string) error {
	path := path.Join(fc.fname,*fname)
	return os.Remove(path)
}
// Функция создаёт файл по запросу
func (fc *FileChecker) CreateFile(fname *string, reader io.ReadCloser) error {
	path := path.Join(fc.fname,*fname)
	file, err := os.Create(path)
	if err == nil {
		var buf [MAXBUF]byte
		bufRead := buf[0:]
		n, err := reader.Read(bufRead)
		if err == nil {
			for n>0  {
				buf2 := buf[0:n]
				 _ , err = file.Write(buf2)
				if err != nil {
					log.Println("Ошибка записи: "+err.Error())
					break
				}
				n, err = reader.Read(bufRead)
				if err != nil && err != io.EOF {
					log.Println("Ошибка чтения: "+err.Error())
				}
			}
		}
		reader.Close()
		file.Close()
	}
	return err
}

func (fc *FileChecker) GetFile(fname *string, writer http.ResponseWriter) error {
	path := path.Join(fc.fname,*fname)
	file,err := os.Open(path)
	defer file.Close()
	if err == nil {
		var buf [MAXBUF]byte
		var n int;
		bufRead := buf[0:]
		n, err = file.Read(bufRead)
		if err == nil {		
			 for n > 0 {
			 	_,err = writer.Write(bufRead[0:n])
			 	if err != nil {
			 					 		
			 	 	break	
			 	}
			 	n, err = file.Read(bufRead)			 	
			 	if err != nil {
			 		if err == io.EOF {
						err = nil
					}		
			 		break	 		
			 	}
			 }
		} 
	}
	return err;
}
// Функция получает текущий список файлов по запросу
func (fc *FileChecker) GetFiles() ([]string, error) {
	return getFiles(&fc.fname)
}
func NewFileChecker(aFname *string) FileChecker {
	return FileChecker{fname:*aFname}
}

// Вспомогательная фунекция - пишет лшибки в лог
func errMsg(msg string, err error) {
	log.Println(msg+err.Error());	
}
