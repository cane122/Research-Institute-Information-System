# Analytics & Activity Logging Implementation

## Overview
Comprehensive activity logging and analytics dashboard for the Research Institute Information System. This implementation tracks all user activities and provides real-time analytics with beautiful visualizations.

## Components Implemented

### 1. Database Schema

#### LogAktivnosti Table
Stores all user activities with JSONB support for flexible metadata:
```sql
CREATE TABLE LogAktivnosti (
    log_id SERIAL PRIMARY KEY,
    korisnik_id INT REFERENCES Korisnici(korisnik_id),
    tip_aktivnosti VARCHAR(50) NOT NULL,
    entitet_tip VARCHAR(50),
    entitet_id INT,
    naziv_entiteta VARCHAR(255),
    opis TEXT,
    ip_adresa VARCHAR(45),
    user_agent TEXT,
    rezultat VARCHAR(20) DEFAULT 'SUCCESS',
    dodatne_informacije JSONB,
    kreiran_datuma TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes created for performance:**
- `korisnik_id` (foreign key)
- `tip_aktivnosti` (activity type)
- `(entitet_tip, entitet_id)` (composite for entity lookup)
- `kreiran_datuma` (time-based queries)
- `rezultat` (success/failure filtering)

#### Database Views

**StatistikaDokumenata** - Document statistics:
- `ukupno_dokumenata` - Total document count
- `novih_dokumenata_mesecno` - New documents this month
- `broj_autora` - Unique authors count
- `broj_projekata_sa_dokumentima` - Projects with documents
- `prosecno_verzija_po_dokumentu` - Average versions per document

**StatistikaAktivnosti** - Activity statistics by type:
- `tip_aktivnosti` - Activity type (UPLOAD, VIEW, DELETE, etc.)
- `broj_aktivnosti` - Total activity count
- `broj_korisnika` - Unique users performing activity
- `danas` - Activities today
- `ove_nedelje` - Activities this week
- `ovog_meseca` - Activities this month
- `poslednja_aktivnost` - Last activity timestamp

**SkornjeAktivnosti** - Recent activity feed (limit 100):
- Joins with Korisnici table for user information
- Shows activity details with user name
- Ordered by most recent first

### 2. Backend Services

#### AnalyticsService (backend/services/analitics_service.go)

**Functions:**

1. **LogActivity(korisnikID *int, req ActivityLogRequest) error**
   - Logs any user activity to database
   - Supports JSONB metadata for flexible data storage
   - Automatically captures timestamp

2. **GetRecentActivity(limit int) ([]SkornjeAktivnosti, error)**
   - Retrieves recent activities from SkornjeAktivnosti view
   - Includes user information
   - Configurable limit (default 100)

3. **GetDocumentStatistics() (*StatistikaDokumenata, error)**
   - Returns aggregated document statistics
   - Used for dashboard summary cards

4. **GetActivityStatistics() ([]StatistikaAktivnosti, error)**
   - Returns activity breakdown by type
   - Shows time-based metrics (today/week/month)

5. **GetDocumentsByType() (map[string]int, error)**
   - Groups documents by type with counts
   - Used for visualization

6. **GetDocumentTrends(days int) ([]map[string]interface{}, error)**
   - Document creation trends over time
   - For time-series charts

7. **GetTopContributors(limit int) ([]map[string]interface{}, error)**
   - Users with most document uploads
   - Ordered by contribution count

### 3. API Functions (main.go)

All analytics functions exposed via Wails:
- `LogActivity(req ActivityLogRequest) error`
- `GetRecentActivity(limit int) ([]SkornjeAktivnosti, error)`
- `GetDocumentStatistics() (*StatistikaDokumenata, error)`
- `GetActivityStatistics() ([]StatistikaAktivnosti, error)`
- `GetDocumentsByType() (map[string]int, error)`
- `GetDocumentTrends(days int) ([]map[string]interface{}, error)`
- `GetTopContributors(limit int) ([]map[string]interface{}, error)`

### 4. Frontend Components

#### DocumentAnalytics.vue
Beautiful analytics dashboard with:

**Statistics Cards (4 cards with gradients):**
1. **Broj dokumenata** - Total documents (purple gradient)
2. **Broj pregleda** - Total views (pink gradient)
3. **Broj izbrisanih dokumenata** - Deleted documents (orange gradient)
4. **Broj novih dokumenata** - New documents this month (teal gradient)

**Recent Activity Section:**
- Live feed of recent user activities
- Color-coded activity icons
- Relative timestamps ("Pre 5 min", "Pre 2 sati")
- User names and activity descriptions
- Smooth hover animations

**Top Contributors Section:**
- Ranked list of top document uploaders
- Shows contribution count per user
- Purple gradient badges for rankings

**Documents by Type:**
- Visual breakdown by document type
- Animated progress bars
- Count and percentage display

**Activity Statistics Grid:**
- Cards for each activity type
- Shows total count and time-based metrics
- Displays: Today / This Week / This Month

**Features:**
- Export PDF button (placeholder for future implementation)
- Responsive design
- Smooth animations and transitions
- Beautiful gradient backgrounds
- Auto-refreshes on mount

#### Activity Logging Integration

**DocumentAdd.vue:**
- Logs UPLOAD activity when new document is uploaded
- Logs EDIT activity when document is updated
- Includes document ID and name in log

**DocumentManagement.vue:**
- Logs VIEW activity when document is previewed
- Logs DELETE activity when document is deleted
- Permission checks before logging

**Auth Store (stores/auth.js):**
- Logs LOGIN activity on successful login
- Logs LOGOUT activity when user logs out
- Includes user ID and username

### 5. Activity Types

The system tracks these activity types:
- **UPLOAD** - Document uploaded
- **EDIT** - Document updated
- **VIEW** - Document viewed/previewed
- **DELETE** - Document deleted
- **LOGIN** - User logged in
- **LOGOUT** - User logged out

Each activity is logged with:
- User ID
- Activity type
- Entity type (DOKUMENT, KORISNIK, etc.)
- Entity ID
- Entity name
- Description
- Result (SUCCESS/FAIL)
- Timestamp

## Usage

### Accessing Analytics Dashboard

1. Navigate to Document Management page
2. Click the "📊 Analytics" button in the header
3. View real-time statistics and activity feed

### How Activity Logging Works

All document operations automatically log activities:

```javascript
// Upload/Edit document
await LogActivity({
  tip_aktivnosti: 'UPLOAD',
  entitet_tip: 'DOKUMENT',
  entitet_id: docId,
  naziv_entiteta: documentName,
  opis: 'Uploaded document: Example.pdf',
  rezultat: 'SUCCESS'
})
```

### Viewing Analytics Data

The dashboard automatically loads:
- Document statistics (total, new, authors, projects)
- Activity statistics (by type with time breakdowns)
- Recent activity feed (last 50 activities)
- Top contributors (top 5 users)
- Documents by type distribution

## Database Queries

### View Recent Activities
```sql
SELECT * FROM SkornjeAktivnosti LIMIT 50;
```

### View Document Statistics
```sql
SELECT * FROM StatistikaDokumenata;
```

### View Activity Statistics
```sql
SELECT * FROM StatistikaAktivnosti;
```

### Custom Activity Query
```sql
SELECT 
    la.*,
    k.korisnicko_ime,
    k.ime || ' ' || k.prezime as korisnik
