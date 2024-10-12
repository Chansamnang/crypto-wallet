package lumberjack

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logMap sync.Map

func InitLumberjackLogger() {
	err := zap.RegisterSink("lumberjack", func(u *url.URL) (zap.Sink, error) {
		// log.Printf("filename: %s %#v \n", u.Path, u)
		if last, ok := logMap.Load(u.Path); ok {
			if sink, ok := last.(zap.Sink); ok {
				return sink, nil
			}
		}

		//log.Println("filename: 2 ", u.Path, u.String())
		path := u.Path
		switch u.Host {
		case ".":
			path = "." + path
		case "..":
			path = ".." + path
		}
		dir := filepath.Dir(path)
		log.Println("filename:", path, dir)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
		// lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
		lumberJackLogger := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    300, // megabytes
			MaxBackups: 30,
			MaxAge:     30, // days
			LocalTime:  true,
		}
		var sink zap.Sink = &writerWrapper{lumberJackLogger}

		logMap.Store(u.Path, sink)

		return sink, nil
	})
	if err != nil {
		return
	}
}

type writerWrapper struct {
	*lumberjack.Logger
}

func (w writerWrapper) Sync() error {
	return nil
}

func (w writerWrapper) Close() error {
	logMap.Delete(w.Filename)
	return w.Logger.Close()
}
