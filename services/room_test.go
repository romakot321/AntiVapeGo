package services

import (
  "testing"

  "antivape/schemas"
  db "antivape/db"
)

func TestCreateRoom(t *testing.T) {
  dbConn, _ := db.InitDatabase("host=localhost port=6432 user=postgres dbname=test password=postgres sslmode=disable")
  roomService := NewRoomService(dbConn)

  schema := schemas.RoomCreateSchema{
    Name: "room",
    ZoneID: 1,
    OwnerID: 1,
  }
  created := roomService.Create(schema)
  if created.Name != schema.Name || created.ID == 0 {
    t.Fatal(created, schema)
  }
}
