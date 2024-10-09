package schemas

type LoginSchema struct {
  Username string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

type RegisterSchema struct {
  Name string `json:"username" binding:"required"`
  Password string `json:"password" binding:"required"`
}

type TokenSchema struct {
  Token string `json:"token" binding:"required"`
}
