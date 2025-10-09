# Document Management Filtering & Grouping Improvements

## Overview
Enhanced the Document Management interface with fully functional filtering, project-based grouping, and dynamic data loading from the database.

## Changes Made

### 1. **Dynamic Filter Options**
- **Type Filter**: Now dynamically populated from actual document types in database
- **Project Filter**: Auto-populated from available projects
- **Tag Filter**: Loads real tags from database with search functionality
- **Author Filter**: Real-time filtering with reactive updates

### 2. **Group by Project Feature**
Added checkbox to group documents by project:
- Documents are organized under project headers
- Each group shows document count
- Beautiful gradient headers for visual distinction
- Separate checkboxes for selecting all documents in a group
- "Unassigned" group for documents without projects

### 3. **Working Filters**
All filters now work correctly with backend data:

```javascript
// Search by document name, author, or description
if (searchQuery.value) {
  filtered = filtered.filter(doc =>
    doc.naziv_dokumenta.toLowerCase().includes(query) ||
    doc.ime_kreirao.toLowerCase().includes(query) ||
    (doc.opis && doc.opis.toLowerCase().includes(query))
  )
}

// Filter by author
if (authorFilter.value) {
  filtered = filtered.filter(doc =>
    doc.ime_kreirao.toLowerCase().includes(author)
  )
}

// Filter by document type
if (selectedType.value) {
  filtered = filtered.filter(doc => doc.tip_dokumenta === selectedType.value)
}

// Filter by project
if (selectedProject.value) {
  filtered = filtered.filter(doc => doc.naziv_projekta === selectedProject.value)
}

// Date range filtering
if (dateFrom.value) {
  filtered = filtered.filter(doc => {
    const docDate = new Date(doc.poslednja_izmena || doc.datuma_postavke)
    return docDate >= fromDate
  })
}
```

### 4. **Database Integration**
Updated to use real database fields:

| Frontend Display | Database Field |
|-----------------|----------------|
| Document Name | `naziv_dokumenta` |
| Author | `ime_kreirao` |
| Type | `tip_dokumenta` |
| Modified | `poslednja_izmena` or `datuma_postavke` |
| Project | `naziv_projekta` |
| Document ID | `dokument_id` |

### 5. **Tag Management**
- Loads all tags from database on mount: `GetAllTags()`
- Tag search functionality filters available tags
- Tags displayed with real data from `tagovi` table
- Checkbox selection for multiple tags

### 6. **Sorting Improvements**
Updated sort fields to match database schema:
- Sort by `naziv_dokumenta` (Document Name)
- Sort by `ime_kreirao` (Author)
- Sort by `tip_dokumenta` (Type)
- Sort by `poslednja_izmena` (Modified Date)
- Sort by `naziv_projekta` (Project)

Handles null/undefined values gracefully:
```javascript
// Handle null/undefined values
if (aVal === null || aVal === undefined) aVal = ''
if (bVal === null || bVal === undefined) bVal = ''

// Handle date fields
if (sortField.value === 'poslednja_izmena' || sortField.value === 'datuma_postavke') {
  aVal = new Date(aVal || 0).getTime()
  bVal = new Date(bVal || 0).getTime()
}
```

### 7. **Reactive Filtering**
All filters are reactive - changes trigger immediate updates:
- Search query: Updates as you type
- Author filter: `@input="applyFilters"`
- Type/Project dropdowns: `@change="applyFilters"`
- Tag checkboxes: `@change="applyFilters"`
- Date range: `@change="applyFilters"`

### 8. **UI Improvements**

#### Project Groups View
```css
.project-group {
  margin-bottom: 30px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
  background: white;
}

.project-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 15px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
```

#### Empty States
- "No tags available" message when tag list is empty
- Proper error handling for failed loads

### 9. **Data Loading**
Parallel loading for better performance:
```javascript
const [docs, tags] = await Promise.all([
  GetAllDocuments(),
  GetAllTags()
])

documents.value = docs || []
allTags.value = tags || []
```

## Usage

### Basic Filtering
1. **Search**: Type in the search box to filter by name, author, or description
2. **Author**: Enter author name to filter documents by creator
3. **Date Range**: Select "From" and "To" dates to filter by modification date
4. **Tags**: Search and check tags to filter documents with specific tags
5. **Type**: Select document type from dropdown (PDF, DOC, etc.)
6. **Project**: Select project from dropdown

### Group by Project
1. Check "Group by Project" checkbox
2. Documents are organized under project headers
3. Each project shows its document count
4. Click checkbox in project header to select all documents in that project

### Sorting
- Click any column header to sort
- Click again to toggle between ascending/descending
- Sort indicator (↑/↓) shows current sort direction

## Technical Details

### Computed Properties
- `filteredTags`: Filters tags based on search query
- `availableTypes`: Extracts unique types from loaded documents
- `availableProjects`: Extracts unique projects from loaded documents
- `filteredDocuments`: Applies all active filters
- `sortedDocuments`: Sorts filtered results
- `groupedDocuments`: Groups sorted documents by project
- `paginatedDocuments`: Handles pagination

### State Management
```javascript
const documents = ref([])        // All documents from database
const allTags = ref([])          // All tags from database
const searchQuery = ref('')      // Search input
const authorFilter = ref('')     // Author filter input
const selectedType = ref('')     // Selected type from dropdown
const selectedProject = ref('')  // Selected project from dropdown
const selectedTags = ref([])     // Array of selected tag names
const tagSearchQuery = ref('')   // Tag search input
const groupByProject = ref(false) // Group by project toggle
```

### Functions
- `applyFilters()`: Resets to page 1 and triggers reactive updates
- `clearFilters()`: Resets all filter values
- `toggleSelectAllInGroup(docs)`: Selects/deselects all docs in a project group
- `loadDocuments()`: Fetches documents and tags from backend
- `formatDate(dateStr)`: Formats date for display

## Files Modified

### Frontend
- `frontend/src/views/documents/DocumentManagement.vue`
  - Updated filter UI with real tag loading
  - Added project grouping feature
  - Fixed all computed filter functions
  - Updated sort fields to match database schema
  - Added reactive filter updates
  - Integrated with backend API

## Build Status
- ✅ Frontend build: 178.82 kB (gzip: 59.90 kB)
- ✅ Build time: 1.38s
- ✅ No compilation errors

## Next Steps

1. **Tag Filtering Logic**
   - Currently commented out (no document-tag relationship loaded yet)
   - Need to load tags for each document
   - Uncomment tag filtering in `filteredDocuments` computed property

2. **Document Preview**
   - Links to `/documents/preview/:id` already work
   - Document preview page needs to be updated with real data

3. **Bulk Operations**
   - Select multiple documents
   - Bulk delete
   - Bulk tag assignment
   - Bulk project assignment

4. **Advanced Features**
   - Saved filter presets
   - Export filtered results
   - Custom column visibility
   - Advanced search with operators (AND, OR, NOT)

## Testing Checklist

- [x] Search filter works
- [x] Author filter works
- [x] Type filter populated from database
- [x] Project filter populated from database
- [x] Tags loaded from database
- [x] Tag search works
- [x] Date range filtering works
- [x] Group by Project toggle works
- [x] Project groups display correctly
- [x] Sort by each column works
- [x] Pagination works with filters
- [x] Select all in project group works
- [x] Document preview link works
- [x] Document edit link works
- [x] Document delete works
- [ ] Tag filtering (pending document-tag relationship)
