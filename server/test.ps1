param([String]$type="wrong_parameter", [String]$package="...", [String]$test=".")

if($type -ne "ut" -and $type -ne "ft" -and $type -ne "all"){
    "test.ps1 -type ut | will run unit tests"
    "test.ps1 -type ut -package package | will run UTies only for specific package"
    "test.ps1 -type ut -test regex | will run test which matches given regex"
    "test.ps1 -type ft | will start application with additional endpoints for functional testing"
    exit 1
}

$Env:PGDB = "postgres"
$Env:PGIP = "localhost"
$Env:PGUSER = "postgres"
$Env:PGPORT = 5432
$Env:PGPASSWORD = "postgres"
$Env:TLS_ENABLED = "false"
$Env:CERT_PATH = ""
$Env:CERT_PRIVATE_KEY_PATH = ""
$Env:ALLOWED_ORIGIN = "*"

function UnitTests(){    
    go clean -testcache
    if($LASTEXITCODE -ne 0){
        "failed to clean cache"
        exit 1
    }
    
    go test -v ./src/$package -run $test
    if($LASTEXITCODE -ne 0){
        "uties failed"
        exit 1
    }
}

function FunctionalTests(){
    go run .\main.go --testing
}

if($type -eq "ut" -or $type -eq "all"){
    "Starting unit tests"
    UnitTests
}

if($type -eq "ft" -or $type -eq "all"){
    "Starting app for functional testing"
    FunctionalTests
}
