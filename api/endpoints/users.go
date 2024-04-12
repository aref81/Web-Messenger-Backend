package endpoints

import (
	"backend/internal/authorize"
	"backend/internal/configs"
	"backend/internal/model"
	"backend/internal/mwares"
	"backend/utils"
	"backend/utils/datasource"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	repo        model.UserRepository
	contactRepo model.ContactRepository
}

func NewUser(repo model.UserRepository, contactRepo model.ContactRepository) *User {
	return &User{
		contactRepo: contactRepo,
		repo:        repo,
	}
}

func (u *User) RegisterUser(c echo.Context) error {
	password := c.FormValue("password")

	if password == "" || c.FormValue("username") == "" || c.FormValue("phone") == "" || c.FormValue("name") == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "fields can not be empty",
		})
	}

	passHash := utils.HashData(password)

	req := model.User{
		Name:         c.FormValue("name"),
		Username:     c.FormValue("username"),
		Password:     passHash,
		Phone:        c.FormValue("phone"),
		IsFirtsLogin: "true",
		IsActive:     "true",
	}

	if err := u.repo.Create(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"err": err.Error(),
			"msg": "Cant create user",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"map": "User created",
	})
}

func (u *User) LoginUser(c echo.Context) error {
	password := c.FormValue("password")
	username := c.FormValue("username")

	passHash := utils.HashData(password)

	users, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		Username: &username,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(users) == 0 {
		return echo.ErrNotFound
	}
	if len(users) > 1 {
		return echo.ErrInternalServerError
	}
	if users[0].Password != passHash {
		return c.JSON(http.StatusNotFound, map[string]string{
			"err": "Unauthorized",
			"msg": "invalid password",
		})
	}

	if err = u.repo.Update(c.Request().Context(), model.User{
		Username: username,
		IsActive: "true",
	}); err != nil {
		return echo.ErrInternalServerError
	}

	token, err := authorize.JwtToken(users[0].UserID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":        token,
		"username":     users[0].Username,
		"userID":       users[0].UserID,
		"isFirstLogin": users[0].IsFirtsLogin,
	})
}

func (u *User) GetUserByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	users, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		ID:       &id,
		Username: nil,
		Phone:    nil,
		IsActive: nil,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(users) == 0 {
		return echo.ErrNotFound
	}
	if len(users) > 1 {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, users[0])
}

func (u *User) UpdateUser(c echo.Context) error {
	username := c.Param("username")

	ph := c.FormValue("phone")
	ps := c.FormValue("password")
	bi := c.FormValue("biography")
	nm := c.FormValue("name")
	pf, err := c.FormFile("profile")
	if err != nil {
		log.Warnln("no profile picture found")
		pf = nil
	}

	if nm != "" {
		if err = u.repo.Update(c.Request().Context(), model.User{
			Username: username,
			Name:     nm,
		}); err != nil {
			return echo.ErrInternalServerError
		}
	}

	if ph != "" {
		if err = u.repo.Update(c.Request().Context(), model.User{
			Username: username,
			Phone:    ph,
		}); err != nil {
			return echo.ErrInternalServerError
		}
	}

	if ps != "" {
		hps := utils.HashData(ps)
		if err = u.repo.Update(c.Request().Context(), model.User{
			Username: username,
			Password: hps,
		}); err != nil {
			return echo.ErrInternalServerError
		}
	}

	if bi != "" {
		if err = u.repo.Update(c.Request().Context(), model.User{
			Username:  username,
			Biography: bi,
		}); err != nil {
			return echo.ErrInternalServerError
		}
	}

	if pf != nil {
		conf, err := configs.LoadConfig()
		if err != nil {
			log.Errorln("cant open config file")
		}

		sess, err := datasource.ConnectS3(conf.S3.AccessKey, conf.S3.SecretKey, conf.S3.Region, conf.S3.Endpoint)
		if err != nil {
			log.Errorln("can not connect to S3")
		}

		PFPpath, err := datasource.UploadS3(sess, pf, conf.S3.Bucket, username)
		if err != nil {
			log.Errorln("can not upload to s3")
		}

		if err = u.repo.Update(c.Request().Context(), model.User{
			Username:       username,
			ProfilePicture: PFPpath,
		}); err != nil {
			return echo.ErrInternalServerError
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "user updated successfuly",
	})
}

