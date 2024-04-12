package main

import (
	"backend/api"
	"backend/internal/configs"
	"backend/internal/model"
	"backend/internal/repositoryImpl/contactRepoImpl"
	"backend/internal/repositoryImpl/groupChatRepoImpl"
	"backend/internal/repositoryImpl/messageRepoImpl"
	"backend/internal/repositoryImpl/userChatRepoImpl"
	"backend/internal/repositoryImpl/userGroupRepoImpl"
	"backend/internal/repositoryImpl/userRepoImpl"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

func main() {
	conf, err := configs.LoadConfig()
	if err != nil {
		logrus.Warnln("error in loading configs")
	}

	port := strconv.Itoa(conf.Database.Port)
	dsn := "host=" + conf.Database.Addr + " user=" + conf.Database.User + " password=" + conf.Database.Password + " dbname=" + conf.Database.DBName + " port=" + port
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("failed to connect database %v", err)
	}

	err = db.AutoMigrate(new(userRepoImpl.UserDTO), new(userChatRepoImpl.UserChatDTO), new(messageRepoImpl.MessageDTO), new(model.GroupDTO),
		new(contactRepoImpl.ContactDTO), new(userGroupRepoImpl.UserGroupDTO), new(groupChatRepoImpl.GroupChatDTO))
	if err != nil {
		return
	}

	api.Run(db, conf)
}
