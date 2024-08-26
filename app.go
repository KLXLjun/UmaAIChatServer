package main

import (
	openai "UmaAIChatServer/API/OpenAI"
	config "UmaAIChatServer/Config"
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/otiai10/openaigo"
	"github.com/sirupsen/logrus"

	nested "github.com/antonfisher/nested-logrus-formatter"
)

var reader = bufio.NewReader(os.Stdin)
var Quit = make(chan os.Signal, 1)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		NoColors:        false,
		TimestampFormat: "2006-01-02 15:04:05",
		FieldsOrder:     []string{},
	})

	ex, err := os.Executable()
	if err != nil {
		logrus.Error("An error occurred while attempting to retrieve the execution directory. 获取程序执行目录失败", err)
		logrus.Error("Please press the Enter key to exit. 按下Enter键退出程序")
		_, _ = reader.ReadString('\n')
		os.Exit(0)
	}
	exPath := filepath.Dir(ex)
	logrus.Info("Execution directory:", exPath)

	if ok, loaderr := config.LoadConfig(exPath); !ok {
		logrus.Error("An error occurred while parsing the configuration file. 解析配置出错", loaderr)
		logrus.Error("Please press the Enter key to exit. 按下Enter键退出程序")
		_, _ = reader.ReadString('\n')
		os.Exit(0)
	}

	openai.InitClient(config.Conf.ChatConfig.Proxy, config.Conf.ChatConfig.APIUrl, config.Conf.ChatConfig.APIKey)
	openai.InitEmotionClient(config.Conf.EmotionConfig.Proxy, config.Conf.EmotionConfig.APIUrl, config.Conf.EmotionConfig.APIKey)

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logrus.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")
			return nil
		},
	}))

	e.Use(middleware.Recover())

	go openai.ChatTask()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "The service is ready.")
	})

	g := e.Group("/v1")

	g.POST("/setSystemPrompt", func(c echo.Context) error {
		prompt := c.FormValue("prompt")
		openai.SetSystemPrompt(prompt)
		openai.ClearToken()
		return c.String(http.StatusOK, "ok")
	})

	g.POST("/clear", func(c echo.Context) error {
		openai.ClearToken()
		return c.String(http.StatusOK, "ok")
	})

	g.POST("/chat", func(c echo.Context) error {
		msg := c.FormValue("prompt")
		emotion := c.FormValue("emotion")

		channel := make(chan openai.ChatResult, 1)

		openai.AddTask(openai.QueuePrompt{
			Emotion:  emotion,
			CallBack: channel,
			PromptGroup: openaigo.Message{
				Role:    openai.ChatMessageRoleUser,
				Content: msg,
			},
		})

		i, ok := <-channel
		if ok {
			defer close(channel)
			return c.JSON(http.StatusOK, i)
		} else {
			return c.JSON(http.StatusOK, "")
		}
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", config.Conf.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
