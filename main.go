package main

import (
	"flag"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type PreProjectTask struct {
	ID              uint64         `gorm:"column:id;primaryKey"`
	ParentID        uint64         `gorm:"column:parent_id"`
	ProjectID       uint64         `gorm:"column:project_id"`
	ColumnID        uint64         `gorm:"column:column_id"`
	DialogID        uint64         `gorm:"column:dialog_id"`
	FlowItemID      uint64         `gorm:"column:flow_item_id"`
	FlowItemName    string         `gorm:"column:flow_item_name"`
	Name            string         `gorm:"column:name"`
	Color           string         `gorm:"column:color"`
	Desc            string         `gorm:"column:desc"`
	StartAt         *time.Time     `gorm:"column:start_at"`
	EndAt           *time.Time     `gorm:"column:end_at"`
	ArchivedAt      *time.Time     `gorm:"column:archived_at"`
	ArchivedUserID  uint64         `gorm:"column:archived_userid"`
	ArchivedFollow  int            `gorm:"column:archived_follow"`
	CompleteAt      *time.Time     `gorm:"column:complete_at"`
	UserID          uint64         `gorm:"column:userid"`
	IsAllVisible    int            `gorm:"column:is_all_visible"`
	PLevel          int            `gorm:"column:p_level"`
	PName           string         `gorm:"column:p_name"`
	PColor          string         `gorm:"column:p_color"`
	Sort            int            `gorm:"column:sort"`
	Loop            string         `gorm:"column:loop"`
	LoopAt          *time.Time     `gorm:"column:loop_at"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	DeletedAt       *gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedUserID   uint64         `gorm:"column:deleted_userid"`
}

type PreProjectTaskUser struct {
	ID         uint64 `gorm:"column:id;primaryKey"`
	ProjectID  uint64 `gorm:"column:project_id"`
	TaskID     uint64 `gorm:"column:task_id"`
	TaskPID    uint64 `gorm:"column:task_pid"`
	UserID     uint64 `gorm:"column:userid"`
	Owner      int8   `gorm:"column:owner"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}


type PreProject struct {
	ID             uint64         `gorm:"column:id;primaryKey"`
	Name           string         `gorm:"column:name"`
	Description    string         `gorm:"column:desc"`
	UserID         uint64         `gorm:"column:userid"`
	Personal       bool           `gorm:"column:personal"`
	UserSimple     string         `gorm:"column:user_simple"`
	DialogID       uint64         `gorm:"column:dialog_id"`
	ArchivedAt     *time.Time     `gorm:"column:archived_at"`
	ArchivedUserID uint64         `gorm:"column:archived_userid"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

type PreProjectUser struct {
	ID        uint64     `gorm:"column:id;primary_key;auto_increment"`
	ProjectID uint64     `gorm:"column:project_id"`
	UserID    uint64     `gorm:"column:userid"`
	Owner     int        `gorm:"column:owner"`
	TopAt     *time.Time `gorm:"column:top_at"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type PreProjectColumn struct {
	ID        uint64 `gorm:"column:id;primaryKey"`
	ProjectID uint64 `gorm:"column:project_id"`
	Name      string `gorm:"column:name"`
	Color     string `gorm:"column:color"`
	Sort      int `gorm:"column:sort"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func RandomString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func newProject(db *gorm.DB,now time.Time, col int)  {
	p := PreProject{
		Name:RandomString(5),
		Description:RandomString(20),
		UserID: 1,
		DialogID: 1,
		UpdatedAt:now,
		CreatedAt:now,
	}
	err := db.Model(&PreProject{}).Create(&p).Error
	if err != nil {
		panic(err)
	}

	var temp PreProject
	db.Model(&PreProject{}).Last(&temp)
	pUser := PreProjectUser{
		ProjectID:temp.ID,
		UserID:temp.UserID,
		Owner:1,
		CreatedAt:&now,
		UpdatedAt:&now,
	}
	proCol := make([]PreProjectColumn,0)
	for j := 0;j < col; j++{
		proCol = append(proCol,PreProjectColumn{
			ProjectID:temp.ID,
			Name:RandomString(6),
			CreatedAt:now,
			UpdatedAt:now,
		})
	}
	err = db.Model(&PreProjectUser{}).Create(&pUser).Error
	err = db.Model(&PreProjectColumn{}).Create(&proCol).Error
	if err != nil {
		panic(err)
	}
}

func main(){

	mysqlPtr := flag.String("mysql", "localhost", "数据库地址")
	portPtr := flag.Int("port", 3306, "数据库端口")
	userPtr := flag.String("user","root","数据库用户")
	passwdPtr := flag.String("password","123456","数据库密码")
	sqlNamePtr := flag.String("sqlname","dootask","数据库名称")
	proNums := flag.Int("p",1,"指定生成项目数量")
	colNums := flag.Int("c",1,"指定每个项目生成列表数量")
	tasksNums := flag.Int("t",1,"指定每个项目生成任务数量")
	flag.Parse()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",*userPtr,*passwdPtr,*mysqlPtr,*portPtr,*sqlNamePtr)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}


	userTasks := make([]PreProjectTaskUser,0)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0;i < *proNums; i++ {
		var projectID uint64
		var columnID []uint64
		var name string
		now := time.Now()
		end := time.Now().Add(time.Hour * 24)
		newProject(db,now,*colNums)
		err = db.Model(&PreProject{}).Where("userid = ?",1).Pluck("id", &projectID).Error
		err = db.Model(&PreProjectColumn{}).Where("project_id = ?",projectID).Pluck("id", &columnID).Error
		err = db.Model(&PreProject{}).Where("userid = ?",1).Pluck("name", &name).Error
		if err != nil {
			panic(err)
		}
		for j := 0; j < *tasksNums; j++ {
			tasks := PreProjectTask{
				ProjectID: projectID,
				ColumnID:columnID[random.Intn(len(columnID))],
				FlowItemID:1,
				FlowItemName:"start|待处理",
				Name: RandomString(6),
				UserID:1,
				IsAllVisible:1,
				Sort:1,
				StartAt:&now, // 设置 StartAt 的值为当前时间
				EndAt:&end, // 设置 EndAt 的值为当前时间往后一天
				UpdatedAt:now,
				CreatedAt:now,
			}

			err = db.Model(&PreProjectTask{}).Create(&tasks).Error
			if err != nil {
				panic(err)
			}

			var temp PreProjectTask
			db.Model(&PreProjectTask{}).Last(&temp)
			userTasks = append(userTasks,PreProjectTaskUser{
				ProjectID:projectID,
				TaskID:temp.ID,
				TaskPID:temp.ID,
				UserID:1,
				Owner:1,
				CreatedAt:now,
				UpdatedAt:now,
			})
		}
		fmt.Println("将于",name,"项目生成",*colNums,"个列表",*tasksNums,"条任务...")
	}

	err = db.Model(&PreProjectTaskUser{}).Create(&userTasks).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("ok!")
}
