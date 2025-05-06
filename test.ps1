param([String]$type="wrong_parameter", [String]$recreate="false")

if($type -ne "ut" -and $type -ne "ft" -and $type -ne "all"){
    "test.ps1 -type ut | will run unit tests"
    "test.ps1 -type ft | will run functional tests"
    "test.ps1 -type ft -recreate true | will recreate virtual env then run functional tests"
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
    
    go test -v ./...
    if($LASTEXITCODE -ne 0){
        "uties failed"
        exit 1
    }
}

function FunctionalTests(){
    if(!(Test-Path ".\.env") -or $recreate -eq "true"){
        Remove-Item -Recurse -Force ".\.env"
        python -m venv .env
        if($LASTEXITCODE -ne 0){
            "Failed to create venv"
            exit 1
        }
    }
    
    .\.env\Scripts\Activate.ps1
    $PIP_VERSION = pip -V
    if(!($PIP_VERSION -Match "\.env")){
        "Failed to activate venv"
        exit 1
    }

    pip install -r test/requirments.txt
    if($LASTEXITCODE -ne 0){
        "Failed to install requirments"
        exit 1
    }

    pytest
    if($LASTEXITCODE -ne 0){
        "Functional tests failed"
        exit 1
    }
}

if($type -eq "ut" -or $type -eq "all"){
    "Starting unit tests"
    UnitTests
}

if($type -eq "ft" -or $type -eq "all"){
    "Starting functional tests"
    FunctionalTests
}
