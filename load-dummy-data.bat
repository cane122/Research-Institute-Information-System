@echo off
echo ======================================
echo    UCITAVANJE DUMMY PODATAKA
echo ======================================
echo.

echo Postavlja se password...
set PGPASSWORD=123

echo Ucitavam dummy podatke u bazu research_institute...
"C:\Program Files\PostgreSQL\17\bin\psql.exe" -U postgres -d research_institute -f "database\dummy_data.sql"

if %ERRORLEVEL% equ 0 (
    echo.
    echo ======================================
    echo ✅ DUMMY PODACI USPESNO UCITANI!
    echo ======================================
    echo.
    echo Ucitano:
    echo - 10 korisnika ^(admin, direktor, rukovodioci, istrazivaci^)
    echo - 5 projekata ^(AI, IoT, Kvantno racunarstvo, Blockchain, Energija^)
    echo - 10 zadataka sa realnim opisima
    echo - 7 dokumenata sa verzijama
    echo - Tagovi, komentari i log aktivnosti
    echo.
    echo Testni korisnici:
    echo   admin/password - Administrator
    echo   researcher1/password - Istrazivac
    echo   proj_manager1/password - Rukovodilac projekta
) else (
    echo.
    echo ❌ GRESKA pri ucitavanju podataka!
    echo Proverite da li je PostgreSQL pokrenut i da li je baza kreirana.
)

echo.
set PGPASSWORD=
echo Pritisnite bilo koji taster za zatvaranje...
pause >nul
