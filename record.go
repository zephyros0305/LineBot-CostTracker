package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Class  string
	Cost   uint64
	Memo   string
	UserID string
}

type StatData struct {
	Class   string
	CostSum uint64
}

func ConvertToRecord(data OperationData) *Record {
	var record Record

	record.Class = data.CostType
	record.Cost = data.Number
	record.Memo = data.Memo

	return &record
}

func (r *Record) Save() bool {
	db := connectDB()
	if db != nil {
		db.Create(r)
		closeDB(db)

		return true
	} else {
		return false
	}
}

func GetLastRecords(num uint) []Record {
	if num > 25 {
		num = 25
	}

	var records []Record

	db := connectDB()
	if db != nil {
		db.Order("id desc").Limit(int(num)).Find(&records)
	}

	return records
}

func GetStatData() []StatData {
	var result []StatData

	db := connectDB()

	if db != nil {
		rows, err := db.Model(&Record{}).Select("class, sum(cost) as costSum").Group("class").Rows()

		if err == nil {
			defer rows.Close()

			for rows.Next() {
				log.Println("rows=", rows)
				var temp StatData
				db.ScanRows(rows, &temp)
				result = append(result, temp)
			}
		} else {
			log.Println("GetStatData err=", err)
		}
	}

	log.Println("StatData=", result)

	return result
}

func connectDB() *gorm.DB {
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Taipei"
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database.")
		return nil
	}

	return db
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("failed to close db.")
	} else {
		sqlDB.Close()
	}
}

func init() {
	db := connectDB()
	if !db.Migrator().HasTable(&Record{}) {
		db.AutoMigrate(&Record{})
	}
	closeDB(db)
}
