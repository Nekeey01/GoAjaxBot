/*
Файл, содержащий вспомогательные функции для работы бота. В основном для обработки данных и отправки сообщений.
Изменять что-либо не советую
*/

package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"mime/multipart"
	"regexp"
	"strings"
)

// Переводим blob данные файла в tgbotapi.FileBytes
func blobToFileBytes(photo *multipart.FileHeader) tgbotapi.FileBytes {
	screenshotContent, err := photo.Open() // Открываем файл для чтения
	if err != nil {
		log.Panic(err)
	}
	screenshotContentToBytes, err := io.ReadAll(screenshotContent) // Вытаскиваем blob данные
	if err != nil {
		log.Panic(err)
	}

	photoFileBytes := tgbotapi.FileBytes{
		Name:  photo.Filename,
		Bytes: screenshotContentToBytes,
	}

	return photoFileBytes
}

// Обработка сообщения с одним документом
//
// Возвращает конфиг для отправки сообщения
func soloDocumentProcessing(args argsForProcessing) tgbotapi.Chattable {

	InputDocumentPhoto := blobToFileBytes(args.media[0])
	msg := tgbotapi.NewDocument(args.chatID, InputDocumentPhoto)

	// Если первый, добавляем текст
	if args.isFirst {
		msg.Caption = args.caption
		msg.ParseMode = parseMode
	}

	if args.msgID != 0 {
		msg.ReplyToMessageID = args.msgID
	}

	return msg
}

// Обработка сообщения с одним фото
//
// Возвращает конфиг для отправки сообщения
func soloPhotoProcessing(args argsForProcessing) tgbotapi.PhotoConfig {

	InputMediaPhoto := blobToFileBytes(args.media[0])
	msg := tgbotapi.NewPhoto(args.chatID, InputMediaPhoto)

	// Если первый, добавляем текст
	if args.isFirst {
		msg.Caption = args.caption
		msg.ParseMode = parseMode
	}

	if args.msgID != 0 {
		msg.ReplyToMessageID = args.msgID
	}

	return msg
}

// Создаем конфиги для нескольких фото или документов
//   - isFirst/isLast нжуны для определения позиции файла
//
// Возвращает конфиг документа/фото
func createInputMedia(fileBytes tgbotapi.FileBytes, isImage bool, isFirst bool, isLast bool, caption string) interface{} {
	if isImage {
		inputMedia := tgbotapi.NewInputMediaPhoto(fileBytes)
		if isFirst {
			inputMedia.ParseMode = parseMode
			inputMedia.Caption = caption
		}
		return inputMedia
	} else {
		inputMedia := tgbotapi.NewInputMediaDocument(fileBytes)
		//	Проверка на последний документ в списке нужна, чтобы текст был ниже всех документов.
		if isLast {
			inputMedia.ParseMode = parseMode
			inputMedia.Caption = caption
		}
		return inputMedia
	}
}

// Обработка списка файлов
//
// Возвращает конфиг для отправки сообщения
func multiMediaProcessing(args argsForProcessing) tgbotapi.MediaGroupConfig {
	var arr []interface{}

	// Проходимся циклом по всем файлам
	for k, v := range args.media {
		photoFileBytes := blobToFileBytes(v)
		isImage := strings.HasPrefix(v.Header.Get("Content-Type"), "image") // Проверка на фото
		arr = append(arr, createInputMedia(photoFileBytes, isImage, k == 0 && args.isFirst, k == len(args.media)-1 && args.isFirst, args.caption))
	}

	newMediaGroup := tgbotapi.NewMediaGroup(args.chatID, arr)

	if args.msgID != 0 {
		newMediaGroup.ReplyToMessageID = args.msgID
	}
	return newMediaGroup
}

// Распределяем файлы на два списка:
//   - Список фото
//   - Список документов
//
// Возвращаются оба списка
func distributionToLists[T []*multipart.FileHeader](media T) (T, T) {
	var photoList T
	var documentList T

	var currentType string

	for _, v := range media {
		currentType = v.Header.Get("Content-Type")

		if strings.HasPrefix(currentType, "image") {
			photoList = append(photoList, v)
		} else {
			documentList = append(documentList, v)
		}
	}

	return photoList, documentList
}

// Экранируем символы для МаркДАУНа
//
// Возваращет экраннированый текст, готовый к отправке
func escapingForMsg(text string) string {
	re := regexp.MustCompile(`[_*\[\]()~>#\+\-=|{}.!]`)
	caption := re.ReplaceAllStringFunc(text, func(x string) string {
		return "\\" + x
	})
	return caption
}

// Главная функция обработки данных.
//
//	1 - Разделяем файлы на два списка (distributionToLists)
//	2 - Проверяем, больше ли картинок, чем 1. Если да - отправляем сообщение
//	3 - Проверяем, есть ли картинка (обязана быть). Если да - отправляем сообщение
//	4 - Создаем сообщение с кнопкой
//	5 - Проверяем, больше ли документов, чем 1. Если да - отправляем сообщение
//	6 - Проверяем, есть ли документов. Если да - отправляем сообщение
//	7 - Отправляем сообщение с кнопками
func chooseMediaSend(bot *tgbotapi.BotAPI, bugData Params, chatID int64) {

	args := argsForProcessing{
		chatID:  chatID,
		msgID:   0,
		media:   bugData.Screenshot,
		caption: makeCaption(bugData),
		isFirst: true,
	}

	photoList, documentList := distributionToLists(args.media)

	if len(photoList) > 1 {
		args.media = photoList
		args.msgID = sendGroupMediaMessage(bot, multiMediaProcessing(args))
		args.isFirst = false
	}

	if len(photoList) == 1 {
		args.media = photoList
		args.msgID = sendOneMediaMessage(bot, soloPhotoProcessing(args))
		args.isFirst = false
	}

	msgWithBtn := tgbotapi.NewMessage(args.chatID, msgTextForReplyButtons) // создаем сообщение с кнопками
	msgWithBtn.ReplyMarkup = createInlineBtn()                             // Добавляем кнопки
	msgWithBtn.ReplyToMessageID = args.msgID                               // указываем сообщение, на которое отвечаем

	if len(documentList) > 1 {
		args.media = documentList
		args.msgID = sendGroupMediaMessage(bot, multiMediaProcessing(args))
		args.isFirst = false
	}

	if len(documentList) == 1 {
		args.media = documentList
		args.msgID = sendOneMediaMessage(bot, soloDocumentProcessing(args))
		args.isFirst = false
	}

	if _, err := bot.Send(msgWithBtn); err != nil {
		panic(err)
	}

}
