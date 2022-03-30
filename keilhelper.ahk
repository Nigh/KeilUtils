
#SingleInstance Ignore
SetWorkingDir(A_ScriptDir)

keilProjectPath:=getKeilProject()
if(keilProjectPath==-1) {
	ExitApp -1
}
oriProject := FileRead(keilProjectPath, "UTF-8-RAW")
state := IniRead("config.ini", "setup", "state", 0)

if A_Args.Length = 0
{
	if state = 0 {
		firstSetup()
		ExitApp -1
	} else if state = 1 {
		ExitApp 0
	} else if state = 2 {
		afterpush()
		ExitApp -2
	}
} else {
	if state = 0 {
		firstSetup()
	}
	for n, param in A_Args
	{
		if(param=="--change") {
			changePath()
		}
		if(param=="--prepush") {
			prepush()
		}
		if(param=="--afterpush") {
			afterpush()
		}
	}
}

ExitApp 0

stateTo(n)
{
	IniWrite n, "config.ini", "setup", "state"
}

getKeilProject()
{
	Loop Files, ".\*.uvprojx" ,"F"
	{
		Return A_LoopFilePath
	}
	Return -1
}

keilProjectUpdate(newProject)
{
	global keilProjectPath
	FileObj := FileOpen(keilProjectPath, 0x201, "UTF-8-RAW")
	FileObj.Write(newProject)
	FileObj.Close()
}

firstSetup()
{
	global
	isSetup := IniRead("config.ini", "setup", "issetup", "")
	if(isSetup=="Y29uZmlnZWQ=") {
		ExitApp 0
	}
	While True
	{
		SDK_path := DirSelect("", 0, "请选择SDK的路径")
		SDK_path := RegExReplace(SDK_path, "\\$")
		if SDK_path = "" {
			ExitApp -1
		}
		if not FileExist(SDK_path "\middlewares\Nationstech\ble_library\ns_ble_stack\arch\compiler.h") {
			MsgBox("无效的SDK路径！请重新选择")
		} else {
			Break
		}
	}
	keilProjectUpdate(StrReplace(oriProject, "$(SDK_PATH)", SDK_path))
	IniWrite SDK_path, "config.ini", "setup", "sdk_path"
	IniWrite "Y29uZmlnZWQ=", "config.ini", "setup", "issetup"
	stateTo(1)
}

changePath()
{
	global
	isSetup := IniRead("config.ini", "setup", "issetup", "")
	if(isSetup!="Y29uZmlnZWQ=") {
		ExitApp -2
	}
	ori_SDK_path := IniRead("config.ini", "setup", "sdk_path", "")
	While True
	{
		SDK_path := DirSelect("", 0, "请选择新的SDK的路径")
		SDK_path := RegExReplace(SDK_path, "\\$")
		if SDK_path = "" {
			ExitApp -1
		}
		if not FileExist(SDK_path "\middlewares\Nationstech\ble_library\ns_ble_stack\arch\compiler.h") {
			MsgBox("无效的SDK路径！请重新选择")
		} else {
			Break
		}
	}
	if(state = 2) {
		afterpush()
	}
	keilProjectUpdate(StrReplace(oriProject, ori_SDK_path, SDK_path))
	IniWrite SDK_path, "config.ini", "setup", "sdk_path"
}

prepush()
{
	global
	ori_SDK_path := IniRead("config.ini", "setup", "sdk_path", "bm90IHNldA==")
	keilProjectUpdate(StrReplace(oriProject, ori_SDK_path, "$(SDK_PATH)"))
	stateTo(2)
}

afterpush()
{
	global
	SDK_path := IniRead("config.ini", "setup", "sdk_path", "bm90IHNldA==")
	keilProjectUpdate(StrReplace(oriProject, "$(SDK_PATH)", SDK_path))
	stateTo(1)
}
