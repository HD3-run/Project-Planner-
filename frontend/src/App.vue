<template>
  <div class="app-wrapper">
    <!-- Ambient Animated Background -->
    <div class="ambient-bg"></div>
    <div class="ambient-glow-1"></div>
    <div class="ambient-glow-2"></div>

    <!-- Sticky Header -->
    <header class="header">
      <div class="brand">
        <div class="brand-icon">E</div>
        <div class="brand-text">
          <h1>ECOMMITRA</h1>
          <p>Product Architecture</p>
        </div>
      </div>

      <div class="header-controls">
        <div class="progress-container" title="Platform Completion">
          <div class="progress-bar-wrap">
            <div class="progress-bar" :style="{ width: completionPercentage + '%' }"></div>
          </div>
          <div class="progress-text">{{ completionPercentage }}% Live</div>
        </div>
        
        <button class="btn-primary" :class="{ active: editMode }" @click="toggleEditMode">
          <svg v-if="!editMode" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
          <svg v-else width="16" height="16" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
          {{ editMode ? 'Lock Architecture' : 'Edit Architecture' }}
        </button>
      </div>
    </header>

    <div class="layout">
      <!-- Sidebar Navigation -->
      <aside class="sidebar">
        <nav class="sidebar-nav">
          <a v-for="section in sections" :key="'nav-'+section.id" class="nav-link" :class="{ active: activeSection === section.id }" @click.prevent="scrollTo(section.id)">
            <span>{{ section.icon }}</span>
            {{ section.title }}
            <span class="nav-count">{{ section.features?.length || 0 }}</span>
          </a>
        </nav>
      </aside>

      <!-- Main Content Area -->
      <main class="main-content">
        
        <div v-if="loading" class="loader-view">
          <div class="spinner"></div>
          <p style="margin-top: 16px; font-weight: 600; color: var(--text-muted);">Syncing with Supabase...</p>
        </div>

        <div v-else-if="error" class="loader-view">
          <p style="color: #ef4444; font-weight: 700;">❌ Connection Error: {{ error }}</p>
        </div>

        <div v-else>
          <!-- Hero Banner -->
          <div class="hero">
            <h2>The Blueprint of <span>ECOMMITRA</span></h2>
            <p>An interactive map of our technical infrastructure, live capabilities, and future roadmap.</p>
          </div>

          <!-- Mind Map Visualization -->
          <div class="visual-map-wrapper" v-if="sections.length">
            <div class="tree-container">
              <div class="tree-root">🌍 PLATFORM CORE</div>
              <div class="tree-branches" :style="'--bcount: ' + sections.length">
                
                <div v-for="sec in sections" :key="'tree-'+sec.id" class="tree-branch">
                  <div class="tree-section-node" :style="`border-top-color: ${sec.color}`" @click="scrollTo(sec.id)">
                    <span>{{ sec.icon }}</span> {{ sec.title }}
                  </div>
                  
                  <div class="tree-features-spine">
                    <div v-for="feat in sec.features" :key="'tfeat-'+feat.id" class="tree-feature-node" :class="feat.status">
                      {{ feat.title }}
                    </div>
                  </div>
                </div>

              </div>
            </div>
          </div>

          <!-- Sections -->
          <div v-for="section in sections" :key="section.id" :id="section.id" class="section">
            
            <div class="section-header">
              <div class="section-title-wrap">
                <div class="section-icon-large" :style="`color: ${section.color}`">{{ section.icon }}</div>
                <div>
                  <h3 class="section-title">{{ section.title }}</h3>
                  <p class="section-desc">{{ section.description || section.desc }}</p>
                </div>
              </div>
            </div>

            <!-- Features -->
            <div v-for="feature in section.features" :key="feature.id" class="feature-card" :class="{ 'is-expanded': expandedCards.includes(feature.id) }">
              
              <div class="card-header" @click="toggleCard(feature.id)">
                <div class="card-icon-box" :style="`color: ${section.color}`">{{ feature.icon }}</div>
                
                <div class="card-info">
                  <div class="card-title-row">
                    <div class="card-title">
                      <input v-if="editMode" v-model="feature.title" class="editable-input" @click.stop @input="debouncedUpdate(feature)" style="width: 250px; padding: 4px 8px;" />
                      <template v-else>{{ feature.title }}</template>
                    </div>
                    
                    <select v-if="editMode" v-model="feature.status" class="editable-select" @click.stop @change="updateStatus(feature)" style="width: 120px; padding: 4px 8px; font-weight: 700;">
                      <option value="live">LIVE</option>
                      <option value="planned">PLANNED</option>
                      <option value="future">FUTURE</option>
                    </select>
                    <span v-else class="status-badge" :class="feature.status">{{ feature.status }}</span>
                  </div>
                  
                  <div class="card-subtitle">
                    <input v-if="editMode" v-model="feature.subtitle" class="editable-input" @click.stop @input="debouncedUpdate(feature)" style="width: 100%; padding: 4px 8px;" />
                    <template v-else>{{ feature.subtitle }}</template>
                  </div>
                </div>

                <button v-if="editMode" @click.stop="deleteFeature(feature.id)" class="btn-primary" style="background: #fee2e2; color: #ef4444; box-shadow: none;">Delete</button>
                
                <div class="card-chevron">
                  <svg width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"></path></svg>
                </div>
              </div>

              <!-- Card Details -->
              <div class="card-body" :style="{ maxHeight: expandedCards.includes(feature.id) ? '2000px' : '0' }">
                <div class="card-body-inner">
                  
                  <!-- Business Impact (Highlight) -->
                  <div class="detail-block highlight" style="margin-bottom: 24px;" v-if="editMode || feature.impact">
                    <div class="detail-label"><svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg> Business Impact</div>
                    <textarea v-if="editMode" v-model="feature.impact" class="editable-textarea" @input="debouncedUpdate(feature)"></textarea>
                    <div v-else class="detail-text">{{ feature.impact }}</div>
                  </div>

                  <div class="details-grid">
                    <!-- How It Works -->
                    <div class="detail-block">
                      <div class="detail-label"><svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"></path></svg> How It Works</div>
                      <textarea v-if="editMode" v-model="feature.how_it_works" class="editable-textarea" @input="debouncedUpdate(feature)"></textarea>
                      <div v-else class="detail-text">{{ feature.how_it_works || 'No description provided.' }}</div>
                    </div>

                    <!-- Approach -->
                    <div class="detail-block">
                      <div class="detail-label"><svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path></svg> Implementation Approach</div>
                      <textarea v-if="editMode" v-model="feature.approach" class="editable-textarea" @input="debouncedUpdate(feature)"></textarea>
                      <div v-else class="detail-text">{{ feature.approach || 'No technical approach provided.' }}</div>
                    </div>
                  </div>

                  <!-- Tech Stack & Capabilities -->
                  <div class="details-grid">
                    <div class="detail-block" v-if="editMode || (feature.tech && feature.tech.length)">
                      <div class="detail-label">Tech Stack</div>
                      <input v-if="editMode" :value="(feature.tech || []).join(', ')" class="editable-input" @input="e => updateTech(feature, e.target.value)" placeholder="Vue, Node, PostgreSQL..." />
                      <div v-else class="tech-stack">
                        <span v-for="t in feature.tech" :key="t" class="tech-pill">{{ t }}</span>
                      </div>
                    </div>

                    <div class="detail-block" v-if="editMode || (feature.capabilities && feature.capabilities.length)">
                      <div class="detail-label">Key Capabilities</div>
                      <textarea v-if="editMode" :value="(feature.capabilities || []).join('\n')" class="editable-textarea" @input="e => updateCapabilities(feature, e.target.value)" placeholder="One per line..."></textarea>
                      <ul v-else class="capability-list">
                        <li v-for="c in feature.capabilities" :key="c">{{ c }}</li>
                      </ul>
                    </div>
                  </div>

                </div>
              </div>
            </div>

            <button v-if="editMode" class="btn-add" @click="addFeature(section.id)">
              + Add Feature to {{ section.title }}
            </button>
          </div>

        </div>
      </main>
    </div>

    <!-- Toast Notification -->
    <div v-if="saveStatus.text" class="save-status" :style="{ color: saveStatus.type === 'error' ? '#ef4444' : '#10b981' }">
      {{ saveStatus.text }}
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { supabase } from './supabase'

