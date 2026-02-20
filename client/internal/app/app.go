package app

type AndrewChatApp struct {
	MainModel *MainModel
}

func NewApp() *AndrewChatApp {
	return &AndrewChatApp{
		MainModel: NewMainModel(),
	}
}
