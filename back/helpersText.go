/*
Файл, включающий в себя функции, взаимодействующие с текстом.
Если надо поменять текст\формат сообщения\кнопок, то менять исключительно здесь
*/
package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
)

// Форматируем текст для редактированного сообщения
func editTextForReplyMsg(text string) string {
	formattedText := strings.Replace(text, "1. URL: ", "1. URL: ```\n", -1)
	formattedText = strings.Replace(formattedText, "\n2. Система: ", "```\n2. Система: ```\n", -1)
	formattedText = strings.Replace(formattedText, "\n3. Время обнаружения: ", "```\n3. Время обнаружения: ```\n", -1)
	formattedText = strings.Replace(formattedText, "\n4. Срочность: ", "```\n4. Срочность: ```\n", -1)
	formattedText = strings.Replace(formattedText, "\n5. Текст ошибки: ", "```\n5. Текст ошибки: ```\n", -1)
	formattedText = strings.Replace(formattedText, "\n6. Trace ошибки: ", "```\n6. Trace ошибки: ```php\n", -1)
	formattedText = strings.Replace(formattedText, "\n7. Шаги воспроизведения: ", "```\n7. Шаги воспроизведения: ```\n", -1)
	formattedText = strings.Replace(formattedText, "\n8. Доп информация: ", "```\n8. Доп информация: ```\n", -1)
	formattedText = strings.Replace(formattedText, "\nSTATUS", "```\nSTATUS", -1)

	return formattedText
}

// Ищем статус в изначальном сообщении и меняем его на callback.Data (то, что отправляет кнопка)
//
// Возвращает изменный текст со статусом и
func regularForMsg(callback_data string, oldMsg string, workerName string) string {
    var updatedText string

	// Вытаскиваем статус из callback.Data
	reStatus := regexp.MustCompile(`STATUS: [^\n]+`)
	newStatus := reStatus.FindString(callback_data)
	updatedText = reStatus.ReplaceAllString(oldMsg, newStatus)

	// Вытаскиваем статус из сообщения, и меняем на полученный.
	// Исключительно для статуса "В работе"
	reWorkerInWork := regexp.MustCompile(`\n*Работу взял - @[^\n]+`)
	newWorkerInWork := reWorkerInWork.FindString(callback_data)
	updatedText = reWorkerInWork.ReplaceAllString(updatedText, newWorkerInWork)

	// Вытаскиваем статус из сообщения, и меняем на полученный.
	// Исключительно для статуса "Выполнено"
	reWorkerComplete := regexp.MustCompile(`\n*Выполнено - @[^\n]+`)
	newWorkerComplete := reWorkerComplete.FindString(callback_data)
	updatedText = reWorkerComplete.ReplaceAllString(updatedText, newWorkerComplete)

	// Вытаскиваем статус из сообщения, и меняем на полученный.
	// Исключительно для статуса "Отменено"
	reCancelComplete := regexp.MustCompile(`\n*Отменено - @[^\n]+`)
	newCancelComplete := reCancelComplete.FindString(callback_data)
	updatedText = reCancelComplete.ReplaceAllString(updatedText, newCancelComplete)

	updatedText += "\n\n"

	if newStatus == statusInWork.data {
		updatedText += "Работу взял - @" + workerName
	} else if newStatus == statusComplete.data {
		updatedText += "Выполнено - @" + workerName
	} else if newStatus == statusNotWork.data {
		updatedText += "Отменено - @" + workerName
	}

	return updatedText
}

// Формируем текст для изначального сообщения
func makeCaption(bugData Params) string {
	// Вставить на релизе
	text := "1. URL: " + bugData.System.bugReportFullUrl + "\n"
	text += "2. Система: ```\n" + bugData.System.bugReportSystem + "```\n" +
		"3. Время обнаружения: ```\n" + bugData.bugReportDateTime + "```\n" +
		"4. Срочность: ```\n" + bugData.bugReportUrgency + "```\n" +
		"5. Текст ошибки:  ```\n" + bugData.bugReportErrorMsg + "```\n" +
		"6. Trace ошибки:  ```php\n" + bugData.bugReportErrorTrace + "```\n" +
		"7. Шаги воспроизведения: ```\n" + bugData.bugReportSteps + "```\n" +
		"8. Доп информация: ```\n" + bugData.bugReportDopInfo + "```\n\n"
	text += statusNotWork.data + "\n"

	return escapingForMsg(text)
}

// Создаем кнопки
func createInlineBtn() tgbotapi.InlineKeyboardMarkup {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(statusComplete.text, statusComplete.data),
			tgbotapi.NewInlineKeyboardButtonData(statusInWork.text, statusInWork.data),
			tgbotapi.NewInlineKeyboardButtonData(statusNotWork.text, statusNotWork.data),
		))
	return keyboard
}
