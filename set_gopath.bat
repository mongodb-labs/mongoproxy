@echo off
set PKG=%cd%\.gopath\src\github.com\mongodbinc-interns\mongoproxy
for %%t in (bsonutil, buffer, convert, log, main, messages, mock, modules, server, tests) do echo d | xcopy %cd%\%%t %PKG%\%%t /Y /E /S
REM copy vendored libraries to GOPATH
for /f %%v in ('dir /b /a:d "%cd%\vendor\src\*"') do echo d | xcopy %cd%\vendor\src\%%v %cd%\.gopath\src\%%v /Y /E /S
set GOPATH=%cd%\.gopath;%cd%\vendor
