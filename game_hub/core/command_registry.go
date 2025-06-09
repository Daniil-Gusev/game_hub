package core

import (
	"sort"
	"strings"
)

type CommandRegistry struct {
	globalCommands  []Command
	localCommands   []Command
	globalAliasMap  map[string]string
	localAliasMap   map[string]string
	globalLocalizer *CommandLocalizer
	localLocalizer  *CommandLocalizer
}

func NewCommandRegistry(globalLocalizer, localLocalizer *CommandLocalizer) *CommandRegistry {
	return &CommandRegistry{
		globalAliasMap:  make(map[string]string),
		localAliasMap:   make(map[string]string),
		globalLocalizer: globalLocalizer,
		localLocalizer:  localLocalizer,
	}
}

func (r *CommandRegistry) LoadGlobalTranslations(filePath string) error {
	return r.globalLocalizer.LoadTranslations(filePath)
}

func (r *CommandRegistry) LoadLocalTranslations(filePath string) error {
	return r.localLocalizer.LoadTranslations(filePath)
}

func (r *CommandRegistry) UpdateAliases() {
	r.updateGlobalAliases()
	r.updateLocalAliases()
}

func (r *CommandRegistry) containsCommand(cmds []Command, cmd Command) bool {
	for _, c := range cmds {
		if c.Id() == cmd.Id() {
			return true
		}
	}
	return false
}

func (r *CommandRegistry) GetName(cmd Command) (string, error) {
	cmdScope := cmd.Scope()
	cmdId := cmd.Id()
	if r.containsCommand(r.localCommands, cmd) {
		if r.localLocalizer.Exists(cmdScope, cmdId) {
			return r.localLocalizer.GetName(cmdScope, cmdId)
		}
		return "", NewAppError(ErrLocalization, "command_not_found", map[string]any{
			"scope": cmdScope,
			"cmd":   cmdId,
		})
	}
	if r.containsCommand(r.globalCommands, cmd) {
		if r.globalLocalizer.Exists(cmdScope, cmdId) {
			return r.globalLocalizer.GetName(cmdScope, cmdId)
		}
		return "", NewAppError(ErrLocalization, "command_not_found", map[string]any{
			"scope": cmdScope,
			"cmd":   cmdId,
		})
	}
	return "", NewAppError(ErrLocalization, "command_not_found", map[string]any{
		"scope": cmdScope,
		"cmd":   cmdId,
	})
}

func (r *CommandRegistry) GetDescription(cmd Command) (string, error) {
	cmdScope := cmd.Scope()
	cmdId := cmd.Id()
	if r.containsCommand(r.localCommands, cmd) {
		if r.localLocalizer.Exists(cmdScope, cmdId) {
			return r.localLocalizer.GetDescription(cmdScope, cmdId)
		}
		return "", NewAppError(ErrLocalization, "command_not_found", map[string]any{
			"scope": cmdScope,
			"cmd":   cmdId,
		})
	}
	if r.containsCommand(r.globalCommands, cmd) {
		if r.globalLocalizer.Exists(cmdScope, cmdId) {
			return r.globalLocalizer.GetDescription(cmdScope, cmdId)
		}
		return "", NewAppError(ErrLocalization, "command_not_found", map[string]any{
			"scope": cmdScope,
			"cmd":   cmdId,
		})
	}
	return "", NewAppError(ErrLocalization, "command_not_found", map[string]any{
		"scope": cmdScope,
		"cmd":   cmdId,
	})
}

func (r *CommandRegistry) GetAliases(cmd Command) ([]string, error) {
	cmdScope := cmd.Scope()
	cmdId := cmd.Id()
	if r.containsCommand(r.localCommands, cmd) {
		if r.localLocalizer.Exists(cmdScope, cmdId) {
			return r.localLocalizer.GetAliases(cmdScope, cmdId)
		}
		return []string{}, NewAppError(ErrLocalization, "command_not_found", map[string]any{
			"scope": cmdScope,
			"cmd":   cmdId,
		})
	}
	if r.containsCommand(r.globalCommands, cmd) {
		if r.globalLocalizer.Exists(cmdScope, cmdId) {
			return r.globalLocalizer.GetAliases(cmdScope, cmdId)
		}
		return []string{}, NewAppError(ErrLocalization, "command_not_found", map[string]any{
			"scope": cmdScope,
			"cmd":   cmdId,
		})
	}
	return []string{}, NewAppError(ErrLocalization, "command_not_found", map[string]any{
		"scope": cmdScope,
		"cmd":   cmdId,
	})
}

func (r *CommandRegistry) updateGlobalAliases() {
	r.globalAliasMap = make(map[string]string)
	for _, cmd := range r.globalCommands {
		cmdScope := cmd.Scope()
		cmdId := cmd.Id()
		if r.globalLocalizer.Exists(cmdScope, cmdId) {
			name, _ := r.globalLocalizer.GetName(cmdScope, cmdId)
			r.globalAliasMap[strings.ToLower(name)] = cmdId
			aliases, _ := r.globalLocalizer.GetAliases(cmdScope, cmdId)
			for _, alias := range aliases {
				r.globalAliasMap[strings.ToLower(alias)] = cmdId
			}
		}
	}
}

