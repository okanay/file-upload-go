package types

type Assets struct {
	ID          int    `json:"id"`
	Creator     string `json:"creator"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
	Size        int64  `json:"size"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateAssetReq struct {
	Creator     string `json:"creator"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
	Size        int64  `json:"size"`
}

type UploadAssetReq struct {
	Description string `json:"description"`
	File        string `json:"file"`
	Size        string `json:"size"`
}
