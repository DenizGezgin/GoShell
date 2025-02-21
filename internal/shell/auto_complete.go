package shell

type AutoComplete interface {
	CompleteCommand() []string
}

type autoComplete struct {
	commandToken string
}

func NewAutoComplete(commandToken string) AutoComplete {
	return &autoComplete{commandToken: commandToken}
}

func (a *autoComplete) CompleteCommand() []string {
	if len(a.commandToken) < 1 {
		return nil
	}

	commandRepository := GetShell().GetCommandRepository()
	matchGroups := commandRepository.GetAllCommandNamesGroupedByPrefix(a.commandToken)

	if len(matchGroups) == 0 {
		return nil
	}

	// If there's only one group, we can complete with the first match
	if len(matchGroups) == 1 {
		matches := matchGroups[0]
		if len(matches) > 0 {
			return []string{matches[0]}
		}
	}

	// If there are multiple groups, return all matches flattened
	var allMatches []string
	for _, group := range matchGroups {
		allMatches = append(allMatches, group...)
	}
	return allMatches
}
