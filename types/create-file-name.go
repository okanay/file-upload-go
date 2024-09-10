package types

type UniqueFileName struct {
	Filename  string `json:"filename"`
	ID        string `json:"id"`
	IdWithExt string `json:"id_with_ext"`
	Type      string `json:"type"`
	Base      string `json:"base"`
}
