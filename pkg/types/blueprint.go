package types

// BluePrint defines blue print db structure
type BluePrint struct {
	ID        int    `gorm:"column:id"`
	BluePrint string `gorm:"column:blue_print"`
	Path      string `gorm:"column:path"`
}

func (b *BluePrint) GetID() int {
	return b.ID
}

func (b *BluePrint) GetBluePrint() string {
	return b.BluePrint
}

func (b *BluePrint) GetPath() string {
	return b.Path
}
