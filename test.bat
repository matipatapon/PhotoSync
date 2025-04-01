psql --version
if %errorlevel% neq 0 exit /b %errorlevel%

SET PGPASSWORD=postgres

psql -U postgres -d postgres -a -f setup_test_db.sql