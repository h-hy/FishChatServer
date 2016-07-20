@echo off

@title RDAWatchServer

@for /f %%i in ('cd') do set PWD=%%i


set ABC=%PWD%

c:\Windows\system32\cmd.exe /c ""C:\Program Files\Git\git-bash.exe" --cd=%PWD%"