const sections = ref([])
const loading = ref(true)
const error = ref(null)
const editMode = ref(false)
const expandedCards = ref([])
const activeSection = ref('')
const saveStatus = ref({ type: '', text: '' })

let saveTimeouts = {}

// Computed stats
const liveCount = computed(() => {
  let count = 0
  sections.value.forEach(s => s.features?.forEach(f => { if(f.status === 'live') count++ }))
  return count
})

const plannedCount = computed(() => {
  let count = 0
  sections.value.forEach(s => s.features?.forEach(f => { if(f.status !== 'live') count++ }))
  return count
})

const completionPercentage = computed(() => {
  const total = liveCount.value + plannedCount.value
  if(total === 0) return 0
  return Math.round((liveCount.value / total) * 100)
})

// Methods
const toggleEditMode = () => { editMode.value = !editMode.value }

const toggleCard = (id) => {
  if (expandedCards.value.includes(id)) {
    expandedCards.value = expandedCards.value.filter(i => i !== id)
  } else {
    expandedCards.value.push(id)
  }
}

const scrollTo = (id) => {
  activeSection.value = id
  const el = document.getElementById(id)
  if(el) {
    const y = el.getBoundingClientRect().top + window.scrollY - 120
    window.scrollTo({top: y, behavior: 'smooth'})
  }
}

