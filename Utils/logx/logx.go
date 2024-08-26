package logx

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

func runfuncname() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
func Error(tag string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"tag": tag, "func": runfuncname()}).Errorln(args)
}
func Warn(tag string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"tag": tag, "func": runfuncname()}).Warnln(args)
}
func Info(tag string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"tag": tag, "func": runfuncname()}).Infoln(args)
}
func Debug(tag string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"tag": tag, "func": runfuncname()}).Debugln(args)
}
func Trace(tag string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"tag": tag, "func": runfuncname()}).Traceln(args)
}

// func Error(tag string,funcx string,args ...interface{}){
// 	logrus.WithFields(logrus.Fields{"tag": tag,"func":funcx}).Errorln(args)
// }
// func Warn(tag string,funcx string,args ...interface{}) {
// 	logrus.WithFields(logrus.Fields{"tag": tag,"func":funcx}).Warnln(args)
// }
// func Info(tag string,funcx string,args ...interface{}) {
// 	logrus.WithFields(logrus.Fields{"tag": tag,"func":funcx}).Infoln(args)
// }
// func Debug(tag string,funcx string,args ...interface{}) {
// 	logrus.WithFields(logrus.Fields{"tag": tag,"func":funcx}).Debugln(args)
// }
// func Trace(tag string,funcx string,args ...interface{}){
// 	logrus.WithFields(logrus.Fields{"tag": tag,"func":funcx}).Traceln(args)
// }
