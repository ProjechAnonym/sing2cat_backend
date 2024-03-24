package ginconfig

import (
	"github.com/glebarez/sqlite"
	"github.com/tidwall/buntdb"

	"gorm.io/gorm"
)
var Db *gorm.DB
var BuntDb *buntdb.DB
func Get_database(){
	Db,_ = gorm.Open(sqlite.Open("sing2cat.db"),&gorm.Config{})
	BuntDb, _= buntdb.Open(":memory:")
	Db.AutoMigrate(&Component{})
}