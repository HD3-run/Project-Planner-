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

const isBaba = computed(() => {
  const adminRole = import.meta.env.VITE_ADMIN_ROLE || 'admin'
  return userRole.value === adminRole
})

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
      const tech = Array.isArray(f.tech) ? f.tech : []
      const q = searchQuery.value.toLowerCase()
      const matchesSearch = !searchQuery.value || 
        s.title.toLowerCase().includes(q) ||
        f.title.toLowerCase().includes(q) ||
        (f.subtitle || '').toLowerCase().includes(q) ||
        (f.impact || '').toLowerCase().includes(q) ||
        tech.some(t => typeof t === 'string' && t.toLowerCase().includes(q))
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
    if (!res.ok) {
      throw new Error(data.error || 'Authentication failed')
    }
    
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
      id: s.id, title: s.title, icon: s.icon, color: s.color, description: s.description,
      features: (s.features || []).map(f => ({
        id: f.id, section_id: f.section_id, title: f.title, status: f.status, subtitle: f.subtitle,
        impact: f.impact, how_it_works: f.how_it_works, approach: f.approach,
        tech: Array.isArray(f.tech) ? f.tech : JSON.parse(f.tech || '[]')
      }))
    }))
  } catch (err) { console.error(err) }
}

// BABA Editing Power
const saveFeature = async (feature) => {
  if (!isBaba.value) {
    console.warn('Save blocked: User is not a BABA.')
    return
  }
  isUpdating.value = true
  console.log('Syncing feature to DB:', feature)
  
  try {
    const res = await apiFetch(`${API_URL}/architecture/feature`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...feature,
        id: Number(feature.id),
        section_id: String(feature.section_id),
        tech: JSON.stringify(feature.tech)
      })
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Save failed')
    console.log('Feature synced successfully:', data)
  } catch (err) { 
    console.error('Persistence Error:', err.message)
    alert('Save Failed: ' + err.message)
  }
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
    impact: '', how_it_works: '', approach: ''
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
  if (editMode.value) return
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

    <!-- Polished Header -->
    <header class="top-nav">
      <div class="header-left">
        <div class="brand-wrap">
          <div class="brand-logo">E</div>
          <div class="brand-info"><h1>ECOMMITRA</h1><p>Product Architecture</p></div>
        </div>
      </div>

      <div v-if="sessionToken" class="header-center">
        <div class="search-and-filter">
          <div class="search-container">
            <svg class="search-icon-svg" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
            <input type="text" v-model="searchQuery" class="search-input" placeholder="Search architecture..." />
          </div>

          <div class="status-filter-group">
            <button v-for="s in ['all', 'live', 'planned', 'future']" :key="s" @click="statusFilter = s" class="filter-pill" :class="{ active: statusFilter === s }">
              {{ s === 'all' ? 'All' : s.charAt(0).toUpperCase() + s.slice(1) }}
            </button>
          </div>
        </div>
      </div>

      <div class="header-right" v-if="sessionToken">
        <div class="action-cluster">
          <button v-if="isBaba" @click="editMode = !editMode" class="btn-edit-toggle" :class="{ active: editMode }">
            <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
            <span>Edit</span>
          </button>

          <div class="view-toggle">
            <button @click="activeView = 'list'" :class="{ active: activeView === 'list' }">List</button>
            <button @click="activeView = 'map'" :class="{ active: activeView === 'map' }">Map</button>
          </div>

          <div class="progress-indicator">
            <div class="prog-bar-bg"><div class="prog-bar-fill" :style="{ width: completionPercentage + '%' }"></div></div>
            <span class="prog-text">{{ completionPercentage }}%</span>
          </div>

          <button @click="handleLogout" class="btn-logout" title="Logout">
            <svg width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
          </button>
        </div>
      </div>
    </header>

    <div v-if="sessionToken" class="layout">
      <!-- Elite Sidebar -->
      <aside class="sidebar">
        <div class="sidebar-section-label">Navigation</div>
        <nav class="sidebar-nav">
          <div v-for="s in filteredSections" :key="'nav-'+s.id" class="nav-item" :class="{ active: activeSection === s.id }" @click="scrollTo(s.id)">
            <span>{{ s.icon }}</span><span class="nav-label">{{ s.title }}</span><span class="nav-badge">{{ s.features.length }}</span>
          </div>
        </nav>
      </aside>

      <!-- Main Elite Content -->
      <main class="main-content">
        <!-- Interactive Architecture Map -->
        <div v-if="activeView === 'map'" class="map-container">
          <div class="map-root">
            <div class="map-core">
              <div class="core-node">ECOMMITRA</div>
            </div>
            <div class="map-branches">
              <div v-for="section in filteredSections" :key="'map-'+section.id" class="map-section">
                <div class="section-node" :style="{ borderColor: section.color }">
                  <span class="node-icon">{{ section.icon }}</span>
                  <span class="node-title">{{ section.title }}</span>
                </div>
                <div class="map-features">
                  <div v-for="feat in section.features" :key="'mapfeat-'+feat.id" class="feat-node" :class="feat.status">
                    {{ feat.title }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

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
              <div v-for="feat in section.features" :key="feat.id || 'new'" class="feature-card" :class="{ editing: editMode, expanded: expandedCards.includes(feat.id) }" @click="toggleCard(feat.id)">
                <div class="card-top">
                  <div class="card-title-group">
                    <div class="feature-icon-box">{{ section.icon }}</div>
                    <input v-if="editMode" v-model="feat.title" @blur="saveFeature(feat)" class="edit-input-title" />
                    <h4 v-else>{{ feat.title }}</h4>
                  </div>
                  <div style="display: flex; align-items: center; gap: 12px;">
                    <select v-if="editMode" v-model="feat.status" @change="saveFeature(feat)" class="edit-select">
                      <option value="live">Live</option><option value="planned">Planned</option><option value="future">Future</option>
                    </select>
                    <span v-else class="status-pill" :class="feat.status">{{ feat.status }}</span>
                    <button v-if="editMode && isBaba" @click.stop="deleteFeature(feat.id, section.id)" class="btn-del">×</button>
                  </div>
                </div>
                <!-- Card Subtitle Edit -->
                <div v-if="editMode" style="margin-top: 12px;">
                  <span class="edit-label">Short Description</span>
                  <textarea v-model="feat.subtitle" @blur="saveFeature(feat)" class="edit-textarea-sub"></textarea>
                </div>
                <p v-else class="feature-subtitle">{{ feat.subtitle }}</p>

                <div v-if="expandedCards.includes(feat.id) || editMode" class="card-details">
                  <div class="detail-section">
                    <span v-if="editMode" class="edit-label">Business Impact</span>
                    <h5 v-else>Business Impact</h5>
                    <textarea v-if="editMode" v-model="feat.impact" @blur="saveFeature(feat)" class="edit-textarea-detail"></textarea>
                    <p v-else class="detail-text">{{ feat.impact || 'Standard platform capability.' }}</p>
                    
                    <div style="margin-top: 20px;">
                      <span v-if="editMode" class="edit-label">Logic & Architecture</span>
                      <h5 v-else>Logic & Architecture</h5>
                      <textarea v-if="editMode" v-model="feat.how_it_works" @blur="saveFeature(feat)" class="edit-textarea-detail"></textarea>
                      <p v-else class="detail-text">{{ feat.how_it_works || 'Logic contained in core services.' }}</p>
                    </div>
                  </div>
                  <div class="detail-section">
                    <span v-if="editMode" class="edit-label">Tech Stack (JSON Array)</span>
                    <h5 v-else>Tech Stack</h5>
                    <input v-if="editMode" :value="JSON.stringify(feat.tech)" @blur="e => { try { feat.tech = JSON.parse(e.target.value); saveFeature(feat) } catch(err){} }" class="edit-input-tech" />
                    <div v-else class="tech-chips">
                      <span v-for="t in feat.tech" :key="t" class="tech-chip">{{ t }}</span>
                    </div>

                    <div style="margin-top: 20px;">
                      <span v-if="editMode" class="edit-label">Implementation Approach</span>
                      <h5 v-else>Implementation Approach</h5>
                      <textarea v-if="editMode" v-model="feat.approach" @blur="saveFeature(feat)" class="edit-textarea-detail"></textarea>
                      <p v-else class="detail-text" style="font-size: 13px;">{{ feat.approach || 'Standard approach.' }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>

    <!-- Auth Modal -->
    <div v-if="showAuthModal" class="modal-backdrop">
      <div class="modal-content">
        <div class="brand-logo" style="margin: 0 auto 24px;">E</div>
        <h3>{{ authMode === 'login' ? 'Authentication' : 'Member Registration' }}</h3>
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
/* Header Optimization */
.top-nav { display: grid; grid-template-columns: 1fr auto 1fr; gap: 20px; }
.header-center { display: flex; justify-content: center; width: 100%; }
.header-right { display: flex; justify-content: flex-end; }
.action-cluster { display: flex; align-items: center; gap: 12px; background: rgba(255,255,255,0.5); padding: 6px; border-radius: 16px; border: 1px solid var(--border-light); }

.view-toggle { display: flex; background: #f1f5f9; padding: 2px; border-radius: 10px; }
.view-toggle button { border: none; background: transparent; padding: 6px 14px; border-radius: 8px; font-size: 12px; font-weight: 800; color: var(--text-muted); cursor: pointer; }
.view-toggle button.active { background: white; color: var(--brand-primary); box-shadow: var(--shadow-sm); }

.progress-indicator { display: flex; align-items: center; gap: 8px; padding: 0 10px; border-left: 1px solid var(--border-light); border-right: 1px solid var(--border-light); }
.prog-bar-bg { width: 60px; height: 6px; background: #f1f5f9; border-radius: 10px; overflow: hidden; }
.prog-bar-fill { height: 100%; background: var(--brand-gradient); border-radius: 10px; }
.prog-text { font-size: 11px; font-weight: 800; color: var(--text-dark); }

/* Map UI */
.map-container { padding: 40px; min-height: 80vh; overflow-x: auto; }
.map-root { display: flex; flex-direction: column; align-items: center; gap: 60px; }
.core-node { background: var(--brand-gradient); color: white; padding: 20px 40px; border-radius: 50px; font-weight: 900; font-size: 24px; box-shadow: var(--brand-glow); }
.map-branches { display: flex; gap: 40px; align-items: flex-start; justify-content: center; flex-wrap: wrap; }
.map-section { display: flex; flex-direction: column; align-items: center; gap: 20px; width: 220px; }
.section-node { width: 100%; background: white; border: 2px solid; padding: 16px; border-radius: 16px; text-align: center; display: flex; flex-direction: column; gap: 8px; box-shadow: var(--shadow-md); }
.node-icon { font-size: 24px; }
.node-title { font-size: 13px; font-weight: 800; color: var(--text-dark); }
.map-features { display: flex; flex-direction: column; gap: 8px; width: 100%; }
.feat-node { background: #f8fafc; padding: 8px 12px; border-radius: 10px; font-size: 11px; font-weight: 700; color: var(--text-base); border: 1px solid var(--border-light); text-align: center; }
.feat-node.live { border-left: 3px solid var(--status-live); }
.feat-node.planned { border-left: 3px solid var(--status-planned); }
</style>
