# Keil_MacroPath_Helper
Use macro path in Keil project

#### Works with `Autohotkey Version 2.0-beta.3`

## Usage

1. Replace SDK path to `$(SDK_PATH)` in `uvprojx` file
2. Put `KeilHelper.exe` in the same directory with `uvprojx` file
3. Run `KeilHelper.exe` to choose the path of your SDK
4. Run `KeilHelper.exe --change` to change the path of your SDK
5. Run `KeilHelper.exe --prepush` to replace the SDK path back to `$(SDK_PATH)` macro
6. Run `KeilHelper.exe --afterpush` to replace the `$(SDK_PATH)` macro back to your SDK path