func (r *CommandRegistry) updateLocalAliases() {
	r.localAliasMap = make(map[string]string)
	for _, cmd := range r.localCommands {
		cmdScope := cmd.Scope()
		cmdId := cmd.Id()
		if r.localLocalizer.Exists(cmdScope, cmdId) {
			name, _ := r.localLocalizer.GetName(cmdScope, cmdId)
			r.localAliasMap[strings.ToLower(name)] = cmdId
			aliases, _ := r.localLocalizer.GetAliases(cmdScope, cmdId)
			for _, alias := range aliases {
				r.localAliasMap[strings.ToLower(alias)] = cmdId
			}
		}
	}
}

func (r *CommandRegistry) findCommandById(cmds []Command, id string) Command {
	for _, cmd := range cmds {
		if cmd.Id() == id {
			return cmd
		}
	}
	return nil
}

func (r *CommandRegistry) FindCommandWithoutLocalization(commandList []Command, translations CommandTranslations) Command {
	for _, cmd := range commandList {
		cmds, scopeExists := translations[cmd.Scope()]
		if !scopeExists {
			return cmd
		}
		if _, exists := cmds[cmd.Id()]; !exists {
			return cmd
		}
	}
	return nil
}

func (r *CommandRegistry) RegisterGlobalCommands(cmds []Command) error {
	if cmd := r.FindCommandWithoutLocalization(cmds, r.globalLocalizer.Translations); cmd != nil {
		return NewAppError(ErrLocalization, "command_localization_not_found", map[string]any{
			"scope": cmd.Scope(),
			"cmd":   cmd.Id(),
		})
	}
	r.globalCommands = cmds
	r.sortCommands(r.globalCommands)
	r.updateGlobalAliases()
	return nil
}

func (r *CommandRegistry) RegisterLocalCommands(cmds []Command) error {
	if cmd := r.FindCommandWithoutLocalization(cmds, r.localLocalizer.Translations); cmd != nil {
		return NewAppError(ErrLocalization, "command_localization_not_found", map[string]any{
			"scope": cmd.Scope(),
			"cmd":   cmd.Id(),
		})
	}
	r.localCommands = cmds
	r.sortCommands(r.localCommands)
	r.updateLocalAliases()
	return nil
}

func (r *CommandRegistry) GetGlobalCommands() []Command {
	return r.globalCommands
}

func (r *CommandRegistry) GetLocalCommands() []Command {
	return r.localCommands
}

func (r *CommandRegistry) sortCommands(cmds []Command) {
	sort.Slice(cmds, func(i, j int) bool {
		nameI, _ := r.GetName(cmds[i])
		nameJ, _ := r.GetName(cmds[j])
		return nameI < nameJ
	})
}

func (r *CommandRegistry) GetCommand(input string) Command {
	input = strings.ToLower(input)
	if cmdId, exists := r.localAliasMap[input]; exists {
		return r.findCommandById(r.localCommands, cmdId)
	}
	if cmdId, exists := r.globalAliasMap[input]; exists {
		return r.findCommandById(r.globalCommands, cmdId)
	}
	if cmd := r.findCommandOrAliasByPrefix(r.localCommands, r.localAliasMap, input); cmd != nil {
		return cmd
	}
	if cmd := r.findCommandOrAliasByPrefix(r.globalCommands, r.globalAliasMap, input); cmd != nil {
		return cmd
	}
	return nil
}

func (r *CommandRegistry) findCommandOrAliasByPrefix(cmds []Command, aliasMap map[string]string, prefix string) Command {
	// 1. Проверяем команды по префиксу
	if cmd := r.findCommandByPrefix(cmds, prefix); cmd != nil {
		return cmd
	}
	// 2. Проверяем алиасы по префиксу
	for alias, cmdId := range aliasMap {
		if strings.HasPrefix(alias, prefix) {
			return r.findCommandById(cmds, cmdId)
		}
	}
	return nil
}

func (r *CommandRegistry) findCommandByPrefix(cmds []Command, prefix string) Command {
	index := sort.Search(len(cmds), func(i int) bool {
		name, _ := r.GetName(cmds[i])
		return name >= prefix
	})
	if index < len(cmds) {
		name, _ := r.GetName(cmds[index])
		if strings.HasPrefix(name, prefix) {
			return cmds[index]
		}
	}
	return nil
}

func (r *CommandRegistry) ParseInput(input string) (Command, []string) {
	if input == "" {
		return nil, []string{}
	}
	args := strings.Fields(input)
	cmdPart := args[0]
	args[0] = strings.ToLower(cmdPart)
	if cmd := r.GetCommand(cmdPart); cmd != nil {
		return cmd, args
	}
	return nil, []string{}
}