FROM LogAktivnosti la
LEFT JOIN Korisnici k ON la.korisnik_id = k.korisnik_id
WHERE tip_aktivnosti = 'UPLOAD'
ORDER BY kreiran_datuma DESC;
```

## Performance Considerations

1. **Indexes** - All frequently queried columns have indexes
2. **Views** - Pre-aggregated statistics for fast queries
3. **Limits** - Recent activity limited to 100 entries
4. **JSONB** - Efficient storage for flexible metadata
5. **Async Logging** - Activity logging doesn't block UI

## Future Enhancements

1. **PDF Export** - Generate PDF reports from analytics
2. **Charts** - Add Chart.js for visualizations:
   - Document trends over time (line chart)
   - Activity distribution (pie chart)
   - Type breakdown (bar chart)
3. **Filters** - Date range pickers for custom reporting
4. **User Analytics** - Per-user activity dashboards
5. **Project Analytics** - Per-project statistics
6. **Email Reports** - Scheduled analytics reports
7. **Real-time Updates** - WebSocket for live activity feed

## Testing

To test the implementation:

1. **Login** - Check LogAktivnosti for LOGIN entry
2. **Upload Document** - Verify UPLOAD activity logged
3. **View Document** - Confirm VIEW activity created
4. **Edit Document** - Check EDIT activity recorded
5. **Delete Document** - Verify DELETE activity logged
6. **Logout** - Confirm LOGOUT activity saved
7. **Analytics Dashboard** - Visit /documents/analytics and verify:
   - Statistics cards show correct counts
   - Recent activity displays logged activities
   - Top contributors list appears
   - Documents by type shows distribution

## Files Modified

### Backend
- `database/schema.sql` - Added LogAktivnosti table and views
- `backend/models/models.go` - Added analytics models
- `backend/services/analitics_service.go` - Complete rewrite
- `main.go` - Exposed analytics API functions

### Frontend
- `frontend/src/views/documents/DocumentAnalytics.vue` - Analytics dashboard
- `frontend/src/views/documents/DocumentAdd.vue` - Added upload/edit logging
- `frontend/src/views/documents/DocumentManagement.vue` - Added view/delete logging
- `frontend/src/stores/auth.js` - Added login/logout logging
- `frontend/src/router.js` - Added analytics route

## Summary

✅ **Complete Implementation:**
- Database schema with optimized indexes
- 3 database views for efficient statistics
- 7 backend analytics functions
- Full API exposure via Wails
- Beautiful analytics dashboard
- Activity logging in all operations
- Login/logout tracking
- Responsive design with animations

The system now tracks all user activities and provides comprehensive analytics with a beautiful, modern interface matching the design requirements.
