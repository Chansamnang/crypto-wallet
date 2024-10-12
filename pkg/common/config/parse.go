package config

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"wallet/pkg/constant"

	"gopkg.in/yaml.v2"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project.
	Root = filepath.Join(filepath.Dir(b), "../..")
)

func readConfig(configFile string) ([]byte, error) {
	b, err := os.ReadFile(configFile)
	if err != nil { // File exists and was read successfully
		return nil, err
	}
	return b, nil
}

func InitConfig(configFile string) error {
	data, err := readConfig(configFile)
	if err != nil {
		return fmt.Errorf("read loacl config file error: %w", err)
	}

	if err = yaml.NewDecoder(bytes.NewReader(data)).Decode(&Config); err != nil {
		return fmt.Errorf("parse loacl  file error: %w", err)
	}

	// configData, err := yaml.Marshal(&Config)
	// fmt.Printf("debug: %s\nconfig:\n%s\n", time.Now(), string(configData))

	// timezone
	initTimezone()

	return nil
}

func configFieldCopy[T any](local **T, remote T) {
	if *local == nil {
		*local = &remote
	}
}

type zkLogger struct{}

func (l *zkLogger) Printf(format string, a ...interface{}) {
	fmt.Printf("zk get config %s\n", fmt.Sprintf(format, a...))
}

func checkFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func findConfigFile(paths []string) (string, error) {
	for _, path := range paths {
		if checkFileExists(path) {
			return path, nil
		}
	}
	return "", fmt.Errorf("configPath not found")
}

func CreateCatalogPath(path string) []string {

	path1 := filepath.Dir(path)
	path1 = filepath.Dir(path1)
	// the parent of  binary file
	pa1 := filepath.Join(path1, constant.ConfigPath)
	path2 := filepath.Dir(path1)
	path2 = filepath.Dir(path2)
	path2 = filepath.Dir(path2)
	// the parent is _output
	pa2 := filepath.Join(path2, constant.ConfigPath)
	path3 := filepath.Dir(path2)
	// the parent is project(default)
	pa3 := filepath.Join(path3, constant.ConfigPath)

	return []string{pa1, pa2, pa3}

}

func findConfigPath(configFile string) (string, error) {
	path := make([]string, 10)

	// First, check the configFile argument
	if configFile != "" {
		if _, err := findConfigFile([]string{configFile}); err != nil {
			return "", errors.New("the configFile argument path is error")
		}
		fmt.Println("configfile:", configFile)
		return configFile, nil
	}

	p1, err := os.Executable()
	if err != nil {
		return "", err
	}

	path = CreateCatalogPath(p1)
	pathFind, err := findConfigFile(path)
	if err == nil {
		return pathFind, nil
	}

	// Forth, use the Default path.
	return "", errors.New("the config.yaml path not found")
}

func FlagParse(env string) (configFile string, logFileName string, ginPort int, hide bool, showVersion bool, err error) {
	flag.StringVar(&configFile, "config", "./config/config.yaml", "Config full path")

	flag.IntVar(&ginPort, "port", 10009, "get ginServerPort from cmd")

	flag.BoolVar(&hide, "hide", false, "hide the ComponentCheck result")

	// Version flag
	flag.BoolVar(&showVersion, "version", false, "show version and exit")

	// log file
	logFileName = fmt.Sprintf("lumberjack://./logs/%s.log", env)
	flag.StringVar(&logFileName, "log_file", logFileName, "log file name")

	flag.Parse()

	configFile, err = findConfigPath(configFile)
	return
}

func getArrPointEnv(key1, key2 string, fallback *[]string) *[]string {
	str1 := getEnv(key1, "")
	str2 := getEnv(key2, "")
	str := fmt.Sprintf("%s:%s", str1, str2)
	if len(str) <= 1 {
		return fallback
	}
	return &[]string{str}
}

func getStringEnv(key1, key2 string, fallback string) string {
	str1 := getEnv(key1, "")
	str2 := getEnv(key2, "")
	str := fmt.Sprintf("%s:%s", str1, str2)
	if len(str) <= 2 {
		return fallback
	}
	return str
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvStringPoint(key string, fallback *string) *string {
	if value, exists := os.LookupEnv(key); exists {
		return &value
	}
	return fallback
}

func getEnvIntPoint(key string, fallback *int64) (*int64, error) {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		temp := int64(val)
		if err != nil {
			return nil, err
		}
		return &temp, nil
	}
	return fallback, nil
}
