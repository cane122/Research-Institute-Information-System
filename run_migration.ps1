# PowerShell script to run database migration for adding resursi column
# This adds the 'resursi' TEXT column to the Zadaci table

$dbHost = "localhost"
$dbPort = "5432"
$dbUser = "postgres"
$dbName = "research_institute"
$dbPassword = "123"

$migrationFile = ".\database\migrations\add_resursi_to_tasks.sql"

Write-Host "Running migration: add_resursi_to_tasks.sql" -ForegroundColor Cyan
Write-Host "  Host: $dbHost"
Write-Host "  Database: $dbName"
Write-Host "  User: $dbUser"

$env:PGPASSWORD = $dbPassword

try {
    psql -h $dbHost -p $dbPort -U $dbUser -d $dbName -f $migrationFile
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`nMigration completed successfully!" -ForegroundColor Green
    } else {
        Write-Host "`nMigration failed with error code: $LASTEXITCODE" -ForegroundColor Red
    }
} catch {
    Write-Host "`nError running migration: $_" -ForegroundColor Red
} finally {
    Remove-Item Env:\PGPASSWORD
}
