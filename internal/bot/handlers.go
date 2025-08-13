package bot

func (b *Bot) handleStart(chatID int64) {
	text := `<b>Привет! Я ДейлиБот - твой помощник на каждый день!</b>

<b>Мои команды:</b>
/help - подробная справка`

	b.sendMessage(chatID, text)
}

func (b *Bot) handleHelp(chatID int64) {
	text := `<b>Справка по командам:</b>

<i>Бот работает на языке Go и использует официальные API</i>`

	b.sendMessage(chatID, text)
}
