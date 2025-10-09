# 🤖 Automatsko Tagovanje Dokumenata - Uputstvo

## Pregled

Sistem omogućava automatsko generisanje tagova (ključnih reči) i opisa za dokumente koristeći OpenAI GPT-3.5 Turbo AI model. Ove funkcionalnosti pomažu u bržoj kategorizaciji i pretrazi dokumenata.

## Kako Funkcioniše

### Auto-Tagging (Tagovanje)
Sistem analizira:
- **Naziv dokumenta** - Glavni naslov dokumenta
- **Tip dokumenta** - Research Paper, PDF, Report, itd.
- **Opis dokumenta** - Detaljni opis sadržaja

Na osnovu ovih informacija, AI model generiše 5 relevantnih tagova koji najbolje opisuju dokument.

### Auto-Description (Generisanje Opisa)
Sistem analizira:
- **Naziv dokumenta** - Glavni naslov dokumenta
- **Tip dokumenta** - Research Paper, PDF, Report, itd.
- **Ime fajla** - Naziv uploadovanog fajla

AI model generiše profesionalan opis od 2-3 rečenice koji objašnjava sadržaj i svrhu dokumenta.

## Korišćenje

### 1. Otvorite Formu za Upload Dokumenta

Kliknite na **"Upload Document"** ili navigirajte na `/documents/add`

### 2. Popunite Osnovne Informacije

```
Document Name: Machine Learning in Medical Imaging
Document Type: Research Paper
```

### 3a. Generišite Opis (NOVO! 🎉)

Kliknite **"🤖 Generate Description"** dugme

Sistem će automatski generisati opis:
```
Description: This research paper explores the application of machine learning 
algorithms in medical imaging analysis. It focuses on automated diagnosis and 
pattern recognition in radiological scans, providing insights into current 
methodologies and future developments in the field.
```

### 3b. Generišite Tagove

Nakon što imate opis, kliknite **"🤖 Generate Tags (AI)"**

Dugme će pokazati:
- **🤖 Generate Tags (AI)** - Kada je spremno
- **🤖 Generating...** - Tokom generisanja (1-3 sekunde)

### 4. Rezultat - Tag Labels

Tagovi će biti prikazani kao vizuelni label-ovi sa mogućnošću uklanjanja:

```
[machine learning ×] [medical imaging ×] [diagnosis ×] [neural networks ×] [radiology ×]
```

## Primeri Korišćenja

### Primer 1: Naučni Rad (Sa Auto-Description)
**Input:**
- Name: "Blockchain Technology in Supply Chain Management"
- Type: Research Paper

**Korak 1 - Generate Description:**
```
Description: This research paper examines the implementation of blockchain 
technology in modern supply chain management systems. It analyzes the benefits 
of distributed ledger technology for tracking products, ensuring transparency, 
and reducing fraud across global supply networks.
```

**Korak 2 - Generate Tags:**
```
Tags: [blockchain ×] [supply chain ×] [distributed ledger ×] [transparency ×] [tracking ×]
```

### Primer 2: Tehnički Dokument
**Input:**
- Name: "REST API Documentation v2.0"
- Type: DOC
- File: "api-docs-v2.pdf"

**Generated Description:**
```
This technical documentation provides comprehensive guidance for implementing 
and using the REST API version 2.0. It includes endpoint specifications, 
authentication methods, and example code snippets for developers.
```

**Generated Tags:**
```
Tags: [REST API ×] [documentation ×] [endpoints ×] [authentication ×] [developer guide ×]
```

### Primer 3: Finansijski Izveštaj
**Input:**
- Name: "Q3 2025 Financial Report"
- Type: Report
- Description: "Quarterly financial analysis and performance metrics"

**Generated Tags:**
```
financial report, Q3 2025, analysis, performance metrics, quarterly
```

## Najbolje Prakse

### ✅ DO - Preporuke

1. **Napišite detaljan opis**
   ```
   ❌ Loše: "Document about AI"
   ✅ Dobro: "Comprehensive analysis of artificial intelligence applications 
              in healthcare, focusing on diagnostic imaging and patient care"
   ```

