
package config

import (
	"log"
	"os"
	
	"ecommitra-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var DB *gorm.DB

// ConnectDatabase initializes the PostgreSQL database connection and runs migrations
func ConnectDatabase() {
	var err error
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set. Please add your Supabase Postgres URI to the .env file.")
	}

	// Open connection to Supabase Postgres
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")

	// 1. Mandatory deduplication of features (Run this FIRST before migrations)
	log.Println("🧹 Step 1: Cleaning up any existing duplicate features...")
	DB.Exec("DELETE FROM features a USING features b WHERE a.id < b.id AND a.title = b.title AND a.section_id = b.section_id")

	// 2. MIGRATE & SEED (Background)
	go func() {
		log.Println("🚀 Step 2: Running Auto Migrations in background...")
		err := DB.AutoMigrate(&models.User{}, &models.Section{}, &models.Feature{}, &models.Session{})
		if err != nil {
			log.Printf("⚠️ Background migration error: %v", err)
			return
		}
		seedInitialData()
		log.Println("✅ Database Migrations & Seeding Complete.")
	}()
}

// seedInitialData ensures the database isn't completely empty on first launch
func seedInitialData() {
	var count int64
	DB.Model(&models.Section{}).Count(&count)
	
	// We only seed if the database is truly fresh (no sections)
	if count > 0 {
		return
	}

	log.Println("🌱 Seeding Master ECOMMITRA Architecture Manifest...")

	sections := []models.Section{
		{ID: "core", Title: "Core Commerce Engine", Color: "#4f46e5", Icon: "⚡", SortOrder: 1, Description: "Foundation layer powering all commerce operations, multi-tenant isolation, and secure authentication."},
		{ID: "inventory", Title: "Intelligent Inventory Matrix", Color: "#10b981", Icon: "📊", SortOrder: 2, Description: "Advanced 4-mode inventory tracking system with real-time sync and dynamic asset generation."},
		{ID: "invoices", Title: "Automated Billing & Invoicing", Color: "#3b82f6", Icon: "📄", SortOrder: 3, Description: "GST-compliant dynamic PDF generation with custom merchant branding."},
		{ID: "reports", Title: "Reports & Analytics Dashboard", Color: "#8b5cf6", Icon: "📈", SortOrder: 4, Description: "Real-time KPI engine with 'Explainable Data' drill-down modals."},
		{ID: "ai", Title: "AI Executive Assistant (LangGraph)", Color: "#f59e0b", Icon: "🤖", SortOrder: 5, Description: "Agentic business intelligence using RAG and LangGraph for complex data reasoning."},
		{ID: "whatsapp", Title: "Integrated WhatsApp Commerce", Color: "#25d366", Icon: "💬", SortOrder: 6, Description: "Automated bot for order tracking and customer engagement."},
		{ID: "sync", Title: "Multi-Platform Sync Engine", Color: "#06b6d4", Icon: "🔄", SortOrder: 7, Description: "Real-time stock and price synchronization across all sales channels."},
		{ID: "bulk", Title: "AI-Powered Bulk Upload", Color: "#6366f1", Icon: "📁", SortOrder: 8, Description: "Intelligent CSV processing with batching and automated image mapping."},
		{ID: "storefront", Title: "Dynamic Catalog Storefronts", Color: "#ec4899", Icon: "🎨", SortOrder: 9, Description: "Template-driven B2C shopping experience with high-performance routing."},
		{ID: "returns", Title: "Returns & Reverse Logistics", Color: "#ef4444", Icon: "📦", SortOrder: 10, Description: "Automated return lifecycle with integrated inventory restocking."},
		{ID: "payment", Title: "High-Resilience Payment Service", Color: "#f97316", Icon: "💳", SortOrder: 11, Description: "Fault-tolerant payment processing microservice with Razorpay integration."},
		{ID: "gst", Title: "GST Verification Service", Color: "#14b8a6", Icon: "✅", SortOrder: 12, Description: "Real-time scraping integration for official GSTIN verification."},
		{ID: "employees", Title: "Employee RBAC Workflow", Color: "#64748b", Icon: "👥", SortOrder: 13, Description: "Tailored internal dashboards with granular status-change permissions."},
		{ID: "assets", Title: "Brand & Asset Management", Color: "#a855f7", Icon: "🖼️", SortOrder: 14, Description: "Self-service portal for professional branding and S3 asset orchestration."},
		{ID: "suppliers", Title: "Supplier & Vendor Portal", Color: "#eab308", Icon: "🤝", SortOrder: 15, Description: "Centralized management of supply chain partners and contact points."},
		{ID: "security", Title: "Security & Observability", Color: "#0f172a", Icon: "🛡️", SortOrder: 16, Description: "Enterprise-grade security layers and real-time log monitoring."},
		{ID: "roadmap", Title: "Future Roadmap", Color: "#94a3b8", Icon: "🚀", SortOrder: 17, Description: "Expanding the ecosystem with Marketplace and shared logistics capabilities."},
	}

	for _, s := range sections {
		DB.Create(&s)
	}

	features := []models.Feature{
		{SectionID: "core", Title: "Secure Auth (Phantom Tokens)", Status: "live", Subtitle: "Indian-compliant phone verification", Impact: "Merchant's operational foundation. Hardened security.", HowItWorks: "Uses stateful session management backed by a secure DB pool.", Approach: "RBAC enforced at the route level.", Tech: `["Node.js", "Express", "PostgreSQL", "TypeScript"]`},
		{SectionID: "inventory", Title: "4-Mode Inventory Tracking", Status: "live", Subtitle: "Stock Only -> Items -> Variants", Impact: "Eliminates 'Ghost Stock' and overselling.", HowItWorks: "Database triggers ensure consistency across parent products.", Approach: "Atomic stock updates with row-level locking (FOR UPDATE).", Tech: `["PostgreSQL Triggers", "AWS S3", "QR Engine"]`},
		{SectionID: "ai", Title: "LangGraph Business Agent", Status: "live", Subtitle: "ReAct-style agent with RAG", Impact: "Turns 'Data' into 'Strategy'. Immediate answers.", HowItWorks: "Python microservice with 30+ specialized tools.", Approach: "Cyclic reasoning to verify facts across data points.", Tech: `["Python", "LangGraph", "ChromaDB", "OpenAI"]`},
		{SectionID: "payment", Title: "Fault-Tolerant Payment Ritual", Status: "live", Subtitle: "Razorpay + Resurrection Logic", Impact: "Protects merchant revenue. 99.9% success rate.", HowItWorks: "Dedicated microservice using the Outbox Pattern.", Approach: "Advisory Locking in PostgreSQL to serialize requests.", Tech: `["Razorpay SDK", "Advisory Locks", "Webhooks"]`},
		{SectionID: "roadmap", Title: "Multi-vendor Marketplace", Status: "planned", Subtitle: "Unified commerce hub", Impact: "Scale to thousands of vendors.", HowItWorks: "Aggregating individual catalogs into a shared storefront.", Approach: "Shared order routing and commission logic.", Tech: `["Next.js", "Redis", "ElasticSearch"]`},
	}

	for _, f := range features {
		DB.Create(&f)
	}
	
	log.Println("✅ Master Seed Complete.")
}
