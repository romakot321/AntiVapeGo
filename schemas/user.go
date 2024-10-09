package schemas

type UserCreateSchema struct {
  Name string `json:"name" binding:"required"`
  Password string `json:"password" binding:"required"`
}

type UserSchema struct {
  ID uint `json:"id" binding:"required"`
  Name string `json:"name" binding:"required"`
  IsSuperuser bool `json:"is_superuser" binding:"required"`
}

type UserUpdateSchema struct {
  Name string `json:"name"`
}

type UserFindSchema struct {}
