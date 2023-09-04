package ping

import "time"

type SwaggerPing struct {
	Example string `json:"example" binging:"required" example:"asdfasdf"`
}

type SwaggerListPing struct {
	Name string    `json:"name" binging:"required" example:"wilson"`
	Age  int       `json:"age" binging:"required" example:"20"`
	Time time.Time `json:"time,omitempty" example:"2023-09-04T03:05:50.692675318+08:00"`
}
