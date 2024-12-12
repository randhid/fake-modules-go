package fakeboard

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/resource"
)

const (
	boardName      = "fake-board"
	emptyBoardName = "empty-board"
)

var (
	Model      = common.FakesFamily.WithModel(boardName)
	EmptyModel = common.FakesFamily.WithModel(emptyBoardName)
)

func init() {
	resource.RegisterComponent(board.API, Model, resource.Registration[board.Board, *Config]{
		Constructor: newFakeBoard,
	})

	resource.RegisterComponent(board.API, EmptyModel, resource.Registration[board.Board, resource.NoNativeConfig]{
		Constructor: newEmptyBoard,
	})
}
