package api

import (
	"backend/api/endpoints"
	"backend/internal/configs"
	"backend/internal/repositoryImpl/contactRepoImpl"
	"backend/internal/repositoryImpl/groupChatRepoImpl"
	"backend/internal/repositoryImpl/groupRepoImpl"
	"backend/internal/repositoryImpl/messageRepoImpl"
	"backend/internal/repositoryImpl/userChatRepoImpl"
	"backend/internal/repositoryImpl/userGroupRepoImpl"
	"backend/internal/repositoryImpl/userRepoImpl"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type repos struct {
	userRepo      *userRepoImpl.Repository
	messageRepo   *messageRepoImpl.Repository
	groupRepo     *groupRepoImpl.Repository
	contactRepo   *contactRepoImpl.Repository
	userGroupRepo *userGroupRepoImpl.Repository
	groupChatRepo *groupChatRepoImpl.Repository
	userChatRepo  *userChatRepoImpl.Repository
}

func initRepos(db *gorm.DB) *repos {
	return &repos{
		userRepo:      userRepoImpl.New(db),
		messageRepo:   messageRepoImpl.New(db),
		groupRepo:     groupRepoImpl.New(db),
		contactRepo:   contactRepoImpl.New(db),
		userGroupRepo: userGroupRepoImpl.New(db),
		groupChatRepo: groupChatRepoImpl.New(db),
		userChatRepo:  userChatRepoImpl.New(db),
	}
}

func Run(db *gorm.DB, conf *configs.Config) {
	repos := initRepos(db)

	e := echo.New()

	hu := endpoints.NewUser(repos.userRepo, repos.contactRepo)
	hc := endpoints.NewUserChat(repos.userChatRepo, repos.messageRepo)
	hg := endpoints.NewGroup(repos.groupRepo, repos.messageRepo, repos.userGroupRepo, repos.groupChatRepo)

	apiGroup := e.Group("/api")

	hu.NewUserHandler(apiGroup)
	hc.NewUserChatHandler(apiGroup)
	hg.NewGroupHandler(apiGroup)

	if err := e.Start(conf.Server.Address + ":" + conf.Server.Port); err != nil {
		logrus.Fatalf("server failed to start %v", err)
	}
}
