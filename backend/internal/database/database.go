package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDatabase 初始化数据库连接
func InitDatabase() (*sql.DB, error) {
	// 确保数据目录存在
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	// 数据库文件路径
	dbPath := filepath.Join(dataDir, "nav-panel.db")

	// 连接数据库
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate 运行数据库迁移
func Migrate(db *sql.DB) error {
	// 启用外键约束
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	// 创建分组表
	groupsSQL := `
	CREATE TABLE IF NOT EXISTS groups (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		color TEXT,
		icon TEXT,
		sort_order INTEGER DEFAULT 0,
		is_collapsed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(groupsSQL); err != nil {
		return err
	}

	// 创建网站表
	websitesSQL := `
	CREATE TABLE IF NOT EXISTS websites (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		icon TEXT,
		description TEXT,
		category TEXT,
		group_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE SET NULL
	);`

	if _, err := db.Exec(websitesSQL); err != nil {
		return err
	}

	// 创建搜索引擎表
	searchEnginesSQL := `
	CREATE TABLE IF NOT EXISTS search_engines (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		icon TEXT,
		placeholder TEXT,
		is_default BOOLEAN DEFAULT FALSE
	);`

	if _, err := db.Exec(searchEnginesSQL); err != nil {
		return err
	}

	// 检查并添加新字段（向后兼容）
	if err := addColumnIfNotExists(db, "search_engines", "icon", "TEXT"); err != nil {
		return err
	}
	if err := addColumnIfNotExists(db, "websites", "sort_order", "INTEGER DEFAULT 0"); err != nil {
		return err
	}
	if err := addColumnIfNotExists(db, "websites", "group_id", "TEXT"); err != nil {
		return err
	}

	// 创建应用设置表
	settingsSQL := `
	CREATE TABLE IF NOT EXISTS app_settings (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		theme TEXT DEFAULT 'light',
		layout_config TEXT,
		background_config TEXT,
		card_style_config TEXT,
		search_config TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(settingsSQL); err != nil {
		return err
	}

	// 插入默认设置（如果不存在）
	defaultSettingsSQL := `
	INSERT OR IGNORE INTO app_settings (id, theme, layout_config, background_config, card_style_config, search_config)
	VALUES (1, 'light', 
		'{"columns":6,"cardSize":"medium","showLabels":true,"gap":20}',
		'{"type":"gradient","value":"linear-gradient(135deg, #667eea 0%, #764ba2 100%)"}',
		'{"borderRadius":12,"opacity":0.9,"shadow":true}',
		'{"enabled":true,"defaultEngineId":"google","openInNewTab":true}'
	);`

	if _, err := db.Exec(defaultSettingsSQL); err != nil {
		return err
	}

	// 插入默认搜索引擎（如果不存在）
	defaultEngines := []struct {
		id, name, url, placeholder string
		isDefault                  bool
	}{
		{"google", "Google", "https://www.google.com/search?q={}", "使用 Google 搜索...", true},
		{"baidu", "百度", "https://www.baidu.com/s?wd={}", "使用百度搜索...", false},
		{"bing", "Bing", "https://www.bing.com/search?q={}", "使用 Bing 搜索...", false},
		{"github", "GitHub", "https://github.com/search?q={}", "在 GitHub 中搜索...", false},
	}

	for _, engine := range defaultEngines {
		engineSQL := `
		INSERT OR IGNORE INTO search_engines (id, name, url, placeholder, is_default)
		VALUES (?, ?, ?, ?, ?);`
		
		if _, err := db.Exec(engineSQL, engine.id, engine.name, engine.url, engine.placeholder, engine.isDefault); err != nil {
			return err
		}
	}

	return nil
}

// addColumnIfNotExists 检查列是否存在，如果不存在则添加
func addColumnIfNotExists(db *sql.DB, tableName, columnName, columnType string) error {
	// 检查列是否存在
	exists, err := columnExists(db, tableName, columnName)
	if err != nil {
		return err
	}
	
	// 如果列不存在，则添加
	if !exists {
		sql := "ALTER TABLE " + tableName + " ADD COLUMN " + columnName + " " + columnType
		_, err := db.Exec(sql)
		return err
	}
	
	return nil
}

// columnExists 检查指定表中的列是否存在
func columnExists(db *sql.DB, tableName, columnName string) (bool, error) {
	query := "PRAGMA table_info(" + tableName + ")"
	rows, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, dfltValue, pk interface{}
		
		err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltValue, &pk)
		if err != nil {
			return false, err
		}
		
		if name == columnName {
			return true, nil
		}
	}
	
	return false, nil
} 