func (u *User) GetUserPfpf(c echo.Context) error {
	username := c.Param("username")

	users, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		Username: &username,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	pfp := users[0].ProfilePicture

	conf, err := configs.LoadConfig()
	if err != nil {
		log.Errorln("cant open config file")
	}

	sess, err := datasource.ConnectS3(conf.S3.AccessKey, conf.S3.SecretKey, conf.S3.Region, conf.S3.Endpoint)
	if err != nil {
		log.Errorln("can not connect to S3")
		return echo.ErrInternalServerError
	}
	_, err = datasource.DownloadS3(sess, conf.S3.Bucket, pfp)
	if err != nil {
		log.Errorln("can not download from s3")
		return echo.ErrInternalServerError
	}

	return c.File("./profile/" + pfp)
}

func (u *User) DeleteUser(c echo.Context) error {
	username := c.Param("username")

	if err := u.repo.Delete(c.Request().Context(), model.UserInterface{
		Username: &username,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "user deleted",
	})
}

func (u *User) GetUserByKey(c echo.Context) error {
	keyword := c.QueryParam("keyword")

	users, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		Username: &keyword,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, users[0])
}

func (u *User) GetUserContacts(c echo.Context) error {
	username := c.Param("username")

	users, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		Username: &username,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	contacts, err := u.contactRepo.Get(c.Request().Context(), model.ContactInterface{
		UserID: &users[0].UserID,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, contacts)
}

func (u *User) NewUserContact(c echo.Context) error {
	username := c.Param("username")

	contactUsername := c.FormValue("username")

	users, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		Username: &username,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}
	contactUsers, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		Username: &contactUsername,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	oldContact, err := u.contactRepo.Get(c.Request().Context(), model.ContactInterface{
		UserID:        &users[0].UserID,
		ContactUserID: &contactUsers[0].UserID,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(oldContact) != 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"msg": "contact already exist",
		})
	}

	if err := u.contactRepo.Create(c.Request().Context(), model.Contact{
		UserID:        users[0].UserID,
		ContactUserID: contactUsers[0].UserID,
		Status:        model.Pending,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"msg": "contact created",
	})
}

func (u *User) DeleteUserContact(c echo.Context) error {
	username := c.Param("username")

	id, err := strconv.ParseUint(c.Param("contactid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	users, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		Username: &username,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if err = u.contactRepo.Delete(c.Request().Context(), model.ContactInterface{
		UserID:        &users[0].UserID,
		ContactUserID: &id,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "deleted contact successfuly",
	})
}

func (u *User) UpdateContact(c echo.Context) error {
	username := c.Param("username")

	id, err := strconv.ParseUint(c.Param("contactid"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest
	}

	users, err := u.repo.Get(c.Request().Context(), model.UserInterface{
		Username: &username,
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if err := u.contactRepo.Update(c.Request().Context(), model.Contact{
		UserID:        users[0].UserID,
		ContactUserID: id,
		Status:        model.Accepted,
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"msg": "contact status updated",
	})
}

func (u *User) UpdateOnlineStatus(c echo.Context) error {
	id := c.Get("userID")
	userid, _ := id.(uint64)

	if err := u.repo.Update(c.Request().Context(), model.User{
		UserID:   userid,
		LastSeen: time.Now(),
		IsActive: "false",
	}); err != nil {
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

func (u *User) NewUserHandler(g *echo.Group) {
	g.POST("/register", u.RegisterUser)
	g.POST("/login", u.LoginUser)

	userGroup := g.Group("/users")
	userGroup.GET("/:userid", u.GetUserByID, mwares.JWTMiddleware)
	userGroup.PATCH("/:username", u.UpdateUser, mwares.JWTMiddleware)
	userGroup.DELETE("/:username", u.DeleteUser, mwares.JWTMiddleware)
	userGroup.GET("/pfp/:username", u.GetUserPfpf, mwares.JWTMiddleware)
	userGroup.GET("", u.GetUserByKey, mwares.JWTMiddleware)
	userGroup.POST("/signout", u.UpdateOnlineStatus, mwares.JWTMiddleware)

	userGroup.GET("/:username/contacts", u.GetUserContacts, mwares.JWTMiddleware)
	userGroup.POST("/:username/contacts", u.NewUserContact, mwares.JWTMiddleware)
	userGroup.DELETE("/:username/contacts/:contactid", u.DeleteUserContact, mwares.JWTMiddleware)
	userGroup.PATCH("/:username/contacts/:contactid", u.UpdateContact, mwares.JWTMiddleware)
}
