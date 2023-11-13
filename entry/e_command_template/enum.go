package e_command_template

import "github.com/goccy/go-json"

type Protocol int

const (
	ProtocolNone Protocol = iota
	Http
	Websocket
	Mqtt
	RedisTopic
)

func (p Protocol) String() string {
	return [...]string{"", "http", "websocket", "mqtt", "redis_topic"}[p]
}

func (p Protocol) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Protocol) UnmarshalJSON(data []byte) error {
	var tp string
	err := json.Unmarshal(data, &tp)
	if err != nil {
		return err
	}
	*p = S2Protocol(&tp)
	return nil
}

type HTTPMethod int

const (
	GET HTTPMethod = iota
	POST
	PATCH
	PUT
	DELETE
)

func (hm HTTPMethod) String() string {
	return [...]string{"GET", "POST", "PATCH", "PUT", "DELETE"}[hm]
}

func (hm HTTPMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(hm.String())
}

func (hm *HTTPMethod) UnmarshalJSON(data []byte) error {
	var th string
	err := json.Unmarshal(data, &th)
	if err != nil {
		return err
	}
	*hm = S2HTTPMethod(&th)
	return nil
}

type BodyType int

const (
	BodyTypeNone BodyType = iota
	Text
	HTML
	XML
	FormData
	XWWWFormUrlencoded
	Json
)

func (b BodyType) String() string {
	return [...]string{"", "text", "html", "xml", "form_data", "x_www_form_urlencoded", "json"}[b]
}

func (b BodyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *BodyType) UnmarshalJSON(data []byte) error {
	var th string
	err := json.Unmarshal(data, &th)
	if err != nil {
		return err
	}
	*b = S2BodyType(&th)
	return nil
}
