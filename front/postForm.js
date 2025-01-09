$(document).ready(function () {
    const ServiceList = {
        'manager.spasskievorota.ru': "Рабочая",
        "managertest.spasskievorota.ru": "Тестовая"
    }

    const Select2UrgencyList = {
        minimumResultsForSearch: Infinity,
        data: [
            {id: 'Blocking', text: 'Блокирующая'},
            {id: 'Urgent', text: 'Срочная'},
            {id: 'Medium', text: 'Средний приоритет'},
            {id: 'Low', text: 'Низкий приоритет'}
        ]
    }

    const DefaultUrgencyList = {
        'Блокирующая': 'Блокирующая',
        'Срочная': 'Срочная',
        'Средний приоритет': 'Средний приоритет',
        'Низкий приоритет': 'Низкий приоритет'
    }

    let bugReportFullUrl_bugReportSystem;
    let full_url;
    let system;
    var dateTime = null;
    var errorMsg = 'VAM JOPA';
    let screenshot = [];
    let urgencyIsFull = false;


    // Чиним нижнюю модалку при закрытии фулл пикчи
    $('body').on('hidden.bs.modal', function () {
        if ($('.modal.show').length > 0) {
            $('body').addClass('modal-open');
        }
    });

    // Тут кароч фон такой красивый, типа реально модалку вызвали
    $(document).on('show.bs.modal', '.modal', function (event) {
        var zIndex = 1040 + (10 * $('.modal:visible').length);
        $(this).css('z-index', zIndex);
        setTimeout(function () {
            $('.modal-backdrop').not('.modal-stack').css('z-index', zIndex - 1).addClass('modal-stack');
        }, 0);
    });


    function initUrgencyList() {
        $("#bugReportUrgency").html('').select2(Select2UrgencyList); // Заполняем
        $('#bugReportUrgency').val('Urgent'); // Значение по умолчанию
        $('#bugReportUrgency').trigger('change');
        urgencyIsFull = true
    }

    function initUrgencyList2() {

        // Заполняем селект
        $.each(DefaultUrgencyList, function (i, item) {
            $('#bugReportUrgency').append($('<option>', {
                value: i,
                text: item
            }));
        });

        $('#bugReportUrgency option[value=Срочная]').attr('selected', true);

        urgencyIsFull = true

    }

    function initbugReportSystemAndUrl() {
        full_url = window.location.href
        // Включить на релизе
        // system = ServiceList[window.location.host]
        system = "Тестовая"
        bugReportFullUrl_bugReportSystem = `${full_url} - ${system}`;

        $('#bugReportFullUrl').val(full_url)
        $('#bugReportSystem').val(system)
    }

    function initDateTime() {
        var today = new Date();
        var date = today.getFullYear() + '-' + (today.getMonth() + 1) + '-' + today.getDate();
        var minute = (today.getMinutes() < 10 ? '0' : '') + today.getMinutes()
        var sec = (today.getSeconds() < 10 ? '0' : '') + today.getSeconds()
        var time = today.getHours() + ":" + minute + ":" + sec;

        dateTime = date + ' ' + time;
    }

    function initErrorMsg() {
        $("#bugReportErrorMsg").val(errorMsg);
    }

    function initAll() {
        console.log('init start')

        if (dateTime == null) {
            initDateTime()
        }

        if (urgencyIsFull === false) {
            // initbugReportUrgencyList() // select2
            initUrgencyList2()   // default select
        }
        initErrorMsg()
        initbugReportSystemAndUrl()

        console.log('init end')
    }


    function getFormData() {

        let formData = new FormData(bugReportForm)
        // formData.append('bugReportSystem', system)
        // formData.append('Full_url', full_url)
        formData.append('bugReportDateTime', dateTime)

        // Заполняем формдату скринами
        filesToUpload.forEach((file, index) => {
            formData.append(`file`, file);
        });
        console.log(formData)

        return formData
    }

    let bugReportModal;
    // Тупа клозед
    $(".closeBugReportModal").on('click', e => {
        bugReportModal.hide()
    });

    // Делаем скрин и открываем модалку
    $("#bugReportButton").on('click', e => {
        console.log('clicked bugReportButton')
        html2canvas(document.body, {
            scrollY: 0,
            scrollX: 0,
            // width: document.body.scrollWidth,
            // height: document.body.scrollHeight,
            // scale: 3,
        }).then(function (canvas) {
            // Преобразуем canvas в Blob
            canvas.toBlob(function (blob) {
                screenshot = [] // очищаем массив, потому что пользователи ебланы
                screenshot.push(blob)
                handleFiles(screenshot)

                bugReportModal = new bootstrap.Modal(document.getElementById('bugReportModal'));
                bugReportModal.show();

                // var item = new ClipboardItem({ "image/png": blob });

                // Копируем изображение в буфер обмена
                // navigator.clipboard.write([item]).then(function() {
                //     // alert("Скриншот скопирован в буфер обмена!");
                // }).catch(function(error) {
                //     console.error("Ошибка копирования в буфер:", error);
                // });
            }, "image/png");
        });
        console.log("make screen")
        initAll()

        // console.log(getFormData())
    });

    $("#submitBugReportBtn").on('click', e => {
        console.log('clicked submitBugReportBtn')
        e.preventDefault();

        let formData = getFormData()
        $.ajax({
            url: "http://localhost:8080/form",
            method: "POST",
            data: formData,
            async: false,
            processData: false,
            contentType: false,
            dataType: 'json',
            success: function (data) {
                console.log(data)
                // $('#submitBugReportBtn').html(data);
            },
        })
    })
});

