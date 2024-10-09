package schemas

import (
  "encoding/json"
)

func SchemaToMap(schema interface{}) map[string]interface{} {
  var ret map[string]interface{}
  inrec, _ := json.Marshal(schema)
  json.Unmarshal(inrec, &ret)
  return ret
}
