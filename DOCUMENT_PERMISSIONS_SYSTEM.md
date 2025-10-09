# Document Permissions System

## Overview
The Research Institute Information System now includes a comprehensive permission-based access control system for documents. Users can only View, Edit, or Delete documents if they have the appropriate privileges.

## Permission Types

### Three Permission Levels:
1. **Read Permission** (`moze_citati`)
   - Allows viewing/previewing document content
   - Shows "View" button in document list

2. **Write Permission** (`moze_menjati`)
   - Allows editing document metadata, tags, and content
   - Shows "Edit" button in document list
   - Can update existing documents

3. **Delete Permission** (`moze_brisati`)
   - Allows permanent deletion of documents
   - Shows "Delete" button in document list

## Automatic Full Access

### Users with automatic full permissions:
1. **Document Creator**
   - User who uploaded/created the document
   - Has read, write, and delete access automatically
   - No explicit permission entry needed in `DozvoleDokumenata` table

2. **System Administrators**
   - Users with `uloga = 'admin'` role
   - Have full access to ALL documents
   - Can view, edit, and delete any document

## Permission Storage

### Database Table: `DozvoleDokumenata`
```sql
CREATE TABLE DozvoleDokumenata (
    dozvola_id SERIAL PRIMARY KEY,
    dokument_id INTEGER REFERENCES Dokumenti(dokument_id),
    korisnik_id INTEGER REFERENCES Korisnici(korisnik_id),
    moze_citati BOOLEAN DEFAULT FALSE,
    moze_menjati BOOLEAN DEFAULT FALSE,
    moze_brisati BOOLEAN DEFAULT FALSE,
    datum_dodele TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Permission Assignment:
- Permissions are set during document upload in `DocumentAdd.vue`
- Managed through checkboxes in the "Document Permissions" section
- Only users with at least one permission enabled are saved to database
- Permissions can be updated when editing documents

## Backend Implementation

### Service Layer: `backend/services/document_service.go`

#### `CheckUserPermission(documentID, userID, permissionType)`:
```go
// Permission check logic:
1. Check if user is document creator (kreirao_korisnik_id)
   → Grant full access

2. Check if user is admin (naziv_uloge = 'admin')
   → Grant full access

3. Query DozvoleDokumenata table for specific permission
   → Return moze_citati / moze_menjati / moze_brisati
```

**Permission Types:** `"read"`, `"write"`, `"delete"`

### API Layer: `main.go`

#### `CheckUserPermission(documentID int, permissionType string)`:
```go
// Frontend-accessible function
// Uses currentUser.KorisnikID from session
// Returns (bool, error)
```

## Frontend Implementation

### Document Management View: `DocumentManagement.vue`

#### Permission Loading:
```javascript
// loadDocuments() function
for each document:
  1. Load document metadata
  2. Load associated tags
  3. Call CheckUserPermission for each permission type:
     - CheckUserPermission(documentID, 'read')
     - CheckUserPermission(documentID, 'write')
     - CheckUserPermission(documentID, 'delete')
  4. Store in document.permissions object
```

#### UI Rendering:
```vue
<!-- Conditional button display -->
<button v-if="doc.permissions?.canRead" @click="previewDocument(doc)">
  View
</button>
<button v-if="doc.permissions?.canWrite" @click="editDocument(doc)">
  Edit
</button>
<button v-if="doc.permissions?.canDelete" @click="deleteDocument(doc)">
  Delete
</button>

<!-- No access message -->
<span v-if="!canRead && !canWrite && !canDelete" class="no-access-text">
  No Access
</span>
```

#### Function Guards:
```javascript
function previewDocument(doc) {
  if (!doc.permissions?.canRead) {
    alert('You do not have permission to view this document.')
    return
  }
  // Navigate to preview...
}

function editDocument(doc) {
  if (!doc.permissions?.canWrite) {
    alert('You do not have permission to edit this document.')
    return
  }
  // Navigate to edit...
}

