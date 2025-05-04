# Unit tests
$Env:PGDB = "postgres"
$Env:PGIP = "localhost"
$Env:PGUSER = "postgres"
$Env:PGPORT = 5432
$Env:PGPASSWORD = "postgres"

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

# Functional tests

python -m venv .env
if($LASTEXITCODE -ne 0){
    "Failed to create venv"
    exit 1
}

.env\Scripts\activate.ps1
if($LASTEXITCODE -ne 0){
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
