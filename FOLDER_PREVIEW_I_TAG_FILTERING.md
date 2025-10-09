# Folder Preview i Filtriranje po Tagovima

## Implementirane Promene

### 1. **Folder Preview Umesto Grupisanja**
Zamenjen sistem grupisanja po projektima sa **folder preview** sistemom gde:
- Projekti su prikazani kao folderi koji mogu da se otvaraju i zatvaraju
- Dokumenti su stavke unutar tih foldera
- Klik na folder header otvara/zatvara folder
- Vizuelna animacija pri otvaranju/zatvaranju

### 2. **Popravljena Pretraga po Tagovima**
Sada radi potpuno funkcionalno:
- Učitavaju se tagovi za svaki dokument
- Filtriranje po tagovima radi u realnom vremenu
- Možete selektovati više tagova odjednom
- Prikazuju se samo dokumenti koji imaju bar jedan od selektovanih tagova

## Kako Radi Folder Preview

### Vizuelni Prikaz
```
📁 AI Research (5)                    ▶  <- Klikni da otvoriš
📂 Climate Studies (3)                ▼  <- Folder otvoren
  ├── 📄 Climate Report 2024
  ├── 📄 Temperature Analysis
  └── 📄 Data Collection Methods
📁 Unassigned (2)                     ▶
```

### Funkcionalnost
1. **Klik na folder header** - Otvara/zatvara folder
2. **Ikonica foldera** - 📁 (zatvoren) ili 📂 (otvoren)
3. **Strelica** - ▶ (zatvoren) ili ▼ (otvoren)
4. **Brojač dokumenata** - Prikazuje koliko dokumenata ima u folderu
5. **Hover efekat** - Folder se pomera udesno i menja boju

### Interakcija
- Kada je folder **zatvoren**: Vide se samo header sa nazivom projekta i brojem dokumenata
- Kada je folder **otvoren**: Prikazuju se svi dokumenti u tabeli sa svim podacima
- **Select all checkbox** - U svakom folderu možete selektovati sve dokumente odjednom

## Filtriranje po Tagovima

### Kako Radi
1. **Učitavanje tagova**: Pri učitavanju dokumenata, za svaki dokument se učitavaju i njegovi tagovi
```javascript
const docsWithTags = await Promise.all(
  docs.map(async (doc) => {
    const docTags = await GetDocumentTags(doc.dokument_id)
    return { ...doc, tagovi: docTags || [] }
  })
)
```

2. **Filtriranje**: Dokumenti se filtriraju ako imaju bar jedan tag koji je selektovan
```javascript
if (selectedTags.value.length > 0) {
  filtered = filtered.filter(doc => {
    return doc.tagovi && doc.tagovi.some(tag => 
      selectedTags.value.includes(tag.naziv_taga)
    )
  })
}
```

3. **Pretraga tagova**: Možete pretražiti tagove u sidebaru pre nego što ih selektujete

### Primer Upotrebe
1. U sidebaru, sekcija "Tags"
2. Kucajte u search box da filtrirate dostupne tagove
3. Checkboxom selektujte tagove koje želite
4. Dokumenti se automatski filtriraju da prikažu samo one sa selektovanim tagovima

## Tehnički Detalji

### State Management
```javascript
const expandedProjects = ref([])  // Lista otvorenih projekata
const selectedTags = ref([])      // Lista selektovanih tagova
const showFolderPreview = ref(false)  // Toggle za folder view
```

### Funkcije

#### toggleProjectFolder(projectName)
Otvara/zatvara folder projekta:
```javascript
function toggleProjectFolder(projectName) {
  const index = expandedProjects.value.indexOf(projectName)
  if (index > -1) {
    expandedProjects.value.splice(index, 1)  // Zatvori
  } else {
    expandedProjects.value.push(projectName)  // Otvori
  }
}
```

#### loadDocuments()
Učitava dokumente sa tagovima:
```javascript
async function loadDocuments() {
  // 1. Učitaj sve dokumente i tagove
  const [docs, tags] = await Promise.all([
    GetAllDocuments(),
    GetAllTags()
  ])
  
  // 2. Za svaki dokument učitaj njegove tagove
  const docsWithTags = await Promise.all(
    docs.map(async (doc) => {
      const docTags = await GetDocumentTags(doc.dokument_id)
      return { ...doc, tagovi: docTags || [] }
    })
  )
  
  documents.value = docsWithTags
}
```

### CSS Animacije

#### Slide Down Animacija
Kada se folder otvara:
```css
@keyframes slideDown {
  from {
    opacity: 0;
    max-height: 0;
  }
  to {
    opacity: 1;
    max-height: 2000px;
  }
}
```

#### Hover Efekat
```css
.folder-header:hover {
  background: linear-gradient(135deg, #5568d3 0%, #6a3f8f 100%);
  transform: translateX(5px);  /* Pomera se udesno */
}
```

## UI/UX Poboljšanja

### 1. Vizuelna Hijerarhija
- **Folder header**: Gradient purple pozadina
- **Dokumenti**: Bela pozadina sa purple bordurom sleva
- **Indent**: Dokumenti su uvučeni 20px udesno
- **Hover**: Dokumenti menjaju boju na hover

### 2. Interaktivnost
- Cursor pointer na folder header
- Smooth animacije (0.3s ease)
- Visual feedback na klik
- Scale efekat na ikonicu foldera

### 3. Responsive Design
- Folder se prilagođava širini ekrana
- Dokumenti u tabeli imaju fiksne kolone
- Actions dugmići ostaju vidljivi

## Kombinovanje Filtera

Svi filteri rade zajedno:
1. **Search** - Pretraga po nazivu, autoru, opisu
2. **Author** - Filter po autoru
3. **Tags** - Filter po tagovima ✅ SADA RADI
4. **Type** - Filter po tipu dokumenta
5. **Project** - Filter po projektu
6. **Date Range** - Filter po datumu

Primer: Možete pretražiti sve PDF dokumente od autora "Marko" sa tagom "AI Research" kreirane u poslednjih 30 dana.

## Build Status
- ✅ Frontend build: 179.66 kB (gzip: 60.11 kB)
- ✅ Build time: 1.39s
- ✅ No compilation errors
- ✅ CSS animations included
- ✅ Tag filtering working
- ✅ Folder preview working

## Upotreba

### Aktiviranje Folder Preview
1. U "Document Management" stranici
2. Checkbox "Folder preview" (gore desno)
3. Dokumenti se grupišu po projektima
4. Kliknite na folder da ga otvorite/zatvorite

### Filtriranje po Tagovima
1. U sidebaru, sekcija "Tags"
2. Pretražite tagove (search box)
3. Selektujte željene tagove checkboxom
4. Dokumenti se automatski filtriraju

### Kombinovanje sa Drugim Filterima
- Možete koristiti folder preview zajedno sa svim filterima
- Filteri se primenjuju pre grupisanja u foldere
- Folderi bez dokumenata (nakon filtriranja) se NE prikazuju

## Sledeći Koraci

- [ ] Dodati animaciju za broj dokumenata u folderu
- [ ] Dodati ikonu za "Expand All" / "Collapse All" foldere
- [ ] Zapamtiti stanje otvorenih foldera (localStorage)
- [ ] Dodati drag & drop za premeštanje dokumenata između projekata
- [ ] Bulk select sa shift+click
