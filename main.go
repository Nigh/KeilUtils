package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

var marcoMap map[string]string
var projectFile string

const version string = "1.1.0"

func getConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	projectFile = viper.GetString("project")
	marcoMap = viper.GetStringMapString("macro")
}

func init() {
	// rand.Seed(time.Now().UnixNano())
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

func projectMacroScan(projectFileName string) []string {
	project, err := os.ReadFile(projectFileName)
	if err != nil {
		panic(fmt.Errorf("fatal error read project file: %s", projectFileName))
	}
	projectStr := string(project)
	r, _ := regexp.Compile(`\$\((\w+?)\)`)
	macros := r.FindAllStringSubmatch(projectStr, -1)
	var macrosSingle []string
	for _, v := range macros {
		var dup bool = false
		for _, b := range macrosSingle {
			if v[1] == b {
				dup = true
				break
			}
		}
		if dup {
			continue
		}
		macrosSingle = append(macrosSingle, v[1])
	}
	return macrosSingle
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Command list:")
		fmt.Println("version - get binary version")
		fmt.Println("init - generate init config")
		fmt.Println("set MARCO_NAME STRING - set a string marco in config")
		fmt.Println("list - check your settings")
		fmt.Println("path2macro - using config replace the STRING to $(MACRO)")
		fmt.Println("macro2path - using config replace the $(MACRO) back to STRING")
		fmt.Println("remove MARCO_NAME - delete the marco set in config")
		fmt.Println("replace FOO BAR - replace all the string FOO to BAR in your project")
		return
	}
	marcoMap = make(map[string]string)
	if len(flag.Args()) >= 1 {
		cmd := flag.Arg(0)
		if cmd == "version" {
			fmt.Println("version " + version)
			return
		}
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
							fmt.Println("found project", "["+f.Name()+"]")
							viper.Set("project", f.Name())
							macros := projectMacroScan(f.Name())
							for _, v := range macros {
								_, ok := marcoMap[v]
								if !ok {
									marcoMap[v] = "$(" + v + ")"
									fmt.Println("found unconfigured macro", "$("+v+")")
								}
							}
							viper.Set("macro", marcoMap)
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
