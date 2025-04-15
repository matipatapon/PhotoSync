SET PGDB=postgres
SET PGIP=localhost
SET PGUSER=postgres
SET PGPORT=5432
SET PGPASSWORD=postgres

psql --version
if %errorlevel% neq 0 exit /b %errorlevel%

psql -U %PGUSER% -d %PGDB% -a -f ./test/init.sql
if %errorlevel% neq 0 exit /b %errorlevel%

go clean -testcache
go test -v ./...
