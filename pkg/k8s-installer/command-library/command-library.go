package command_library

type commandStage uint8

const (
	os commandStage = iota
)

type commandBook map[commandStage][]commandAndParser

type commandAndParser struct {
	command   command
	parser    outputParser
	condition condition
}

type command string

type outputParser func([]byte, interface{}) error

// condition означает востребованность выполнения операции
type condition uint8

const (
	sufficient condition = iota // После команды с таким condition исполнение дальнейшей цепочки остановится
	required                    // Команда обязательна к выполнению
	anyway                      // Безразличен результат выполнения команды
)
