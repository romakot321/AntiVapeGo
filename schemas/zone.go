package schemas

type ZoneCreateSchema struct {
  Name string `json:"name" binding:"required"`
  OwnerID uint `json:"owner_id" binding:"required"`
}

type ZoneSchema struct {
  ID uint `json:"id" binding:"required"`
  Name string `json:"name" binding:"required"`
  OwnerID uint `json:"owner_id" binding:"required"`
  Rooms []RoomSchema `json:"rooms`
}

type ZoneUpdateSchema struct {
  Name *string `json:"name,omitempty"`
}

type ZoneFindSchema struct {
  OwnerID *uint `json:"owner_id,omitempty"`
}
