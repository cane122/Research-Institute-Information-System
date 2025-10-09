$env:PGPASSWORD = "123"
psql -h localhost -U postgres -d research_institute -c "ALTER TABLE Dokumenti ADD COLUMN IF NOT EXISTS kljucne_reci TEXT;"
Write-Host "Kolona kljucne_reci je uspešno dodata!"
