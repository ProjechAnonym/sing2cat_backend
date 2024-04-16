package ginconfig

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/tidwall/buntdb"

	"gorm.io/gorm"
)
var Db *gorm.DB
var BuntDb *buntdb.DB
func Get_database(){
	project_dir,_ := Get_value("project_dir")
	Db,_ = gorm.Open(sqlite.Open(fmt.Sprintf("%s/sing2cat.db",project_dir)),&gorm.Config{})
	BuntDb, _= buntdb.Open(":memory:")
	Db.AutoMigrate(&Component{})
}