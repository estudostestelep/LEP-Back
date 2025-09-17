package server

import (
	"encoding/json"
	"fmt"
	"lep/handler"
	"lep/repositories/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResourceUsers struct {
	handler *handler.Handlers
}

type IServerUsers interface {
	ServiceGetUser(c *gin.Context)
	ServiceGetUserByGroup(c *gin.Context)
	ServiceCreateUser(c *gin.Context)
	ServiceUpdateUser(c *gin.Context)
	ServiceDeleteUser(c *gin.Context)
}

func (r *ResourceUsers) ServiceGetUser(c *gin.Context) {
	id := c.Param("id")

	number, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Erro ao converter a string para inteiro:", err)
		return
	}

	resp, err := r.handler.HandlerUser.GetUser(number)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao obter o usuário")
		return
	}

	if resp == nil {
		c.String(http.StatusNotFound, "Usuário não encontrado")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *ResourceUsers) ServiceGetUserByGroup(c *gin.Context) {
	id := c.Param("id")
	resp, err := r.handler.HandlerUser.GetUserByGroup(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao obter o usuário")
		return
	}

	if resp == nil {
		c.String(http.StatusNotFound, "Usuário não encontrado")
		return
	}

	responseJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Error")
	}
	c.String(http.StatusOK, string(responseJSON))
}

func (r *ResourceUsers) ServiceCreateUser(c *gin.Context) {
	var newUser models.User
	err := c.BindJSON(&newUser)
	if err != nil {
		c.String(http.StatusBadRequest, "Erro ao decodificar dados do usuário")
		return
	}

	err = r.handler.HandlerUser.CreateUser(&newUser)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao criar o usuário")
		return
	}

	c.String(http.StatusCreated, "Usuário criado com sucesso")
}

func (r *ResourceUsers) ServiceUpdateUser(c *gin.Context) {

	var updatedUser models.User
	err := c.BindJSON(&updatedUser)
	if err != nil {
		fmt.Println("err", err)
		c.String(http.StatusBadRequest, "Erro ao decodificar dados do usuário")
		return
	}

	err = r.handler.HandlerUser.UpdateUser(&updatedUser)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao atualizar o usuário")
		return
	}

	c.String(http.StatusOK, "Usuário atualizado com sucesso")
}

func (r *ResourceUsers) ServiceDeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := r.handler.HandlerUser.DeleteUser(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao excluir o usuário")
		return
	}

	c.String(http.StatusOK, "Usuário excluído com sucesso")
}

func NewSourceServerUsers(handler *handler.Handlers) IServerUsers {
	return &ResourceUsers{handler: handler}
}
