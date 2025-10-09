<template>
  <div class="analytics-container">
    <div class="analytics-header">
      <div class="header-left">
        <button @click="goBack" class="back-btn">
          <span class="icon">←</span>
          Back
        </button>
        <h1>📊 Analitika Dokumenata</h1>
      </div>
      <button @click="exportPDF" class="export-btn">
        <span class="icon">📄</span>
        Export PDF
      </button>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <div class="stat-card documents">
        <div class="stat-icon">📄</div>
        <div class="stat-content">
          <h3>Broj dokumenata</h3>
          <p class="stat-value">{{ stats.ukupno_dokumenata || 0 }}</p>
        </div>
      </div>

      <div class="stat-card views">
        <div class="stat-icon">👁️</div>
        <div class="stat-content">
          <h3>Broj pregleda</h3>
          <p class="stat-value">{{ viewCount || 0 }}</p>
        </div>
      </div>

      <div class="stat-card deleted">
        <div class="stat-icon">🗑️</div>
        <div class="stat-content">
          <h3>Broj izbrisanih dokumenata</h3>
          <p class="stat-value">{{ deletedCount || 0 }}</p>
        </div>
      </div>

      <div class="stat-card new">
        <div class="stat-icon">✨</div>
        <div class="stat-content">
          <h3>Broj novih dokumenata</h3>
          <p class="stat-value">{{ stats.novih_dokumenata_mesecno || 0 }}</p>
        </div>
      </div>
    </div>

    <!-- Recent Activity and Additional Stats -->
    <div class="content-grid">
      <!-- Recent Activity -->
      <div class="activity-section">
        <h2>Skornje Aktivnosti</h2>
        <div class="activity-list">
          <div v-if="recentActivity.length === 0" class="no-activity">
            Nema aktivnosti
          </div>
          <div v-for="activity in recentActivity" :key="activity.log_id" class="activity-item">
            <div class="activity-icon" :class="getActivityTypeClass(activity.tip_aktivnosti)">
              {{ getActivityIcon(activity.tip_aktivnosti) }}
            </div>
            <div class="activity-details">
              <div class="activity-title">
                <strong>{{ activity.korisnik_ime }}</strong>
                <span class="activity-type">{{ formatActivityType(activity.tip_aktivnosti) }}</span>
              </div>
              <div class="activity-description">{{ activity.opis }}</div>
              <div class="activity-time">{{ formatTime(activity.kreiran_datuma) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Top Contributors -->
      <div class="contributors-section">
        <h2>Top Autori</h2>
        <div class="contributors-list">
          <div v-for="(contributor, index) in topContributors" :key="index" class="contributor-item">
            <div class="contributor-rank">{{ index + 1 }}</div>
            <div class="contributor-info">
              <div class="contributor-name">{{ contributor.ime }}</div>
              <div class="contributor-stats">{{ contributor.broj_uploadovanih }} dokumenata</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Documents by Type -->
    <div class="type-section">
      <h2>Dokumenti po tipu</h2>
      <div class="type-list">
        <div v-for="(count, type) in documentsByType" :key="type" class="type-item">
          <div class="type-header">
            <span class="type-name">{{ type }}</span>
            <span class="type-count">{{ count }}</span>
          </div>
          <div class="type-bar">
            <div class="type-progress" :style="{ width: getTypePercentage(count) + '%' }"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Activity Statistics -->
    <div class="activity-stats-section">
      <h2>Statistika Aktivnosti</h2>
      <div class="activity-stats-grid">
        <div v-for="stat in activityStats" :key="stat.tip_aktivnosti" class="activity-stat-card">
          <h4>{{ formatActivityType(stat.tip_aktivnosti) }}</h4>
          <div class="activity-stat-value">{{ stat.broj_aktivnosti }}</div>
          <div class="activity-stat-details">
            <span>Danas: {{ stat.danas }}</span>
            <span>Ove nedelje: {{ stat.ove_nedelje }}</span>
            <span>Ovog meseca: {{ stat.ovog_meseca }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  GetDocumentStatistics, 
  GetRecentActivity, 
  GetActivityStatistics,
  GetDocumentsByType,
  GetTopContributors
} from '../../../wailsjs/go/main/App'
import jsPDF from 'jspdf'
import 'jspdf-autotable'

const router = useRouter()
const stats = ref({})
const recentActivity = ref([])
const activityStats = ref([])
const documentsByType = ref({})
const topContributors = ref([])
const viewCount = ref(0)
const deletedCount = ref(0)

const goBack = () => {
  router.push('/documents')
}

const loadData = async () => {
  try {
    // Load document statistics
    const docStats = await GetDocumentStatistics()
    stats.value = docStats || {}
    console.log('Document Stats:', stats.value)

    // Load recent activity
    const activity = await GetRecentActivity(50)
    recentActivity.value = activity || []
    console.log('Recent Activity:', recentActivity.value)

    // Load activity statistics
    const actStats = await GetActivityStatistics()
    activityStats.value = actStats || []
    console.log('Activity Stats:', activityStats.value)
    
    // Calculate view and delete counts from activity stats
    const viewStat = actStats.find(s => s.tip_aktivnosti === 'VIEW')
    const deleteStat = actStats.find(s => s.tip_aktivnosti === 'DELETE')
    viewCount.value = viewStat ? viewStat.broj_aktivnosti : 0
    deletedCount.value = deleteStat ? deleteStat.broj_aktivnosti : 0

    // Load documents by type
    const byType = await GetDocumentsByType()
    documentsByType.value = byType || {}
    console.log('Documents by Type:', documentsByType.value)

    // Load top contributors
    const contributors = await GetTopContributors(5)
    topContributors.value = contributors || []
    console.log('Top Contributors:', topContributors.value)
  } catch (error) {
    console.error('Failed to load analytics:', error)
  }
}

const getActivityTypeClass = (type) => {
  const classes = {
    'UPLOAD': 'upload',
    'EDIT': 'edit',
    'DELETE': 'delete',
    'VIEW': 'view',
    'LOGIN': 'login',
    'LOGOUT': 'logout'
  }
  return classes[type] || 'default'
}

const getActivityIcon = (type) => {
  const icons = {
    'UPLOAD': '⬆️',
    'EDIT': '✏️',
    'DELETE': '🗑️',
    'VIEW': '👁️',
    'LOGIN': '🔓',
    'LOGOUT': '🔒'
  }
  return icons[type] || '📝'
}

const formatActivityType = (type) => {
  const types = {
    'UPLOAD': 'Upload',
    'EDIT': 'Izmena',
    'DELETE': 'Brisanje',
    'VIEW': 'Pregled',
    'LOGIN': 'Prijava',
    'LOGOUT': 'Odjava'
  }
  return types[type] || type
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date
  
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return 'Upravo sada'
  if (minutes < 60) return `Pre ${minutes} min`
  if (hours < 24) return `Pre ${hours} sati`
  if (days < 7) return `Pre ${days} dana`
  
  return date.toLocaleDateString('sr-RS')
}

const getTypePercentage = (count) => {
  if (!documentsByType.value || typeof documentsByType.value !== 'object') {
    return 0
  }
  const total = Object.values(documentsByType.value).reduce((sum, c) => sum + (c || 0), 0)
  return total > 0 ? (count / total) * 100 : 0
}

const exportPDF = async () => {
  try {
    // Validate data before generating PDF
    if (!stats.value || Object.keys(stats.value).length === 0) {
      alert('Nema podataka za generisanje PDF-a. Molimo sacekajte da se podaci ucitaju.')
      return
    }

    // Helper function to convert Serbian Cyrillic/Latin characters to ASCII
    const toASCII = (str) => {
      if (!str) return ''
      return str
        .replace(/č/g, 'c').replace(/Č/g, 'C')
        .replace(/ć/g, 'c').replace(/Ć/g, 'C')
        .replace(/ž/g, 'z').replace(/Ž/g, 'Z')
        .replace(/š/g, 's').replace(/Š/g, 'S')
        .replace(/đ/g, 'd').replace(/Đ/g, 'D')
    }

    const doc = new jsPDF()
    
    // Set up colors
    const primaryColor = [102, 126, 234] // #667eea
    const secondaryColor = [118, 75, 162] // #764ba2
    const darkColor = [30, 41, 59] // #1e293b
    const lightColor = [148, 163, 184] // #94a3b8
    
    let yPosition = 20
    
    // ========== HEADER ==========
    // Add gradient background effect (using rectangles)
    doc.setFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
    doc.rect(0, 0, 210, 50, 'F')
    
    // Title
    doc.setTextColor(255, 255, 255)
    doc.setFontSize(24)
    doc.setFont(undefined, 'bold')
    doc.text('ANALITICKI IZVESTAJ', 105, 25, { align: 'center' })
    
    // Subtitle
    doc.setFontSize(12)
    doc.setFont(undefined, 'normal')
    const reportDate = new Date().toLocaleDateString('sr-RS', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    })
    doc.text(`Generisano: ${reportDate}`, 105, 35, { align: 'center' })
    
    yPosition = 60
    
    // ========== EXECUTIVE SUMMARY ==========
    doc.setTextColor(darkColor[0], darkColor[1], darkColor[2])
    doc.setFontSize(16)
    doc.setFont(undefined, 'bold')
    doc.text('Izvrsni Rezime', 20, yPosition)
    
    yPosition += 10
    
    // Summary box with border
    doc.setDrawColor(primaryColor[0], primaryColor[1], primaryColor[2])
    doc.setLineWidth(0.5)
    doc.rect(20, yPosition, 170, 45, 'S')
    
    yPosition += 8
    
    doc.setFontSize(10)
    doc.setFont(undefined, 'normal')
    doc.setTextColor(darkColor[0], darkColor[1], darkColor[2])
    
    const summaryText = `Ovaj izvestaj pruza sveobuhvatan pregled aktivnosti u sistemu za upravljanje dokumentima. ` +
      `Analizirano je ukupno ${stats.value.ukupno_dokumenata || 0} dokumenata sa ${viewCount.value || 0} pregleda. ` +
      `Sistem belezi aktivnosti ${activityStats.value.length} razlicitih tipova, sa ukupno ` +
      `${recentActivity.value.length} skorasnjih aktivnosti.`
    
    const splitSummary = doc.splitTextToSize(summaryText, 160)
    doc.text(splitSummary, 25, yPosition)
    
    yPosition += 55
    
    // ========== KEY METRICS ==========
    doc.setFontSize(16)
    doc.setFont(undefined, 'bold')
    doc.text('Kljucne Metrike', 20, yPosition)
    
    yPosition += 10
    
    // Metrics in a grid
    const metrics = [
      { label: 'Ukupno Dokumenata', value: stats.value.ukupno_dokumenata || 0, color: primaryColor },
      { label: 'Broj Pregleda', value: viewCount.value || 0, color: [245, 87, 108] },
      { label: 'Izbrisanih Dokumenata', value: deletedCount.value || 0, color: [250, 112, 154] },
      { label: 'Novih Dokumenata', value: stats.value.novih_dokumenata_mesecno || 0, color: [48, 207, 208] }
    ]
    
    const boxWidth = 40
    const boxHeight = 25
    const spacing = 5
    
    metrics.forEach((metric, index) => {
      const x = 20 + (index * (boxWidth + spacing))
      
      // Colored box
      doc.setFillColor(metric.color[0], metric.color[1], metric.color[2])
      doc.roundedRect(x, yPosition, boxWidth, boxHeight, 3, 3, 'F')
      
      // Value
      doc.setTextColor(255, 255, 255)
      doc.setFontSize(18)
      doc.setFont(undefined, 'bold')
      doc.text(metric.value.toString(), x + boxWidth/2, yPosition + 12, { align: 'center' })
      
      // Label
      doc.setFontSize(7)
      doc.setFont(undefined, 'normal')
      const labelLines = doc.splitTextToSize(metric.label, boxWidth - 4)
      doc.text(labelLines, x + boxWidth/2, yPosition + 18, { align: 'center' })
    })
    
    yPosition += boxHeight + 15
    
    // ========== ACTIVITY STATISTICS ==========
    doc.setTextColor(darkColor[0], darkColor[1], darkColor[2])
    doc.setFontSize(16)
    doc.setFont(undefined, 'bold')
    doc.text('Statistika Aktivnosti', 20, yPosition)
    
    yPosition += 5
    
    // Activity table
    const activityTableData = activityStats.value.map(stat => [
      formatActivityType(stat.tip_aktivnosti || ''),
      (stat.broj_aktivnosti || 0).toString(),
      (stat.danas || 0).toString(),
      (stat.ove_nedelje || 0).toString(),
      (stat.ovog_meseca || 0).toString()
    ])
    
    doc.autoTable({
      startY: yPosition,
      head: [['Tip Aktivnosti', 'Ukupno', 'Danas', 'Ova Nedelja', 'Ovaj Mesec']],
      body: activityTableData,
      theme: 'striped',
      headStyles: {
        fillColor: primaryColor,
        textColor: [255, 255, 255],
        fontStyle: 'bold',
        fontSize: 10
      },
      styles: {
        fontSize: 9,
        cellPadding: 4
      },
      alternateRowStyles: {
        fillColor: [248, 250, 252]
      },
      margin: { left: 20, right: 20 }
    })
    
    yPosition = doc.lastAutoTable.finalY + 15
    
    // ========== TOP CONTRIBUTORS ==========
    if (yPosition > 250) {
      doc.addPage()
      yPosition = 20
    }
    
    doc.setFontSize(16)
    doc.setFont(undefined, 'bold')
    doc.text('Top Autori', 20, yPosition)
    
    yPosition += 5
    
    const contributorsTableData = topContributors.value.map((contributor, index) => [
      (index + 1).toString(),
      toASCII(contributor.ime) || 'N/A',
      (contributor.broj_uploadovanih || 0).toString()
    ])
    
    doc.autoTable({
      startY: yPosition,
      head: [['Rang', 'Korisnik', 'Broj Dokumenata']],
      body: contributorsTableData,
      theme: 'striped',
      headStyles: {
        fillColor: secondaryColor,
        textColor: [255, 255, 255],
        fontStyle: 'bold',
        fontSize: 10
      },
      styles: {
        fontSize: 9,
        cellPadding: 4
      },
      alternateRowStyles: {
        fillColor: [248, 250, 252]
      },
      margin: { left: 20, right: 20 },
      columnStyles: {
        0: { cellWidth: 20, halign: 'center' },
        2: { halign: 'center' }
      }
    })
    
    yPosition = doc.lastAutoTable.finalY + 15
    
    // ========== DOCUMENTS BY TYPE ==========
    if (yPosition > 250) {
      doc.addPage()
      yPosition = 20
    }
    
    doc.setFontSize(16)
    doc.setFont(undefined, 'bold')
    doc.text('Dokumenti po Tipu', 20, yPosition)
    
    yPosition += 5
    
    const typeTableData = Object.entries(documentsByType.value).map(([type, count]) => {
      const percentage = getTypePercentage(count)
      return [toASCII(type) || 'N/A', (count || 0).toString(), (percentage || 0).toFixed(1) + '%']
    })
    
    doc.autoTable({
      startY: yPosition,
      head: [['Tip Dokumenta', 'Broj', 'Procenat']],
      body: typeTableData,
      theme: 'striped',
      headStyles: {
        fillColor: [245, 87, 108],
        textColor: [255, 255, 255],
        fontStyle: 'bold',
        fontSize: 10
      },
      styles: {
        fontSize: 9,
        cellPadding: 4
      },
      alternateRowStyles: {
        fillColor: [248, 250, 252]
      },
      margin: { left: 20, right: 20 },
      columnStyles: {
        1: { halign: 'center' },
        2: { halign: 'center' }
      }
    })
    
    // ========== RECENT ACTIVITY ==========
    doc.addPage()
    yPosition = 20
    
    doc.setFontSize(16)
    doc.setFont(undefined, 'bold')
    doc.text('Skornje Aktivnosti', 20, yPosition)
    
    doc.setFontSize(9)
    doc.setFont(undefined, 'normal')
    doc.setTextColor(lightColor[0], lightColor[1], lightColor[2])
    doc.text(`Prikazano: ${Math.min(50, recentActivity.value.length)} najnovijih aktivnosti`, 20, yPosition + 5)
    
    yPosition += 10
    
    const activityTableDataDetailed = recentActivity.value.slice(0, 50).map(activity => {
      let dateString = 'N/A'
      if (activity.kreiran_datuma) {
        try {
          const date = new Date(activity.kreiran_datuma)
          if (!isNaN(date.getTime())) {
            dateString = date.toLocaleDateString('sr-RS') + ' ' + date.toLocaleTimeString('sr-RS', { 
              hour: '2-digit', 
              minute: '2-digit' 
            })
          }
        } catch (e) {
          console.error('Invalid date:', activity.kreiran_datuma, e)
        }
      }
      return [
        formatActivityType(activity.tip_aktivnosti || ''),
        toASCII(activity.korisnik_ime) || 'N/A',
        toASCII(activity.naziv_entiteta) || 'N/A',
        dateString
      ]
    })
    
    doc.autoTable({
      startY: yPosition,
      head: [['Aktivnost', 'Korisnik', 'Entitet', 'Vreme']],
      body: activityTableDataDetailed,
      theme: 'striped',
      headStyles: {
        fillColor: [48, 207, 208],
        textColor: [255, 255, 255],
        fontStyle: 'bold',
        fontSize: 9
      },
      styles: {
        fontSize: 8,
        cellPadding: 3
      },
      alternateRowStyles: {
        fillColor: [248, 250, 252]
      },
      margin: { left: 20, right: 20 },
      columnStyles: {
        0: { cellWidth: 30 },
        1: { cellWidth: 40 },
        2: { cellWidth: 60 },
        3: { cellWidth: 40 }
      }
    })
    
    // ========== FOOTER ON ALL PAGES ==========
    const pageCount = doc.internal.getNumberOfPages()
    
    for (let i = 1; i <= pageCount; i++) {
      doc.setPage(i)
      
      // Footer line
      doc.setDrawColor(lightColor[0], lightColor[1], lightColor[2])
      doc.setLineWidth(0.5)
      doc.line(20, 285, 190, 285)
      
      // Footer text
      doc.setTextColor(lightColor[0], lightColor[1], lightColor[2])
      doc.setFontSize(8)
      doc.setFont(undefined, 'normal')
      doc.text('Research Institute Information System', 20, 290)
      doc.text(`Strana ${i} od ${pageCount}`, 190, 290, { align: 'right' })
    }
    
    // Save the PDF
    const filename = `Analytics_Report_${new Date().toISOString().split('T')[0]}.pdf`
    doc.save(filename)
    
  } catch (error) {
    console.error('Failed to generate PDF:', error)
    alert('Greška pri generisanju PDF-a: ' + error.message)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.analytics-container {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.analytics-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.analytics-header h1 {
  font-size: 2rem;
  color: #1e293b;
  margin: 0;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  background: linear-gradient(135deg, #30cfd0 0%, #330867 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: transform 0.2s, box-shadow 0.2s;
}

.back-btn:hover {
  transform: translateX(-2px);
  box-shadow: 0 4px 12px rgba(48, 207, 208, 0.4);
}

.export-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: transform 0.2s, box-shadow 0.2s;
}

.export-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.stat-card.documents {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-card.views {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.stat-card.deleted {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  color: white;
}

.stat-card.new {
  background: linear-gradient(135deg, #30cfd0 0%, #330867 100%);
  color: white;
}

.stat-icon {
  font-size: 3rem;
  opacity: 0.9;
}

.stat-content h3 {
  font-size: 0.875rem;
  font-weight: 500;
  margin: 0 0 0.5rem 0;
  opacity: 0.9;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0;
}

.content-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

@media (max-width: 968px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}

.activity-section, .contributors-section, .type-section, .activity-stats-section {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.activity-section h2, .contributors-section h2, .type-section h2, .activity-stats-section h2 {
  font-size: 1.25rem;
  color: #1e293b;
  margin: 0 0 1.5rem 0;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 0.75rem;
}

.activity-list {
  max-height: 500px;
  overflow-y: auto;
}

.no-activity {
  text-align: center;
  padding: 2rem;
  color: #94a3b8;
}

.activity-item {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  border-bottom: 1px solid #f1f5f9;
  transition: background-color 0.2s;
}

.activity-item:hover {
  background-color: #f8fafc;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
  flex-shrink: 0;
}

.activity-icon.upload {
  background: #dbeafe;
}

.activity-icon.edit {
  background: #fef3c7;
}

.activity-icon.delete {
  background: #fee2e2;
}

.activity-icon.view {
  background: #e0e7ff;
}

.activity-icon.login, .activity-icon.logout {
  background: #dcfce7;
}

.activity-details {
  flex: 1;
}

.activity-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
}

.activity-type {
  background: #f1f5f9;
  padding: 0.125rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  color: #64748b;
}

.activity-description {
  color: #64748b;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.activity-time {
  color: #94a3b8;
  font-size: 0.75rem;
}

.contributors-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.contributor-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  background: #f8fafc;
  border-radius: 8px;
}

.contributor-rank {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  flex-shrink: 0;
}

.contributor-info {
  flex: 1;
}

.contributor-name {
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 0.25rem;
}

.contributor-stats {
  color: #64748b;
  font-size: 0.875rem;
}

.type-section {
  margin-bottom: 2rem;
}

.type-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.type-item {
  padding: 0.75rem;
}

.type-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.type-name {
  font-weight: 600;
  color: #1e293b;
}

.type-count {
  color: #64748b;
  font-weight: 500;
}

.type-bar {
  height: 8px;
  background: #f1f5f9;
  border-radius: 4px;
  overflow: hidden;
}

.type-progress {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s ease;
}

.activity-stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}

.activity-stat-card {
  padding: 1rem;
  background: #f8fafc;
  border-radius: 8px;
  border-left: 4px solid #667eea;
}

.activity-stat-card h4 {
  margin: 0 0 0.5rem 0;
  color: #1e293b;
  font-size: 0.875rem;
}

.activity-stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #667eea;
  margin-bottom: 0.5rem;
}

.activity-stat-details {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.activity-stat-details span {
  font-size: 0.75rem;
  color: #64748b;
}
</style>
