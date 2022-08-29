package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Result struct {
	Code    int         `json:"rc"`
	Message string      `json:"rd"`
	Data    interface{} `json:"data"`
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func LoggerToFile() gin.HandlerFunc {
	logFilePath := viper.GetString("logger.path")
	logFileName := viper.GetString("logger.name")

	fileName := path.Join(logFilePath, logFileName)

	_, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("err", err)
	}

	logger := logrus.New()

	logWriter, err := rotatelogs.New(
		fileName+".%Y%m%d.log",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Duration(viper.GetInt("logger.max_age")*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(viper.GetInt("logger.rotation_time")*24)*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		// DisableColors:   true,
		// PrettyPrint: true,
	})

	logger.AddHook(lfHook)

	return func(c *gin.Context) {

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		startTime := time.Now()

		c.Next()

		responseBody := bodyLogWriter.body.String()

		var responseCode string
		var responseMsg string
		var responseData interface{}

		if responseBody != "" {
			res := Result{}
			err := json.Unmarshal([]byte(responseBody), &res)
			if err == nil {
				responseCode = strconv.Itoa(res.Code)
				responseMsg = res.Message
				responseData = res.Data
			}
		}

		endTime := time.Now()

		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}

		logger.WithFields(logrus.Fields{
			"res_time":   endTime.Sub(startTime),
			"__method":   c.Request.Method,
			"__uri":      c.Request.RequestURI,
			"_proto":     c.Request.Proto,
			"_useragent": c.Request.UserAgent(),
			"_req":       c.Request.PostForm.Encode(),
			"_ip":        c.ClientIP(),

			"res_status_code": c.Writer.Status(),
			"res_rc":          responseCode,
			"res_rd":          responseMsg,
			"res_data":        responseData,
		}).Info("HTTP")
	}
}

func Log(w string) func(string, interface{}) {

	// only debug devel
	if strings.ToLower(viper.GetString("ENVIRONMENT")) == "production" {
		return func(s string, r interface{}) {}
	}

	log := logrus.New()
	log.Level = logrus.TraceLevel
	log.Out = os.Stdout

	logFilePath := viper.GetString("logger.path")
	logFileName := viper.GetString("logger.name")

	fileName := path.Join(logFilePath, logFileName)

	_, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		// log.Info("Failed to log to file, using default stderr")
	}

	logWriter, err := rotatelogs.New(
		fileName+".%Y%m%d.log",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Duration(viper.GetInt("logger.max_age")*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(viper.GetInt("logger.rotation_time")*24)*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		// FullTimestamp:   true,
		// DisableColors:   true,
		PrettyPrint: true,
	})

	log.AddHook(lfHook)

	if strings.ToLower(w) == "info" {
		return func(message string, data interface{}) {
			log.WithFields(logrus.Fields{
				"log": data,
			}).Info(message)
		}
	} else if strings.ToLower(w) == "debug" {
		return func(message string, data interface{}) {
			log.WithFields(logrus.Fields{
				"log": data,
			}).Debug(message)
		}
	} else if strings.ToLower(w) == "warning" {
		return func(message string, data interface{}) {
			log.WithFields(logrus.Fields{
				"log": data,
			}).Warning(message)
		}
	} else if strings.ToLower(w) == "error" {
		return func(message string, data interface{}) {
			log.WithFields(logrus.Fields{
				"log": data,
			}).Error(message)
		}
	} else {
		return func(message string, data interface{}) {
			log.Info("type tidak tersedia")
		}
	}
}

func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
