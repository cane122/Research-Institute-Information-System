import { 
  GetProjectsByCurrentUser,
  GetProjectByID,
  CreateNewProject,
  UpdateProject,
  DeleteProject,
  GetProjectMembers,
  AddProjectMember,
  RemoveProjectMember,
  CompleteProject,
  GetProjectAnalytics,
  GetAllProjects
} from '../../wailsjs/go/main/App.js'

/**
 * Fetches all projects for the currently logged-in user
 */
export async function fetchUserProjects() {
  const items = await GetProjectsByCurrentUser()
  return mapProjectsToUI(items)
}

/**
 * Fetches all projects in the system (for admins)
 */
export async function fetchAllProjects() {
  const items = await GetAllProjects()
  return mapProjectsToUI(items)
}

/**
 * Maps backend project data to UI-friendly format
 */
function mapProjectsToUI(items) {
  return (items || []).map(p => ({
    id: p.projekat_id,
    name: p.naziv_projekta,
    description: p.opis || '',
    progress: calculateProgress(p.broj_zadataka),
    status: mapStatus(p.status),
    leader: p.rukovodilac_ime || '',
    leaderId: p.rukovodilac_id || null,
    deadline: p.datum_zavrsetka ? formatDateForInput(p.datum_zavrsetka) : null,
    created: p.datum_pocetka ? formatDateForInput(p.datum_pocetka) : null,
    updated: p.datum_zavrsetka ? formatDateForInput(p.datum_zavrsetka) : null,
    teamCount: p.broj_clanova || 0,
    taskCount: p.broj_zadataka || 0,
    workflowId: p.radni_tok_id || null,
    team: []
  }))
}

/**
 * Maps backend status to UI status
 */
function mapStatus(backendStatus) {
  const statusMap = {
    'aktivan': 'active',
    'aktivni': 'active', 
    'završen': 'completed',
    'zavrsen': 'completed',
    'na čekanju': 'on-hold',
    'na cekanju': 'on-hold',
    'otkazan': 'cancelled'
  }
  const normalized = (backendStatus || 'aktivan').toLowerCase()
  return statusMap[normalized] || 'active'
}

/**
 * Calculates project progress based on task count
 */
function calculateProgress(taskCount) {
  // Simplified progress calculation - in real app would need completed vs total tasks
  if (!taskCount) return 0
  return Math.min(100, Math.round(taskCount * 5)) // Placeholder logic
}

/**
 * Formats date for input fields (YYYY-MM-DD)
 */
function formatDateForInput(dateStr) {
  if (!dateStr) return null
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return null
  return d.toISOString().split('T')[0]
}

/**
 * Formats date for API (RFC3339)
 */
function formatDateForAPI(dateStr) {
  if (!dateStr) return null
  return new Date(dateStr).toISOString()
}

/**
 * Creates a new project
 */
export async function createProject(projectData) {
  const payload = {
    naziv_projekta: projectData.name || projectData.naziv_projekta,
    opis: projectData.description || projectData.opis || '',
    datum_pocetka: projectData.startDate ? formatDateForAPI(projectData.startDate) : formatDateForAPI(new Date()),
    datum_zavrsetka: projectData.deadline ? formatDateForAPI(projectData.deadline) : null,
    radni_tok_id: projectData.workflowId || null,
    clanovi_tima: projectData.teamMembers || []
  }
  
  await CreateNewProject(payload)
}

/**
 * Updates an existing project
 */
export async function updateProject(projectId, projectData) {
  const payload = {
    projekat_id: projectId,
    naziv_projekta: projectData.name,
    opis: projectData.description || '',
    datum_pocetka: projectData.created ? formatDateForAPI(projectData.created) : null,
    datum_zavrsetka: projectData.deadline ? formatDateForAPI(projectData.deadline) : null,
    status: projectData.status === 'active' ? 'aktivan' : projectData.status,
    rukovodilac_id: projectData.leaderId || null,
    radni_tok_id: projectData.workflowId || null
  }
  
  await UpdateProject(projectId, payload)
}

/**
 * Deletes a project
 */
export async function deleteProject(projectId) {
  await DeleteProject(projectId)
}

/**
 * Fetches detailed project information
 */
export async function fetchProjectDetails(projectId) {
  const project = await GetProjectByID(projectId)
  const members = await GetProjectMembers(projectId)
  
  return {
    ...mapProjectsToUI([project])[0],
    team: members.map(m => ({
      id: m.korisnik_id,
      name: `${m.ime || ''} ${m.prezime || ''}`.trim() || m.korisnicko_ime,
      username: m.korisnicko_ime,
      role: m.naziv_uloge || ''
    }))
  }
}

/**
 * Adds a team member to a project
 */
export async function addTeamMember(projectId, userId) {
  await AddProjectMember(projectId, userId)
}

/**
 * Removes a team member from a project
 */
export async function removeTeamMember(projectId, userId) {
  await RemoveProjectMember(projectId, userId)
}

/**
 * Marks a project as completed
 */
export async function completeProject(projectId) {
  await CompleteProject(projectId)
}

/**
 * Fetches project analytics
 */
export async function fetchProjectAnalytics(projectId) {
  const analytics = await GetProjectAnalytics(projectId)
  
  return {
    totalTasks: analytics.total_tasks || 0,
    completedTasks: analytics.completed_tasks || 0,
    inProgressTasks: analytics.in_progress_tasks || 0,
    pendingTasks: analytics.pending_tasks || 0,
    completionPercentage: analytics.completion_percentage || 0,
    overdueTasks: analytics.overdue_tasks || 0,
    teamMemberCount: analytics.team_member_count || 0,
    tasksByPhase: analytics.tasks_by_phase || []
  }
}
