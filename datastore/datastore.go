package datastore

import (
	"github.com/priykumar/notification-service/model"
)

type DataStore interface {
	GetTemplate(string) *model.Template
	PutTemplate(model.Template) error
	PutNotification(model.Notification)
	GetNotification(string) model.Notification
}

type DB struct {
	type_Template map[string]model.Template
	id_Notify     map[string]model.Notification
}

func InitialiseDB() *DB {
	return &DB{
		type_Template: make(map[string]model.Template),
		id_Notify:     make(map[string]model.Notification),
	}
}

// Get existing template
func (d *DB) GetTemplate(tName string) *model.Template {
	if resp, exist := d.type_Template[tName]; exist {
		return &resp
	}
	return nil
}

// Insert template in DB
func (d *DB) PutTemplate(t model.Template) error {
	d.type_Template[t.Name] = t
	return nil
}

// Insert notification in DB
func (d *DB) PutNotification(n model.Notification) {
	d.id_Notify[n.ID] = n
}

// Get existing template for given notification ID
func (d *DB) GetNotification(id string) model.Notification {
	return d.id_Notify[id]
}
