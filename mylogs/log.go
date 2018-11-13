package mylogs

import (
	"github.com/op/go-logging"
	"os"
	"time"
	"io/ioutil"
	"strings"
	"fmt"
)

var (
	Log = logging.MustGetLogger("mylog")
	logf *os.File
	lognameLeft string
	runMode string
	logdir string
)

const (
	consoleFormat = `%{color}%{time:15:04:05.000} [%{callpath}] [%{shortfile}] ▶ %{level:.4s}%{color:reset} %{message}`
	fileFormat = `%{time:15:04:05.000} [%{callpath}] [%{shortfile}] ▶ %{level:.4s}  %{message}`

)

func LogInit(dir, name, mode string) (err error) {
	lognameLeft = name
	runMode = mode
	logdir = dir

	consolelog := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stdout, "", 0),
		logging.MustStringFormatter(consoleFormat))
	backlogs := []logging.Backend{consolelog}
	if runMode == "prod" || runMode == "test" {
		logName := fmt.Sprintf("%s%s.log", lognameLeft, time.Now().Format("2006-01-02"))
		logF, err := os.OpenFile(logdir + "/" + logName, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Println("日志文件打开失败！ err: ", err.Error())
				return err
			}
			logF, err = os.OpenFile(logdir + "/" + logName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("日志文件打开失败！ err: ", err.Error())
				return err
			}
		}
		logf = logF

		filelog := logging.AddModuleLevel(
			logging.NewBackendFormatter(
				logging.NewLogBackend(logF, "", 0),
				logging.MustStringFormatter(fileFormat)))
		filelog.SetLevel(logging.INFO, "")
		backlogs = append(backlogs, filelog)

	}

	logging.SetBackend(backlogs...)
	return nil
}

//按天数生成日志文件
func NewLog() {


	if runMode == "prod" || runMode == "test" {
		consolelog := logging.NewBackendFormatter(
			logging.NewLogBackend(os.Stdout, "", 0),
			logging.MustStringFormatter(consoleFormat))
		backlogs := []logging.Backend{consolelog}
		logName := fmt.Sprintf("%s%s.log", lognameLeft, time.Now().Format("2006-01-02"))
		logF, err := os.OpenFile(logdir + "/" + logName, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Println("日志文件打开失败！ err: ", err.Error())
				return
			}
			logF, err = os.OpenFile(logdir + "/" + logName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("日志文件打开失败！ err: ", err.Error())
				return
			}
		}

		logf.Close()
		logf = logF

		filelog := logging.AddModuleLevel(
			logging.NewBackendFormatter(
				logging.NewLogBackend(logF, "", 0),
				logging.MustStringFormatter(fileFormat)))
		filelog.SetLevel(logging.INFO, "")
		backlogs = append(backlogs, filelog)

		logging.SetBackend(backlogs...)
	}
}

//清除历史日志文件
func CleanHistoryLog(leftday int) {
	if leftday < 1 {
		return
	}
	dir, err := ioutil.ReadDir(logdir)
	if err != nil {
		return
	}
	PthSep := string(os.PathSeparator)
	date := time.Now().Add(-((time.Duration(leftday) * 24 * time.Hour)))
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasPrefix(fi.Name(), lognameLeft) {
			if tmpt, err := time.ParseInLocation("2006-01-02", fi.Name()[len(lognameLeft):len(lognameLeft)+10], time.Local); err == nil {
				if date.After(tmpt) {
					os.Remove(logdir + PthSep + fi.Name())
				}
			}
		}
	}
}