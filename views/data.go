package views

import "log"

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
	AlertMsgGeneric = "Something went wrong!!! ERROR!!!!"
)

type PublicError interface {
	error
	Public() string
}

type Alert struct {
	Level   string
	Message string
}
type Data struct {
	Alert *Alert
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// func (d *Data) SetAlert(err error) {
// 	d.Alert = &Alert{
// 	  Level:   AlertLvlError,
// 	  Message: err.Error(),
// 	}
//   }

//   func (d *Data) SetAlert(err error) {
// 	var msg string
// 	if err is public {
// 	  // Public() would return the public error message
// 	  msg = err.Public()
// 	} else {
// 	  msg = AlertMsgGeneric
// 	}
// 	d.Alert = &Alert{
// 	  Level:   AlertLvlError,
// 	  Message: msg,
// 	}
//   }
