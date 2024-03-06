package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type macro struct {
	Content   string `json:"content"`
	Desc      string `json:"desc"`
	CheckPath string `json:"checkpath"`
}
type macromap map[string]macro

type fileConfig struct {
	Path      string `json:"path"`
	CheckPath string `json:"checkpath"`
}

type config struct {
	Macro macromap     `json:"macro"`
	Files []fileConfig `json:"files"`
}

var cfg config

const version string = "2.0.0"

const macroPrefix = "$("
const macroSuffix = ")"

func fileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	if !fileExist("macros.json") {
		fmt.Println("macros.json" + " not exist, create it")
		os.WriteFile("macros.json", []byte("{}"), 0644)
	}
	if !fileExist("macros-local.json") {
		fmt.Println("macros-local.json" + " not exist, create from macros.json")
		data, _ := os.ReadFile("macros.json")
		os.WriteFile("macros-local.json", data, 0644)
	}

	data, err := os.ReadFile("macros-local.json")
	if err != nil {
		fmt.Println("Read marcos.json error")
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	json.Unmarshal(data, &cfg)
}
func cfgSave() {
	data, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println("Error encoding:", err)
		panic(fmt.Errorf("fatal error marshal config file: %s", err))
	} else {
		os.WriteFile("macros-local.json", data, 0644)
		cfgPush()
	}
}

// 将本地配置清洗后覆盖至公共配置
func cfgPush() {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	err := enc.Encode(cfg)
	if err != nil {
		fmt.Println("Error encoding:", err)
		return
	}
	var globalCfg config
	err = dec.Decode(&globalCfg)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	for k, v := range globalCfg.Macro {
		v.Content = ""
		globalCfg.Macro[k] = v
	}

	data, err := json.Marshal(globalCfg)
	if err != nil {
		fmt.Println("Error encoding:", err)
		panic(fmt.Errorf("fatal error marshal config file: %s", err))
	} else {
		os.WriteFile("macros.json", data, 0644)
	}
}

