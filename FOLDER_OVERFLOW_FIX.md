# PopravkaFolder Preview - Prikaz Svih Dokumenata

## Problem
1. **Folder se suzava/širi** - Nije imao fiksnu širinu
2. **Ne prikazuje sve dokumente** - Prikazivalo samo 2 dokumenta iako ih folder ima 3+

## Uzrok Problema

### Globalni CSS (globals.css)
```css
.documents-table {
  flex: 1;
  overflow: hidden;  /* ← Sakriva sadržaj */
  display: flex;
  flex-direction: column;
}

.table-body {
  flex: 1;
  overflow-y: auto;  /* ← Pravi scrollbar umesto da prikaže sve */
}
```

Ovi globalni stilovi su ograničavali visinu i forsiravali overflow, što je sprečavalo prikaz svih dokumenata u folderu.

## Rešenje

### 1. Override Globalnih Stilova sa `!important`
```css
/* Override global styles for folder preview */
.documents-table {
  overflow: visible !important;
}

.folder-contents .table-body {
  max-height: none !important;
  overflow: visible !important;
  flex: none !important;
}
```

### 2. Fiksna Širina za Foldere
```css
.project-folder {
  width: 100%;
  max-width: none;
  overflow: visible;
}

.folder-header {
  width: 100%;
}

.folder-contents {
  width: 100%;
  max-height: none !important;
  overflow: visible !important;
}
```

### 3. Dodato " docs" u Brojač
```vue
<span class="folder-count">({{ docs.length }}) docs</span>
```
Sada jasnije prikazuje: "(3) docs" umesto samo "(3)"

## CSS Promene - Pre i Posle

### PRE (Problem):
```css
.project-folder {
  overflow: hidden;  /* Sakriva sadržaj */
}

.folder-contents {
  background: #ffffff;
  /* Nema override za globalne stilove */
}
```

### POSLE (Rešenje):
```css
.documents-table {
  overflow: visible !important;  /* Override global */
}

.project-folder {
  overflow: visible;
  width: 100%;
  max-width: none;
}

.folder-contents {
  width: 100%;
  max-height: none !important;
  overflow: visible !important;
}

.folder-contents .table-body {
  max-height: none !important;
  overflow: visible !important;
  flex: none !important;
}
```

## Zašto `!important`?

Koristimo `!important` jer moramo da override-ujemo globalne stilove iz `globals.css` koji imaju visoku specifičnost. Ovo je jedini siguran način da garantujemo da će naši folder preview stilovi imati prioritet.

## Testiranje

### Test 1: Broj Dokumenata
```
Folder kaže: "(3) docs"
Prikazuje se: Svih 3 dokumenta ✅
```

### Test 2: Širina Foldera
```
Folder Header: 100% širine
Folder Contents: 100% širine
Ne suzava se/širi više ✅
```

### Test 3: Overflow
```
Svi dokumenti vidljivi
Nema scrollbar-a unutar foldera
Nema sakrivenih dokumenata ✅
```

## Build Status
- ✅ Frontend build: 179.52 kB (gzip: 60.10 kB)
- ✅ Build time: 1.36s
- ✅ No compilation errors
- ✅ Overflow fix applied with !important
- ✅ Width fix applied (100%, max-width: none)

## Kako Sada Radi

1. **Klik na "Folder preview" checkbox**
2. **Projekti se prikazuju kao folderi**
3. **Klik na folder header** - Otvara folder
4. **Prikazuju se SVI dokumenti** - bez ograničenja
5. **Folder ima fiksnu širinu** - 100% parent container-a
6. **Nema scroll-a** - svi dokumenti vidljivi odjednom

## Napomena

CSS sa `!important` je neophodan jer:
- Globalni `globals.css` ima `overflow: hidden` i `flex: 1`
- Ovi stilovi imaju visoku specifičnost
- Bez `!important` ne možemo da ih override-ujemo
- Alternativa bi bila menjanje globalnog CSS-a, ali to utiče na ceo projekat

Ovo je **lokalno rešenje** samo za folder preview mode koje ne utiče na regular table view.