function deleteDocument(docId) {
  const doc = documents.value.find(d => d.dokument_id === docId)
  if (!doc || !doc.permissions?.canDelete) {
    alert('You do not have permission to delete this document.')
    return
  }
  // Confirm and delete...
}
```

## Permission Assignment Workflow

### During Document Upload:

1. **User selects file and fills metadata**
   - Document name, author, description, type, language
   - Tags and keywords

2. **User configures permissions**
   - See list of all system users
   - Check boxes for each user:
     - ☐ Read
     - ☐ Write
     - ☐ Delete
   - Search/filter users by name

3. **Document uploaded**
   - Document saved to `Dokumenti` table
   - `kreirao_korisnik_id` set to current user
   - Returns `dokument_id`

4. **Permissions saved**
   - For each user with at least one permission checked:
     ```javascript
     SetDocumentPermission({
       dokument_id: documentId,
       korisnik_id: user.id,
       moze_citati: user.permissions.read,
       moze_menjati: user.permissions.write,
       moze_brisati: user.permissions.delete
     })
     ```

### During Document Edit:

1. **Load existing permissions**
   - `GetDocumentPermissions(documentId)` on mount
   - Pre-check boxes for users who already have access

2. **Modify permissions**
   - User can add new permissions
   - User can remove existing permissions
   - User can modify permission levels

3. **Save updated permissions**
   - Add/update: `SetDocumentPermission()`
   - Remove: `RemoveDocumentPermission()` for unchecked users

## Security Features

### Frontend Security:
1. **UI Hiding**: Buttons hidden if no permission
2. **Function Guards**: Alert shown if user tries to bypass
3. **Permission Checks**: Before navigation to edit/preview/delete

### Backend Security:
1. **Session Validation**: Checks `currentUser` is authenticated
2. **Database Queries**: Verifies permissions in database
3. **Creator Check**: Automatic access for document creators
4. **Admin Check**: Full access for administrators
5. **SQL Protection**: Parameterized queries prevent injection

### Defense in Depth:
- Frontend checks provide UX feedback
- Backend enforces actual security
- Both layers validate permissions independently

## User Experience

### Scenario 1: Full Access (Creator/Admin)
```
User sees: [View] [Edit] [Delete]
User can: View, modify, and delete document
```

### Scenario 2: Read-Only Access
```
User sees: [View]
User can: Preview document content only
```

### Scenario 3: Read + Write Access
```
User sees: [View] [Edit]
User can: View and modify document, but not delete
```

### Scenario 4: No Access
```
User sees: "No Access" (gray italic text)
User can: See document exists, but cannot interact
```

### Scenario 5: Attempted Unauthorized Action
```
User action: Clicks disabled/hidden button or uses devtools
System response: Alert → "You do not have permission to [action] this document."
```

## Permission Matrix

| User Type | Read | Write | Delete | Notes |
|-----------|------|-------|--------|-------|
| Creator | ✅ | ✅ | ✅ | Automatic full access |
| Admin | ✅ | ✅ | ✅ | Automatic full access |
| Granted User | ✅ | ⚠️ | ⚠️ | Based on assigned permissions |
| Other User | ❌ | ❌ | ❌ | No access by default |

⚠️ = Permission must be explicitly granted

## API Reference

### Backend Functions (Wails-exposed):

#### CheckUserPermission
```go
func (a *App) CheckUserPermission(documentID int, permissionType string) (bool, error)
```
- **Parameters:**
  - `documentID`: Document to check
  - `permissionType`: "read", "write", or "delete"
- **Returns:** `(hasPermission bool, error)`

#### GetDocumentPermissions
```go
func (a *App) GetDocumentPermissions(documentID int) ([]models.DocumentPermissionResponse, error)
```
- **Returns:** List of all users with permissions on the document

#### SetDocumentPermission
```go
func (a *App) SetDocumentPermission(perm models.DocumentPermissionRequest) error
```
- **Parameters:** Permission object with dokument_id, korisnik_id, and three boolean flags

#### RemoveDocumentPermission
```go
func (a *App) RemoveDocumentPermission(documentID int, userID int) error
```
- **Parameters:** Document ID and User ID to remove permissions for

### Frontend Usage:

```javascript
// Check permission
const { CheckUserPermission } = window.go.main.App
const canEdit = await CheckUserPermission(123, 'write')

