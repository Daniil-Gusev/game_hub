package main

import (
	"fmt"
	"game_hub/app"
	"game_hub/config"
	"game_hub/core"
	"game_hub/games"
)

func main() {
	availableGames := games.AvailableGames()
	lm := core.NewLocalizationManager("")
	appMessageLocalizer := core.NewMessageLocalizer(lm)
	console, err := core.NewStdReadlineConsole()
	if err != nil {
		fmt.Errorf("Failed to initialize console: %v\r\n", err)
		return
	}
	defer console.Close()
	uiCtx := &core.UiContext{
		Console:             console,
		Validator:           core.InputValidator{},
		ErrHandler:          core.NewLocalizedErrorHandler(appMessageLocalizer),
		CommandRegistry:     core.NewCommandRegistry(core.NewCommandLocalizer(lm), core.NewCommandLocalizer(lm)),
		LocalizationManager: lm,
		AppLocalizer:        appMessageLocalizer,
		GameLocalizer:       core.NewMessageLocalizer(lm),
		StateLocalizer:      core.NewStateLocalizer(lm),
	}
	cfg, err := config.NewConfig(core.AppName)
	if err != nil {
		uiCtx.DisplayError(err)
		return
	}
	appCtx := &core.AppContext{
		Config:       cfg,
		StateStack:   core.NewStateStack(),
		Game:         nil,
		AppIsRunning: true,
		GoToMenu:     false,
	}
	lm.SetLanguage(appCtx.Config.Language.CurrentLanguage)
	if err := uiCtx.AppLocalizer.LoadTranslations(appCtx.Config.Paths.CoreTranslationsPath()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.StateLocalizer.LoadTranslations(appCtx.Config.Paths.CoreStatesPath()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.LoadGlobalTranslations(appCtx.Config.Paths.CoreGlobalCommandsPath()); err != nil {
		return
	}
	if err := uiCtx.CommandRegistry.RegisterGlobalCommands(core.DefaultGlobalCommands()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.LoadLocalTranslations(appCtx.Config.Paths.CoreLocalCommandsPath()); err != nil {
		return
	}
	startState := core.State(app.NewGameSelectionMenu(availableGames))
	currentState := startState
	appCtx.StateStack.Push(currentState)
	if _, err := currentState.Init(appCtx, uiCtx); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.RegisterLocalCommands(currentState.GetCommands()); err != nil {
		uiCtx.DisplayError(err)
	}

	for appCtx.AppIsRunning {
		currentState.Display(appCtx, uiCtx)
		input, inputErr := "", error(nil)
		if currentState.RequiresInput() {
			input, inputErr = uiCtx.Console.Read()
			uiCtx.DisplayError(inputErr)
			if appErr, ok := inputErr.(*core.AppError); ok && appErr.Code == core.ErrEOF {
				currentState = &core.ExitState{}
				_, err := currentState.Init(appCtx, uiCtx)
				uiCtx.DisplayError(err)
				appCtx.StateStack.Push(currentState)
				currentState.Display(appCtx, uiCtx)
			}
		}
		nextState, err := uiCtx.HandleInput(input, appCtx)
		uiCtx.DisplayError(err)
		if appErr, ok := err.(*core.AppError); ok && appErr.Code == core.ErrStateStack {
			currentState = startState
			_, err := currentState.Init(appCtx, uiCtx)
			uiCtx.DisplayError(err)
			appCtx.StateStack.Push(currentState)
			currentState.Display(appCtx, uiCtx)
		}
		if appCtx.GoToMenu {
			currentState = startState
			appCtx.StateStack.Clear()
			appCtx.StateStack.Push(currentState)
			if _, err := currentState.Init(appCtx, uiCtx); err != nil {
				uiCtx.DisplayError(err)
				return
			}
			if err := uiCtx.CommandRegistry.RegisterLocalCommands(currentState.GetCommands()); err != nil {
				uiCtx.DisplayError(err)
			}
			appCtx.GoToMenu = false
			continue
		}
		if nextState != currentState {
			if nextState == startState {
				appCtx.StateStack.Clear()
			}
			appCtx.StateStack.Push(nextState)
			currentState = nextState
			newState, err := currentState.Init(appCtx, uiCtx)
			uiCtx.DisplayError(err)
			if err != nil && (newState != currentState) {
				appCtx.StateStack.Pop()
				appCtx.StateStack.Push(newState)
				currentState = newState
			}
			if err := uiCtx.CommandRegistry.RegisterLocalCommands(currentState.GetCommands()); err != nil {
				uiCtx.DisplayError(err)
			}
		}
		uiCtx.DisplayMessage()
	}
}
