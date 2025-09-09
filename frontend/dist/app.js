// Research Institute System Frontend
class App {
    constructor() {
        this.currentUser = null;
        this.currentSection = 'projects';
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.showLoginScreen();
    }

    setupEventListeners() {
        // Login form
        document.getElementById('loginForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleLogin();
        });

        // Logout button
        document.getElementById('logoutBtn').addEventListener('click', () => {
            this.handleLogout();
        });

        // Navigation
        document.querySelectorAll('.nav-item').forEach(item => {
            item.addEventListener('click', (e) => {
                e.preventDefault();
                const section = item.getAttribute('data-section');
                this.switchSection(section);
            });
        });

        // Create project button
        document.getElementById('createProjectBtn').addEventListener('click', () => {
            this.showModal('createProjectModal');
        });

        // Create user button
        document.getElementById('createUserBtn').addEventListener('click', () => {
            this.showModal('createUserModal');
        });

        // Create project form
        document.getElementById('createProjectForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleCreateProject();
        });

        // Create user form
        document.getElementById('createUserForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleCreateUser();
        });

        // Generate password button
        document.getElementById('generatePasswordBtn').addEventListener('click', () => {
            this.generateTempPassword();
        });

        // Modal close buttons
        document.querySelectorAll('.close').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const modal = e.target.closest('.modal');
                this.closeModal(modal.id);
            });
        });

        // Close modal on outside click
        window.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal')) {
                this.closeModal(e.target.id);
            }
        });
    }

    showLoginScreen() {
        document.getElementById('loginScreen').classList.add('active');
        document.getElementById('dashboardScreen').classList.remove('active');
    }

    showDashboard() {
        document.getElementById('loginScreen').classList.remove('active');
        document.getElementById('dashboardScreen').classList.add('active');
        this.updateUserInterface();
        this.loadUserProjects();
    }

    showFirstTimeSetup(user) {
        // Hide login screen
        document.getElementById('loginScreen').classList.remove('active');
        
        // Create first-time setup modal
        this.createFirstTimeSetupModal(user);
    }

    createFirstTimeSetupModal(user) {
        // Remove existing modal if it exists
        const existingModal = document.getElementById('firstTimeSetupModal');
        if (existingModal) {
            existingModal.remove();
        }

        // Create modal HTML
        const modalHTML = `
            <div id="firstTimeSetupModal" class="modal show">
                <div class="modal-content">
                    <div class="modal-header">
                        <h2>Prva prijava - Postavite lozinku</h2>
                    </div>
                    <div class="modal-body">
                        <p>Dobrodošli ${user.ime} ${user.prezime}!</p>
                        <p>Ovo je vaša prva prijava u sistem. Molimo vas da postavite novu lozinku.</p>
                        <form id="firstTimeSetupForm">
                            <div class="form-group">
                                <label for="newPassword">Nova lozinka:</label>
                                <input type="password" id="newPassword" name="newPassword" required minlength="6">
                            </div>
                            <div class="form-group">
                                <label for="confirmPassword">Potvrdite lozinku:</label>
                                <input type="password" id="confirmPassword" name="confirmPassword" required minlength="6">
                            </div>
                            <div id="passwordError" class="error-message" style="display: none;"></div>
                            <div class="form-actions">
                                <button type="submit" class="btn btn-primary">Postavite lozinku</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        `;

        // Add modal to body
        document.body.insertAdjacentHTML('beforeend', modalHTML);

        // Store user data for later use
        this.firstTimeUser = user;

        // Add event listener for form submission
        document.getElementById('firstTimeSetupForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleFirstTimePasswordSetup();
        });
    }

    async handleFirstTimePasswordSetup() {
        const newPassword = document.getElementById('newPassword').value;
        const confirmPassword = document.getElementById('confirmPassword').value;
        const errorDiv = document.getElementById('passwordError');

        // Reset error
        errorDiv.style.display = 'none';
        errorDiv.textContent = '';

        // Validate passwords
        if (newPassword !== confirmPassword) {
            errorDiv.textContent = 'Lozinke se ne poklapaju!';
            errorDiv.style.display = 'block';
            return;
        }

        if (newPassword.length < 6) {
            errorDiv.textContent = 'Lozinka mora imati najmanje 6 karaktera!';
            errorDiv.style.display = 'block';
            return;
        }

        try {
            // Call backend to complete first-time setup
            console.log('FirstTimeUser object:', this.firstTimeUser);
            const response = await window.go.main.App.CompleteFirstTimeSetup(
                this.firstTimeUser.korisnicko_ime || this.firstTimeUser.KorisnickoIme, 
                newPassword
            );

            if (response.success) {
                // Remove modal
                document.getElementById('firstTimeSetupModal').remove();
                
                // Set current user and show dashboard
                this.currentUser = this.firstTimeUser;
                this.showDashboard();
                
                this.showSuccessMessage('Lozinka je uspešno postavljena! Dobrodošli u sistem!');
            } else {
                errorDiv.textContent = response.message || 'Greška pri postavljanju lozinke';
                errorDiv.style.display = 'block';
            }
        } catch (error) {
            console.error('First time setup error:', error);
            errorDiv.textContent = 'Greška pri postavljanju lozinke: ' + error.message;
            errorDiv.style.display = 'block';
        }
    }

    updateUserInterface() {
        if (!this.currentUser) return;

        // Update welcome message
        const welcomeSpan = document.getElementById('userWelcome');
        welcomeSpan.textContent = `Dobrodošli, ${this.currentUser.ime} ${this.currentUser.prezime}`;

        // Show/hide role-specific elements
        const adminElements = document.querySelectorAll('.admin-only');
        const managerElements = document.querySelectorAll('.manager-only');

        // Map role IDs to role names (based on your database schema)
        const getRoleName = (ulogaId) => {
            switch (ulogaId) {
                case 1: return 'Administrator';
                case 2: return 'Rukovodilac projekta';
                case 3: return 'Istraživač';
                case 4: return 'Organizator projekta';
                default: return 'Nepoznato';
            }
        };

        const roleName = getRoleName(this.currentUser.ulogaId);

        adminElements.forEach(el => {
            el.style.display = roleName === 'Administrator' ? 'block' : 'none';
        });

        managerElements.forEach(el => {
            el.style.display = 
                (roleName === 'Rukovodilac projekta' || 
                 roleName === 'Administrator') ? 'block' : 'none';
        });
    }

    switchSection(section) {
        // Update navigation
        document.querySelectorAll('.nav-item').forEach(item => {
            item.classList.remove('active');
        });
        document.querySelector(`[data-section="${section}"]`).classList.add('active');

        // Update content
        document.querySelectorAll('.content-section').forEach(sec => {
            sec.classList.remove('active');
        });
        document.getElementById(`${section}Section`).classList.add('active');

        this.currentSection = section;

        // Load section-specific data
        switch (section) {
            case 'projects':
                this.loadUserProjects();
                break;
            case 'documents':
                this.loadDocuments();
                break;
            case 'documentation':
                this.loadDocumentation();
                break;
            case 'users':
                this.loadUsers();
                break;
        }
    }

    async handleLogin() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorDiv = document.getElementById('loginError');
    const statusDiv = document.getElementById('loginStatus');
    const loginBtn = document.querySelector('#loginForm button[type="submit"]');
    // Reset status and error
    if (statusDiv) {
        statusDiv.textContent = '';
        statusDiv.style.display = 'none';
    }
    errorDiv.textContent = '';
    errorDiv.style.display = 'none';
    loginBtn.disabled = true;
    loginBtn.textContent = 'Prijavljivanje...';
    console.log('Login attempt:', { username });

    try {
        // Call Wails backend
        const response = await window.go.main.App.Login(username, password);
        console.log('Login response:', response);

        if (response.success) {
            if (response.message === 'FIRST_TIME_LOGIN') {
                // Handle first-time login - show password setup form
                this.showFirstTimeSetup(response.user);
            } else {
                this.currentUser = response.user;
                this.showDashboard();
                errorDiv.textContent = '';
                errorDiv.style.display = 'none';
                this.showSuccessMessage('Uspešno ste se prijavili!');
                if (statusDiv) {
                    statusDiv.textContent = 'Uspešno ste se prijavili!';
                    statusDiv.style.color = 'green';
                    statusDiv.style.display = 'block';
                }
            }
        } else {
            errorDiv.textContent = response.message || 'Neispravno korisničko ime ili lozinka.';
            errorDiv.style.display = 'block';
            errorDiv.classList.add('error-visible');
            this.showErrorMessage('Neuspešno logovanje! Proverite podatke.');
            if (statusDiv) {
                statusDiv.textContent = response.message || 'Neuspešno logovanje! Proverite podatke.';
                statusDiv.style.color = 'red';
                statusDiv.style.display = 'block';
            }
        }
    } catch (error) {
        console.error('Login error:', error);
        if (error && error.message) {
            errorDiv.textContent = `Greška: ${error.message}`;
        } else if (error && error.toString) {
            errorDiv.textContent = `Greška: ${error.toString()}`;
        } else {
            errorDiv.textContent = 'Greška pri povezivanju sa serverom. Proverite internet konekciju ili pokušajte ponovo.';
        }
        errorDiv.style.display = 'block';
        errorDiv.classList.add('error-visible');
        this.showErrorMessage('Neuspešno logovanje! Greška u konekciji ili serveru.');
        if (statusDiv) {
            statusDiv.textContent = 'Neuspešno logovanje! Greška u konekciji ili serveru.';
            statusDiv.style.color = 'red';
            statusDiv.style.display = 'block';
        }
    }
    loginBtn.disabled = false;
    loginBtn.textContent = 'Prijavi se';
}

    async handleLogout() {
        try {
            await window.go.main.App.Logout();
            this.currentUser = null;
            this.showLoginScreen();
            
            // Clear form
            document.getElementById('loginForm').reset();
        } catch (error) {
            console.error('Logout error:', error);
        }
    }

    async handleCreateProject() {
        const formData = new FormData(document.getElementById('createProjectForm'));
        
        const project = {
            name: formData.get('name'),
            description: formData.get('description'),
            startDate: formData.get('startDate') || null,
            endDate: formData.get('endDate') || null
        };

        try {
            await window.go.main.App.CreateProject(project);
            this.closeModal('createProjectModal');
            this.loadUserProjects();
            this.showSuccessMessage('Projekat je uspešno kreiran');
        } catch (error) {
            this.showErrorMessage('Greška pri kreiranju projekta: ' + error.message);
        }
    }

    async handleCreateUser() {
        const formData = new FormData(document.getElementById('createUserForm'));
        
        const user = {
            username: formData.get('username'),
            email: formData.get('email'),
            firstName: formData.get('firstName'),
            lastName: formData.get('lastName'),
            roleId: parseInt(formData.get('roleId'))
        };

        const tempPassword = formData.get('tempPassword');

        try {
            await window.go.main.App.CreateUser(user, tempPassword);
            this.closeModal('createUserModal');
            this.loadUsers();
            this.showSuccessMessage('Korisnik je uspešno kreiran');
        } catch (error) {
            this.showErrorMessage('Greška pri kreiranju korisnika: ' + error.message);
        }
    }

    async loadUserProjects() {
        const projectsList = document.getElementById('projectsList');
        
        try {
            const projects = await window.go.main.App.GetUserProjects();
            
            if (projects && projects.length > 0) {
                projectsList.innerHTML = projects.map(project => `
                    <div class="project-card">
                        <h3>${project.name}</h3>
                        <p>${project.description || 'Nema opisa'}</p>
                        <div class="project-meta">
                            <span>Rukovodilac: ${project.manager ? project.manager.firstName + ' ' + project.manager.lastName : 'N/A'}</span>
                            <span class="project-status ${project.status.toLowerCase()}">${project.status}</span>
                        </div>
                    </div>
                `).join('');
            } else {
                projectsList.innerHTML = '<div class="text-center"><p>Nemate dodeljene projekte.</p></div>';
            }
        } catch (error) {
            projectsList.innerHTML = '<div class="text-center"><p>Greška pri učitavanju projekata.</p></div>';
            console.error('Load projects error:', error);
        }
    }

    async loadDocuments() {
        const documentsList = document.getElementById('documentsList');
        documentsList.innerHTML = '<div class="text-center"><p>Funkcionalnost upravljanja dokumentima će biti implementirana uskoro.</p></div>';
    }

    async loadDocumentation() {
        const documentationContent = document.getElementById('documentationContent');
        documentationContent.innerHTML = '<div class="text-center"><p>Funkcionalnost projektne dokumentacije će biti implementirana uskoro.</p></div>';
    }

    async loadUsers() {
        const usersList = document.getElementById('usersList');
        
        try {
            const users = await window.go.main.App.GetAllUsers();
            
            if (users && users.length > 0) {
                usersList.innerHTML = users.map(user => `
                    <div class="list-item">
                        <div class="item-info">
                            <h4>${user.firstName} ${user.lastName}</h4>
                            <div class="item-meta">
                                ${user.username} | ${user.email} | ${user.role ? user.role.name : 'N/A'} | ${user.status}
                            </div>
                        </div>
                        <div class="item-actions">
                            <button class="btn btn-secondary btn-sm" onclick="app.editUser(${user.id})">Izmeni</button>
                            <button class="btn btn-danger btn-sm" onclick="app.resetPassword(${user.id})">Resetuj lozinku</button>
                        </div>
                    </div>
                `).join('');
            } else {
                usersList.innerHTML = '<div class="text-center"><p>Nema korisnika u sistemu.</p></div>';
            }
        } catch (error) {
            usersList.innerHTML = '<div class="text-center"><p>Greška pri učitavanju korisnika.</p></div>';
            console.error('Load users error:', error);
        }
    }

    generateTempPassword() {
        const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        let password = '';
        for (let i = 0; i < 12; i++) {
            password += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        document.getElementById('tempPassword').value = password;
    }

    showModal(modalId) {
        document.getElementById(modalId).classList.add('show');
        
        // Generate password for user creation modal
        if (modalId === 'createUserModal') {
            this.generateTempPassword();
        }
    }

    closeModal(modalId) {
        document.getElementById(modalId).classList.remove('show');
        
        // Clear forms
        const form = document.querySelector(`#${modalId} form`);
        if (form) {
            form.reset();
        }
    }

    showSuccessMessage(message) {
        const statusDiv = document.getElementById('loginStatus');
        if (statusDiv) {
            statusDiv.textContent = message;
            statusDiv.style.color = 'green';
            statusDiv.style.display = 'block';
        }
    }

    showErrorMessage(message) {
        const statusDiv = document.getElementById('loginStatus');
        if (statusDiv) {
            statusDiv.textContent = message;
            statusDiv.style.color = 'red';
            statusDiv.style.display = 'block';
        }
    }

    editUser(userId) {
        // Implement user editing functionality
        alert('Funkcionalnost izmene korisnika će biti implementirana uskoro');
    }

    resetPassword(userId) {
        if (confirm('Da li ste sigurni da želite da resetujete lozinku za ovog korisnika?')) {
            // Implement password reset functionality
            alert('Funkcionalnost resetovanja lozinke će biti implementirana uskoro');
        }
    }
}

