package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/mongo-experiments/go/pkg/api"
	"github.com/mongo-experiments/go/pkg/models"
)

type UserService interface {
	GetUsers() ([]api.User, error)
	CreateUser(user api.User) (*api.User, error)
}

type UserController struct {
	UserService UserService
}

func NewUserController(users UserService) *UserController {
	return &UserController{
		UserService: users,
	}
}

// Routes for users.
func (c *UserController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", c.List)
	r.Post("/", c.Create)
	return r
}

// List renders all the users.
func (c *UserController) List(w http.ResponseWriter, r *http.Request) {

	list, err := c.UserService.GetUsers()
	if err != nil {
		CheckError(err, w, r)
	}
	res := &models.Users{}
	for _, user := range list {
		res.Users = append(res.Users, models.ToResponseUser(&user))
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, res)
	return
}

// Create stores a new user.
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {

	data := &models.UserPayload{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}

	newUser := api.User{
		Username:      data.Username,
		Name:          data.Username,
		Active:        *data.Active,
		AnimesWatched: *data.AnimesWatched,
	}

	res, err := c.UserService.CreateUser(newUser)
	if err != nil {
		CheckError(err, w, r)
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, models.ToResponseUser(res))
	return
}