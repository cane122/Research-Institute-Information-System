# PowerShell script to update the PostgreSQL database
# This script executes the dummy_data.sql file to populate the database.

# --- Configuration ---
# Please update these variables if your database connection details are different.
$dbHost = "localhost"
$dbPort = "5432"
$dbUser = "postgres"
$dbName = "research_institute"
$dbPassword = "123" # IMPORTANT: Change this to your actual PostgreSQL password

# Path to the SQL files
$schemaFile = ".\database\schema.sql"
$dummyFile = ".\database\dummy_data.sql"

# --- Execution ---
Write-Host "Starting database update..."
Write-Host "  Host: $dbHost"
Write-Host "  Database: $dbName"
Write-Host "  User: $dbUser"
Write-Host "  SQL File: $sqlFile"

# Set the password as an environment variable for the psql command
$env:PGPASSWORD = $dbPassword


# Drop all tables (CASCADE)
$dropCmd = @"
DO $$ DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END $$;
"@

try {
    # Drop all tables
    psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -c $dropCmd
    Write-Host "All tables dropped." -ForegroundColor Yellow

    # Recreate schema
    psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -f $schemaFile
    Write-Host "Schema recreated." -ForegroundColor Yellow

    # Load dummy data
    psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -f $dummyFile
    Write-Host "Dummy data loaded." -ForegroundColor Green
} catch {
    Write-Host "An error occurred during the database update." -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
} finally {
    # Clean up the environment variable
    Remove-Item env:PGPASSWORD
}
