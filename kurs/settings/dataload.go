package settings

import "time"

type Settings struct{
	Timeout time.Time `json: "timeout" toml: "timeout"`
	LogAccesPath string `json: "logaccess"`
	LogErrPath string `json "logerr"`

}

func createSettings(){
	timeout:= time.Microsecond*10
	s:=Settings{
		timeout: time.Millisecond*10,
		logAccesPath: "access.log",
		logErrPath: "err.log",
	}
	

}