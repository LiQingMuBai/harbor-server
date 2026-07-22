package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	URI      string `json:"uri"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
	PoolSize int    `json:"poolsize"`
	Clinet   *mongo.Client
	DBHandle *mongo.Database
	Session  *mongo.Session
}

func (m *MongoDB) CreateClient() {
	uri := m.URI
	if uri == "" {
		uri = fmt.Sprintf("mongodb://%s:%d", m.Host, m.Port)
	}

	m_options := options.Client().ApplyURI(uri)

	m_options.SetConnectTimeout(100 * time.Second)
	m_options.SetHeartbeatInterval(3 * time.Second)
	m_options.SetSocketTimeout(3 * time.Second)
	//m_options.SetMaxConnIdleTime(10 * time.Second)
	//m_options.SetMaxPoolSize(200)
	//m_options.SetDirect(true)
	m_options.SetRetryReads(true)
	m_options.SetRetryWrites(true)

	if m.PoolSize > 0 {
		m_options.SetMaxPoolSize(uint64(m.PoolSize))
	}
	c, err := mongo.Connect(context.TODO(), m_options)
	if err != nil {
		panic(err)
	}
	m.Clinet = c

	m.DBHandle = m.Clinet.Database(m.DBName)
	//s, _ := m.Clinet.StartSession()
	//m.Session = &s

}

func (m *MongoDB) InsertData(table string, data interface{}) *mongo.InsertOneResult {
	opts := new(options.InsertOneOptions)
	opts.SetBypassDocumentValidation(true)
	rs, _ := m.DBHandle.Collection(table).InsertOne(context.TODO(), data, opts)

	return rs
}
func (m *MongoDB) InsertMany(table string, data []interface{}) *mongo.InsertManyResult {
	opts := new(options.InsertManyOptions)
	opts.SetBypassDocumentValidation(true)
	opts.SetOrdered(false)
	rs, _ := m.DBHandle.Collection(table).InsertMany(context.TODO(), data)

	return rs
}

func (m *MongoDB) GetOne(table string, condtion interface{}, sort interface{}) primitive.M {
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()
	opts := options.FindOneOptions{}

	if sort != nil {
		opts.Sort = sort
	}

	rs := m.DBHandle.Collection(table).FindOne(context.TODO(), condtion, &opts)
	if rs.Err() != nil {
		return nil
	}
	var crs bson.D
	err := rs.Decode(&crs)
	if err != nil {
		return nil
	}

	return crs.Map()
}
func (m *MongoDB) GetList(table string, condition interface{}, sort interface{}, limit int64) []primitive.M {

	opts := options.FindOptions{}
	//opts.SetAllowDiskUse(true)
	if sort != nil {
		opts.SetSort(sort)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}

	rs, err := m.DBHandle.Collection(table).Find(context.TODO(), condition, &opts)

	if err != nil {
		return nil
	}
	if rs.Err() != nil {
		return nil
	}
	defer rs.Close(context.TODO())

	list := make([]primitive.M, 0)

	//err = rs.All(context.TODO(), &list)
	//if err != nil {
	//return nil
	//}
	for rs.Next(context.TODO()) {
		var result bson.D
		err := rs.Decode(&result)
		if err != nil {
			return nil
		}
		list = append(list, result.Map())
	}

	return list

}
func (m *MongoDB) Update(table string, data interface{}, condition interface{}) error {
	_, err := m.DBHandle.Collection(table).UpdateOne(context.TODO(), condition, bson.M{"$set": data})

	return err
}
func (m *MongoDB) FindAndReplace(table string, data interface{}, condition interface{}) {

	one := m.GetOne(table, condition, nil)
	if one != nil {
		//m.DBHandle.Collection(table).FindOneAndReplace(context.TODO(), condition, data)
		m.Update(table, data, condition)
	} else {
		m.InsertData(table, data)
	}

}