2. **Koristite specifičan naziv**
   ```
   ❌ Loše: "Report 2025"
   ✅ Dobro: "Annual Cybersecurity Threat Assessment Report 2025"
   ```

3. **Izaberite odgovarajući tip dokumenta**
   - Research Paper → za naučne radove
   - Report → za izveštaje
   - PDF/DOC → za opšte dokumente

### ❌ DON'T - Izbegavati

1. **Nemojte koristiti AI bez opisa**
   - Opis je ključan za kvalitetne tagove
   - Minimum 10-20 reči u opisu

2. **Nemojte koristiti generičke nazive**
   - "Document1", "File", "New" → loši nazivi
   - Koristite smislene i opisne nazive

## Troubleshooting

### Problem: "Please enter a document name first!"

**Rešenje:** Prvo popunite polje "Document Name" pre generisanja tagova.

### Problem: "OpenAI API key not configured"

**Rešenje:** 
1. Proverite da li je `OPENAI_API_KEY` postavljen u environment variables
2. Restartujte aplikaciju
3. Pogledajte `AI_FEATURES.md` za detaljne instrukcije

### Problem: "No tags were generated"

**Mogući uzroci:**
- API key je nevažeći
- Nedovoljan kredit na OpenAI nalogu
- Opis dokumenta je previše kratak

**Rešenje:**
1. Proverite OpenAI account balance
2. Dodajte detaljniji opis dokumenta
3. Proverite da API key počinje sa `sk-`

### Problem: Tagovi nisu relevantni

**Rešenje:**
1. Napišite detaljniji i precizniji opis
2. Uključite specifične termine i koncepte
3. Pokušajte ponovo sa boljim opisom

## Tehnički Detalji

### API Poziv

```javascript
const tags = await GenerateTagsFromText(
  documentName,    // Naziv dokumenta
  description,     // Opis dokumenta
  documentType,    // Tip dokumenta
  5               // Maksimalan broj tagova
)
```

### Performanse

- **Vreme odgovora:** 1-3 sekunde
- **Cena po pozivu:** ~$0.001 (1/10 centa)
- **Rate limit:** 3 zahteva po minuti (free tier)

### Model Informacije

- **Model:** GPT-3.5 Turbo
- **Temperature:** 0.3 (fokusiran na konzistentnost)
- **Max tokens:** 100
- **Context window:** ~12,000 karaktera

## Napredne Opcije

### Prilagođavanje Broja Tagova

U kodu možete promeniti broj generisanih tagova:

```javascript
const tags = await GenerateTagsFromText(
  documentName,
  description,
  documentType,
  10  // Generiši 10 tagova umesto 5
)
```

### Bulk Tagging

Za tagovanje više dokumenata odjednom, koristite skriptu:

```javascript
for (const doc of documents) {
  const tags = await GenerateTagsFromText(
    doc.name,
    doc.description,
    doc.type,
    5
  )
  doc.keywords = tags.join(', ')
  await saveDocument(doc)
  
  // Pauza da izbegnemo rate limiting
  await sleep(1000)
}
```

## FAQ

**Q: Da li mogu ručno editovati AI-generisane tagove?**
A: Da! Tagovi se popunjavaju u Keywords polje koje možete slobodno menjati.

**Q: Koliko kosta mesečno korišćenje?**
A: Za tipičnu upotrebu (100 dokumenata mesečno): ~$0.10 (10 centi)

**Q: Da li AI tagovanje radi offline?**
A: Ne, potrebna je internet konekcija za OpenAI API pozive.

**Q: Mogu li koristiti drugi AI model?**
A: Da, možete promeniti model u `llm_service.go` (npr. GPT-4, Claude, itd.)

## Podrška

Za dodatnu pomoć:
- Pogledajte `AI_FEATURES.md` za kompletnu dokumentaciju
- Proverite [OpenAI Status](https://status.openai.com/) za API status
- Kontaktirajte tim za podršku

---

💡 **Pro Tip:** Što je opis detaljniji, AI će generisati preciznije i korisnije tagove!
