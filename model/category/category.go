package category

type Category struct {
	CategoryID string `json:"category_id"`
	Name string	`json:"name"`
	Description string `json:"description"`
	FullPath string `json:"full_path"`
	Price string `json:"price"`
	VendorID string `json:"vendor_id"`
	ParentID string	`json:"parent_id"`
	UpdatedAt int `json:"updated_at"`
	CreatedAt int `json:"created_at"`
	HasChild bool `json:"has_child"`
	Archived bool `json:"archived"`
}


