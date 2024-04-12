package endpoints

import (
	"backend/internal/model"
	"backend/internal/mwares"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strconv"
)

type Group struct {
	repo          model.GroupRepository
	messageRepo   model.MessageRepository
	userGroupRepo model.UserGroupRepository
	groupChatRepo model.GroupChatRepository
}

func NewGroup(repo model.GroupRepository, messageRepo model.MessageRepository, userGroupRepo model.UserGroupRepository,
	groupChatRepo model.GroupChatRepository) *Group {
	return &Group{
		groupChatRepo: groupChatRepo,
		messageRepo:   messageRepo,
		userGroupRepo: userGroupRepo,
		repo:          repo,
	}
}

func (g *Group) NewGroup(c echo.Context) error {
	id := c.Get("userID")
	creatorid, _ := id.(uint64)

	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "group name can not be empty",
		})
	}

	logrus.Warnln("uc", creatorid)

	groupID, err := g.repo.Create(c.Request().Context(), model.Group{
		Creator:     creatorid,
		Name:        name,
		Description: c.FormValue("description"),
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if err = g.userGroupRepo.Create(c.Request().Context(), model.UserGroup{
		GroupID: groupID,
		UserID:  creatorid,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"msg": "group created",
	})
}

func (g *Group) GetGroupData(c echo.Context) error {
	groupid, err := strconv.ParseUint(c.Param("groupid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	group, users, err := g.userGroupRepo.GetGroupWithUserGroups(c.Request().Context(), groupid)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(group) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "no groups found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"group": group,
		"users": users,
	})
}

func (g *Group) GetGroups(c echo.Context) error {
	id := c.Get("userID")
	userID, _ := id.(uint64)

	groups, err := g.userGroupRepo.Get(c.Request().Context(), model.UserGroupInterface{
		UserID: &userID,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(groups) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, groups)
}

func (g *Group) DeleteGroup(c echo.Context) error {
	id := c.Get("userID")
	creatorid, _ := id.(uint64)

	groupID, err := strconv.ParseUint(c.Param("groupid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	groups, err := g.repo.Get(c.Request().Context(), model.GroupInterface{
		ID:        &groupID,
		CreatorID: &creatorid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(groups) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "group not found",
		})
	}

	if err = g.repo.Delete(c.Request().Context(), model.GroupInterface{
		ID:        &groupID,
		CreatorID: &creatorid,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "group deleted",
	})
}

func (g *Group) AddUserToGroup(c echo.Context) error {
	cid := c.Get("userID")
	creatorid, _ := cid.(uint64)

	groupID, err := strconv.ParseUint(c.Param("groupid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	groups, err := g.repo.Get(c.Request().Context(), model.GroupInterface{
		ID:        &groupID,
		CreatorID: &creatorid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(groups) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "group not found",
		})
	}

	id, err := strconv.ParseUint(c.FormValue("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err = g.userGroupRepo.Create(c.Request().Context(), model.UserGroup{
		GroupID: groupID,
		UserID:  id,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"msg": "user added",
	})
}

func (g *Group) DeleteUserFromGroup(c echo.Context) error {
	cid := c.Get("userID")
	creatorid, _ := cid.(uint64)

	groupID, err := strconv.ParseUint(c.Param("groupid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	groups, err := g.repo.Get(c.Request().Context(), model.GroupInterface{
		ID:        &groupID,
		CreatorID: &creatorid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(groups) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "group not found",
		})
	}

	id, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err = g.userGroupRepo.Delete(c.Request().Context(), model.UserGroupInterface{
		UserID:  &id,
		GroupID: &groupID,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "user deleted",
	})
}

func (g *Group) NewGroupMessage(c echo.Context) error {
	id := c.Get("userID")
	uid, _ := id.(uint64)

	userid, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	if userid != uid {
		return echo.ErrUnauthorized
	}

	groupID, err := strconv.ParseUint(c.Param("groupid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	users, err := g.userGroupRepo.Get(c.Request().Context(), model.UserGroupInterface{
		GroupID: &groupID,
		UserID:  &uid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(users) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "you are not part of this group",
		})
	}

	messageContent := c.FormValue("content")
	if messageContent == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "message body can not be empty",
		})
	}

	messageID, err := g.messageRepo.Create(c.Request().Context(), model.Message{
		Content:  messageContent,
		ChatID:   groupID,
		SenderID: uid,
		Type:     model.TypeGP,
		IsRead:   "false",
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if err = g.groupChatRepo.Create(c.Request().Context(), model.GroupChat{
		GroupID:   groupID,
		MessageID: messageID,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"msg": "message sent",
	})
}

func (g *Group) DeleteGroupMessage(c echo.Context) error {
	id := c.Get("userID")
	uid, _ := id.(uint64)

	groupid, err := strconv.ParseUint(c.Param("groupid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	users, err := g.userGroupRepo.Get(c.Request().Context(), model.UserGroupInterface{
		GroupID: &groupid,
		UserID:  &uid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(users) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "you are not part of this group",
		})
	}

	messageid, err := strconv.ParseUint(c.Param("messageid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	if err = g.messageRepo.Delete(c.Request().Context(), model.MessageInterface{
		ID: &messageid,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	if err = g.groupChatRepo.Delete(c.Request().Context(), model.GroupChatInterface{
		ID:      &messageid,
		GroupID: &groupid,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"msg": "message deleted",
	})
}

func (g *Group) GetGroupMessages(c echo.Context) error {
	id := c.Get("userID")
	uid, _ := id.(uint64)

	groupID, err := strconv.ParseUint(c.Param("groupid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	users, err := g.userGroupRepo.Get(c.Request().Context(), model.UserGroupInterface{
		GroupID: &groupID,
		UserID:  &uid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(users) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "you are not part of this group",
		})
	}

	count, err := strconv.ParseUint(c.Param("count"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	chatType := model.TypeGP

	messages, err := g.messageRepo.GetDto(c.Request().Context(), model.MessageInterface{
		ChatID: &groupID,
		Type:   &chatType,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})

	if count > uint64(len(messages)) {
		count = uint64(len(messages))
	}

	return c.JSON(http.StatusOK, messages[:count])
}

func (g *Group) NewGroupHandler(gr *echo.Group) {
	GroupsGroup := gr.Group("/groups")

	GroupsGroup.POST("", g.NewGroup, mwares.JWTMiddleware)
	GroupsGroup.GET("/allgroups", g.GetGroups, mwares.JWTMiddleware)
	GroupsGroup.GET("/:groupid", g.GetGroupData, mwares.JWTMiddleware)
	GroupsGroup.DELETE("/:groupid", g.DeleteGroup, mwares.JWTMiddleware)
	GroupsGroup.POST("/:groupid", g.AddUserToGroup, mwares.JWTMiddleware)
	GroupsGroup.DELETE("/:groupid/:userid", g.DeleteUserFromGroup, mwares.JWTMiddleware)
	GroupsGroup.POST("/:groupid/message/:userid", g.NewGroupMessage, mwares.JWTMiddleware)
	GroupsGroup.DELETE("/:groupid/message/:messageid", g.DeleteGroupMessage, mwares.JWTMiddleware)
	GroupsGroup.GET("/:groupid/message/:count", g.GetGroupMessages, mwares.JWTMiddleware)
}
