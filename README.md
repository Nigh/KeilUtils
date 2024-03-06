# MacroUtil

> previous `KeilUtil`, changed name to `MacroUtil` at v2.0.0

This project started out as a way to use macro paths in Keil projects, since Keil does not support macro paths.
本项目一开始是为了在Keil项目中使用宏路径，因为Keil并不支持宏路径。

This project was created to allow IDEs that do not support macro paths to continue to work in different environments.
为了让这些不支持宏路径的IDE能够在不同的环境下继续运作，这个项目应运而生。

## Usage

- `MacroUtil.exe version` show binary version
- `MacroUtil.exe init` get init info
- `MacroUtil.exe set MARCO_NAME STRING` set a global marco in config
- `MacroUtil.exe list` list all your settings
- `MacroUtil.exe path2macro` replace the STRING to `$(MACRO)` in your `filelist`
- `MacroUtil.exe macro2path` replace the `$(MACRO)` back to STRING in your `filelist`
- `MacroUtil.exe remove MARCO_NAME` delete the global marco set in config
- `MacroUtil.exe add FILE_PATH` add a file to your `filelist`
- `MacroUtil.exe ignore FILE_PATH` remove a file in your `filelist`
- `MacroUtil.exe replace foo bar` replace all the string `foo` to `bar` in your files
