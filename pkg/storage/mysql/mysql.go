package mysql

import (
	"fmt"

	"github.com/you06/chaos-commander/config"
	"github.com/you06/chaos-commander/pkg/storage/basic"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/juju/errors"
)

type driver struct {
	db *gorm.DB
}

// Connect create database connection
func Connect(config *config.Database) (basic.Driver, error) {
	fmt.Println(config)
	connect := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open("mysql", connect)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &driver{
		db,
	}, nil
}

// Close connection
func (d *driver) Close() error {
	return errors.Trace(d.db.Close())
}

// Save insert or update data
func (d *driver) Save(model interface{}) error {
	return errors.Trace(d.db.Save(model).Error)
}

// FindOne get one matched result
func (d *driver) FindOne(model interface{}, state string, cond ...interface{}) error {
	return errors.Trace(d.db.Where(state, cond...).First(model).Error)
}

// Find get all matched results
func (d *driver) Find(model interface{}, state string, cond ...interface{}) error {
	return errors.Trace(d.db.Where(state, cond...).Find(model).Error)
}

// Update data
func (d *driver) Update(model interface{}, state string, cond []interface{}, update map[string]interface{}) error {
	return errors.Trace(d.db.Model(model).Where(state, cond...).Updates(update).Error)
}

// UpdateAll update all data
func (d *driver) UpdateAll(model interface{}, update map[string]interface{}) error {
	return errors.Trace(d.db.Model(model).Updates(update).Error)
}
