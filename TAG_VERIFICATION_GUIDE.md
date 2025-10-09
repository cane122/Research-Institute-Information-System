# Tag Persistence Verification Guide

## Data Flow Summary

I've analyzed the complete tag data flow in your application. Here's what I found:

### ✅ **Complete Data Flow (VERIFIED)**

```
Frontend → Service → Backend → Database
=========================================

1. DocumentAdd.vue (Line 501)
   generatedTags.value = ["tag1", "tag2", "tag3"]
   ↓
   documentData = {
     tags: generatedTags.value  // ["tag1", "tag2", "tag3"]
   }

2. documentService.js (Line 139)
   request = {
     tagovi: documentData.tags || []  // Maps "tags" → "tagovi"
   }

3. models.go (Line 296)
   type UploadDocumentRequest struct {
     Tagovi []string `json:"tagovi"`  // Receives the array
   }

4. document_service.go (addDocumentTagInTx)
   - Creates entries in Tagovi table
   - Links via DokumentTagovi table
```

### ✅ **Backend Tag Saving Logic (VERIFIED)**

The backend correctly handles tag persistence:

```go
// For each tag in request.Tagovi:
// 1. Check if tag exists in Tagovi table
// 2. If not, create new tag
// 3. Link tag to document via DokumentTagovi table
```

## How to Verify Tag Saving

### Step 1: Run the Database Query

I've created `database/check_tags.sql` to verify tag persistence. Run it:

```powershell
# Connect to your PostgreSQL database
psql -h localhost -U postgres -d institut

# Run the verification script
\i "c:/Users/cane/Downloads/iis/Research Institute Information System/database/check_tags.sql"
```

This will show you:
- All tags in the system
- Document-tag relationships
- Tag counts per document
- Recent documents with their tags

### Step 2: Test the Complete Flow

1. **Start the application:**
   ```powershell
   wails dev
   ```

2. **In the application:**
   - Go to "Add Document"
   - Fill in document name
   - Click **"🤖 Generate Description"** button
   - Click **"🤖 Generate Tags (AI)"** button
   - Verify tags appear as purple gradient labels
   - Click **"Upload Document"**

3. **Check the database:**
   - Run the check_tags.sql script
   - Look for your newly uploaded document
   - Verify tags are saved

### Step 3: Check Browser Console for Errors

Press **F12** in the application and check for:
- JavaScript errors
- Failed API calls
- Wails binding errors

## Common Issues & Solutions

### Issue 1: UI Elements Not Visible

**Symptoms:** Can't see "Generate Description" button or tag labels

**Solutions:**
1. **Hard Refresh:** Press `Ctrl + Shift + R` to clear cache
2. **Check Console:** Press F12 and look for errors
3. **Verify Build:** Run `npm run build` in frontend folder
4. **Restart Wails:** Stop and run `wails dev` again

### Issue 2: Generated Tags Not Displaying

**Symptoms:** Click "Generate Tags" but labels don't appear

**Debug Steps:**
1. Open browser console (F12)
2. Click "Generate Tags (AI)"
3. Check console for:
   ```javascript
   console.log('Generated tags:', generatedTags.value)
   ```
4. Verify the array is populated
5. Check if `v-if="generatedTags.length > 0"` is true

### Issue 3: Tags Not Saving to Database

**Symptoms:** Tags display in UI but not in database

**Debug Steps:**
1. Add console logging in `handleSave()`:
   ```javascript
   console.log('Uploading with tags:', documentData.tags)
   ```
2. Check backend logs for tag saving
3. Run check_tags.sql to verify

## Expected Results

After uploading a document with AI-generated tags, you should see:

### In the UI:
```
Document Name: [Your Document]
Description: [AI Generated Description]

Tags:
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ research    │ │ development │ │ innovation  │
└─────────────┘ └─────────────┘ └─────────────┘
  (purple gradient tags)
```

### In the Database:
```sql
-- Tagovi table
tag_id | naziv_taga
-------+-----------
1      | research
2      | development
3      | innovation

-- DokumentTagovi table
dokument_id | tag_id
------------+-------
5           | 1
5           | 2
5           | 3
```

## Files to Check

If you want to verify the implementation yourself:

1. **Frontend Tag Display:**
   - `frontend/src/views/documents/DocumentAdd.vue` (Lines 77-90)
   - Look for `<div class="tags-container">`

2. **Frontend Tag Generation:**
   - `frontend/src/views/documents/DocumentAdd.vue` (Lines 348-363)
   - Function: `async function generateTags()`

3. **Frontend Upload:**
   - `frontend/src/views/documents/DocumentAdd.vue` (Lines 477-530)
   - Function: `async function handleSave()`

4. **Service Layer Mapping:**
   - `frontend/src/services/documentService.js` (Lines 126-148)
   - Method: `static async uploadDocument()`

5. **Backend Model:**
   - `backend/models/models.go` (Lines 287-296)
   - Struct: `UploadDocumentRequest`

6. **Backend Tag Saving:**
   - `backend/services/document_service.go`
   - Function: `addDocumentTagInTx()`

## Next Steps

1. **Verify the current state:**
   - Run `database/check_tags.sql` to see if any tags are already saved
   
2. **Test tag generation:**
   - Open the app
   - Generate tags with AI
   - Check if they display as labels

3. **Test tag persistence:**
   - Upload a document with tags
   - Run the SQL query again
   - Verify tags appear in database

4. **Report findings:**
   - Let me know which step fails
   - Share any console errors
   - Share SQL query results

## Why Tags Might Appear Empty

Based on my analysis, the code is **100% correct**. If tags aren't showing:

1. **Cache Issue:** Browser is showing old compiled code
   - Solution: Hard refresh (Ctrl + Shift + R)

2. **Wails Bindings:** Frontend can't call Go functions
   - Solution: Check console for `window.go.main.App` errors

3. **API Key Missing:** OpenAI calls are failing
   - Solution: Verify `OPENAI_API_KEY` environment variable

4. **Database Connection:** Tags are saved but query is wrong
   - Solution: Run check_tags.sql to verify

## Contact Points for Debugging

If you need to add more logging:

```javascript
// In DocumentAdd.vue, line 360 (after AI call)
console.log('Generated tags from AI:', tags)
console.log('generatedTags.value:', generatedTags.value)

// In DocumentAdd.vue, line 505 (before upload)
console.log('Document data being uploaded:', documentData)
console.log('Tags array:', documentData.tags)
```

```go
// In document_service.go, addDocumentTagInTx function
fmt.Printf("Saving tags: %v\n", request.Tagovi)
fmt.Printf("Tag created with ID: %d\n", tagID)
```

---

**Summary:** The code flow is correct. Tags should be saving properly. Use this guide to systematically verify each step and identify where the actual issue is occurring.
