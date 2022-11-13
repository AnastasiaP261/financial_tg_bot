package messages

var (
	ErrTxtUnknownCommand  = "Не знаю эту команду"
	ErrTxtInvalidInput    = "Кажется, вы ошиблись при вводе команды. Введите /help, чтобы посмотреть шаблоны команд"
	ErrTxtInvalidCurrency = "Вы можете выбрать только одну из следующих валют: RUB, USD, EUR, CNY. Пожалуйста, введите команду /currency заново с одной из доступных валют"
	ErrTxtInvalidStatus   = "Не верный статус, попробуйте заново"

	ScsTxtPurchaseAdded        = "Трата добавлена"
	ScsTxtCategoryCreated      = "Категория создана"
	ScsTxtCategoryAddedToUser  = "Категория добавлена вам"
	ScsTxtCategoryAddSelected  = "Вы выбрали создание новой категории. Создайте категорию с помощью команды /category, а затем введите трату заново"
	ScsTxtCurrencyChanged      = "Ваша основная валюта изменена"
	ScsTxtLimitChanged         = "Лимит установлен. Для того, чтобы сбросить лимит, отправьте \"/limit -1\""
	ScsTxtReportRequestCreated = "Отчет готовится..."

	ButtonTxtCreateCategory = "Создать категорию"
)
