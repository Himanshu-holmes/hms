package model

// APIError provides a standard structure for JSON error responses.
type APIError struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"` // Can be map[string]string for validation errors, or a simple string
}

// PaginationParams defines common query parameters for paginated API endpoints.
// `form` tags are used by Gin to bind query parameters.
// `default` tag sets a default value if the parameter is not provided.
type PaginationParams struct {
	Limit  int `form:"limit,default=10" validate:"omitempty,min=1,max=100"` // Added validation
	Offset int `form:"offset,default=0" validate:"omitempty,min=0"`        // Added validation
}

// PaginatedResponse provides a standard structure for responses that include paginated data.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`        // The actual list of items
	Total      int64       `json:"total"`       // Total number of items available across all pages
	Limit      int         `json:"limit"`       // The limit used for this page
	Offset     int         `json:"offset"`      // The offset used for this page
	Page       int         `json:"page,omitempty"` // Current page number (calculated: offset/limit + 1)
	TotalPages int         `json:"total_pages,omitempty"` // Total pages (calculated: ceil(total/limit))
}