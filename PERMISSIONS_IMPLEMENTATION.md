# Document Permissions/Sharing Implementation

## Overview
Implemented a complete document permissions/sharing system using the existing `DozvoleDokumenata` database table. Users can now assign granular read/write/delete permissions to team members when uploading documents.

## Database Schema
The system uses the existing `DozvoleDokumenata` table:
```sql
CREATE TABLE DozvoleDokumenata (
    dozvola_id SERIAL PRIMARY KEY,
    dokument_id INT REFERENCES Dokumenti(dokument_id) ON DELETE CASCADE,
    korisnik_id INT REFERENCES Korisnici(korisnik_id) ON DELETE CASCADE,
    moze_citati BOOLEAN DEFAULT FALSE,
    moze_menjati BOOLEAN DEFAULT FALSE,
    moze_brisati BOOLEAN DEFAULT FALSE,
    UNIQUE(dokument_id, korisnik_id)
);
```

## Backend Implementation

### Models (`backend/models/models.go`)
Added two new structs for permission management:

1. **DocumentPermissionRequest** - Input structure for setting permissions
   - `dokument_id` - Document ID
   - `korisnik_id` - User ID
   - `moze_citati` - Can read permission
   - `moze_menjati` - Can edit permission
   - `moze_brisati` - Can delete permission

2. **DocumentPermissionResponse** - Output structure with user details
   - All fields from request
   - `korisnicko_ime` - Username
   - `ime` - First name
   - `prezime` - Last name

### Service Layer (`backend/services/document_service.go`)

#### 1. GetDocumentPermissions(documentID int)
- Returns all permissions for a document with user details
- Uses JOIN query to include username, first name, last name
- Orders results by username

#### 2. SetDocumentPermission(req DocumentPermissionRequest)
- Upserts permission (INSERT...ON CONFLICT DO UPDATE)
- Allows updating existing permissions
- Atomic operation

#### 3. RemoveDocumentPermission(documentID, userID int)
- Deletes specific user permission for a document
- Returns error if permission doesn't exist

#### 4. CheckUserPermission(documentID, userID, permissionType string)
- Validates if user has specific permission
- Permission types: "read", "write", "delete"
- Returns boolean

#### 5. UploadDocument() - Modified
- **Changed return type from `error` to `(int, error)`**
- Now returns the uploaded document ID
- Necessary for linking permissions to newly created documents

### API Layer (`main.go`)
All 5 functions exposed to frontend via Wails:
- Authentication check on all functions (`currentUser == nil`)
- Proper error messages in Serbian
- CheckUserPermission() automatically uses current user's ID

## Frontend Implementation

### User Loading (`DocumentAdd.vue`)
Modified `onMounted()` to load users from database:
```javascript
const { GetAllUsers } = window.go.main.App
const allUsers = await GetAllUsers()

users.value = (allUsers || []).map(user => ({
  id: user.korisnik_id,
  name: user.korisnicko_ime || `${user.ime} ${user.prezime}`.trim(),
  permissions: {
    read: false,
    write: false,
    delete: false
  }
}))
```

### UI Structure
The permissions UI includes:
- **Search box** - Filter users by name (`userSearch`)
- **User table** - Shows all users with permission checkboxes
- **Checkboxes** - Read (R), Write (W), Delete (D) for each user
- **Reactive filtering** - `filteredUsers` computed property

### Permission Saving (`handleSave()`)
After successful document upload:
```javascript
const documentId = await DocumentService.uploadDocument(documentData, file)

const { SetDocumentPermission } = window.go.main.App
for (const user of users.value) {
  if (user.permissions.read || user.permissions.write || user.permissions.delete) {
    await SetDocumentPermission({
      dokument_id: documentId,
      korisnik_id: user.id,
      moze_citati: user.permissions.read,
      moze_menjati: user.permissions.write,
      moze_brisati: user.permissions.delete
    })
  }
}
```

### Document Service (`documentService.js`)
Modified `uploadDocument()` to return document ID:
```javascript
const documentId = await UploadDocument(request, fileData, file.name)
return documentId
```

## Flow Diagram

```
User fills form → Selects file → Sets permissions → Clicks "Sačuvaj dokument"
                                                              ↓
                                            1. Upload document (returns ID)
                                                              ↓
                                            2. Loop through users with permissions
                                                              ↓
                                            3. Call SetDocumentPermission for each
                                                              ↓
                                            4. Redirect to /documents
```

## Testing Checklist

- [ ] Load users in Team Permissions section
- [ ] Search users by name
- [ ] Check/uncheck permission boxes
- [ ] Upload document with permissions
- [ ] Verify permissions saved in database:
  ```sql
  SELECT d.*, k.korisnicko_ime 
  FROM DozvoleDokumenata d
  JOIN Korisnici k ON d.korisnik_id = k.korisnik_id
  WHERE dokument_id = <uploaded_document_id>;
  ```
- [ ] Test with multiple users
- [ ] Test with all permission combinations (R only, RW, RWD, etc.)

## Error Handling

1. **User loading fails** - Console error, empty user list
2. **Permission save fails** - Console error, continues with other users
3. **Document upload fails** - Alert shown, no permissions saved
4. **Not authenticated** - Backend returns "niste prijavljeni"

## Next Steps

1. **Document Management View**
   - Show who has access to each document
   - Edit permissions on existing documents
   - Call GetDocumentPermissions() to display current state

2. **Permission Enforcement**
   - Check permissions before allowing document view/edit/delete
   - Use CheckUserPermission() before operations
   - Hide UI elements based on permissions

3. **Bulk Operations**
   - Assign same permissions to multiple users at once
   - Copy permissions from another document
   - Permission templates/roles

## Files Modified

### Backend
- `backend/models/models.go` - Added permission models
- `backend/services/document_service.go` - Added 4 permission functions, modified UploadDocument
- `main.go` - Exposed permission functions to frontend

### Frontend
- `frontend/src/views/documents/DocumentAdd.vue` - User loading, permission saving
- `frontend/src/services/documentService.js` - Return document ID from upload

### Build
- Wails bindings regenerated successfully
- Frontend build: 175.73 kB (gzip: 59.19 kB) ✓
- Backend compilation successful ✓
