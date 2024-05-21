package pkg

import (
	"os"

	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

func InitLog() {
	logLevel := getLoggerLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logLevel)

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		DisableColors:   true,
	})

	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    20, // megabytes
		MaxBackups: 14,
		MaxAge:     14, // days
		Compress:   true,
	})
}

func getLoggerLevel(level string) log.Level {
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

func LogRequest(c *gin.Context) {
	nodeEnv := os.Getenv("ENV")
	if nodeEnv == constants.DEV {
		// Print console log
		log.SetOutput(os.Stdout)
	}

	requestId := c.GetHeader("X-Request-ID")
	if requestId == "" {
		requestId = uuid.New().String()
	}
	c.Set("requestId", requestId)
	clientIP := c.ClientIP()

	log.WithFields(log.Fields{
		"requestId": requestId,
		"clientIP":  clientIP,
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
		"params":    c.Request.URL.Query(),
	}).Info("Incoming request")

	c.Next()
}
