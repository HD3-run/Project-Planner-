<script setup>
import { ref, onMounted, computed } from 'vue'

const API_URL = window.location.hostname === 'localhost' ? 'http://localhost:8080/api' : '/api'

// UI & Auth State
const sections = ref([])
const activeSection = ref(null)
const expandedCards = ref([])
const showAuthModal = ref(false)
const authMode = ref('login')
const authEmail = ref('')
const authPassword = ref('')
const authError = ref('')
const searchQuery = ref('')
const statusFilter = ref('all')
const activeView = ref('list')
const editMode = ref(false)
const isUpdating = ref(false)

// Token State
const sessionToken = ref(localStorage.getItem('auth_token') || null)
const refreshToken = ref(localStorage.getItem('auth_refresh_token') || null)
const userRole = ref(localStorage.getItem('auth_role') || 'user')

const isBaba = computed(() => userRole.value === 'BABA')

// Stats
const completionPercentage = computed(() => {
  if (!sections.value.length) return 0
  const allFeatures = sections.value.flatMap(s => s.features || [])
  if (!allFeatures.length) return 0
  const liveCount = allFeatures.filter(f => f.status === 'live').length
  return Math.round((liveCount / allFeatures.length) * 100)
})

const filteredSections = computed(() => {
  return sections.value.map(s => ({
    ...s,
    features: (s.features || []).filter(f => {
      const matchesSearch = !searchQuery.value || 
        f.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        (f.tech || []).some(t => t.toLowerCase().includes(searchQuery.value.toLowerCase()))
      const matchesStatus = statusFilter.value === 'all' || f.status === statusFilter.value
      return matchesSearch && matchesStatus
    })
  })).filter(s => s.features.length > 0 || !searchQuery.value)
})

// API Logic
const apiFetch = async (url, options = {}) => {
  if (!options.headers) options.headers = {}
  if (sessionToken.value) options.headers['Authorization'] = `Bearer ${sessionToken.value}`

  let res = await fetch(url, options)

  if (res.status === 401 && refreshToken.value) {
    try {
      const refreshRes = await fetch(`${API_URL}/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken.value })
      })
      const data = await refreshRes.json()
      if (!refreshRes.ok) throw new Error('Refresh failed')
      sessionToken.value = data.access_token
      userRole.value = data.role
      localStorage.setItem('auth_token', data.access_token)
      localStorage.setItem('auth_role', data.role)
      options.headers['Authorization'] = `Bearer ${sessionToken.value}`
      res = await fetch(url, options)
    } catch (err) { handleLogout() }
  }
  return res
}

const handleAuth = async () => {
  authError.value = ''
  try {
    const res = await fetch(`${API_URL}/${authMode.value}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value, password: authPassword.value })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Auth failed')
    
    sessionToken.value = data.access_token
    refreshToken.value = data.refresh_token
    userRole.value = data.role
    localStorage.setItem('auth_token', data.access_token)
    localStorage.setItem('auth_refresh_token', data.refresh_token)
    localStorage.setItem('auth_role', data.role)
    
    showAuthModal.value = false
    loadData()
  } catch (err) { authError.value = err.message }
}

const handleLogout = () => {
  sessionToken.value = null
  refreshToken.value = null
  userRole.value = 'user'
  localStorage.clear()
  showAuthModal.value = true 
}

const loadData = async () => {
  if (!sessionToken.value) return
  try {
    const res = await apiFetch(`${API_URL}/architecture`)
    const data = await res.json()
    sections.value = data.map(s => ({
      id: s.ID, title: s.Title, icon: s.Icon, color: s.Color, description: s.Description,
      features: (s.Features || []).map(f => ({
        id: f.ID, section_id: f.SectionID, title: f.Title, status: f.Status, subtitle: f.Subtitle,
        impact: f.Impact, how_works: f.HowItWorks, approach: f.Approach,
        tech: Array.isArray(f.Tech) ? f.Tech : JSON.parse(f.Tech || '[]')
      }))
    }))
  } catch (err) { console.error(err) }
}

// BABA Editing Power
const saveFeature = async (feature) => {
  if (!isBaba.value) return
  isUpdating.value = true
  try {
    await apiFetch(`${API_URL}/architecture/feature`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...feature,
        tech: JSON.stringify(feature.tech)
      })
    })
  } catch (err) { console.error(err) }
  finally { isUpdating.value = false }
}

const addFeature = (sectionId) => {
  if (!isBaba.value) return
  const section = sections.value.find(s => s.id === sectionId)
  const newFeat = {
    section_id: sectionId,
    title: 'New Feature',
    status: 'planned',
    tech: ['Node.js'],
    subtitle: 'Describe this feature...',
    impact: '', how_works: '', approach: ''
  }
  section.features.push(newFeat)
  editMode.value = true
}

