# KeilUtils

> Use macro path in Keil project

Because `keil` does not support macro paths, we will encounter path problems when opening `keil` projects in different environments.
And paths that cannot be represented by variables can cause a lot of trouble when people collaborate using version control tools like git.
To solve this problem, I have developed this tool.

## Usage

**At first, put `KeilUtils.exe` in the same directory with `uvprojx` file**

- `KeilUtils.exe version` to get binary version
- `KeilUtils.exe init` to generate init config
- `KeilUtils.exe set MARCO_NAME STRING` to set a string marco in config
- `KeilUtils.exe list` to check your settings
- `KeilUtils.exe path2macro` using config replace the STRING to `$(MACRO)`
- `KeilUtils.exe macro2path` using config replace the `$(MACRO)` back to STRING
- `KeilUtils.exe remove MARCO_NAME` to delete the marco set before
- `KeilUtils.exe replace foo bar` to replace all the string `foo` to `bar` in your project
