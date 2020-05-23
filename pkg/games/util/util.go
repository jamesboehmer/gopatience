package util

type UndoAction struct {
	Function func(...interface{}) error
	Args     []interface{}
}

func (action *UndoAction) call() error {
	return action.Function(action.Args...)
}

type Undoable struct {
	UndoStack []UndoAction
}

func (undoable *Undoable) Undo() error {
	if len(undoable.UndoStack) > 0 {
		action := undoable.UndoStack[len(undoable.UndoStack)-1]
		undoable.UndoStack = undoable.UndoStack[:len(undoable.UndoStack)-1]
		return action.call()
	}
	return nil
}
