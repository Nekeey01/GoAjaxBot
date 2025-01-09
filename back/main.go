package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cors "github.com/rs/cors/wrapper/gin"
	"net/url"
)

// Точка входа в проект.
// Инициализирует бота.
// Запускает startGin и listenUpd
func main() {

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	startGin(bot, chatID)
	listenUpd(bot, chatID)

}

// Запуск gin'а с параметрами ginUrl и ginPort
func startGin(bot *tgbotapi.BotAPI, chatID int64) {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST(ginUrl, postForm(bot, chatID))

	go func() {
		err := r.RunTLS(ginPort, "certificate.crt", "certificate.key")
		if err != nil {
			log.Panic("У вас GIN не запустился, кекв")
		}
	}()
}

// Обработка входящего запроса от клиента
// Запускает getReq, convertReq и postTG
func postForm(bot *tgbotapi.BotAPI, chatID int64) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		req, files := getReq(c)                  // Получаем данные
		completeParams := convertReq(req, files) // Пихаем их в структуру

		// Возвращаем json ответ клиенту
		c.JSON(200, gin.H{
			"POST FORM": "True",
			"PARAMS":    completeParams,
		})

		postTG(bot, chatID, completeParams) // Начинаем работу бота
	}
	return gin.HandlerFunc(fn)
}

// Получение данных от клиента.
//
// Возвращает список строковых данных и список файлов
func getReq(ctx *gin.Context) (url.Values, []*multipart.FileHeader) {
	req, err := ctx.MultipartForm()

	if err != nil {
		panic("Не получилось обработать параметры от клиента(")
	}
	values := req.Value
	file := req.File["file"]

	return values, file
}

// Преобразование полученых параметров из getReq в структуру Params
//
// Возвращается готовая структура
func convertReq(reqData url.Values, files []*multipart.FileHeader) Params {

	return Params{
		System: SystemInfo{
			bugReportFullUrl: reqData.Get("bugReportFullUrl"),
			bugReportSystem:  reqData.Get("bugReportSystem"),
		},
		bugReportDateTime: reqData.Get("bugReportDateTime"),
		bugReportUrgency:  reqData.Get("bugReportUrgency"),
		bugReportErrorMsg: reqData.Get("bugReportErrorMsg"),
		bugReportSteps:    reqData.Get("bugReportSteps"),
		Screenshot:        files,
		bugReportDopInfo:  reqData.Get("bugReportDopInfo"),
		bugReportErrorTrace: reqData.Get("bugReportErrorTrace"),
	}
}

// Прослушиваем нажатия кнопок в боте в горутине.
func listenUpd(bot *tgbotapi.BotAPI, chatID int64) { //args argsForProcessing) {
u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go func(update tgbotapi.Update) {
			if update.CallbackQuery != nil {

				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data) // ловим коллбек
				if _, err := bot.Request(callback); err != nil {
					print(err)
				}

				// Получаем id сообщения, на которое был сделан реплай сообщения с кнопками
				msgId := update.CallbackQuery.Message.ReplyToMessage.MessageID
				// Меняем текст
				updatedText := regularForMsg(update.CallbackQuery.Data, update.CallbackQuery.Message.ReplyToMessage.Caption, update.CallbackQuery.From.UserName)

				// Создаем новое сообщение
				editmsg := tgbotapi.NewEditMessageCaption(chatID, msgId, escapingForMsg(editTextForReplyMsg(updatedText)))
				editmsg.ParseMode = parseMode

				_, err := bot.Send(editmsg)
				if err != nil {
					print(err)
				}

			}
		}(update)
	}

}

// Отправляем несколько фото\документов.
//
// Возвращает id сообщения
func sendGroupMediaMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MediaGroupConfig) int {
	message, err := bot.SendMediaGroup(msg)
	if err != nil {
		panic(err)
	}

	return message[0].MessageID
}

// Отправляем одно фото/документ.
//
// Возвращает id сообщения
func sendOneMediaMessage(bot *tgbotapi.BotAPI, msg tgbotapi.Chattable) int {
	message, err := bot.Send(msg)
	if err != nil {
		panic(err)
	}

	return message.MessageID
}

// Функция-посредник.
// Запускает chooseMediaSend
func postTG(bot *tgbotapi.BotAPI, chatID int64, bugData Params) {
	chooseMediaSend(bot, bugData, chatID)
}
