package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var marcoMap map[string]string
var projectFile string

func getConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	projectFile = viper.GetString("project")
	marcoMap = viper.GetStringMapString("macro")
}

func init() {
	rand.Seed(time.Now().UnixNano())
	// viper.SetDefault("state", "0")
	viper.SetDefault("macro", map[string]string{})
	viper.SetDefault("project", "")
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SafeWriteConfigAs("config.json")
}

func projectReplace(fn func(string) string) {
	project, err := os.ReadFile(projectFile)
	if err != nil {
		panic(fmt.Errorf("fatal error read project file: %s", projectFile))
	}
	projectStr := string(project)
	projectStr = fn(projectStr)
	os.WriteFile(projectFile, []byte(projectStr), 0644)
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		// TODO: help info
		return
	}
	if len(flag.Args()) >= 1 {
		cmd := flag.Arg(0)
		if cmd != "init" {
			getConfig()
		}
		switch cmd {
		case "init":
			filepath.Walk(".",
				func(path string, f os.FileInfo, err error) error {
					if !f.IsDir() {
						ext := filepath.Ext(path)
						if ext == ".uvprojx" {
							fmt.Println("found project", f.Name())
							viper.Set("project", f.Name())
							viper.WriteConfig()
							return errors.New("found project")
						}
					}
					return nil
				})
		case "list":
			fmt.Println("project", projectFile)
			for k, v := range marcoMap {
				k = strings.ToUpper(k)
				fmt.Println("$(" + k + ")=" + v)
			}
		case "path2macro":
			projectReplace(func(in string) string {
				for k, v := range marcoMap {
					k = strings.ToUpper(k)
					fmt.Println("replace", v, "to", "$("+k+")")
					in = strings.Replace(in, v, "$("+k+")", -1)
				}
				return in
			})
		case "macro2path":
			projectReplace(func(in string) string {
				for k, v := range marcoMap {
					k = strings.ToUpper(k)
					fmt.Println("replace", "$("+k+")", "to", v)
					in = strings.Replace(in, "$("+k+")", v, -1)
				}
				return in
			})
		case "set":
			if len(flag.Args()) >= 3 {
				marco := flag.Arg(1)
				path := flag.Arg(2)
				fmt.Println("$(" + marco + ")=" + path)
				marcoMap[marco] = path
				viper.Set("macro", marcoMap)
				viper.WriteConfig()
			}
		case "remove":
			if len(flag.Args()) >= 2 {
				for i := 0; i < len(flag.Args())-1; i++ {
					k := strings.ToUpper(flag.Arg(i + 1))
					fmt.Println("try to delete macro", "$("+k+")")
					delete(marcoMap, k)
				}
				viper.Set("macro", marcoMap)
				viper.WriteConfig()
			}
		case "replace":
			if len(flag.Args()) >= 3 {
				from := flag.Arg(1)
				to := flag.Arg(2)
				projectReplace(func(in string) string {
					fmt.Println("replace", from, "to", to)
					in = strings.Replace(in, from, to, -1)
					return in
				})
			}
		}
	}
}