const deleteFeature = async (id, sectionId) => {
  if (!isBaba.value || !confirm('Delete this feature?')) return
  try {
    await apiFetch(`${API_URL}/architecture/feature?id=${id}`, { method: 'DELETE' })
    const section = sections.value.find(s => s.id === sectionId)
    section.features = section.features.filter(f => f.id !== id)
  } catch (err) { console.error(err) }
}

const toggleCard = (id) => {
  if (editMode.value) return // Don't collapse during editing
  expandedCards.value = expandedCards.value.includes(id) 
    ? expandedCards.value.filter(i => i !== id) 
    : [...expandedCards.value, id]
}

const scrollTo = (id) => {
  activeSection.value = id
  activeView.value = 'list'
  setTimeout(() => {
    const el = document.getElementById(id)
    if(el) window.scrollTo({ top: el.getBoundingClientRect().top + window.scrollY - 120, behavior: 'smooth' })
  }, 50)
}

onMounted(() => {
  if (!sessionToken.value) showAuthModal.value = true
  else loadData()
})
</script>

<template>
  <div class="app-container">
    <div class="ambient-bg"><div class="glow-1"></div><div class="glow-2"></div></div>

    <header class="top-nav">
      <div class="brand-wrap">
        <div class="brand-logo">E</div>
        <div class="brand-info"><h1>ECOMMITRA</h1><p>Product Architecture</p></div>
      </div>

      <div v-if="sessionToken" class="search-container">
        <svg class="search-icon-svg" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
        <input type="text" v-model="searchQuery" class="search-input" placeholder="Search architecture..." />
      </div>

      <div class="header-actions" v-if="sessionToken">
        <!-- BABA Edit Toggle -->
        <button v-if="isBaba" @click="editMode = !editMode" class="btn-edit-toggle" :class="{ active: editMode }">
          <svg width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
          <span>{{ editMode ? 'Finish Editing' : 'Edit Architecture' }}</span>
        </button>

        <div class="view-toggle">
          <button @click="activeView = 'list'" class="view-btn" :class="{ active: activeView === 'list' }">List</button>
          <button @click="activeView = 'map'" class="view-btn" :class="{ active: activeView === 'map' }">Map</button>
        </div>

        <div class="progress-ring-wrap">
          <div class="ring-bar-bg"><div class="ring-bar-fill" :style="{ width: completionPercentage + '%' }"></div></div>
          <span class="ring-text">{{ completionPercentage }}% LIVE</span>
        </div>

        <button @click="handleLogout" class="btn-logout" title="Logout"><svg width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg></button>
      </div>
    </header>

    <div v-if="sessionToken" class="layout">
      <aside class="sidebar">
        <div class="sidebar-section-label">Navigation</div>
        <nav class="sidebar-nav">
          <div v-for="s in filteredSections" :key="'nav-'+s.id" class="nav-item" :class="{ active: activeSection === s.id }" @click="scrollTo(s.id)">
            <span>{{ s.icon }}</span><span class="nav-label">{{ s.title }}</span><span class="nav-badge">{{ s.features.length }}</span>
          </div>
        </nav>
      </aside>

      <main class="main-content">
        <!-- List View -->
        <div v-if="activeView === 'list'">
          <div v-for="section in filteredSections" :key="section.id" :id="section.id" class="section-block">
            <div class="section-header">
              <div class="section-title-wrap">
                <div class="section-icon-dot" :style="{ backgroundColor: section.color }"></div>
                <h3>{{ section.title }}</h3>
                <button v-if="editMode && isBaba" @click="addFeature(section.id)" class="btn-add-feat">+ Add Feature</button>
              </div>
              <p class="detail-text" style="margin-top: 8px; opacity: 0.7;">{{ section.description }}</p>
            </div>

            <div class="features-grid">
              <div v-for="feat in section.features" :key="feat.id || 'new'" class="feature-card" :class="{ editing: editMode }" @click="toggleCard(feat.id)">
                <!-- Card Header -->
                <div class="card-top">
                  <div class="card-title-group">
                    <div class="feature-icon-box">{{ section.icon }}</div>
                    <input v-if="editMode" v-model="feat.title" @blur="saveFeature(feat)" class="edit-input-title" />
                    <h4 v-else>{{ feat.title }}</h4>
                  </div>
                  
                  <div style="display: flex; align-items: center; gap: 12px;">
                    <select v-if="editMode" v-model="feat.status" @change="saveFeature(feat)" class="edit-select">
                      <option value="live">Live</option>
                      <option value="planned">Planned</option>
                      <option value="future">Future</option>
                    </select>
                    <span v-else class="status-pill" :class="feat.status">{{ feat.status }}</span>
                    
                    <button v-if="editMode && isBaba" @click.stop="deleteFeature(feat.id, section.id)" class="btn-del">×</button>
                  </div>
                </div>

                <!-- Card Subtitle -->
                <textarea v-if="editMode" v-model="feat.subtitle" @blur="saveFeature(feat)" class="edit-textarea-sub"></textarea>
                <p v-else class="feature-subtitle">{{ feat.subtitle }}</p>

                <!-- Expanded Details (Always show if editing) -->
                <div v-if="expandedCards.includes(feat.id) || editMode" class="card-details">
                  <div class="detail-section">
                    <h5>Business Impact</h5>
                    <textarea v-if="editMode" v-model="feat.impact" @blur="saveFeature(feat)" class="edit-textarea-detail"></textarea>
                    <p v-else class="detail-text">{{ feat.impact || 'Standard platform capability.' }}</p>
                    
                    <h5 style="margin-top: 20px;">Logic & Architecture</h5>
                    <textarea v-if="editMode" v-model="feat.how_works" @blur="saveFeature(feat)" class="edit-textarea-detail"></textarea>
                    <p v-else class="detail-text">{{ feat.how_works || 'Logic contained in core services.' }}</p>
                  </div>
                  
                  <div class="detail-section">
                    <h5>Tech Stack (JSON Array)</h5>
                    <input v-if="editMode" :value="JSON.stringify(feat.tech)" @blur="e => { feat.tech = JSON.parse(e.target.value); saveFeature(feat) }" class="edit-input-tech" />
                    <div v-else class="tech-chips">
                      <span v-for="t in feat.tech" :key="t" class="tech-chip">{{ t }}</span>
                    </div>

                    <h5 style="margin-top: 20px;">Implementation Approach</h5>
                    <textarea v-if="editMode" v-model="feat.approach" @blur="saveFeature(feat)" class="edit-textarea-detail"></textarea>
                    <p v-else class="detail-text" style="font-size: 13px;">{{ feat.approach || 'Standard approach.' }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>

    <!-- Elite Auth Modal -->
    <div v-if="showAuthModal" class="modal-backdrop">
      <div class="modal-content">
        <div class="brand-logo" style="margin: 0 auto 24px;">E</div>
        <h3>{{ authMode === 'login' ? 'BABA Authentication' : 'Member Registration' }}</h3>
        <p>Access the ECOMMITRA architecture blueprints.</p>
        <div v-if="authError" class="auth-error">{{ authError }}</div>
        <form @submit.prevent="handleAuth">
          <div class="input-group"><label>Email</label><input type="email" v-model="authEmail" class="elite-input" required /></div>
          <div class="input-group"><label>Password</label><input type="password" v-model="authPassword" class="elite-input" required /></div>
          <button type="submit" class="btn-primary">{{ authMode === 'login' ? 'Unlock Access' : 'Register' }}</button>
        </form>
        <div style="margin-top: 24px; text-align: center;"><a href="#" @click.prevent="authMode = authMode === 'login' ? 'signup' : 'login'" style="color: var(--brand-primary); font-size: 14px; font-weight: 700;">{{ authMode === 'login' ? 'Switch to Registration' : 'Return to Login' }}</a></div>
      </div>
    </div>
  </div>
</template>

<style>
.btn-edit-toggle { display: flex; align-items: center; gap: 8px; background: white; border: 1px solid var(--border-light); padding: 8px 16px; border-radius: 12px; font-weight: 800; font-size: 13px; color: var(--text-dark); cursor: pointer; transition: 0.3s; }
.btn-edit-toggle.active { background: var(--brand-primary); color: white; border-color: var(--brand-primary); box-shadow: var(--brand-glow); }
.btn-add-feat { margin-left: auto; background: var(--status-live); color: white; border: none; padding: 6px 12px; border-radius: 8px; font-weight: 800; font-size: 11px; cursor: pointer; }
.edit-input-title { flex: 1; border: 1px solid var(--brand-primary); padding: 4px 8px; border-radius: 6px; font-size: 18px; font-weight: 700; font-family: inherit; }
.edit-select { border: 1px solid var(--border-light); padding: 4px; border-radius: 6px; font-weight: 800; font-size: 11px; }
.edit-textarea-sub { width: 100%; border: 1px solid var(--border-light); border-radius: 6px; padding: 8px; margin-bottom: 12px; font-family: inherit; font-size: 14px; }
.edit-textarea-detail { width: 100%; border: 1px solid var(--border-light); border-radius: 6px; padding: 8px; min-height: 80px; font-family: inherit; font-size: 13px; }
.edit-input-tech { width: 100%; border: 1px solid var(--border-light); border-radius: 6px; padding: 4px 8px; font-size: 12px; font-family: 'JetBrains Mono', monospace; }
.btn-del { background: #fef2f2; color: #ef4444; border: 1px solid #fee2e2; width: 24px; height: 24px; border-radius: 6px; font-weight: 800; cursor: pointer; }
.auth-error { background: #fef2f2; color: #ef4444; padding: 12px; border-radius: 8px; margin-bottom: 20px; font-size: 13px; font-weight: 700; text-align: center; }
</style>
