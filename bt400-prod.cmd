@echo off
setlocal enabledelayedexpansion

set ymd=%DATE:/=%

set JOBNAME=batch-bt400-%ymd%
set JOBQUEUE=batch-queue
set JOBDEF=batch-bt400:1
set HOST=db.xxxxxxxxxxx.ap-northeast-1.rds.amazonaws.com
set PORT=5432
set DB=database
set USER=batchdb
set PASS=
set UPDUSER=upd@domain
set S3=s3://dittos/sample/%ymd%
::% must write as %% in string

set /p indb="input database name:(%DB%)"
set /p inpass="input password:(%PASS%)"
set /p ins3="input output s3 path:(%S3%)"
if [%indb%] neq [] (
DB=indb
)
if [%inpass%] neq [] (
PASS=inpass
)
if [%ins3%] neq [] (
S3=ins3
)
echo dump the database %DB% to %S3%
set /p yes= "are you sure?(Y/n):"
if [%yes%] neq [] (
if /i [%yes%] neq [y] (
goto :EOF
)
)

submit.exe --job-name="%JOBNAME%" --job-queue="%JOBQUEUE%" --job-definition="%JOBDEF%" --parameters="--host=%HOST%,--port=%PORT%,--database=%DB%,--user=%USER%,--password=%PASS%,--extra=?,--filename=?,--s3=%S3%,--updateUser=%UPDUSER%" --wait
:: %errorlevel%

endlocal

exit /b 0
