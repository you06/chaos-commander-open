package basic

// Driver defines storage model
type Driver interface {
	Close() error
	Save(model interface{}) error
	FindOne(model interface{}, state string, cond ...interface{}) error
	Find(model interface{}, state string, cond ...interface{}) error
	Update(model interface{}, state string, cond []interface{}, update map[string]interface{}) error
	UpdateAll(model interface{}, update map[string]interface{}) error
}
