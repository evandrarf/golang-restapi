package web

type ProductUpdateRequest struct {
	Id          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,max=200"`
	Description string `json:"description" validate:"required"`
}