const showSaved = () => { 
  saveStatus.value = { type: 'saved', text: '✅ Changes Saved' }
  setTimeout(() => { saveStatus.value = { type: '', text: '' } }, 2000)
}
const showError = () => { 
  saveStatus.value = { type: 'error', text: '❌ Failed to save' }
  setTimeout(() => { saveStatus.value = { type: '', text: '' } }, 3000)
}

const loadData = async () => {
  try {
    const [secRes, featRes] = await Promise.all([
      supabase.from('sections').select('*').order('sort_order'),
      supabase.from('features').select('*').order('sort_order')
    ])
    if (secRes.error) throw secRes.error
    if (featRes.error) throw featRes.error

    sections.value = secRes.data.map(s => ({
      ...s,
      features: featRes.data.filter(f => f.section_id === s.id).map(f => ({
        ...f,
        tech: typeof f.tech === 'string' ? JSON.parse(f.tech) : (f.tech || []),
        capabilities: typeof f.capabilities === 'string' ? JSON.parse(f.capabilities) : (f.capabilities || [])
      }))
    }))
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

// Updates
const debouncedUpdate = (feature) => {
  if (saveTimeouts[feature.id]) clearTimeout(saveTimeouts[feature.id])
  saveTimeouts[feature.id] = setTimeout(() => { saveFeatureToDB(feature) }, 800)
}

const updateStatus = (feature) => saveFeatureToDB(feature)

const updateTech = (feature, val) => {
  feature.tech = val.split(',').map(s => s.trim()).filter(Boolean)
  debouncedUpdate(feature)
}

const updateCapabilities = (feature, val) => {
  feature.capabilities = val.split('\n').filter(s => s.trim())
  debouncedUpdate(feature)
}

const saveFeatureToDB = async (feature) => {
  const dbUpdates = {
    title: feature.title, subtitle: feature.subtitle, status: feature.status,
    how_it_works: feature.how_it_works, approach: feature.approach, impact: feature.impact,
    tech: JSON.stringify(feature.tech || []), capabilities: JSON.stringify(feature.capabilities || [])
  }
  const { error } = await supabase.from('features').update(dbUpdates).eq('id', feature.id)
  if (error) { console.error(error); showError(); } else { showSaved(); }
}

const addFeature = async (sectionId) => {
  const newFeat = {
    section_id: sectionId, title: 'New Feature', icon: '✨', status: 'planned',
    subtitle: 'Describe the feature...', capabilities: '[]', tech: '[]', sort_order: 99
  }
  const { error } = await supabase.from('features').insert(newFeat)
  if (error) { showError(); return; }
  showSaved(); loadData()
}

const deleteFeature = async (featureId) => {
  if (!confirm('Are you sure you want to delete this feature?')) return
  const { error } = await supabase.from('features').delete().eq('id', featureId)
  if (error) { showError(); return; }
  showSaved(); loadData()
}

onMounted(() => { loadData() })
</script>
