# Load Analytics Dummy Data Script
# This script populates the LogAktivnosti table with realistic test data

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Analytics Dummy Data Loader" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# PostgreSQL connection parameters
$env:PGHOST = "localhost"
$env:PGPORT = "5432"
$env:PGDATABASE = "research_institute"
$env:PGUSER = "postgres"
$env:PGPASSWORD = "123"

Write-Host "Connecting to database: $env:PGDATABASE" -ForegroundColor Yellow
Write-Host ""

# Run the analytics dummy data script
Write-Host "Loading analytics test data..." -ForegroundColor Green
$scriptPath = Join-Path $PSScriptRoot "database\analytics_dummy_data.sql"

try {
    psql -h $env:PGHOST -p $env:PGPORT -U $env:PGUSER -d $env:PGDATABASE -f $scriptPath
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✓ Analytics dummy data loaded successfully!" -ForegroundColor Green
        Write-Host ""
        
        Write-Host "Activity Statistics:" -ForegroundColor Cyan
        psql -h $env:PGHOST -p $env:PGPORT -U $env:PGUSER -d $env:PGDATABASE -c "SELECT tip_aktivnosti, COUNT(*) as broj FROM LogAktivnosti GROUP BY tip_aktivnosti ORDER BY broj DESC;"
        
        Write-Host ""
        Write-Host "Summary:" -ForegroundColor Cyan
        psql -h $env:PGHOST -p $env:PGPORT -U $env:PGUSER -d $env:PGDATABASE -c "SELECT 'Total Activities' as metric, COUNT(*)::TEXT as value FROM LogAktivnosti UNION ALL SELECT 'Unique Users', COUNT(DISTINCT korisnik_id)::TEXT FROM LogAktivnosti UNION ALL SELECT 'Unique Documents', COUNT(DISTINCT entitet_id)::TEXT FROM LogAktivnosti WHERE entitet_tip = 'DOKUMENT';"
    } else {
        Write-Host ""
        Write-Host "✗ Error loading analytics data" -ForegroundColor Red
    }
} catch {
    Write-Host ""
    Write-Host "✗ Error: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "Press any key to continue..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
