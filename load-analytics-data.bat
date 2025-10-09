@echo off
echo Loading Analytics Dummy Data...
echo.

REM Set PostgreSQL connection parameters
set PGHOST=localhost
set PGPORT=5432
set PGDATABASE=research_institute
set PGUSER=postgres
set PGPASSWORD=123

REM Run the analytics dummy data script
echo Loading analytics test data...
psql -h %PGHOST% -p %PGPORT% -U %PGUSER% -d %PGDATABASE% -f database\analytics_dummy_data.sql

echo.
echo Analytics dummy data loaded successfully!
echo.
echo Statistics:
psql -h %PGHOST% -p %PGPORT% -U %PGUSER% -d %PGDATABASE% -c "SELECT tip_aktivnosti, COUNT(*) as broj FROM LogAktivnosti GROUP BY tip_aktivnosti ORDER BY broj DESC;"

pause
