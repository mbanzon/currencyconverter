package server

var (
	server   *Server
	runError error
)

func init() {
	var err error
	server, err = New()
	if err != nil {
		panic(err)
	}

	go func() {
		runError = server.Run()
		if runError != nil {
			panic(runError)
		}
	}()
}
