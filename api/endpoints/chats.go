package endpoints

import (
	"backend/internal/model"
	"backend/internal/mwares"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"strconv"
)

type Chat struct {
	repo        model.UserChatRepository
	messageRepo model.MessageRepository
}

func NewUserChat(repo model.UserChatRepository, messageRepo model.MessageRepository) *Chat {
	return &Chat{
		messageRepo: messageRepo,
		repo:        repo,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ch *Chat) NewChat(c echo.Context) error {
	id := c.Get("userID")
	userid, _ := id.(uint64)

	reciverid, err := strconv.ParseUint(c.FormValue("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
		UserID:     &userid,
		ReceiverID: &reciverid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(chats) != 0 {
		return c.JSON(http.StatusCreated, map[string]string{
			"msg": "this chat already exists",
		})
	}

	if err := ch.repo.Create(c.Request().Context(), model.Chat{
		UserID:     userid,
		ReceiverID: reciverid,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"msg": "chat created",
	})
}

func (ch *Chat) GetChat(c echo.Context) error {
	chatid, err := strconv.ParseUint(c.Param("chatid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
		ID: &chatid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(chats) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, chats)
}

func (ch *Chat) GetChats(c echo.Context) error {
	id := c.Get("userID")
	userid, _ := id.(uint64)

	chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
		UserID: &userid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	reciverChats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
		ReceiverID: &userid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	allChats := append(chats, reciverChats...)

	type response struct {
		Chat          model.Chat `json:"chat"`
		UnreadMessage int        `json:"unreadMessage"`
	}

	if len(allChats) == 0 {
		return echo.ErrNotFound
	}

	res := make([]response, len(allChats))
	cond := "false"
	for i, v := range allChats {
		messages, _ := ch.messageRepo.Get(c.Request().Context(), model.MessageInterface{
			IsRead: &cond,
			ChatID: &v.ChatID,
		})
		res[i] = response{
			Chat:          v,
			UnreadMessage: len(messages),
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (ch *Chat) GetChatsSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	id := c.Get("userID")
	userid, _ := id.(uint64)

	sendUpdatedChats := func() error {
		chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
			UserID: &userid,
		})
		if err != nil {
			return echo.ErrInternalServerError
		}

		receiverChats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
			ReceiverID: &userid,
		})
		if err != nil {
			return echo.ErrInternalServerError
		}

		allChats := append(chats, receiverChats...)

		type response struct {
			Chat          model.Chat `json:"chat"`
			UnreadMessage int        `json:"unreadMessage"`
		}

		if len(allChats) == 0 {
			return nil
		}

		res := make([]response, len(allChats))
		cond := "false"
		for i, v := range allChats {
			messages, _ := ch.messageRepo.Get(c.Request().Context(), model.MessageInterface{
				IsRead: &cond,
				ChatID: &v.ChatID,
			})
			res[i] = response{
				Chat:          v,
				UnreadMessage: len(messages),
			}
		}

		return ws.WriteJSON(res)
	}

	for {
		err = sendUpdatedChats()
		if err != nil {
			return err
		}

		type msg struct {
			Message string `json:"message"`
		}
		var m msg
		err = ws.ReadJSON(&m)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return err
			}
			break
		}

		if m.Message == "new" {
			err = sendUpdatedChats()
			if err != nil {
				return err
			}
		}
		if m.Message == "exit" {
			break
		}
	}

	return nil
}

func (ch *Chat) DeleteChat(c echo.Context) error {
	id := c.Get("userID")
	userid, _ := id.(uint64)

	chatid, err := strconv.ParseUint(c.Param("chatid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
		ID: &chatid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(chats) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "this chat does not exist",
		})
	}
	if chats[0].UserID != userid && chats[0].ReceiverID != userid {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "can not access this chat",
		})
	}

	if err = ch.repo.Delete(c.Request().Context(), model.ChatInterface{
		ID: &chatid,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "chat deleted",
	})
}

func (ch *Chat) DeleteChatMessage(c echo.Context) error {
	id := c.Get("userID")
	userid, _ := id.(uint64)

	chatid, err := strconv.ParseUint(c.Param("chatid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	messageid, err := strconv.ParseUint(c.Param("messageid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
		ID: &chatid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(chats) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "this chat does not exist",
		})
	}
	if chats[0].UserID != userid && chats[0].ReceiverID != userid {
		return c.JSON(http.StatusNotFound, map[string]string{
			"map": "can not access this chat",
		})
	}

	if err := ch.messageRepo.Delete(c.Request().Context(), model.MessageInterface{
		ID: &messageid,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "message deleted",
	})
}

func (ch *Chat) NewChatMessage(c echo.Context) error {
	id := c.Get("userID")
	senderid, _ := id.(uint64)

	chatID, err := strconv.ParseUint(c.Param("chatid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
		ID: &chatID,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(chats) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "this chat does not exist",
		})
	}
	if chats[0].UserID != senderid && chats[0].ReceiverID != senderid {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "can not access this chat",
		})
	}

	messageContent := c.FormValue("content")
	if messageContent == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "message content can not be empty",
		})
	}

	if _, err := ch.messageRepo.Create(c.Request().Context(), model.Message{
		ChatID:   chatID,
		SenderID: senderid,
		Type:     model.TypePV,
		IsRead:   "false",
		Content:  messageContent,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"msg": "message sent",
	})
}

