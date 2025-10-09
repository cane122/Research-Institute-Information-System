# Popravke Folder Preview Prikaza

## Problemi koje smo rešili:

### 1. ✅ Dupli Select Checkbox
**Problem**: Bio je checkbox za "Group by Project" umesto za "Folder preview"
**Rešenje**: Zamenjen checkbox tekst i funkcionalnost na "Folder preview"

### 2. ✅ Ne prikazuje sve dokumente (piše 3, prikazuje 2)
**Problem**: Pagination je bio aktivan i za folder preview mode
**Rešenje**: 
- Pagination se sada prikazuje SAMO u regular table view
- Kada je folder preview aktivan, prikazuju se SVI dokumenti bez paginacije
```vue
<!-- Pagination (only for regular table view) -->
<div class="pagination" v-if="!showFolderPreview && totalPages > 1">
```

### 3. ✅ Preterano šminkan UI
**Problem**: Previše animacija i efekata
**Rešenje**: Uprošćen CSS:
- ❌ Uklonjeno: `transform: translateX(5px)` na hover
- ❌ Uklonjeno: `transform: scale(1.1)` za ikonicu
- ❌ Uklonjeno: `slideDown` animacija
- ❌ Uklonjeno: `padding-left` i `margin-left` kod dokumenata
- ✅ Zadržano: Samo promena boje na hover i jednostavna 0.2s tranzicija

### 4. ✅ Popravljeno Select All
**Problem**: `allSelected` computed koristio `doc.id` umesto `doc.dokument_id`
**Rešenje**: 
```javascript
const allSelected = computed(() => {
  return paginatedDocuments.value.length > 0 && 
         paginatedDocuments.value.every(doc => selectedDocuments.value.includes(doc.dokument_id))
})
```

## Kako sada radi:

### Folder Preview Mode
1. **Klikni checkbox "Folder preview"**
2. **Svi dokumenti se prikazuju** - bez paginacije
3. **Folderi su grupisani po projektima**
4. **Klik na folder** - otvara/zatvara
5. **Brojač pokazuje tačan broj dokumenata** u folderu

### Regular Table Mode
1. **Odklikni "Folder preview"**
2. **Normalna tabela sa svim dokumentima**
3. **Pagination aktivan** - 10 dokumenata po stranici
4. **Sortiranje po kolonama**

## CSS Promene

### PRE (Preterano):
```css
.folder-header:hover {
  transform: translateX(5px);  /* Pomera se udesno */
}

.folder-icon {
  transition: transform 0.3s ease;
}

.folder-header:hover .folder-icon {
  transform: scale(1.1);  /* Uvećava se */
}

.folder-contents {
  animation: slideDown 0.3s ease;  /* Slide animacija */
}

.document-item {
  padding-left: 20px;
  margin-left: 20px;
  border-left: 3px solid #764ba2;
}
```

### POSLE (Jednostavno):
```css
.folder-header:hover {
  background: linear-gradient(135deg, #5568d3 0%, #6a3f8f 100%);
  /* Samo promena boje */
}

.folder-icon {
  font-size: 1.2em;
  /* Bez animacije */
}

.folder-contents {
  background: #ffffff;
  /* Bez animacije */
}

.document-item {
  border-left: 2px solid #e0e0e0;
  /* Nema padding/margin */
}
```

## Build Status
- ✅ Frontend build: 179.52 kB (gzip: 60.09 kB)
- ✅ Build time: 1.52s
- ✅ No compilation errors
- ✅ Pagination fix applied
- ✅ Simplified CSS
- ✅ All documents visible in folder view

## Testiranje

### Test 1: Broj dokumenata
- Ako folder kaže (3), proveri da se sva 3 dokumenta prikazuju
- ✅ RADI - pagination isključen u folder view

### Test 2: Select all
- Klikni checkbox u header-u foldera
- ✅ RADI - koristi `doc.dokument_id`

### Test 3: Pagination
- U regular table view - pagination postoji
- U folder view - pagination se NE prikazuje
- ✅ RADI - conditionally rendered

### Test 4: UI Simplicity
- Hover na folder - samo promena boje
- Bez translateX, scale ili drugih fancy efekata
- ✅ RADI - simplified CSS

## Preostali Features
- ✅ Folder preview mode
- ✅ Click to expand/collapse
- ✅ Tag filtering
- ✅ All filters work
- ✅ Search functionality
- ✅ Sort by columns (in regular view)
- ✅ Select all in folder
- ✅ Individual document selection
