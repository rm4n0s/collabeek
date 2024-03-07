package isauthenticated

type IsAuthenticatedHandler struct {
	tokenSecret []byte
}

type IsAuthenticatedForm struct {
	Token string `validate:"required,min=1"`
}
