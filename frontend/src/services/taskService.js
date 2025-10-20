import { 
  GetProjectsByCurrentUser,
  GetTasksByProject,
  GetTasksByCurrentUser,
  GetTaskByID,
  CreateTask,
  UpdateTask,
  DeleteTask,
  MoveTaskToPhase,
  GetWorkflowPhases,
  GetProjectMembers,
  GetConditionsByPhase,
  GetConditionAssessmentsByTask,
  CreateOrUpdateConditionAssessment,
  GetConditionFulfillmentStatus,
  GetTaskComments,
  AddTaskComment,
  DeleteTaskComment
} from '../../wailsjs/go/main/App.js'

/**
 * Fetches all projects for the current user
 */
export async function fetchUserProjects() {
  try {
    const projects = await GetProjectsByCurrentUser()
    return (projects || []).map(p => ({
      id: p.projekat_id,
      name: p.naziv_projekta,
      workflowId: p.radni_tok_id,
      leaderId: p.rukovodilac_id
    }))
  } catch (error) {
    console.error('Error fetching user projects:', error)
    return []
  }
}

/**
 * Fetches phases for a specific workflow
 */
export async function fetchWorkflowPhases(workflowId) {
  try {
    if (!workflowId) return []
    const phases = await GetWorkflowPhases(workflowId)
    return (phases || []).map(p => ({
      id: p.faza_id,
      name: p.naziv_faze,
      order: p.redosled,
      workflowId: p.radni_tok_id
    }))
  } catch (error) {
    console.error('Error fetching workflow phases:', error)
    return []
  }
}

/**
 * Fetches tasks for a specific project
 */
export async function fetchProjectTasks(projectId) {
  try {
    if (!projectId) return []
    const tasks = await GetTasksByProject(projectId)
    return (tasks || []).map(mapTaskToUI)
  } catch (error) {
    console.error('Error fetching project tasks:', error)
    return []
  }
}

/**
 * Fetches all tasks for current user
 */
export async function fetchCurrentUserTasks() {
  try {
    const tasks = await GetTasksByCurrentUser()
    return (tasks || []).map(mapTaskToUI)
  } catch (error) {
    console.error('Error fetching user tasks:', error)
    return []
  }
}

/**
 * Fetches project members
 */
export async function fetchProjectMembers(projectId) {
  try {
    if (!projectId) return []
    const members = await GetProjectMembers(projectId)
    return (members || []).map(m => ({
      id: m.korisnik_id,
      name: `${m.ime || ''} ${m.prezime || ''}`.trim() || m.korisnicko_ime,
      username: m.korisnicko_ime
    }))
  } catch (error) {
    console.error('Error fetching project members:', error)
    return []
  }
}

/**
 * Maps backend task to UI format
 */
function mapTaskToUI(task) {
  return {
    id: task.zadatak_id,
    title: task.naziv_zadatka,
    description: task.opis || '',
    projectId: task.projekat_id,
    projectName: task.naziv_projekta || '',
    phaseId: task.faza_id,
    phaseName: task.naziv_faze || '',
    assigneeId: task.dodeljen_korisniku_id,
    assigneeName: task.dodeljen_korisniku || '',
    deadline: task.rok ? formatDateForInput(task.rok) : null,
    priority: mapPriority(task.prioritet),
    progress: task.progres || 0,
    resources: task.resursi || '',
    created: task.kreiran_datuma ? formatDateForInput(task.kreiran_datuma) : null
  }
}

/**
 * Maps backend priority to UI priority
 */
function mapPriority(priority) {
  if (!priority) return 'medium'
  const p = priority.toLowerCase()
  if (p === 'visok' || p === 'high') return 'high'
  if (p === 'nizak' || p === 'low') return 'low'
  return 'medium'
}

/**
 * Maps UI priority to backend priority
 */
