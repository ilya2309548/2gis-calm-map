package db

import (
	"2gis-calm-map/api/config"
	"2gis-calm-map/api/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.Config) {
	database, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	DB = database

	log.Println("Database connected")

	// MIGRATION: автоматически создаёт таблицы, если их нет
	if err := DB.AutoMigrate(&model.User{}, &model.UserParams{}, &model.Organization{}, &model.OrganizationParams{}, &model.OrganizationComment{}); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	// Удаляем уникальный индекс по owner_id, если он был создан ранее (чтобы админы могли иметь несколько организаций)
	// Предполагаем стандартное имя индекса gorm: idx_organizations_owner_id
	if DB.Migrator().HasIndex(&model.Organization{}, "idx_organizations_owner_id") {
		if err := DB.Exec("DROP INDEX IF EXISTS idx_organizations_owner_id;").Error; err != nil {
			log.Println("warn: failed to drop old unique index idx_organizations_owner_id:", err)
		}
	}
	// Создаём неуникальный индекс для ускорения выборок по owner_id (если ещё нет)
	if !DB.Migrator().HasIndex(&model.Organization{}, "idx_org_owner_id") {
		if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_org_owner_id ON organizations(owner_id);").Error; err != nil {
			log.Println("warn: failed to create index idx_org_owner_id:", err)
		}
	}

	log.Println("Database connected, migrated, indexes adjusted")
}
