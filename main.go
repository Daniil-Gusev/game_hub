package main

import (
	"game_hub/app"
	"game_hub/core"
	"game_hub/games"
	"os"
	"os/signal"
)

func main() {
	availableGames := games.AvailableGames()
	gameCtx := &core.GameContext{
		StateStack:   core.NewStateStack(),
		Game:         nil,
		AppIsRunning: true,
	}
	lm := core.NewLocalizationManager("ru")
	appLocalizer := core.NewMessageLocalizer(lm)
	uiCtx := &core.UiContext{
		Console:         core.NewStdConsole(),
		Validator:       core.InputValidator{},
		ErrHandler:      core.NewLocalizedErrorHandler(appLocalizer),
		CommandRegistry: core.NewCommandRegistry(core.NewCommandLocalizer(lm), core.NewCommandLocalizer(lm)),
		AppLocalizer:    appLocalizer,
		GameLocalizer:   core.NewMessageLocalizer(lm),
		StateLocalizer:  core.NewStateLocalizer(lm),
	}
	if err := uiCtx.AppLocalizer.LoadTranslations("./core/translations.json"); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.StateLocalizer.LoadTranslations("./core/states.json"); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.LoadGlobalTranslations("./core/global_commands.json"); err != nil {
		return
	}
	if err := uiCtx.CommandRegistry.RegisterGlobalCommands(core.DefaultGlobalCommands()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.LoadLocalTranslations("./core/local_commands.json"); err != nil {
		return
	}
	startState := core.State(app.NewGameSelectionMenu(availableGames))
	currentState := startState
	gameCtx.StateStack.Push(currentState)
	if _, err := currentState.Init(gameCtx, uiCtx); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.RegisterLocalCommands(currentState.GetCommands()); err != nil {
		uiCtx.DisplayError(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		uiCtx.Console.Write("\r\n")
		currentState = &core.ExitState{}
		if _, err := currentState.Init(gameCtx, uiCtx); err != nil {
			uiCtx.DisplayError(err)
		}
		currentState.Display(gameCtx, uiCtx)
		currentState.Handle(gameCtx, uiCtx, "")
		os.Exit(0)
	}()

	for gameCtx.AppIsRunning {
		currentState.Display(gameCtx, uiCtx)
		input, inputErr := "", error(nil)
		if currentState.RequiresInput() {
			input, inputErr = uiCtx.Console.Read()
			uiCtx.DisplayError(inputErr)
			if appErr, ok := inputErr.(*core.AppError); ok && appErr.Code == core.ErrEOF {
				currentState = &core.ExitState{}
				_, err := currentState.Init(gameCtx, uiCtx)
				uiCtx.DisplayError(err)
				gameCtx.StateStack.Push(currentState)
				currentState.Display(gameCtx, uiCtx)
			}
		}
		nextState, err := uiCtx.HandleInput(input, gameCtx)
		uiCtx.DisplayError(err)
		if appErr, ok := err.(*core.AppError); ok && appErr.Code == core.ErrStateStack {
			currentState = startState
			_, err := currentState.Init(gameCtx, uiCtx)
			uiCtx.DisplayError(err)
			gameCtx.StateStack.Push(currentState)
			currentState.Display(gameCtx, uiCtx)
		}
		if _, ok := currentState.(*core.GameExitState); ok {
			currentState = startState
			gameCtx.StateStack.Clear()
			gameCtx.StateStack.Push(currentState)
			if _, err := currentState.Init(gameCtx, uiCtx); err != nil {
				uiCtx.DisplayError(err)
				return
			}
			if err := uiCtx.CommandRegistry.RegisterLocalCommands(currentState.GetCommands()); err != nil {
				uiCtx.DisplayError(err)
			}
			continue
		}
		if nextState != currentState {
			if nextState == startState {
				gameCtx.StateStack.Clear()
			}
			gameCtx.StateStack.Push(nextState)
			currentState = nextState
			newState, err := currentState.Init(gameCtx, uiCtx)
			uiCtx.DisplayError(err)
			if err != nil && (newState != currentState) {
				gameCtx.StateStack.Pop()
				gameCtx.StateStack.Push(newState)
				currentState = newState
			}
			if err := uiCtx.CommandRegistry.RegisterLocalCommands(currentState.GetCommands()); err != nil {
				uiCtx.DisplayError(err)
			}
		}
		uiCtx.DisplayMessage()
	}
}
