# Deployment Guide

## Development Deployment

### 1. Frontend Development Server

```bash
cd frontend
npm run dev
```
Aplikacija će biti dostupna na: http://localhost:5173/

### 2. Wails Development

```bash
wails dev
```

## Production Deployment

### 1. Frontend Build

```bash
cd frontend
npm run build
```

### 2. Wails Production Build

```bash
wails build
```

Executable fajl će biti kreiran u `build/bin/` direktorijumu.

### 3. Web Deployment (opciono)

Za deploy frontend-a kao web aplikacije:

1. Build frontend:
```bash
cd frontend
npm run build
```

2. Deploy `dist/` folder na web server (Apache, Nginx, etc.)

3. Konfigurišite web server za SPA routing:

**Nginx konfiguracija:**
```nginx
server {
    listen 80;
    server_name yourdomain.com;
    root /path/to/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # Gzip compression
    gzip on;
    gzip_types text/css application/javascript application/json;
}
```

**Apache konfiguracija (.htaccess):**
```apache
RewriteEngine On
RewriteBase /
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule . /index.html [L]

# Enable gzip compression
<IfModule mod_deflate.c>
    AddOutputFilterByType DEFLATE text/plain
    AddOutputFilterByType DEFLATE text/html
    AddOutputFilterByType DEFLATE text/xml
    AddOutputFilterByType DEFLATE text/css
    AddOutputFilterByType DEFLATE application/xml
    AddOutputFilterByType DEFLATE application/xhtml+xml
    AddOutputFilterByType DEFLATE application/rss+xml
    AddOutputFilterByType DEFLATE application/javascript
    AddOutputFilterByType DEFLATE application/x-javascript
</IfModule>
```

## Docker Deployment (budući razvoj)

### Dockerfile (frontend)
```dockerfile
FROM node:18-alpine as build

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci --only=production

COPY frontend/ .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  frontend:
    build: .
    ports:
      - "80:80"
    depends_on:
      - backend
      
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:password@db:5432/research_db
    depends_on:
      - db
      
  db:
    image: postgres:14
    environment:
      POSTGRES_DB: research_db
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/schema.sql:/docker-entrypoint-initdb.d/schema.sql

volumes:
  postgres_data:
```

## Performance Optimizations

### Frontend Optimizations

1. **Code Splitting:**
```javascript
// router.js
const Dashboard = () => import('./views/Dashboard.vue')
const Projects = () => import('./views/Projects.vue')
```

2. **Asset Optimization:**
```javascript
// vite.config.js
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia']
        }
      }
    }
  }
})
```

3. **PWA Setup:**
```bash
npm install -D vite-plugin-pwa
```

### Backend Optimizations

1. **Gzip Compression**
2. **Database Connection Pooling**
3. **Caching Strategy**
4. **Rate Limiting**

## Monitoring and Logging

### Frontend Monitoring
- Error tracking sa Sentry
- Performance monitoring
- User analytics

### Backend Monitoring  
- Application logs
- Database performance
- API response times
- Health checks

## Security Considerations

1. **HTTPS Only** za production
2. **CSP Headers** za security
3. **Rate Limiting** za API
4. **Input Validation** na backend-u
5. **JWT Token Expiration**
6. **CORS Configuration**

## Backup Strategy

1. **Database Backup:**
```bash
pg_dump -h localhost -U username -d research_db > backup.sql
```

2. **File Upload Backup:**
```bash
rsync -av uploads/ backup/uploads/
```

3. **Automated Backups:**
- Daily database backups
- Weekly full system backups
- Offsite backup storage
