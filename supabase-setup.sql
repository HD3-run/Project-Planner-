-- ECOMMITRA Product Architecture – Supabase Schema
-- Run this in your Supabase SQL Editor

-- 1. Sections table
CREATE TABLE IF NOT EXISTS sections (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  icon TEXT DEFAULT '📦',
  description TEXT DEFAULT '',
  color TEXT DEFAULT '#3b82f6',
  sort_order INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- 2. Features table
CREATE TABLE IF NOT EXISTS features (
  id SERIAL PRIMARY KEY,
  section_id TEXT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  icon TEXT DEFAULT '📦',
  status TEXT DEFAULT 'live' CHECK (status IN ('live','planned','future')),
  tag_extra TEXT,
  subtitle TEXT DEFAULT '',
  how_it_works TEXT,
  approach TEXT,
  tech JSONB DEFAULT '[]'::jsonb,
  capabilities JSONB DEFAULT '[]'::jsonb,
  flow JSONB DEFAULT '[]'::jsonb,
  impact TEXT,
  sort_order INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

-- 3. Auto-update timestamps
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = now(); RETURN NEW; END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS sections_updated ON sections;
CREATE TRIGGER sections_updated BEFORE UPDATE ON sections
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

DROP TRIGGER IF EXISTS features_updated ON features;
CREATE TRIGGER features_updated BEFORE UPDATE ON features
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- 4. Enable RLS and secure it for authenticated users only
ALTER TABLE sections ENABLE ROW LEVEL SECURITY;
ALTER TABLE features ENABLE ROW LEVEL SECURITY;

-- Drop existing policies to prevent 'already exists' errors
DROP POLICY IF EXISTS "Allow public read sections" ON sections;
DROP POLICY IF EXISTS "Allow public write sections" ON sections;
DROP POLICY IF EXISTS "Allow public read features" ON features;
DROP POLICY IF EXISTS "Allow public write features" ON features;

-- Create secure policies: Public can read, but only authenticated users can write
CREATE POLICY "Allow public read sections" ON sections FOR SELECT USING (true);
CREATE POLICY "Allow authenticated write sections" ON sections FOR ALL USING (auth.role() = 'authenticated') WITH CHECK (auth.role() = 'authenticated');

CREATE POLICY "Allow public read features" ON features FOR SELECT USING (true);
CREATE POLICY "Allow authenticated write features" ON features FOR ALL USING (auth.role() = 'authenticated') WITH CHECK (auth.role() = 'authenticated');

-- 5. Indexes
CREATE INDEX IF NOT EXISTS idx_features_section ON features(section_id);
CREATE INDEX IF NOT EXISTS idx_features_sort ON features(sort_order);

-- ========== SEED DATA ==========

INSERT INTO sections (id, title, icon, description, color, sort_order) VALUES
('core',         'Core System',              '⚡', 'Foundation layer powering all operations',           '#3b82f6', 1),
('inventory',    'Inventory System',          '📊', '4-mode tracking engine with real-time sync',         '#10b981', 2),
('invoices',     'Invoice System',            '🧾', 'Automated and manual invoice generation',            '#f59e0b', 3),
('reports',      'Reports & Analytics',       '📈', 'Multi-dimensional business intelligence',            '#8b5cf6', 4),
('ai',           'AI System',                 '🤖', 'Natural language interface to platform data',         '#06b6d4', 5),
('whatsapp',     'WhatsApp Commerce',         '💚', 'Chat-based sales channel',                           '#25d366', 6),
('marketplace',  'Multi-Platform Sync',       '🔗', 'Cross-marketplace inventory control hub',            '#f97316', 7),
('ai-bulk',      'AI Bulk Upload',            '🧠', 'Accept any format, AI maps it automatically',        '#f43f5e', 8),
('storefront-v2','Custom Storefronts',        '🎨', 'From hosted page to brand identity platform',        '#ec4899', 9),
('extras',       'Additional Planned Features','🚀','Analytics intelligence and loyalty systems',          '#a855f7', 10)
ON CONFLICT (id) DO NOTHING;

INSERT INTO features (section_id, title, icon, status, subtitle, how_it_works, tech, capabilities, flow, impact, sort_order) VALUES
('core', 'Order Management', '📦', 'live',
 'Full lifecycle order processing with multi-source support',
 'Express REST API (backend/orders/) → PostgreSQL (oms.orders). Session-authenticated routes handle CRUD. Orders support manual entry, CSV bulk upload, and storefront checkout. Each order is merchant-scoped via merchant_id.',
 '["Express Router","PostgreSQL","Session Auth","CSV Parser","WebSocket"]',
 '["Manual order entry via dashboard UI","Bulk order creation via CSV upload with schema validation","Order source tagging (manual, storefront, whatsapp, etc.)","Full lifecycle: Pending → Confirmed → Shipped → Delivered","Returns & refund flow with inventory auto-restoration","Real-time WebSocket status broadcasts to merchant room"]',
 '["Order Created","→","Payment Verified","→","Stock Deducted","→","Invoice Generated","→","Shipped","→","Delivered"]',
 'Handles 100% of revenue flow. Razorpay webhook + advisory locking ensures zero double-payments.', 1),

('core', 'Storefront System', '🏪', 'live',
 'Public catalog with checkout and Razorpay payment integration',
 'backend/public-catalog/ serves product data via public (no-auth) routes. Frontend PublicCatalog.tsx renders browsable catalog per merchant. Payment via dedicated payment-service/ microservice with Razorpay provider.',
 '["React SPA","Public REST API","Razorpay SDK","Payment Microservice","S3 Images"]',
 '["Public catalog per merchant at /catalog/{merchant_id}","Product browsing with images from S3","Template-based UI with light/dark mode + brand customization","Checkout flow with Razorpay payment gateway","Webhook-driven payment confirmation with idempotency keys","Advisory lock-based race condition prevention"]',
 '["Browse Catalog","→","Add to Cart","→","Checkout","→","Razorpay Pay","→","Webhook Confirms","→","Order Created"]',
 'Merchant''s online sales channel. Payment resilience hardened with PostgreSQL advisory locks to prevent double-initiation.', 2),

('core', 'Bulk CSV Operations', '📄', 'live',
 'Schema-based CSV processing for products and stock updates',
 'backend/inventory.ts exposes /upload-csv and /update-stock-csv. Uses multer disk storage to prevent RAM exhaustion on large uploads. processProductCSV() handles batch processing (500 items/batch) with image file co-upload support.',
 '["Multer (Disk Storage)","CSV Parser","Batch Processing","S3 Upload","Lambda Barcode"]',
 '["Product CSV upload with variant & attribute support","Stock update CSV for bulk quantity changes","Batch processing: 500 items per batch to prevent timeouts","Co-upload images alongside CSV (matched by filename)","Auto-barcode generation via generateProductBarcodeValue()","Detailed error reporting per row"]',
 '[]',
 'Current limitation: Users must follow a predefined CSV schema including JSON fields for attributes/variants. This is the key friction point that AI Bulk Upload (planned) will solve.', 3),

('inventory', '4-Mode Inventory Tracking', '🔢', 'live',
 'Flexible complexity: Stock Only → Items → Variants → Items+Variants',
 'Product.track_items flag + product_variants table control mode. DB trigger auto_create_inventory_items generates serial/batch items when track_items=true. inventory_variants table links variant-level stock. LPAD-based serial generation patched for 4000+ item batches.',
 '["PostgreSQL Triggers","Product Variants","Inventory Variants","Serial Generation","Batch Tracking"]',
 '["Mode A – Stock Only: Simple quantity tracking (track_items=false, no variants)","Mode B – Items Only: Serial/batch number tracking per unit (track_items=true)","Mode C – Variants: Size/Color/etc. variants with per-variant stock","Mode D – Items + Variants: Full serial tracking per variant combination","Auto-creation of inventory items via PostgreSQL trigger","Configurable batch sizes and manufacturing dates"]',
 '[]',
 'Covers every merchant type — from simple stock counters to fashion brands needing size×color×serial tracking.', 1),

('inventory', 'Real-Time Stock Sync', '🔄', 'live',
 'WebSocket-based live inventory updates across all connected clients',
 'Socket.IO rooms scoped to merchant:{merchantId}. Every stock mutation emits events like inventory-product-added, inventory-price-updated. Frontend listens and patches local state without reload.',
 '["Socket.IO","Merchant Rooms","Event Emitters","Cache Invalidation"]',
 '["Live stock count updates across dashboard tabs","Storefront stock availability auto-refresh","Price change broadcasts (cost + selling)","New product notifications to all merchant users","Cache invalidation via invalidateUserCache() on every mutation"]',
 '[]',
 'Eliminates stale data. Multiple staff members see the same stock counts simultaneously.', 2),

('invoices', 'Invoice Generation', '📃', 'live',
 'Auto-generated on payment, with manual creation support',
 'backend/invoices/ module with dedicated queries/services. Auto-triggers on payment confirmation webhook. Linked to order data. Supports GST calculations with inclusive/exclusive/no_gst pricing modes.',
 '["Express Router","PostgreSQL","GST Calculator","PDF Generation"]',
 '["Auto invoice generation on successful payment","Manual invoice creation from dashboard","GST-aware: inclusive, exclusive, and no_gst pricing modes","HSN code support for tax compliance","Linked with order and customer data","AI chatbot can fetch invoice details via natural language"]',
 '[]',
 'Handles Indian GST compliance automatically. Merchants don''t need separate billing software.', 1),

('reports', 'Reports Dashboard', '📊', 'live',
 'Product, sales, customer, channel, and financial analytics',
 'backend/reports.ts with extracted services: salesService, kpiService, productService, customerService, channelService, financialService. Cached with cacheMiddleware (60-300s TTL).',
 '["Service Layer Architecture","Chart.js Components","Cache Middleware","CSV Export"]',
 '["Sales reports with day/week/month grouping","KPI dashboard: AOV, conversion rates, revenue metrics","Top selling, unsold, and dead stock analysis","Category performance breakdown","Customer analytics: lifetime value, location, retention","Channel performance (manual vs storefront)","Financial P&L and tax breakdown reports","CSV export for sales data"]',
 '[]',
 'Full BI suite built-in. No need for external analytics tools for most merchants.', 1),

('ai', 'AI Chatbot (Read-Only)', '💬', 'live',
 'Natural language queries across orders, inventory, invoices, and KPIs',
 'Python microservice (python-chatbot/) with RAG architecture. Node backend/chatbot.ts proxies requests with session cookie forwarding. Agent uses 15+ specialized tools.',
 '["Python FastAPI","LLM Agent","RAG Pipeline","15+ API Tools","Session Proxy"]',
 '["Show me order #123 → fetches full order details","List low stock products → queries inventory API","What are my KPIs? → retrieves dashboard metrics","Show invoice for order 456 → fetches invoice data","Order status distribution → returns status breakdown","Supports orders, inventory, items, invoices, returns, and reports"]',
 '[]',
 'Merchants can query their entire business data in plain English. 15+ tools covering every data domain.', 1),

('ai', 'AI Write Layer', '✍️', 'planned',
 'Natural language → system actions (create orders, products, etc.)',
 NULL,
 '[]',
 '["Create order for Rahul, 2 units of XYZ → POST /api/orders","Add product: Shirt, ₹500, size M → POST /api/inventory/add-product","Update stock for SKU-001 to 50 → PATCH /api/inventory/:id/stock","Confirmation prompt before any write action","Full audit trail of AI-initiated actions"]',
 '[]',
 'Removes UI friction completely. Non-technical merchants can operate entirely via chat.', 2),

('whatsapp', 'WhatsApp Automation', '📱', 'planned',
 'Convert WhatsApp into a full sales + order channel',
 NULL,
 '[]',
 '["Sync product catalog to WhatsApp Business","Auto-reply chatbot for product queries and pricing","Direct order placement via chat conversation","Orders auto-created in ECOMMITRA with source=whatsapp","Stock availability checks in real-time","Order confirmation & tracking via WhatsApp messages"]',
 '[]',
 'Removes friction for non-tech customers. Leverages existing merchant behavior — India''s #1 communication channel.', 1),

('marketplace', 'Cross-Marketplace Engine', '🌐', 'future',
 'Central inventory control across Meesho, Amazon, and more',
 NULL,
 '[]',
 '["Bi-directional inventory sync with external marketplaces","Stock=10 → ECOMMITRA order → stock=8 → Meesho auto-updated","Reverse: Meesho order → ECOMMITRA stock decremented","Order aggregation dashboard across all channels","Unified fulfillment layer","Overselling prevention via atomic stock operations"]',
 '["Stock Change","→","Event Emitted","→","Sync Adapters","→","Marketplace APIs","→","Stock Updated"]',
 'Transforms ECOMMITRA from a tool into a control hub. Eliminates overselling across platforms.', 1),

('ai-bulk', 'AI-Powered Flexible Bulk Upload', '🪄', 'planned',
 'Accept any CSV/Excel format — AI interprets and maps fields automatically',
 NULL,
 '[]',
 '["Accept any CSV, Excel, or tabular format","AI reads column headers and sample data","Auto-maps: Product Name|Size|Price → name, variants, selling_price","Variant detection: recognizes size/color columns as variant attributes","Attribute structuring: converts flat columns to JSON attributes","Error correction: fixes common data issues before import","Preview step: shows mapped data before committing"]',
 '[]',
 'Removes the #1 onboarding friction. Makes bulk operations effortless for any merchant.', 1),

('storefront-v2', 'Custom Domains + Builder', '🏗️', 'future',
 'Custom subdomains, template builder, and SEO enablement',
 NULL,
 '[]',
 '["Custom subdomains: merchantname.ecommitra.com","Template builder with drag-and-drop layout customization","Existing: light/dark mode, brand color, logo customization","SEO-optimized merchant-specific URLs","Indexable pages with proper meta tags","Mobile-responsive storefront layouts"]',
 '[]',
 'Competes with Shopify-level branding. Moves ECOMMITRA from tool → platform.', 1),

('extras', 'AI Analytics Engine', '📉', 'future',
 'Dedicated intelligence layer for trend analysis and insights',
 NULL,
 '[]',
 '["Sales trend analysis with anomaly detection","Product-level performance insights","Customer behavior pattern recognition","Proactive alerts: Sales dropped 20% vs last week","Correlation analysis across data domains"]',
 '[]',
 'Turns data into actionable intelligence without merchants needing analytics expertise.', 1),

('extras', 'Loyalty & Rewards System', '🎁', 'future',
 'Gift cards, discount codes, reward points, and tier-based loyalty',
 NULL,
 '[]',
 '["Gift card issuance and redemption","Discount codes with usage limits and expiry","Points-per-purchase reward system","Tier-based loyalty (Bronze → Silver → Gold)","POS integration: give coupon at physical store","Online redemption on storefront checkout"]',
 '[]',
 'Drives repeat purchases. Bridges offline and online commerce — key for Indian merchants.', 2),

('extras', 'Advanced Visitor Analytics', '👁️', 'future',
 'Storefront visitor tracking, heatmaps, and conversion funnels',
 NULL,
 '[]',
 '["Visitor count and unique visitor tracking","Time spent on product pages","Click heatmaps on storefront","Conversion funnel: Visit → View → Cart → Purchase","Traffic source attribution"]',
 '[]',
 'Gives merchants Hotjar-level insights into their storefront performance.', 3);
