package main

import (
	"fmt"
	"game_hub/app"
	"game_hub/config"
	"game_hub/core"
	"game_hub/games"
	"os"
)

func main() {
	availableGames := games.AvailableGames()
	cfg, err := config.NewConfig(core.AppName)
	if err != nil {
		fmt.Printf("Failed to initialize Configuration: %v\r\n", err)
		return
	}
	appCtx := &core.AppContext{
		Config:         cfg,
		StateStack:     core.NewStateStack(),
		Game:           nil,
		AvailableGames: availableGames,
		AppIsRunning:   true,
		GoToMenu:       false,
	}
	console, err := core.NewStdReadlineConsole()
	if err != nil {
		fmt.Printf("Failed to initialize console: %v\r\n", err)
		return
	}
	defer func() {
		if consoleErr := console.Close(); consoleErr != nil && err == nil {
			err = consoleErr
		}
	}()
	lm, err := core.NewLocalizationManager(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize localization manager: %v\r\n", err)
		return
	}
	appMessageLocalizer := core.NewMessageLocalizer(lm)
	errorHandler := core.NewLocalizedErrorHandler(appMessageLocalizer)
	logger := core.NewStdLogger(os.Stdout, errorHandler)
	lm.SetLogger(logger)
	uiCtx := &core.UiContext{
		Console:             console,
		Validator:           core.InputValidator{},
		ErrorHandler:        errorHandler,
		Logger:              logger,
		CommandRegistry:     core.NewCommandRegistry(core.NewCommandLocalizer(lm), core.NewCommandLocalizer(lm)),
		LocalizationManager: lm,
		AppLocalizer:        appMessageLocalizer,
		GameLocalizer:       core.NewMessageLocalizer(lm),
		StateLocalizer:      core.NewStateLocalizer(lm),
	}
	if err := uiCtx.AppLocalizer.LoadTranslations(appCtx.Config.Paths.CoreTranslationsPath()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.StateLocalizer.LoadTranslations(appCtx.Config.Paths.CoreStatesPath()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.LoadGlobalTranslations(appCtx.Config.Paths.CoreGlobalCommandsPath()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.RegisterGlobalCommands(core.DefaultGlobalCommands()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	if err := uiCtx.CommandRegistry.LoadLocalTranslations(appCtx.Config.Paths.CoreLocalCommandsPath()); err != nil {
		uiCtx.DisplayError(err)
		return
	}
	startState := core.State(&app.StartState{})
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
