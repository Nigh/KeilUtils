
SetWorkingDir(A_ScriptDir)

keilProjectPath:=getKeilProject()
if(keilProjectPath==-1) {
	ExitApp -1
}
oriProject := FileRead(keilProjectPath, "UTF-8")

if A_Args.Length = 0
{
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
	newProject := StrReplace(oriProject, "$(SDK_PATH)", SDK_path)
	FileObj := FileOpen(keilProjectPath, 0x201, "UTF-8")
	FileObj.Write(newProject)
	FileObj.Close()
	IniWrite SDK_path, "config.ini", "setup", "sdk_path"
	IniWrite "Y29uZmlnZWQ=", "config.ini", "setup", "issetup"
} else {
	for n, param in A_Args
	{
		if(param=="--change") {
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
			newProject := StrReplace(oriProject, ori_SDK_path, SDK_path)
			FileObj := FileOpen(keilProjectPath, 0x201, "UTF-8")
			FileObj.Write(newProject)
			FileObj.Close()
			IniWrite SDK_path, "config.ini", "setup", "sdk_path"
		}
		if(param=="--prepush") {
			ori_SDK_path := IniRead("config.ini", "setup", "sdk_path", "bm90IHNldA==")
			newProject := StrReplace(oriProject, ori_SDK_path, "$(SDK_PATH)")
			FileObj := FileOpen(keilProjectPath, 0x201, "UTF-8")
			FileObj.Write(newProject)
			FileObj.Close()
		}
		if(param=="--afterpush") {
			SDK_path := IniRead("config.ini", "setup", "sdk_path", "bm90IHNldA==")
			newProject := StrReplace(oriProject, "$(SDK_PATH)", SDK_path)
			FileObj := FileOpen(keilProjectPath, 0x201, "UTF-8")
			FileObj.Write(newProject)
			FileObj.Close()
		}
	}
}

ExitApp 0

getKeilProject()
{
	Loop Files, ".\*.uvprojx" ,"F"
	{
		Return A_LoopFilePath
	}
	Return -1
}
