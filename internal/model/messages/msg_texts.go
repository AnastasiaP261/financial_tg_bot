package messages

var (
	ErrTxtUnknownCommand      = "Не знаю эту команду"
	ErrTxtInvalidInput        = "Кажется, вы ошиблись при вводе команды. Введите /help, чтобы посмотреть шаблоны команд"
	ErrTxtCategoryDoesntExist = "Кажется, такой категории еще нет. Вы можете создать ее с помощью команды /category"
	ErrTxtInvalidCurrency     = "Вы можете выбрать только одну из следующих валют: RUB, USD, EUR, CNY. Пожалуйста, введите команду /currency заново с одной из доступных валют"

	ScsTxtPurchaseAdded   = "Трата добавлена"
	ScsTxtCategoryAdded   = "Категория добавлена"
	ScsTxtCurrencyChanged = "Ваша основная валюта изменена"
)
