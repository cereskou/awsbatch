@echo off

set GOARCH=amd64
set GOOS=windows
set CGO_ENABLED=0

set MODULE=submit.exe

echo clean old module ...
go clean
if exist resource.syso (
    del resource.syso
)

if not exist go.mod (
    echo golang mod init ...
    go mod init
)

echo generate resource ...
go generate

set filetime=
if exist created.txt (
set /p filetime=< created.txt
)

echo build ...
::go build -ldflags "-s -w -extldflags -static -X 'main.version=%version%'" -a -i .
go build -ldflags "-s -w -extldflags -static" -a -i -o %MODULE% .
if %errorlevel% equ 0 (
        if "%filetime%"=="" (
            echo done.
        ) else (
            echo set LastWriteTime to %filetime% ...
            powershell Set-ItemProperty "%~dp0\%MODULE%" -Name LastWriteTime -Value '%filetime%'
            echo set CreationTime  to %filetime% ...
            powershell Set-ItemProperty "%~dp0\%MODULE%" -Name CreationTime  -Value '%filetime%' 
            echo done.
        )
) else (
    echo failed.
)
