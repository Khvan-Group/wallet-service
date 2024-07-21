package api

import (
	"encoding/json"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	users "wallet-service/internal/common/model"
	"wallet-service/internal/models"
)

// Save
// @Summary Создать/Обновить кошелек пользователя
// @ID save-wallet
// @Accept json
// @Produce json
// @Param input body models.Wallet true "Информация о кошельке"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /wallets [post]
// @Security ApiKeyAuth
func (a *API) Save(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	var input models.Wallet

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	updateErr := a.Wallets.Service.Save(input)
	if updateErr != nil {
		errors.HandleError(w, updateErr)
	}

	w.WriteHeader(http.StatusOK)
}

// FindByUser
// @Summary Получить кошелек пользователя
// @ID find-wallet-by-user
// @Accept json
// @Produce json
// @Param username path string true "Логин пользователя"
// @Success 200 {object} models.Wallet "Кошелек пользователя"
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /wallets/{username} [get]
// @Security ApiKeyAuth
func (a *API) FindByUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	user := users.JwtUser{
		Login: username,
	}

	wallet, findUser := a.Wallets.Service.FindByUser(user)
	if findUser != nil {
		errors.HandleError(w, findUser)
	}

	response, err := json.Marshal(wallet)
	if err != nil {
		panic(err)
	}

	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

// Delete
// @Summary Удалить кошелек пользователя
// @ID delete-wallet
// @Accept json
// @Produce json
// @Param username path string true "Логин пользователя"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /wallets/{username} [delete]
// @Security ApiKeyAuth
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	user := users.JwtUser{
		Login: username,
	}

	if err := a.Wallets.Service.Delete(user); err != nil {
		errors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getJwtUser(r *http.Request) users.JwtUser {
	return users.JwtUser{
		Login: utils.ToString(context.Get(r, "login")),
		Role:  utils.ToString(context.Get(r, "role")),
	}
}
