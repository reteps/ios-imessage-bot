# ios-imessage-bot


uses remote messages on a jailbroken ipod to send/receive messsges


add a processing functions (`processingFuncs`)
+ `func(*Chat, Message) Chat`
+ `func(*Chat, Message) (string, Chat)`


add a command (`funcMap`)
+ `func(*Chat, Message) (string, error)`
+ `func(*Chat, Message) (string, *Chat, error)`
+ `func(*Chat, Message) (*Chat, error)`

data structure inside main.go

