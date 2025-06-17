param([String]$type="wrong_parameter", [String]$package="...", [String]$test=".")

if($type -ne "ut" -and $type -ne "clean-run" -and $type -ne "all"){
    "test.ps1 -type ut | will run unit tests"
    "test.ps1 -type ut -package package | will run UTies only for specific package"
    "test.ps1 -type ut -test regex | will run test which matches given regex"
    "test.ps1 -type clean-run | will clear data then start application"
    exit 1
}

$Env:PGDB = "postgres"
$Env:PGIP = "localhost"
$Env:PGUSER = "postgres"
$Env:PGPORT = 5432
$Env:PGPASSWORD = "postgres"

function UnitTests(){
    psql --version
    if($LASTEXITCODE -ne 0){
        "psql not installed"
        exit 1
    }
    
    psql -U $Env:PGUSER -d $Env:PGDB -a -f ./test/init.sql
    if($LASTEXITCODE -ne 0){
        "failed to initialize database"
        exit 1
    }
    
    go clean -testcache
    if($LASTEXITCODE -ne 0){
        "failed to clean cache"
        exit 1
    }
    
    Clear-Host
    "please wait ..."
    go test -v ./src/$package -run $test
    if($LASTEXITCODE -ne 0){
        "uties failed"
        exit 1
    }
}

function CleanRun(){    
    psql -U $Env:PGUSER -d $Env:PGDB -a -f ./test/init.sql
    if($LASTEXITCODE -ne 0){
        "failed to initialize database"
        exit 1
    }
    Clear-Host

    go run .\main.go
}

if($type -eq "ut" -or $type -eq "all"){
    "Starting unit tests"
    UnitTests
}

if($type -eq "clean-run" -or $type -eq "all"){
    "Clean run"
    CleanRun
}
