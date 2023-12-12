package models

type TemplateCreateReq struct {
	TemplateName string `json:"template_name"`
}

type TemplateUpdateReq struct {
	Id           string    `json:"id"`
	TemplateName string `json:"template_name"`
}

type TemplateGetReq struct {
	Id string `json:"id"`
}

type TemplateFindReq struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	OrderByCreatedAt uint64 `json:"order_by_created_at"`
	Search           string `json:"search"`
}

type TemplateDeleteReq struct {
	Id string `json:"id"`
}

type TemplateFindResponse struct {
	Templates []*TemplateResponse `json:"templates"`
	Count     int                 `json:"count"`
}

type TemplateResponse struct {
	Id           string    `json:"id"`
	TemplateName string `json:"template_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}