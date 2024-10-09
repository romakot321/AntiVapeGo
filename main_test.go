package main

import (
  "testing"
  "net/http/httptest"
  "net/http"
  "encoding/json"
  "bytes"
  "bufio"
  "io"
  "fmt"
  "strconv"

  "github.com/stretchr/testify/assert"
	"github.com/gofiber/fiber/v2"
)

type testCase struct {
  description string
  route string
  expectedCode int
  method string
  body map[string]interface{}
}

func doRequest(t *testing.T, app *fiber.App, test testCase, token interface{}) *http.Response {
  encodedBody, err := json.Marshal(test.body)
  assert.NoError(t, err)
  body := bytes.NewBuffer(encodedBody)
  bodyReader := bufio.NewReader(body)

  req := httptest.NewRequest(test.method, test.route, bodyReader)
  req.ContentLength = int64(len(encodedBody))
  req.Header.Set("Content-Type", "application/json")
  if token != nil {
    req.Header.Set("Authorization", token.(string))
  }
  resp, err := app.Test(req, -1)
  assert.NoError(t, err)
  return resp
}

func doRequestReturningJson(app *fiber.App, test testCase, token interface{}) (map[string]interface{}, error) {
  var respMap map[string]interface{}
  var r interface{}
  encodedBody, err := json.Marshal(test.body)
  if err != nil {
    return respMap, err
  }
  body := bytes.NewBuffer(encodedBody)
  bodyReader := bufio.NewReader(body)

  req := httptest.NewRequest(test.method, test.route, bodyReader)
  req.ContentLength = int64(len(encodedBody))
  req.Header.Set("Content-Type", "application/json")
  if token != nil {
    req.Header.Set("Authorization", token.(string))
  }
  resp, err := app.Test(req, -1)
  if err != nil {
    return respMap, err
  }
  respBody, _ := io.ReadAll(resp.Body)
  if test.expectedCode != resp.StatusCode {
    return respMap, fmt.Errorf("Unexpected status code. Want: %d, have: %d", test.expectedCode, resp.StatusCode)
  }
  err = json.Unmarshal(respBody, &r)
  respMap = r.(map[string]interface{})
  if err != nil {
    return respMap, err
  }
  return respMap, nil
}

func createZone(app *fiber.App, name string, ownerID uint, token string) (map[string]interface{}, error) {
  zone := make(map[string]interface{}, 2)
  zone["name"] = name 
  zone["owner_id"] = ownerID
  test := testCase{"zone create", "/zone", 201, "POST", zone}
  return doRequestReturningJson(app, test, token)
}

func createRoom(app *fiber.App, name string, ownerID uint, zoneID uint, token string) (map[string]interface{}, error) {
  room := make(map[string]interface{}, 2)
  room["name"] = name 
  room["owner_id"] = ownerID
  room["zone_id"] = zoneID
  test := testCase{"room create", "/room", 201, "POST", room}
  return doRequestReturningJson(app, test, token)
}

func generateToken(t *testing.T, app *fiber.App, username, password string) string {
  account := make(map[string]interface{}, 2)
  account["username"] = username
  account["password"] = password
  test := testCase{"Login", "/auth/login", 200, "POST", account}

  resp, err := doRequestReturningJson(app, test, nil)
  assert.NoError(t, err)
  return "Bearer " + resp["data"].(map[string]interface{})["token"].(string) 
}

func TestAuth(t *testing.T) {
  t.Parallel()
  account := make(map[string]interface{}, 2)
  account["username"] = "testuser"
  account["password"] = "123456"
  app := InitApp()
  
  defer func(t *testing.T) {
    token := generateToken(t, app, account["username"].(string), account["password"].(string))
    test := testCase{"get me", "/auth/me", 200, "GET", nil}
    resp := doRequest(t, app, test, token)
    assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
    var respMap map[string]json.RawMessage
    respBody, _ := io.ReadAll(resp.Body)
    err := json.Unmarshal(respBody, &respMap)
    assert.NoError(t, err)

    userID := string(respMap["id"])
    test = testCase{"delete user", "/user/" + userID, 204, "DELETE", nil}
    resp = doRequest(t, app, test, token)
    assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
  }(t)

  tests := []testCase{
    {
      "Test register",
      "/auth/register",
      200,
      "POST",
      account,
    },
    {
      "Test login",
      "/auth/login",
      200,
      "POST",
      account,
    },
  }

  for _, test := range tests {
    resp := doRequest(t, app, test, nil)
    assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
  }
}

