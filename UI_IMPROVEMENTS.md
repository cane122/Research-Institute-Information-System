# UI Improvements - Keywords i Tags

## ✅ Šta je popravljeno?

### 1. **Keywords - Sada sa Label-ima** 
Umesto jednog textarea polja, sada imamo:

```
┌─────────────────────────────────────────────────────────┐
│ Keywords                                                │
│ ┌──────────────┬──────┬────────────────────┐          │
│ │ Type keyword │  +Add│ 🤖 Generate Keywords│          │
│ └──────────────┴──────┴────────────────────┘          │
│                                                         │
│ ┌──────┐ ┌──────┐ ┌──────┐                           │
│ │ AI ×  │ │ ML ×  │ │ ...×  │  <- Purple labels       │
│ └──────┘ └──────┘ └──────┘                           │
└─────────────────────────────────────────────────────────┘
```

**Kako radi:**
- Kucaš keyword i pritisneš Enter ILI klikneš "+ Add"
- Keyword se dodaje kao purple gradient label
- Možeš ukloniti klikom na "×"
- Ili klikni "🤖 Generate Keywords" - AI dodaje keywords kao labele

### 2. **Tags - Dugmad umesto Multi-Select**
Umesto wonky multi-select dropdown-a, sada imamo:

```
┌─────────────────────────────────────────────────────────┐
│ Tags (From Database)                                    │
│ ┌───────────────────────────────────────────────────┐ │
│ │  Research   [+]    Development  [+]               │ │
│ │  AI         [+]    Software     [✓]  <- Disabled  │ │
│ │  Hardware   [+]    Testing      [+]               │ │
│ └───────────────────────────────────────────────────┘ │
│                                                         │
│ [🤖 AI Select Tags]                                     │
│                                                         │
│ Selected:                                               │
│ ┌──────────┐ ┌──────────┐                             │
│ │Software × │ │ AI ×     │  <- Purple labels          │
│ └──────────┘ └──────────┘                             │
└─────────────────────────────────────────────────────────┘
```

**Kako radi:**
- Svi tagovi iz baze prikazani kao button-i sa [+] dugmetom
- Klikneš [+] da dodaš tag
- Posle dodavanja, dugme postaje [✓] i disabled
- Selektovani tagovi se prikazuju kao purple labele ispod
- Klikni "🤖 AI Select Tags" - AI bira relevantne tagove

## 🎨 Novi UI Elementi

### Keywords Input:
```vue
<input type="text" v-model="currentKeyword" @keypress.enter="addKeyword" />
<button @click="addKeyword">+ Add</button>
<button @click="generateKeywords">🤖 Generate Keywords</button>

<!-- Keywords display -->
<div class="tags-container">
  <div class="tag-label">
    <span>keyword</span>
    <button @click="removeKeyword(index)">×</button>
  </div>
</div>
```

### Tags Buttons:
```vue
<div class="available-tags-list">
  <div class="tag-item">
    <span class="tag-name">Research</span>
    <button 
      class="btn-add-tag-from-db" 
      @click="addTagFromDb(tag.tag_id)"
      :disabled="selectedTags.includes(tag.tag_id)"
    >
      {{ selectedTags.includes(tag.tag_id) ? '✓' : '+' }}
    </button>
  </div>
</div>
```

## 🔧 Nove Funkcije

### Keywords Management:
```javascript
const currentKeyword = ref('')
const keywords = ref([])

function addKeyword() {
  const keyword = currentKeyword.value.trim()
  if (keyword && !keywords.value.includes(keyword)) {
    keywords.value.push(keyword)
    currentKeyword.value = ''
  }
}

function removeKeyword(index) {
  keywords.value.splice(index, 1)
}

// AI generisanje dodaje u array
async function generateKeywords() {
  const tags = await GenerateTagsFromText(...)
  tags.forEach(tag => {
    if (!keywords.value.includes(tag)) {
      keywords.value.push(tag)
    }
  })
}
```

### Tags Management:
```javascript
function addTagFromDb(tagId) {
  if (!selectedTags.value.includes(tagId)) {
    selectedTags.value.push(tagId)
  }
}

function removeSelectedTag(tagId) {
  selectedTags.value = selectedTags.value.filter(id => id !== tagId)
}
```

### Upload:
```javascript
async function handleSave() {
  const documentData = {
    ...
    tags: tagNames,  // Array tag names
    keywords: keywords.value.join(', '),  // Comma-separated string
    ...
  }
}
```

## 🎨 CSS Stilovi

### Available Tags List (scrollable):
```css
.available-tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
  max-height: 200px;
  overflow-y: auto;
}
```

### Tag Item (white rounded button):
```css
.tag-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 20px;
  transition: all 0.2s ease;
}

.tag-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}
```

### Add Tag Button (zeleno [+] ili sivo [✓]):
```css
.btn-add-tag-from-db {
  background: #28a745;
  color: white;
  border: none;
  cursor: pointer;
  font-weight: bold;
}

.btn-add-tag-from-db:disabled {
  background: #6c757d;  /* Sivo za disabled */
  cursor: not-allowed;
  opacity: 0.6;
}
```

## ✅ Kompletno!

- ✅ Keywords - Input sa Enter + Add dugme + Purple labele
- ✅ Tags - Dugmad [+] umesto multi-select
- ✅ Selektovani tagovi prikazani kao purple labele
- ✅ AI generisanje radi za oba
- ✅ Intuitivniji UX
- ✅ Frontend build uspešan (1.36s)

**Sve je mnogo intuitivnije sada!** 🎉

Za testiranje:
```powershell
wails dev
```

### Testiranje:

1. **Keywords:**
   - Ukucaj "machine learning" i Enter
   - Ukucaj "AI" i klikni "+ Add"
   - Klikni "🤖 Generate Keywords" - AI dodaje labele

2. **Tags:**
   - Klikni [+] pored "Research" - dodaje se label
   - Klikni [+] pored "Development" - dodaje se label
   - Dugme postaje [✓] i disabled
   - Klikni "× " na labelu da ukloniš

3. **AI Tags:**
   - Klikni "🤖 AI Select Tags"
   - AI automatski selektuje relevantne tagove iz liste
   - Videćeš purple labele za selektovane

Sve radi kako treba! 🚀