func macroReplace(fn func(string) string) {
	for _, v := range cfg.Files {
		fmt.Printf("-------------------------\nFile:[%s]\n", v.Path)
		if !fileExist(v.Path) {
			fmt.Println("[ERROR] file not exist")
			os.Exit(2)
		}
		content, err := os.ReadFile(v.Path)
		if err != nil {
			panic(fmt.Errorf("fatal error read file: %s", v))
		}
		contentStr := string(content)
		contentStr = fn(contentStr)
		os.WriteFile(v.Path, []byte(contentStr), 0644)
	}
	fmt.Println("-------------------------")
}
func MACROSTR(s string) string {
	return macroPrefix + strings.ToUpper(s) + macroSuffix
}
func cfgCheck() (string, bool) {
	result := false
	output := "  macros:\n"
	outputErr := "  errors:\n"
	for k, v := range cfg.Macro {
		if len(v.Content) == 0 {
			outputErr += fmt.Sprintf("  --[ERROR] macro NOT set: %s // %s\n", MACROSTR(k), v.Desc)
			result = true
		} else {
			if len(v.CheckPath) > 0 {
				fullcheckpath := v.Content + "\\" + v.CheckPath
				if _, err := os.Stat(fullcheckpath); os.IsNotExist(err) {
					outputErr += fmt.Sprintf(`    [ERROR] INVALID macro %s = "%s"`, MACROSTR(k), v.Content)
					if len(v.Desc) > 0 {
						outputErr += fmt.Sprintf(` // %s`, v.Desc)
					}
					outputErr += "\n     reason: "
					outputErr += fmt.Sprintln("Checkpath NOT exist:", fullcheckpath)
					result = true
					continue
				}
			}
			output += fmt.Sprintf(`  - %s = "%s"`, MACROSTR(k), v.Content)
			if len(v.Desc) > 0 {
				output += fmt.Sprintf(` // %s`, v.Desc)
			}
			output += fmt.Sprintln("")
		}
	}
	if result {
		output += "\n" + outputErr
		output += fmt.Sprintln("[ERROR] macros have not been correctly set, check macros-local.json")
	}

	return output, result
}
func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Command list:")
		fmt.Println("version - show binary version")
		fmt.Println("push - sync local macro to global config")
		fmt.Println("list - list all your settings")
		fmt.Println("set MARCO_NAME STRING - set a global marco in config")
		fmt.Println("remove MARCO_NAME - delete the global marco set in config")
		fmt.Println("add FILE_PATH - add a file to your filelist")
		fmt.Println("ignore FILE_PATH - remove a file in your filelist")
		fmt.Println("path2macro - replace the macro content to `$(MACRO)` in your filelist")
		fmt.Println("macro2path - replace the `$(MACRO)` to macro content in your filelist")
		fmt.Println("replace FOO BAR - replace \"FOO\" to \"BAR\" in everyfile in filelist")
		return
	}
	if len(flag.Args()) >= 1 {
		cmd := flag.Arg(0)
		cmds := []string{"path2macro", "macro2path"}
		for _, v := range cmds {
			if cmd == v {
				if s, b := cfgCheck(); b {
					fmt.Println(s)
					os.Exit(1)
					return
				}
			}
		}
	}
	if len(flag.Args()) >= 1 {
		cmd := flag.Arg(0)
		switch cmd {
		case "version":
			fmt.Println("version " + version)
		case "push":
			cfgSave()
			fmt.Println("config pushed")
		case "list":
			fmt.Println("filelist:")
			for _, v := range cfg.Files {
				if len(v.Path) > 0 {
					fmt.Println("  - " + v.Path)
				}
			}
			fmt.Println("macrolist:")
			s, _ := cfgCheck()
			fmt.Println(s)
		case "set":
			if len(flag.Args()) >= 3 {
				marco := strings.ToUpper(flag.Arg(1))
				path := flag.Arg(2)

				fmt.Println(MACROSTR(marco) + "=" + path)
				cfg.Macro[marco] = macro{Content: path}
				cfgSave()
			}
		case "remove":
			if len(flag.Args()) >= 2 {
				for i := 0; i < len(flag.Args())-1; i++ {
					k := strings.ToUpper(flag.Arg(i + 1))
					fmt.Println("delete macro", MACROSTR(k))
					delete(cfg.Macro, k)
				}
				cfgSave()
			}
		case "add":
			if len(flag.Args()) >= 2 {
				for i := 0; i < len(flag.Args())-1; i++ {
					f := flag.Arg(i + 1)
					fmt.Println("add file", f)
					cfg.Files = append(cfg.Files, fileConfig{Path: f})
				}
				cfgSave()
			}
		case "ignore":
			if len(flag.Args()) >= 2 {
				for i := 0; i < len(flag.Args())-1; i++ {
					f := flag.Arg(i + 1)
					fmt.Println("remove file", f)
					files := make([]fileConfig, 0)
					for i, v := range cfg.Files {
						if filepath.Base(v.Path) == filepath.Base(f) {
							fmt.Println(v.Path + " removed")
							cfg.Files = append(cfg.Files[:i], cfg.Files[i+1:]...)
						} else {
							files = append(files, v)
						}
					}
					cfg.Files = files
				}
				cfgSave()
			}
		case "path2macro":
			macroReplace(func(in string) string {
				for k, v := range cfg.Macro {
					k = strings.ToUpper(k)
					fmt.Printf(`replace "%s" to "%s"`+"\n", v.Content, MACROSTR(k))
					in = strings.Replace(in, v.Content, MACROSTR(k), -1)
				}
				return in
			})
		case "macro2path":
			macroReplace(func(in string) string {
				for k, v := range cfg.Macro {
					k = strings.ToUpper(k)
					fmt.Printf(`replace "%s" to "%s"`+"\n", MACROSTR(k), v.Content)
					in = strings.Replace(in, MACROSTR(k), v.Content, -1)
				}
				return in
			})
		case "replace":
			if len(flag.Args()) >= 3 {
				from := flag.Arg(1)
				to := flag.Arg(2)
				macroReplace(func(in string) string {
					fmt.Printf(`replace "%s" to "%s"`+"\n", from, to)
					in = strings.Replace(in, from, to, -1)
					return in
				})
			}
		}
	}
}
