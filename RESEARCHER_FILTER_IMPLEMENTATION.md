# Implementacija Filtriranja Zadataka za Istraživače

## Pregled

Implementiran je sistem filtriranja zadataka koji omogućava da korisnici sa ulogom **Istraživač** (Istrazivac) vide samo zadatke koji su njima dodeljeni na projektima gde učestvuju.

## Izmene

### 1. Auth Store (`frontend/src/stores/auth.js`)

Dodat je novi computed property `isResearcher` koji proverava da li je korisnik istraživač:

```javascript
const isResearcher = computed(() => 
  user.value?.uloga === 'researcher' || 
  user.value?.naziv_uloge === 'Istrazivac'
)
```

**Napomena:** Provera uključuje oba polja jer:
- `naziv_uloge` je vrednost iz baze ('Istrazivac', 'Administrator', 'Rukovodilac projekta')
- `uloga` može biti korišćen u nekim delovima koda kao shorthand

### 2. Tasks View (`frontend/src/views/Tasks.vue`)

Modifikovan je `filteredTasks` computed property da automatski filtrira zadatke:

```javascript
// If user is a researcher, only show tasks assigned to them
if (authStore.isResearcher && authStore.user) {
  const currentUserId = authStore.user.korisnik_id || authStore.user.korisnikID
  filtered = filtered.filter(t => t.assigneeId == currentUserId)
}
```

## Kako radi

### Za Istraživače (Istrazivac):
1. Kada istraživač otvori stranicu sa zadacima i izabere projekat
2. Sistem automatski filtrira sve zadatke projekta
3. Prikazuju se **samo zadaci gde je `assigneeId` jednak ID-u trenutnog korisnika**
4. Istraživač ne može videti zadatke drugih članova tima

### Za Ostale Uloge:
- **Administrator** - vidi sve zadatke
- **Rukovodilac projekta** - vidi sve zadatke projekta
- **Organizator projekta** - vidi sve zadatke

## Baza Podataka

Uloge su definisane u tabeli `Uloge`:

```sql
CREATE TABLE Uloge (
    uloga_id SERIAL PRIMARY KEY,
    naziv_uloge VARCHAR(50) UNIQUE NOT NULL
);

-- Podrazumevane uloge:
INSERT INTO Uloge (naziv_uloge) VALUES 
('Administrator'),
('Rukovodilac projekta'),
('Istrazivac'),
('Organizator projekta');
```

## Testiranje

### Test Scenario 1: Istraživač sa zadacima
1. Prijavite se kao korisnik sa ulogom 'Istrazivac'
2. Otvorite stranicu Zadaci
3. Izaberite projekat gde imate dodeljene zadatke
4. **Očekivani rezultat:** Vidite samo svoje zadatke (gde ste vi assignee)

### Test Scenario 2: Istraživač bez zadataka
1. Prijavite se kao korisnik sa ulogom 'Istrazivac'
2. Otvorite stranicu Zadaci
3. Izaberite projekat gde nemate dodeljenih zadataka
4. **Očekivani rezultat:** Vidite poruku "Nema zadataka" u svakoj koloni

### Test Scenario 3: Admin/Manager
1. Prijavite se kao admin ili rukovodilac projekta
2. Otvorite stranicu Zadaci
3. Izaberite projekat
4. **Očekivani rezultat:** Vidite sve zadatke u projektu

## Debug Mode

U kodu je dodata debug konzolna ispis koji se može videti u Developer Tools:

```javascript
console.log('🔍 Researcher filter active:', {
  userId: currentUserId,
  naziv_uloge: authStore.user.naziv_uloge,
  totalTasks: filtered.length
})
```

Ovo pomaže u debugovanju i proveri da filter radi kako treba.

## Sigurnosna Napomena

⚠️ **Važno:** Ovo je frontend filtriranje. Za potpunu sigurnost, potrebno je implementirati i backend filtriranje u API-ju koji vraća zadatke. Frontend filtriranje može biti zaobiđeno ako korisnik ima pristup developer tools-ima.

## Buduća Poboljšanja

1. **Backend validacija** - Dodati proveru uloge na serveru
2. **Role-based API** - Kreirati različite endpointe za različite uloge
3. **Audit log** - Pratiti pokušaje pristupa zadacima koji nisu dodeljeni korisniku
4. **Notifications** - Obaveštavati istraživače kada im se dodeli novi zadatak

## Related Files

- `frontend/src/stores/auth.js` - Auth store sa `isResearcher` computed property
- `frontend/src/views/Tasks.vue` - Tasks view sa filterom
- `backend/models/models.go` - User model sa `NazivUloge` poljem
- `database/schema.sql` - Definicija uloga u bazi

## Changelog

- **2025-10-20**: Inicijalna implementacija researcher filtera
  - Dodat `isResearcher` computed u auth store
  - Dodata logika filtriranja u Tasks.vue
  - Dodati debug console logs
