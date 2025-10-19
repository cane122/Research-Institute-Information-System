import { GetUserProjects, CreateProject } from '../../wailsjs/go/main/App.js'

export async function fetchUserProjects() {
  const items = await GetUserProjects()
  // Map backend fields to UI shape
  return (items || []).map(p => ({
    id: p.projekat_id,
    name: p.naziv_projekta,
    description: p.opis || '',
    progress: p.broj_zadataka ? Math.min(100, Math.round((p.broj_zadataka > 0 ? (p.broj_zadataka /* placeholder */) : 0) * 10)) : 0,
    status: (p.status || 'Aktivan').toLowerCase() === 'aktivan' ? 'active' : (p.status || '').toLowerCase(),
    leader: p.rukovodilac_ime || '',
    leaderId: p.rukovodilac_id || '',
    deadline: p.datum_zavrsetka ? new Date(p.datum_zavrsetka).toISOString().split('T')[0] : '',
    created: p.datum_pocetka ? new Date(p.datum_pocetka).toISOString().split('T')[0] : '',
    updated: p.datum_zavrsetka ? new Date(p.datum_zavrsetka).toISOString().split('T')[0] : '',
    team: []
  }))
}

export async function createProject(name, description, deadline) {
  const payload = {
    naziv_projekta: name,
    opis: description,
    datum_zavrsetka: deadline ? new Date(deadline) : null
  }
  await CreateProject(payload)
}