func (ch *Chat) NewChatMessageWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	id := c.Get("userID")
	senderid, _ := id.(uint64)

	chatid, err := strconv.ParseUint(c.Param("chatid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	for {
		var incomingMessage struct {
			Content string `json:"content"`
			Stat    string `json:"stat"`
		}

		err = ws.ReadJSON(&incomingMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return err
			}
			break
		}

		if incomingMessage.Stat == "exit" {
			break
		}

		chatID := chatid
		messageContent := incomingMessage.Content

		chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
			ID: &chatID,
		})
		if err != nil {
			return echo.ErrInternalServerError
		}

		if len(chats) == 0 {
			ws.WriteMessage(websocket.TextMessage, []byte("This chat does not exist"))
			continue
		}
		if chats[0].UserID != senderid && chats[0].ReceiverID != senderid {
			ws.WriteMessage(websocket.TextMessage, []byte("Cannot access this chat"))
			continue
		}

		if messageContent == "" {
			ws.WriteMessage(websocket.TextMessage, []byte("Message content cannot be empty"))
			continue
		}

		if _, err := ch.messageRepo.Create(c.Request().Context(), model.Message{
			ChatID:   chatID,
			SenderID: senderid,
			Type:     model.TypePV,
			IsRead:   "false",
			Content:  messageContent,
		}); err != nil {
			return echo.ErrInternalServerError
		}

		ws.WriteMessage(websocket.TextMessage, []byte("Message sent"))
	}

	return nil
}

func (ch *Chat) GetMessageByCount(c echo.Context) error {
	id := c.Get("userID")
	userid, _ := id.(uint64)

	chatid, err := strconv.ParseUint(c.Param("chatid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
		ID: &chatid,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(chats) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "this chat does not exist",
		})
	}
	if chats[0].UserID != userid && chats[0].ReceiverID != userid {
		return c.JSON(http.StatusNotFound, map[string]string{
			"msg": "can not access this chat",
		})
	}

	count, err := strconv.ParseUint(c.Param("count"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	chatType := model.TypePV

	messages, err := ch.messageRepo.GetDto(c.Request().Context(), model.MessageInterface{
		ChatID: &chatid,
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

func (ch *Chat) GetMessageByCountWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	id := c.Get("userID")
	userid, _ := id.(uint64)

	chatid, err := strconv.ParseUint(c.Param("chatid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	for {
		var requestData struct {
			Count uint64 `json:"count"`
			Stat  string `json:"stat"`
		}

		err = ws.ReadJSON(&requestData)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return err
			}
			break
		}

		if requestData.Stat == "exit" {
			break
		}

		chatid := chatid
		count := requestData.Count

		chats, err := ch.repo.Get(c.Request().Context(), model.ChatInterface{
			ID: &chatid,
		})
		if err != nil {
			return echo.ErrInternalServerError
		}

		if len(chats) == 0 {
			ws.WriteMessage(websocket.TextMessage, []byte("This chat does not exist"))
			continue
		}
		if chats[0].UserID != userid && chats[0].ReceiverID != userid {
			ws.WriteMessage(websocket.TextMessage, []byte("Cannot access this chat"))
			continue
		}

		chatType := model.TypePV

		messages, err := ch.messageRepo.GetDto(c.Request().Context(), model.MessageInterface{
			ChatID: &chatid,
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

		err = ws.WriteJSON(messages[:count])
		if err != nil {
			return err
		}
	}

	return nil
}

func (ch *Chat) NewUserChatHandler(g *echo.Group) {
	chatGroup := g.Group("/chats")

	chatGroup.POST("", ch.NewChat, mwares.JWTMiddleware)
	chatGroup.GET("", ch.GetChats, mwares.JWTMiddleware)
	chatGroup.GET("/ws", ch.GetChatsSocket, mwares.JWTMiddleware)
	chatGroup.GET("/:chatid", ch.GetChat, mwares.JWTMiddleware)
	chatGroup.DELETE("/:chatid", ch.DeleteChat, mwares.JWTMiddleware)
	chatGroup.POST("/:chatid/message", ch.NewChatMessage, mwares.JWTMiddleware)
	chatGroup.GET("/:chatid/message/ws", ch.NewChatMessageWs, mwares.JWTMiddleware)
	chatGroup.DELETE("/:chatid/message/:messageid", ch.DeleteChatMessage, mwares.JWTMiddleware)
	chatGroup.GET("/:chatid/message/:count", ch.GetMessageByCount, mwares.JWTMiddleware)
	chatGroup.GET("/:chatid/message/get/ws", ch.GetMessageByCountWs, mwares.JWTMiddleware)
}
