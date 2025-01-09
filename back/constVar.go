/*
Константы и структуры. По сути это некий "конфиг файл"
*/
package main

import "mime/multipart"

// Информация про url и тип системы
type SystemInfo struct {
	bugReportFullUrl string // Полная урла
	bugReportSystem  string // Тип системы
}

// Структура, содержащая информацию о всех полях, которые пришли от клиента
type Params struct {
	System            SystemInfo              `form:"system"`     // SystemInfo
	bugReportDateTime string                  `form:"datetime"`   // Дата и время
	bugReportUrgency  string                  `form:"urgency"`    // Срочность
	bugReportErrorMsg string                  `form:"ErrorMsg"`   // Текстовое сообщение ошибки
	bugReportSteps    string                  `form:"steps"`      // Шаги для воспроизведения
	Screenshot        []*multipart.FileHeader `form:"screenshot"` // Набор файлов. Фото\документы
	bugReportDopInfo  string                  `form:"dopInfo"`    // Дополнительная информация
	bugReportErrorTrace string                  `form:"ErrorTrace"` // Дополнительная информация
}

// Информация для обработки и отправки сообщений
type argsForProcessing struct {
	chatID  int64                   // id Чата Берется из chatID
	msgID   int                     // id Сообщения
	media   []*multipart.FileHeader // Список файлов
	caption string                  // Подпись к файлам
	isFirst bool                    // Проверка на то, является ли сообщение первым отправленным
}

// Параметры для gin
const (
	ginUrl  = "/debug-form" // url, который прослушивает gin
	ginPort = ":6969" // порт, который прослушивает gin
)

// Параметры для бота
const (
	token                  = "1111:aaaaa"
	chatID                 = int64(1111)
	msgTextForReplyButtons = "Новая ошибка - поставьте статус задачи" // Текст для сообщения, к которому прикрепляются кнопки
	parseMode              = "MarkdownV2"
)

// Структура для статуса
type statusTextData struct {
	text string // Текст на кнопке
	data string // Передаваемые данные в callBack. Так же отображаются в сообщениях
}

var statusComplete = statusTextData{
	text: "Выполнено✅",
	data: "STATUS: ВЫПОЛНЕНО ✅",
}

var statusInWork = statusTextData{
	text: "В работе⚠️",
	data: "STATUS: В РАБОТЕ ⚠️",
}

var statusNotWork = statusTextData{
	text: "Не сделан❌",
	data: "STATUS: НЕ СДЕЛАН ❌",
}
