@echo off 

set FileName=server

if not exist .\build (
    mkdir .\build
)

if exist .\build\%FileName%.exe (
    del .\build\%FileName%.exe
)

set GOOS=windows
set GOARCH=amd64

@echo on
go build -o .\build\%FileName%.exe

@echo off
if exist .\build\%FileName%.exe (
    pushd .\build
    .\%FileName%.exe
    popd
)