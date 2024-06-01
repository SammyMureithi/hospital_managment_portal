package response

import "hospital_management_portal/models"

type Response struct {
	OK     bool        `json:"ok"`
	Status string      `json:"status"`
	Message string     `json:"message"`
	User   models.User `json:"user"`
}
