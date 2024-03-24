package ginconfig

type Component struct {
	ID        uint64                 `gorm:"primaryKey" json:"id"`
	Url       string                 `json:"url" gorm:"not null;type:varchar(255)"`
	Name      string                 `json:"name" gorm:"not null;type:varchar(20)"`
	Class     string                 `json:"class" gorm:"not null;type:varchar(20)"`
	Icon      string                 `json:"icon" gorm:"not null;type:varchar(255)"`
	Data      map[string]interface{} `json:"data" gorm:"-"`
	Gorm_data string                 `json:"gorm_data" gorm:"type:json"`
	Icon_path string                 `json:"icon_path" gorm:"type:varchar(255)"`
}