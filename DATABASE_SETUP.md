# Research Institute Information System

## Database Setup Instructions

The application requires PostgreSQL to function properly. Follow these steps to set up the database:

### 1. Install PostgreSQL
- Download and install PostgreSQL from: https://www.postgresql.org/download/windows/
- During installation, remember the password you set for the `postgres` user

### 2. Create Database
Open PostgreSQL command line (psql) or pgAdmin and run:
```sql
CREATE DATABASE research_institute;
```

### 3. Import Schema
Navigate to the project directory and run:
```bash
psql -U postgres -d research_institute -f database/schema.sql
```

### 4. Configure Environment Variables (Optional)
You can set these environment variables to customize the database connection:
- `DB_HOST` - Database host (default: localhost:5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: password)
- `DB_NAME` - Database name (default: research_institute)

### 5. Run the Application
```bash
# Build the application
wails build

# Run the application
./build/bin/research-institute-system.exe
```

## Default Database Connection
If no environment variables are set, the application will try to connect using:
- Host: localhost:5432
- User: postgres
- Password: password
- Database: research_institute

## Troubleshooting
If you see database connection errors:
1. Make sure PostgreSQL service is running
2. Check if the database exists
3. Verify the username and password
4. Ensure the port 5432 is available
5. Check Windows Firewall settings

The application will continue to run without database connectivity, but most features will not work properly.