func TestZone(t *testing.T) {
  t.Parallel()
  zone := make(map[string]interface{})
  zone["name"] = "Test zone"
  zone["owner_id"] = 1
  app := InitApp()
  token := generateToken(t, app, "user", "password")

  createTest := testCase{
    "Test zone create",
    "/zone",
    201,
    "POST",
    zone,
  }
  respMap, err := doRequestReturningJson(app, createTest, token)
  assert.NoError(t, err)
  zone["id"] = strconv.Itoa(int(respMap["id"].(float64)))

  takeTest := testCase{
    "Test zone take",
    "/zone/" + zone["id"].(string),
    200,
    "GET",
    nil,
  }
  resp := doRequest(t, app, takeTest, token)
  respBody, _ := io.ReadAll(resp.Body)
  assert.Equalf(t, takeTest.expectedCode, resp.StatusCode, takeTest.description + " " + string(respBody))

  deleteTest := testCase{
    "Test zone delete",
    "/zone/" + zone["id"].(string),
    204,
    "DELETE",
    nil,
  }
  resp = doRequest(t, app, deleteTest, token)
  assert.Equalf(t, deleteTest.expectedCode, resp.StatusCode, deleteTest.description)
}

func TestRoom(t *testing.T) {
  t.Parallel()
  app := InitApp()
  token := generateToken(t, app, "user", "password")
  zone, err := createZone(app, "zone with room", 1, token)
  assert.NoError(t, err)

  room := make(map[string]interface{})
  room["name"] = "Test room"
  room["owner_id"] = 1
  room["zone_id"] = zone["id"].(float64)

  createTest := testCase{
    "Test room create",
    "/room",
    201,
    "POST",
    room,
  }
  resp := doRequest(t, app, createTest, token)
  respBody, _ := io.ReadAll(resp.Body)
  assert.Equalf(t, createTest.expectedCode, resp.StatusCode, createTest.description + " " + string(respBody))
  var respMap map[string]json.RawMessage
  err = json.Unmarshal(respBody, &respMap)
  assert.NoError(t, err)
  room["id"] = string(respMap["id"])

  takeTest := testCase{
    "Test room take",
    "/room/" + room["id"].(string),
    200,
    "GET",
    nil,
  }
  resp = doRequest(t, app, takeTest, token)
  respBody, _ = io.ReadAll(resp.Body)
  assert.Equalf(t, takeTest.expectedCode, resp.StatusCode, takeTest.description + " " + string(respBody))

  deleteTest := testCase{
    "Test room delete",
    "/room/" + room["id"].(string),
    204,
    "DELETE",
    nil,
  }
  resp = doRequest(t, app, deleteTest, token)
  assert.Equalf(t, deleteTest.expectedCode, resp.StatusCode, deleteTest.description)
}

func TestSensor(t *testing.T) {
  t.Parallel()
  app := InitApp()
  token := generateToken(t, app, "user", "password")
  room, err := createRoom(app, "room with sensor", 1, 1, token)
  assert.NoError(t, err)

  sensor := make(map[string]interface{})
  sensor["name"] = "Test sensor"
  sensor["owner_id"] = 1
  sensor["room_id"] = room["id"].(float64)

  createTest := testCase{
    "Test sensor create",
    "/sensor",
    200,
    "POST",
    sensor,
  }
  sensor, err = doRequestReturningJson(app, createTest, token)
  sensor["id"] = strconv.Itoa(int(sensor["id"].(float64)))
  assert.NoError(t, err)

  takeTest := testCase{
    "Test sensor take",
    "/sensor/" + sensor["id"].(string),
    200,
    "GET",
    nil,
  }
  resp := doRequest(t, app, takeTest, token)
  respBody, _ := io.ReadAll(resp.Body)
  assert.Equalf(t, takeTest.expectedCode, resp.StatusCode, takeTest.description + " " + string(respBody))

  deleteTest := testCase{
    "Test sensor delete",
    "/sensor/" + sensor["id"].(string),
    204,
    "DELETE",
    nil,
  }
  resp = doRequest(t, app, deleteTest, token)
  assert.Equalf(t, deleteTest.expectedCode, resp.StatusCode, deleteTest.description)
}
