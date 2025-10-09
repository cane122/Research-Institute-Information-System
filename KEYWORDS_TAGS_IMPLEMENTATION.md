# Keywords i Tags - Implementacija

## ✅ Šta je urađeno?

### 1. **Odvojeni Keywords i Tags u formi**

#### **Keywords (Slobodan unos)**
- Polje: Textarea za unos bilo čega
- Lokacija: Basic Information sekcija
- Dugme: "🤖 Generate Keywords" - AI generise slobodne keywords
- Baza: Kolona `kljucne_reci` u tabeli `Dokumenti`
- Format: String razdvojen zarezima

#### **Tags (Samo iz baze)**
- Polje: Multi-select dropdown
- Lokacija: Basic Information sekcija (posle Keywords)
- Dugme: "🤖 AI Select Tags" - AI **bira samo postojeće tagove iz baze**
- Baza: Tabela `Tagovi` + `DokumentTagovi` (many-to-many)
- Format: Array of tag IDs

### 2. **Uklonjen Labels iz Additional Metadata**
- Sada Additional Metadata sadrži samo:
  - ISO Number
  - Source URL

### 3. **Inteligentno AI biranje tagova**

AI generisanje tagova radi ovako:

```javascript
1. Učitavanje svih tagova iz baze (onMounted)
   └─> GetAllTags() -> availableTags[]

2. Klik na "🤖 AI Select Tags"
   └─> Šalje AI-u context sa dostupnim tagovima
       "Available tags in database: Research, Development, AI, Machine Learning..."
   
3. AI sugeriše tagove

4. Frontend match-uje sugestije sa bazom (case-insensitive)
   └─> Samo tagovi koji postoje se dodaju u selectedTags[]

5. Korisnik vidi selektovane tagove kao purple label-e
```

## 📂 Izmenjeni fajlovi

### Backend:
- ✅ `backend/models/models.go` - Dodato `KljucneReci` polje
- ✅ `backend/services/document_service.go` - Dodato `GetAllTags()` i ažuriran INSERT
- ✅ `main.go` - Eksponirana `GetAllTags()` funkcija

### Frontend:
- ✅ `frontend/src/views/documents/DocumentAdd.vue` - Potpuno reorganizovan UI
- ✅ `frontend/src/services/documentService.js` - Dodato `kljucne_reci` polje

### Database:
- ✅ `database/schema.sql` - Dodato `kljucne_reci TEXT`
- ✅ `add_keywords_column.ps1` - Migration script izvršen

## 🎨 UI Promene

### Keywords Polje:
```vue
<textarea v-model="documentInfo.keywords" rows="2">
  Enter keywords separated by commas
</textarea>
<button @click="generateKeywords">
  🤖 Generate Keywords
</button>
```

### Tags Polje:
```vue
<select multiple v-model="selectedTags" size="4">
  <option v-for="tag in availableTags" :value="tag.tag_id">
    {{ tag.naziv_taga }}
  </option>
</select>
<button @click="generateTags">
  🤖 AI Select Tags
</button>

<!-- Selected tags display -->
<div class="selected-tags-display">
  <div class="tag-label" v-for="tagId in selectedTags">
    {{ getTagName(tagId) }}
    <button @click="removeSelectedTag(tagId)">×</button>
  </div>
</div>
```

## 🔧 Nove JavaScript funkcije

### 1. **loadAvailableTags()** (onMounted)
```javascript
const { GetAllTags } = window.go.main.App
const tags = await GetAllTags()
availableTags.value = tags || []
```

### 2. **generateKeywords()**
```javascript
// Koristi GenerateTagsFromText za slobodne keywords
const tags = await GenerateTagsFromText(...)
documentInfo.value.keywords = tags.join(', ')
```

### 3. **generateTags()**
```javascript
// Kreira context sa dostupnim tagovima
const availableTagNames = availableTags.value.map(t => t.naziv_taga).join(', ')
const contextText = `Available tags in database: ${availableTagNames}`

// AI sugeriše
const aiSuggestedTags = await GenerateTagsFromText(contextText, ...)

// Match sa bazom (case-insensitive)
aiSuggestedTags.forEach(suggestedTag => {
  const matchedTag = availableTags.value.find(
    t => t.naziv_taga.toLowerCase() === suggestedTag.toLowerCase()
  )
  if (matchedTag) {
    selectedTags.value.push(matchedTag.tag_id)
  }
})
```

### 4. **handleSave()** - Ažurirano
```javascript
// Konvertuj tag IDs u nazive
const tagNames = selectedTags.value.map(tagId => {
  const tag = availableTags.value.find(t => t.tag_id === tagId)
  return tag ? tag.naziv_taga : null
}).filter(name => name !== null)

const documentData = {
  ...
  tags: tagNames,  // Tagovi iz baze
  keywords: documentInfo.value.keywords || '',  // Slobodni keywords
  ...
}
```

## 📊 Data Flow

```
Frontend                    Backend                    Database
========                    =======                    ========

Keywords Input          →   kljucne_reci (string)  →   Dokumenti.kljucne_reci
"AI, ML, research"                                      "AI, ML, research"

Tags Multi-select       →   tagovi ([]string)      →   1. Tagovi table
[1, 3, 5]                   ["Research", "AI"]          2. DokumentTagovi 
                                                            (many-to-many)
```

## 🎯 Kako testirati?

### Test 1: Keywords (slobodno)
1. Otvori Add Document
2. Unesi bilo šta u Keywords: "quantum physics, research, 2025"
3. Upload
4. Provera:
```sql
SELECT kljucne_reci FROM Dokumenti ORDER BY dokument_id DESC LIMIT 1;
```

### Test 2: Tags (iz baze)
1. Odaberi tagove iz multi-select dropdown-a
2. Upload
3. Provera:
```sql
SELECT t.naziv_taga 
FROM DokumentTagovi dt
JOIN Tagovi t ON dt.tag_id = t.tag_id
WHERE dt.dokument_id = (SELECT MAX(dokument_id) FROM Dokumenti);
```

### Test 3: AI Keywords generisanje
1. Unesi Document Name: "Machine Learning Research Paper"
2. Klik "🤖 Generate Keywords"
3. Trebalo bi da dobije: "machine learning, research, artificial intelligence, ..."

### Test 4: AI Tags selekcija
1. Unesi Document Name: "Software Development Project"
2. Klik "🤖 AI Select Tags"
3. AI će odabrati samo tagove koji postoje u bazi (npr. "Development", "Software")

## 🎨 Stilovi

### Purple gradient tag labels:
```css
.tag-label {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6px 12px;
  border-radius: 16px;
}
```

### AI buttons (Keywords):
```css
.btn-generate-keywords {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}
```

### Multi-select styling:
```css
.tags-multiselect option:checked {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}
```

## ✅ Kompletno!

- ✅ Keywords polje - slobodan unos
- ✅ Tags polje - samo iz baze
- ✅ AI generisanje keywords-a
- ✅ AI biranje postojećih tagova
- ✅ Backend API spremna
- ✅ Baza ažurirana
- ✅ Frontend build uspešan
- ✅ Wails bindings generisani

**Frontend je kompajliran i spreman za testiranje!** 🚀

Za pokretanje:
```powershell
cd "c:\Users\cane\Downloads\iis\Research Institute Information System"
wails dev
```
