
# Keil Macro Path
> Use macro path in Keil project

Because `keil` does not support macro paths, we will encounter path problems when opening `keil` projects in different environments.
And paths that cannot be represented by variables can cause a lot of trouble when people collaborate using version control tools like git.
To solve this problem, I have developed this tool.

## Usage

0. Put `KeilUtils.exe` in the same directory with `uvprojx` file
1. Run `KeilUtils.exe init` to generate init config
2. Run `KeilUtils.exe set ORIGIN_PATH MARCO_NAME` to replace your origin path to marco path
3. Run `KeilUtils.exe list` to check your settings
4. Run `KeilUtils.exe path2macro` to replace the path to `$(MACRO)` macro
5. Run `KeilUtils.exe macro2path` to replace the macro path back to your origin path
6. Run `KeilUtils.exe remove MARCO_NAME` to delete the marco set before
7. Run `KeilUtils.exe replace foo bar` to replace all the string `foo` to `bar` in your project

