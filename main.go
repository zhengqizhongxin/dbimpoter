package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
)

type ValueDB struct {
	ID      int    `gorm:"primary_key" json:"id"`
	Family  string `json:"family"`
	Package string `json:"package"`
	Value   string `json:"value"`
	CVE     string `json:"cve"`
}

type Res struct {
	ID      int    `gorm:"primary_key" json:"id"`
	CveID   string `json:"cve_id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Global struct {
	Config struct {
		Database string `yaml:"database"`
		Server   string `yaml:"server"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"config"`
}

func loadConfig() *Global {
	content, err := ioutil.ReadFile("configure.yaml")
	if err != nil {
		log.Println("open Err:", err)
	}
	output := &Global{}
	err = yaml.Unmarshal(content, &output)
	if err != nil {
		log.Println("conf parse Err:", err)
	}
	return output
}

func getDSN(g *Global) string {
	var dsn string
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", g.Config.Username, g.Config.Password, g.Config.Server, g.Config.Database)
	return dsn
}

func dbInit(g *Global) error {
	var err error
	var sqlDB *sql.DB

	dsn := getDSN(g)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("gorm.Open err: %v", err)
		return err
	}

	iDB = db

	sqlDB, err = iDB.DB()
	if err != nil {
		log.Fatalf("DB.Setup Err: %v", err)
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	err = iDB.AutoMigrate(&ValueDB{})
	if err != nil {
		log.Fatalf("ValueDB Table AutoMigrate Err: %v", err)
		return err
	}

	return err
}

func dbInsert(r Res, db *gorm.DB, famliy string) error {
	var err error
	err = db.Create(&ValueDB{
		CVE:     r.CveID,
		Package: r.Name,
		Value:   r.Version,
		Family:  famliy,
	}).Error
	return err
}

var iDB *gorm.DB
var ubuntuQuery = "select cve_id,name,version from packages,debians,advisories WHERE packages.definition_id = debians.definition_id and debians.definition_id = advisories.definition_id"
var debainQuery = "select cve_id,name,version from packages,debians WHERE packages.definition_id = debians.definition_id"
var redhatQuery = "select cve_id,name,version from packages,advisories,cves WHERE packages.definition_id = advisories.definition_id and advisories.id = cves.advisory_id"

func main() {
	g := loadConfig()
	fmt.Println(getDSN(g))
	err := dbInit(g)
	if err != nil {
		log.Println("Err:", err)
	}
	var res []Res
	//db, err := gorm.Open(sqlite.Open("ubuntu.db"), &gorm.Config{})
	//if err != nil {
	//	log.Println("sqlite Err:", err)
	//}
	//db.Raw(ubuntuQuery).Scan(&res)
	//for _, r := range res {
	//	if len(r.Version) >0 {
	//		err := dbInsert(r,iDB,"Debian")
	//		if err != nil{
	//			log.Println("Insert Err:",err)
	//		}
	//	}
	//}
	//db, err := gorm.Open(sqlite.Open("debian.db"), &gorm.Config{})
	//if err != nil {
	//	log.Println("sqlite Err:", err)
	//}
	//db.Raw(debainQuery).Scan(&res)
	//for _, r := range res {
	//	if len(r.Version) >0 {
	//		err := dbInsert(r,iDB,"Debian")
	//		if err != nil{
	//			log.Println("Insert Err:",err)
	//		}
	//	}
	//}
	db, err := gorm.Open(sqlite.Open("redhat.db"), &gorm.Config{})
	if err != nil {
		log.Println("sqlite Err:", err)
	}
	db.Raw(redhatQuery).Scan(&res)
	for _, r := range res {
		if len(r.Version) > 0 {
			err := dbInsert(r, iDB, "Redhat")
			if err != nil {
				log.Println("Insert Err:", err)
			}
		}
	}
}
