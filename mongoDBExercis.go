package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	client     *mongo.Client
	err        error
	database   *mongo.Database
	collection *mongo.Collection
	cond       *FindByJobName
	cursor     *mongo.Cursor
	record     *LogRecord
	delCond    *DeleteCond
	delResult  *mongo.DeleteResult
)

func init() {
	// 建立连接
	if client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://10.70.93.216:27017")); err != nil {
		fmt.Println(err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		fmt.Println(err)
		return
	}
	// 选择数据库my_db
	database = client.Database("cron")
	// 选择表my_collection
	collection = database.Collection("log")
}

// 任务执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}
type LogRecord struct {
	JobName   string    `bson:"jobName"`   // 任务名称
	Command   string    `bson:"command"`   // shell 命令
	Err       string    `bson:"err"`       // 脚本错误
	Content   string    `bson:"content"`   // 脚本输出
	TimePoint TimePoint `bson:"timePoint"` // 执行时间点
}

// 定义查询条件
type FindByJobName struct {
	JobName string `bson:"jobName"` // 任务名称
}

// startTime小于某时间
// {$lt:xxxx}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// 定义删除条件:{timePoint.startTime:{$lt:xxxx}}
// 只需记住: struct中的bson对应json的键,值对应结构体字段的值
type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	// 插入单条数据
	//insertOneItem()
	// 插入多条数据
	//insertManyItems()
	// 查看
	//findItems()
	// 删除
	delItems()
}
func delItems() {
	// 构造删除条件
	delCond = &DeleteCond{BeforeCond: TimeBeforeCond{Before: time.Now().Unix()}}
	// 执行删除操作
	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("删除的数量:", delResult.DeletedCount)
}
func findItems() {
	// 构造查询条件: 按照jobName字段过滤,
	cond := &FindByJobName{JobName: "job10"}
	// 查询
	if cursor, err = collection.Find(context.TODO(), cond); err != nil {
		fmt.Println(err)
		return
	}
	// 关闭cursor
	defer cursor.Close(context.TODO())
	// 遍历结果集
	for cursor.Next(context.TODO()) {
		record = &LogRecord{}
		// 反序列化到bson到对象
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("结果:", *record)
	}
}
func insertManyItems() {
	// 插入记录(bson)
	var (
		record2, record3 LogRecord
		insertManyResult *mongo.InsertManyResult
		docIds           []interface{}
		docId            primitive.ObjectID
	)
	record2 = LogRecord{
		JobName: "job10",
		Command: "echo record2",
		Err:     "",
		Content: "Hello record2",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}
	record3 = LogRecord{
		JobName: "job10",
		Command: "echo record3",
		Err:     "",
		Content: "Hello record3",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}
	records := []interface{}{record2, record3}
	if insertManyResult, err = collection.InsertMany(context.TODO(), records); err != nil {
		fmt.Println(err)
		return
	}
	docIds = insertManyResult.InsertedIDs
	for _, _docId := range docIds {
		docId = _docId.(primitive.ObjectID)
		fmt.Println("批量插入的IDs:", docId)
	}
}
func insertOneItem() {
	// 插入记录(bson)
	var (
		result *mongo.InsertOneResult
		docId  primitive.ObjectID
	)
	record = &LogRecord{
		JobName: "job10",
		Command: "echo Hello",
		Err:     "",
		Content: "Hello",
		TimePoint: TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}
	// insert one
	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println(err)
		return
	}

	// _id: 默认生成一个全局位置ID,objectId:12字节的二进制
	docId = result.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID:", docId)
}
