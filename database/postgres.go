package database

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

var (
	Dbcon    *gorm.DB
	Errdb    error
	dbuser   string
	dbpass   string
	dbname   string
	dbaddres string
	dbport   string
	dbdebug  bool
	dbtype   string
	sslmode  string
)

func init() {

	dbtype = beego.AppConfig.DefaultString("db.type", "POSTGRES")
	dbuser = beego.AppConfig.DefaultString("db.postgres.user", "ottoagcfg")
	dbpass = beego.AppConfig.DefaultString("db.postgres.pass", "dTj*&56$es")
	dbname = beego.AppConfig.DefaultString("db.postgres.name", "ottoagdb")
	dbaddres = beego.AppConfig.DefaultString("db.postgres.addres", "")
	dbport = beego.AppConfig.DefaultString("db.postgres.port", "8086")
	sslmode = beego.AppConfig.DefaultString("db.postgres.sslmode", "disable")
	dbdebug = beego.AppConfig.DefaultBool("db.postgres.debug", false)
	if DbOpen() != nil {
		fmt.Println("Can Open db Postgres")
	}
}

// DbOpen ...
func DbOpen() error {
	args := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbaddres, dbport, dbuser, dbpass, dbname, sslmode)
	Dbcon, Errdb = gorm.Open("postgres", args)
	fmt.Println("isi postgres sett ", args)
	if Errdb != nil {
		logs.Error("open db Err ", Errdb)
		return Errdb
	}
	if errping := Dbcon.DB().Ping(); errping != nil {
		return errping
	}
	fmt.Println("Database connected [", dbaddres, "] [", dbname, "] [", dbuser, "] !")
	return nil
}

// GetDbCon ...
func GetDbCon() *gorm.DB {

	if errping := Dbcon.DB().Ping(); errping != nil {
		logs.Error("Db Not Connect test Ping :", errping)
		errping = nil
		if errping = DbOpen(); errping != nil {
			logs.Error("try to connect again but error :", errping)
		}
	}
	Dbcon.LogMode(dbdebug)
	return Dbcon
}
