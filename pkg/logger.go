package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
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
		log.SetOutput(os.Stdout) // Print console log in development environment
	}

	requestId := c.GetHeader("X-Request-ID")
	if requestId == "" {
		requestId = uuid.New().String()
	}
	c.Set("requestId", requestId)
	clientIP := c.ClientIP()

	var requestBody interface{}
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			requestBody = formatJSON(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body for further processing
		}
	}

	c.Next()

	status := c.Writer.Status()
	var message string
	var stack string

	if len(c.Errors) > 0 {
		err := c.Errors.Last()
		message = err.Err.Error()
		stack = err.Error() // Contains the full stack trace
	}

	log.WithFields(log.Fields{
		"requestId": requestId,
		"clientIP":  clientIP,
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
		"params":    c.Request.URL.Query(),
		"status":    status,
		"body":      requestBody,
		"message":   message,
		"stack":     stack,
	}).Info("Incoming request")
}

func formatJSON(data []byte) interface{} {
	var parsedData map[string]interface{}
	err := json.Unmarshal(data, &parsedData)
	if err != nil {
		return string(data)
	}
	return parsedData
}
