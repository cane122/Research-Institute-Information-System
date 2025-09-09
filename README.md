# Istraživačko-razvojni centar - Informacioni sistem

Desktop aplikacija za upravljanje projektima, dokumentima i dokumentacijom u istraživačko-razvojnom centru.

## 📋 Pregled

Sistem obuhvata tri glavna podsistema:

1. **📁 Upravljanje dokumentima** - Organizacija i upravljanje dokumentima sa meta-podacima
2. **📋 Priprema projektne dokumentacije** - Životni ciklus projektne dokumentacije
3. **🚀 Realizacija projekata** - Kreiranje, upravljanje i praćenje projekata i zadataka

## 🛠️ Tehnologije

- **Backend**: Go (Wails v2 framework)
- **Frontend**: HTML, CSS, JavaScript
- **Baza podataka**: PostgreSQL
- **Desktop aplikacija**: Wails v2

## Features

### User Management
- Role-based access control (Administrator, Project Manager, Researcher, Project Organizer)
- User authentication and authorization
- Profile management

### Project Management
- Create and manage research projects
- Assign team members to projects
- Define project workflows and phases
- Task management with Kanban-style interface
- Progress tracking and analytics

### Document Management
- Upload and organize documents in folders
- Advanced metadata management
- Document versioning
- Tag-based organization
- Advanced search capabilities
- Access control and permissions

### Project Documentation
- Workflow-based documentation management
- Phase tracking for document preparation
- Version control for project documents
- Analytics and reporting

## Prerequisites

To run this application, you need:

1. **Go** (version 1.21 or higher)
2. **Node.js** (for frontend development)
3. **Wails CLI** 
4. **PostgreSQL** (version 12 or higher)

## Installation

### 1. Install Go
Download and install Go from [https://golang.org/dl/](https://golang.org/dl/)

### 2. Install Wails
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 3. Install PostgreSQL
Install PostgreSQL and create a database:
```sql
CREATE DATABASE research_institute;
```

### 4. Setup Database Schema
Run the database schema:
```bash
psql -U username -d research_institute -f database/schema.sql
```

### 5. Configure Database Connection
Update the database connection string in `main.go`:
```go
db, err := sql.Open("postgres", "postgres://username:password@localhost/research_institute?sslmode=disable")
```

## Running the Application

### Development Mode
```bash
wails dev
```

### Build for Production
```bash
wails build
```

## Project Structure

```
research-institute-system/
├── app/                    # Main application logic
├── backend/               # Go backend code
│   ├── models/           # Data models
│   ├── repositories/     # Database repositories
│   └── services/         # Business logic services
├── frontend/             # Web frontend
│   └── dist/            # Built frontend files
│       ├── index.html   # Main HTML file
│       ├── styles.css   # Stylesheet
│       └── app.js       # JavaScript application
├── database/             # Database schema and migrations
│   └── schema.sql       # PostgreSQL schema
├── main.go              # Main application entry point
├── go.mod               # Go module dependencies
└── README.md            # This file
```

## User Roles and Permissions

### Administrator
- Manage all users and roles
- Access to all system functions
- User creation and management
- System configuration

### Project Manager (Rukovodilac projekta)
- Create and manage projects
- Assign team members
- Define workflows and phases
- Access project analytics

### Researcher (Istraživač)
- Work on assigned tasks
- Upload and manage documents
- Comment on tasks
- Request phase changes

### Project Organizer (Organizator projekta)
- View project progress
- Access reports and analytics
- Read-only access to most functions

## Default Login Credentials

For development and testing:
- **Username**: admin
- **Password**: admin

## Database Schema

The application uses a PostgreSQL database with the following main tables:

- `Uloge` - User roles
- `Korisnici` - System users
- `Projekti` - Research projects
- `Zadaci` - Project tasks
- `Dokumenti` - Document management
- `RadniTokovi` - Workflow definitions
- `Faze` - Workflow phases

See `database/schema.sql` for the complete schema.

## API Endpoints (Wails Bindings)

The following Go methods are exposed to the frontend:

- `Login(username, password string)` - User authentication
- `Logout()` - User logout
- `GetCurrentUser()` - Get current user info
- `CreateUser(user, tempPassword)` - Create new user (Admin only)
- `GetAllUsers()` - List all users (Admin only)
- `GetUserProjects()` - Get user's projects
- `CreateProject(project)` - Create new project

## Development

### Frontend Development
The frontend is built with vanilla HTML, CSS, and JavaScript. All files are in the `frontend/dist/` directory.

### Backend Development
The backend is organized into:
- **Models**: Data structures (`backend/models/`)
- **Repositories**: Database access layer (`backend/repositories/`)
- **Services**: Business logic (`backend/services/`)

### Adding New Features
1. Define models in `backend/models/`
2. Create repository methods in `backend/repositories/`
3. Implement business logic in `backend/services/`
4. Add Wails bindings in `main.go`
5. Update frontend in `frontend/dist/`

## Security Features

- Password hashing using Argon2
- Role-based access control
- Session management
- Input validation
- SQL injection prevention

## Troubleshooting

### Database Connection Issues
1. Ensure PostgreSQL is running
2. Check database credentials in `main.go`
3. Verify database exists and schema is loaded

### Build Issues
1. Ensure Go and Wails are properly installed
2. Run `go mod tidy` to resolve dependencies
3. Check Wails version compatibility

### Frontend Issues
1. Check browser console for JavaScript errors
2. Verify Wails runtime is available
3. Use development mode for debugging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is developed for educational and research purposes.

## Support

For issues and questions, please refer to the project documentation or create an issue in the repository.