function mapPriorityToBackend(priority) {
  const priorityMap = {
    'high': 'visok',
    'medium': 'srednji',
    'low': 'nizak'
  }
  return priorityMap[priority] || 'srednji'
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
 * Creates a new task
 */
export async function createTask(taskData) {
  try {
    const request = {
      projekat_id: taskData.projectId,
      naziv_zadatka: taskData.title,
      opis: taskData.description || '',
      dodeljen_korisniku_id: taskData.assigneeId || null,
      rok: taskData.deadline ? formatDateForAPI(taskData.deadline) : null,
      prioritet: mapPriorityToBackend(taskData.priority || 'medium'),
      faza_id: taskData.phaseId || null,
      resursi: taskData.resources || ''
    }
    
    await CreateTask(request)
    return { success: true }
  } catch (error) {
    console.error('Error creating task:', error)
    return { success: false, error: error.message }
  }
}

/**
 * Updates an existing task
 */
export async function updateTask(taskId, updates) {
  try {
    const request = {}
    
    if (updates.title !== undefined) {
      request.naziv_zadatka = updates.title
    }
    if (updates.description !== undefined) {
      request.opis = updates.description
    }
    if (updates.assigneeId !== undefined) {
      request.dodeljen_korisniku_id = updates.assigneeId
    }
    if (updates.deadline !== undefined) {
      request.rok = updates.deadline ? formatDateForAPI(updates.deadline) : null
    }
    if (updates.priority !== undefined) {
      request.prioritet = mapPriorityToBackend(updates.priority)
    }
    if (updates.progress !== undefined) {
      request.progres = updates.progress
    }
    if (updates.phaseId !== undefined) {
      request.faza_id = updates.phaseId
    }
    if (updates.resources !== undefined) {
      request.resursi = updates.resources
    }
    
    await UpdateTask(taskId, request)
    return { success: true }
  } catch (error) {
    console.error('Error updating task:', error)
    return { success: false, error: error.message }
  }
}

/**
 * Moves task to a different phase
 */
export async function moveTaskToPhase(taskId, newPhaseId) {
  try {
    await MoveTaskToPhase(taskId, newPhaseId)
    return { success: true }
  } catch (error) {
    console.error('Error moving task to phase:', error)
    return { success: false, error: error.message }
  }
}

/**
 * Deletes a task
 */
export async function deleteTask(taskId) {
  try {
    await DeleteTask(taskId)
    return { success: true }
  } catch (error) {
    console.error('Error deleting task:', error)
    return { success: false, error: error.message }
  }
}

/**
 * Fetches conditions for a specific phase
 */
export async function fetchConditionsByPhase(phaseId) {
  try {
    if (!phaseId) return []
    const conditions = await GetConditionsByPhase(phaseId)
    return (conditions || []).map(c => ({
      id: c.uslov_id,
      phaseId: c.faza_id,
      description: c.opis,
      criteria: c.kriterijum,
      phaseName: c.naziv_faze,
      created: c.kreiran_datuma
    }))
  } catch (error) {
    console.error('Error fetching conditions:', error)
    return []
  }
}

/**
 * Fetches condition assessments for a specific task
 */
export async function fetchConditionAssessmentsByTask(taskId) {
  try {
    if (!taskId) return []
    const assessments = await GetConditionAssessmentsByTask(taskId)
    return (assessments || []).map(a => ({
      id: a.procena_id,
      taskId: a.zadatak_id,
      conditionId: a.uslov_id,
      fulfilled: a.ispunjen,
      note: a.napomena,
      evaluatedBy: a.promenio_korisnik_id,
      evaluatedAt: a.datum_procene,
      conditionDescription: a.opis_uslova,
      conditionCriteria: a.kriterijum_uslova,
      taskName: a.naziv_zadatka,
      evaluatorName: a.ime_korisnika
    }))
  } catch (error) {
    console.error('Error fetching condition assessments:', error)
    return []
  }
}

/**
 * Creates or updates a condition assessment
 */
export async function updateConditionAssessment(taskId, conditionId, fulfilled, note = '') {
  try {
    const assessment = {
      zadatak_id: taskId,
      uslov_id: conditionId,
      ispunjen: fulfilled,
      napomena: note
    }
    
    await CreateOrUpdateConditionAssessment(assessment)
    return { success: true }
  } catch (error) {
    console.error('Error updating condition assessment:', error)
    return { success: false, error: error.message }
  }
}

/**
 * Fetches condition fulfillment status for a task
 */
export async function fetchConditionFulfillmentStatus(taskId) {
  try {
    if (!taskId) return null
    const status = await GetConditionFulfillmentStatus(taskId)
    return status
  } catch (error) {
    console.error('Error fetching condition fulfillment status:', error)
    return null
  }
}

/**
 * Fetches comments for a task
 */
export async function fetchTaskComments(taskId) {
  try {
    if (!taskId) return []
    const comments = await GetTaskComments(taskId)
    return (comments || []).map(c => ({
      id: c.komentar_id,
      taskId: c.zadatak_id,
      userId: c.korisnik_id,
      text: c.tekst_komentara,
      createdAt: c.datuma_kreiranja,
      userName: c.ime_korisnika || 'Unknown User'
    }))
  } catch (error) {
    console.error('Error fetching task comments:', error)
    return []
  }
}

/**
 * Adds a comment to a task
 */
export async function addTaskComment(taskId, comment) {
  try {
    if (!taskId || !comment) {
      throw new Error('Task ID and comment text are required')
    }
    await AddTaskComment(taskId, comment)
    return true
  } catch (error) {
    console.error('Error adding task comment:', error)
    throw error
  }
}

/**
 * Deletes a comment from a task
 */
export async function deleteTaskComment(commentId) {
  try {
    if (!commentId) {
      throw new Error('Comment ID is required')
    }
    await DeleteTaskComment(commentId)
    return true
  } catch (error) {
    console.error('Error deleting task comment:', error)
    throw error
  }
}
