package ninja

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Context provides some helper methods for a given w & r
type Context struct {
	w http.ResponseWriter
	r *http.Request
}

// NewContext makes a new context from w & r
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w: w,
		r: r,
	}
}

// Body places the json body into the given object
func (x *Context) Body(obj interface{}) error {
	return Body(x.r, obj)
}

// Body reads the body from r into v
func Body(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

// WriteJSON encodes v as json and writes it to the response
func (x *Context) WriteJSON(v interface{}) {
	err := WriteJSON(x.w, v)
	if err != nil {
		x.Error(err, http.StatusInternalServerError)
		return
	}
}

// WriteJSON encodes v as json and writes it to w
func WriteJSON(w http.ResponseWriter, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

// Error writes err and s to the response
func (x *Context) Error(err error, s int) {
	Error(x.w, err, s)
}

// Error writes err and s to w
func Error(w http.ResponseWriter, err error, s int) {
	http.Error(w, err.Error(), s)
}
