@echo off
REM Batch script for processing a large number of German words
REM Usage: batch_words.bat [input_file] [output_folder] [workers] [timeout]

setlocal enabledelayedexpansion

REM Default values
set INPUT_FILE=%1
set OUTPUT_FOLDER=%2
if "%OUTPUT_FOLDER%"=="" set OUTPUT_FOLDER=output
set WORKERS=%3
if "%WORKERS%"=="" set WORKERS=5
set TIMEOUT=%4
if "%TIMEOUT%"=="" set TIMEOUT=30

REM Check if input file is provided
if "%INPUT_FILE%"=="" (
    echo Error: Input file not specified.
    echo Usage: batch_words.bat [input_file] [output_folder] [workers] [timeout]
    exit /b 1
)

REM Check if input file exists
if not exist "%INPUT_FILE%" (
    echo Error: Input file %INPUT_FILE% not found.
    exit /b 1
)

REM Create output folder if it doesn't exist
if not exist "%OUTPUT_FOLDER%" mkdir "%OUTPUT_FOLDER%"

echo Processing words from %INPUT_FILE% with %WORKERS% workers and %TIMEOUT% seconds timeout...
echo Results will be saved to %OUTPUT_FOLDER%

REM Process words in batches of 10
set BATCH_SIZE=10
set COUNTER=0
set BATCH_WORDS=

for /f "tokens=*" %%w in (%INPUT_FILE%) do (
    set WORD=%%w
    set WORD=!WORD: =!
    if not "!WORD!"=="" (
        set BATCH_WORDS=!BATCH_WORDS! !WORD!
        set /a COUNTER+=1
        
        if !COUNTER! GEQ %BATCH_SIZE% (
            echo Processing batch: !BATCH_WORDS!
            go run main.go batch -w %WORKERS% -t %TIMEOUT% -p "%OUTPUT_FOLDER%/" !BATCH_WORDS!
            set COUNTER=0
            set BATCH_WORDS=
            timeout /t 2 > nul
        )
    )
)

REM Process remaining words
if not "!BATCH_WORDS!"=="" (
    echo Processing final batch: !BATCH_WORDS!
    go run main.go batch -w %WORKERS% -t %TIMEOUT% -p "%OUTPUT_FOLDER%/" !BATCH_WORDS!
)

echo All words processed. Results saved to %OUTPUT_FOLDER% 