# Edit Document Feature - Implementation Summary

## Overview
Added full edit/update functionality to the Document Management system. Users can now click the "Edit" button on any document to modify its metadata, tags, permissions, and optionally replace the file.

## Changes Made

### 1. Frontend Route Integration
**File:** `frontend/src/views/documents/DocumentAdd.vue`

- Added `useRoute` from vue-router to read query parameters
- Added `isEditMode` and `documentId` reactive refs to track edit state
- Page title dynamically changes: "Upload Document" → "Edit Document"

### 2. Document Loading on Mount
**Function:** `onMounted()`

When `route.query.id` is present:
1. Sets `isEditMode = true`
2. Loads document details via `GetDocumentByID(documentId)`
3. Populates form fields:
   - Document name
   - Author
   - Description
   - Type (PDF, DOC, etc.)
   - Language
4. Loads associated tags via `GetDocumentTags(documentId)`
5. Loads existing permissions via `GetDocumentPermissions(documentId)`
6. Updates user permissions checkboxes

### 3. File Selection - Optional in Edit Mode
**Changes:**
- Label shows: "Select File (optional - leave empty to keep current file)"
- Placeholder text in edit mode: "No new file selected - keeping current file"
- `canUpload` computed property updated:
  - **Upload mode:** Requires file + name + author
  - **Edit mode:** Only requires name + author (file optional)

### 4. Save/Update Logic
**Function:** `handleSave()`

**Upload Mode** (new document):
- Requires file selection
- Calls `DocumentService.uploadDocument()`
- Saves permissions for selected users

**Edit Mode** (existing document):
- File selection is optional
- Calls `UpdateDocument(documentId, updateRequest)` from backend
- Updates document metadata
- If new file selected: includes file content in update
- Updates permissions:
  - Adds/updates permissions for users with checkboxes enabled
  - Removes permissions for users with all checkboxes disabled
- Success message: "Document updated successfully!"

### 5. UI Updates
**Button Text:**
- Upload mode: "Upload Document" → "Uploading..."
- Edit mode: "Update Document" → "Updating..."

**Progress Overlay:**
- Same progress bar works for both upload and update
- Message adapts based on mode

### 6. Permission Management in Edit Mode
**Logic:**
- Loads existing permissions from `DozvoleDokumenata` table
- Pre-checks permission checkboxes for users who already have access
- On save:
  - `SetDocumentPermission()` - adds/updates permissions
  - `RemoveDocumentPermission()` - removes permissions when all unchecked

## Backend Integration

### Used Backend Functions:
1. `GetDocumentByID(documentID int)` - Fetch document details
2. `GetDocumentTags(documentID int)` - Fetch associated tags
3. `GetDocumentPermissions(documentID int)` - Fetch existing permissions
4. `UpdateDocument(documentID int, req UploadDocumentRequest)` - Update document
5. `SetDocumentPermission(perm DocumentPermissionRequest)` - Save permissions
6. `RemoveDocumentPermission(docID, userID int)` - Remove permissions

## User Flow

### Editing a Document:
1. User navigates to Document Management
2. Clicks "Edit" button on any document
3. Redirected to: `/documents/add?id=123`
4. Form loads with existing data:
   - Document name, author, description pre-filled
   - Tags pre-selected
   - User permissions pre-checked
5. User can modify:
   - Any metadata field
   - Tags selection
   - User permissions
   - Optionally upload new file
6. Clicks "Update Document"
7. Backend updates document record
8. Permissions updated in database
9. Redirected back to Document Management

## File Handling

### Upload Mode:
- File is **required**
- File content uploaded to server
- Stored in `Dokumenti.fajl_sadrzaj`

### Edit Mode:
- File is **optional**
- If no file selected: keeps existing file
- If new file selected: replaces old file with new content
- File path updated in `putanja_fajla` field

## Error Handling

### Validation:
- Name and Author required in both modes
- File required only in upload mode
- Alert shown if validation fails

### Update Errors:
- Try-catch around `UpdateDocument()` call
- Error message shown to user: "Update failed: [error]"
- Upload progress reset on error

### Permission Errors:
- Each permission update wrapped in try-catch
- Continues with other users if one fails
- Logs errors to console

## Database Changes
**No schema changes required** - uses existing:
- `Dokumenti` table for document updates
- `DozvoleDokumenata` table for permissions
- `TagoviDokumenata` junction table for tags

## Testing Checklist

### Edit Mode:
- ✅ Load existing document data
- ✅ Pre-populate form fields
- ✅ Pre-select tags
- ✅ Pre-check user permissions
- ✅ Update without changing file
- ✅ Update with new file
- ✅ Update metadata only
- ✅ Add new permissions
- ✅ Remove existing permissions
- ✅ Show success message
- ✅ Redirect after update

### Upload Mode (unchanged):
- ✅ Require file selection
- ✅ Create new document
- ✅ Save permissions
- ✅ Show success message

## Security Considerations

1. **Authorization**: Backend should verify user has `moze_menjati` permission before allowing updates
2. **Validation**: Backend validates document exists before updating
3. **File Size**: Existing 80MB limit applies to replacement files
4. **SQL Injection**: Using parameterized queries in backend

## Future Enhancements

1. **Version History**: Show previous versions when editing
2. **Change Tracking**: Highlight which fields were modified
3. **Audit Log**: Record who edited what and when
4. **Collaborative Editing**: Lock document while being edited
5. **Draft Mode**: Save changes without publishing
6. **File Preview**: Show current file before replacing

## Related Files

### Frontend:
- `frontend/src/views/documents/DocumentAdd.vue` - Edit UI
- `frontend/src/views/documents/DocumentManagement.vue` - Edit button
- `frontend/src/router.js` - Route definition

### Backend:
- `main.go` - UpdateDocument, GetDocumentByID, etc.
- `backend/services/document_service.go` - Business logic
- `backend/models/models.go` - Data structures

### Database:
- `database/schema.sql` - Table definitions
- `DozvoleDokumenata` - Permissions table
- `Dokumenti` - Documents table