// Load permissions for document
const { GetDocumentPermissions } = window.go.main.App
const permissions = await GetDocumentPermissions(123)

// Grant permission
const { SetDocumentPermission } = window.go.main.App
await SetDocumentPermission({
  dokument_id: 123,
  korisnik_id: 456,
  moze_citati: true,
  moze_menjati: true,
  moze_brisati: false
})

// Remove permission
const { RemoveDocumentPermission } = window.go.main.App
await RemoveDocumentPermission(123, 456)
```

## Testing Checklist

### Document Upload:
- ✅ Creator automatically has full access
- ✅ Admin automatically has full access
- ✅ Assigned users see their granted permissions
- ✅ Unassigned users see "No Access"

### Document Edit:
- ✅ Creator can edit any document they created
- ✅ Admin can edit any document
- ✅ Users with write permission can edit
- ✅ Users without write permission see alert

### Document Delete:
- ✅ Creator can delete their documents
- ✅ Admin can delete any document
- ✅ Users with delete permission can delete
- ✅ Users without delete permission see alert

### Permission Assignment:
- ✅ Permissions saved during upload
- ✅ Permissions loaded during edit
- ✅ Permissions updated when modified
- ✅ Removed permissions deleted from database

### UI Behavior:
- ✅ Buttons hidden when no permission
- ✅ "No Access" shown when no permissions
- ✅ Alerts shown for unauthorized attempts
- ✅ Proper icons/styling for permission levels

## Database Queries

### Check if user is creator or admin:
```sql
SELECT d.kreirao_korisnik_id, k.naziv_uloge
FROM dokumenti d
JOIN korisnici u ON u.korisnik_id = $2
LEFT JOIN uloge k ON u.uloga_id = k.uloga_id
WHERE d.dokument_id = $1
```

### Check specific permission:
```sql
-- For read permission
SELECT moze_citati 
FROM dozvoledokumenata 
WHERE dokument_id = $1 AND korisnik_id = $2

-- For write permission
SELECT moze_menjati 
FROM dozvoledokumenata 
WHERE dokument_id = $1 AND korisnik_id = $2

-- For delete permission
SELECT moze_brisati 
FROM dozvoledokumenata 
WHERE dokument_id = $1 AND korisnik_id = $2
```

### Get all permissions for document:
```sql
SELECT dd.dozvola_id, dd.dokument_id, dd.korisnik_id,
       dd.moze_citati, dd.moze_menjati, dd.moze_brisati,
       k.korisnicko_ime, k.ime, k.prezime
FROM dozvoledokumenata dd
JOIN korisnici k ON dd.korisnik_id = k.korisnik_id
WHERE dd.dokument_id = $1
```

## Future Enhancements

### Potential Improvements:
1. **Role-Based Permissions**: Group permissions by user roles
2. **Team Permissions**: Grant access to entire teams/departments
3. **Time-Limited Access**: Permissions that expire after a period
4. **Permission Templates**: Predefined permission sets for common scenarios
5. **Audit Log**: Track when permissions are granted/revoked
6. **Inheritance**: Folders with permissions inherited by documents
7. **Notification System**: Alert users when they receive document access
8. **Permission Requests**: Users can request access to documents
9. **Bulk Permission Management**: Set permissions for multiple documents at once
10. **Advanced Sharing**: Share documents via email with temporary access links

## Related Files

### Backend:
- `backend/services/document_service.go` - Permission checking logic
- `backend/models/models.go` - Data structures
- `main.go` - API endpoints

### Frontend:
- `frontend/src/views/documents/DocumentManagement.vue` - Document list with permission-based UI
- `frontend/src/views/documents/DocumentAdd.vue` - Permission assignment during upload/edit
- `frontend/src/stores/auth.js` - User authentication state

### Database:
- `database/schema.sql` - Table definitions
- `DozvoleDokumenata` table - Permission storage