// Global functions for modal management
function closeModal(modalId) {
    app.closeModal(modalId);
}

// Initialize app when DOM is loaded
let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new App();
});

// Wails runtime check
window.addEventListener('DOMContentLoaded', () => {
    // Check if Wails runtime is available
    if (typeof window.go === 'undefined') {
        console.warn('Wails runtime not available. Running in development mode.');
        
        // Mock backend for development
        window.go = {
            main: {
                App: {
                    Login: async (username, password) => {
                        // Mock login for development
                        if (username === 'admin' && password === 'admin') {
                            return {
                                success: true,
                                user: {
                                    id: 1,
                                    username: 'admin',
                                    firstName: 'Admin',
                                    lastName: 'User',
                                    email: 'admin@example.com',
                                    role: { id: 1, name: 'Administrator' }
                                }
                            };
                        }
                        return { success: false, message: 'Neispravno korisničko ime ili lozinka' };
                    },
                    Logout: async () => { return true; },
                    GetUserProjects: async () => {
                        return [
                            {
                                id: 1,
                                name: 'Primer projekta 1',
                                description: 'Ovo je primer projekta za testiranje',
                                status: 'Aktivan',
                                manager: { firstName: 'Marko', lastName: 'Petrović' }
                            },
                            {
                                id: 2,
                                name: 'Primer projekta 2',
                                description: 'Drugi primer projekta',
                                status: 'Završen',
                                manager: { firstName: 'Ana', lastName: 'Jovanović' }
                            }
                        ];
                    },
                    CreateProject: async (project) => { return true; },
                    GetAllUsers: async () => {
                        return [
                            {
                                id: 1,
                                username: 'admin',
                                firstName: 'Admin',
                                lastName: 'User',
                                email: 'admin@example.com',
                                status: 'aktivan',
                                role: { id: 1, name: 'Administrator' }
                            },
                            {
                                id: 2,
                                username: 'manager1',
                                firstName: 'Marko',
                                lastName: 'Petrović',
                                email: 'marko@example.com',
                                status: 'aktivan',
                                role: { id: 2, name: 'Rukovodilac projekta' }
                            }
                        ];
                    },
                    CreateUser: async (user, tempPassword) => { return true; },
                    CompleteFirstTimeSetup: async (username, newPassword) => {
                        console.log('Mock: First time setup completed for', username);
                        return { success: true };
                    }
                }
            }
        };
    }
});
