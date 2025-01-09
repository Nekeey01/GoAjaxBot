    const fileInput = document.getElementById('bugReportInputFile');
    const bugReportDropzone = document.getElementById('bugReportDropzone');
    const previewArea = document.getElementById('preview-thumbnails');

    let filesToUpload = []; // Список файлов для отправки

    const bugPreviewImageModal = new bootstrap.Modal(document.getElementById('bugPreviewImageModal'));
    const bugFullImageElement = document.getElementById('bug-full-image');

    // Обработчик для файлового ввода
    fileInput.addEventListener('change', (event) => {
        handleFiles(event.target.files);
    });

    // Обработчик для дропа файлов
    bugReportDropzone.addEventListener('dragover', (event) => {
        event.preventDefault();
    });

    // Обработчик дорп зоны
    bugReportDropzone.addEventListener('drop', (event) => {
        event.preventDefault();
        const files = event.dataTransfer.files;
        handleFiles(files);
    });

    // Обработка вставки файлов из буфера обмена
    document.addEventListener('paste', (event) => {
        const items = event.clipboardData.items;
        const files = [];

        for (const item of items) {
            if (item.kind === 'file') {
                files.push(item.getAsFile());
            }
        }
        handleFiles(files);
    });

    // Функция для обработки файлов и создания миниатюр
    function handleFiles(files) {
        Array.from(files).forEach(file => {
            const fileType = file.type;

            const thumbnailContainer = document.createElement('div');
            thumbnailContainer.classList.add('thumbnail-container');

            if (fileType.startsWith('image/')) {
                const reader = new FileReader();
                reader.onload = (event) => {
                    const imgElement = document.createElement('img');
                    imgElement.src = event.target.result;
                    imgElement.addEventListener('click', () => openBugPreviewImageModal(imgElement.src));
                    // imgElement.addEventListener('click', () => openImageModal(imgElement.src));
                    thumbnailContainer.appendChild(imgElement);
                    createFileNameLabel(file, thumbnailContainer);
                };
                reader.readAsDataURL(file);
            } else {
                const icon = createFileIcon(file);
                icon.addEventListener('click', () => downloadFile(file));
                thumbnailContainer.appendChild(icon);
                createFileNameLabel(file, thumbnailContainer);
            }

            const deleteIcon = document.createElement('span');
            deleteIcon.classList.add('delete-icon');
            deleteIcon.innerHTML = '&times;'; // Using HTML entity for multiplication sign as delete icon
            deleteIcon.onclick = () => {
                thumbnailContainer.remove(); // Remove the thumbnail container
                filesToUpload = filesToUpload.filter(f => f !== file); // Update the filesToUpload array
            };
            thumbnailContainer.appendChild(deleteIcon);


            previewArea.appendChild(thumbnailContainer);
            filesToUpload.push(file); // Добавляем файл в список для загрузки
        });
    }

    // Функция для создания иконки для файлов
    function createFileIcon(file) {
        const iconDiv = document.createElement('div');
        iconDiv.classList.add('file-icon');

        const fileType = file.name.split('.').pop();

        switch (fileType) {
            case 'docx':
            case 'doc':
                iconDiv.innerHTML = '<i class="fas fa-file-word"></i>';
                iconDiv.style.backgroundColor = '#2b579a'; // Синий для документов Word
                break;
            case 'txt':
                iconDiv.innerHTML = '<i class="fas fa-file-alt"></i>';
                iconDiv.style.backgroundColor = '#4caf50'; // Зеленый для текстовых файлов
                break;
            case 'pdf':
                iconDiv.innerHTML = '<i class="fas fa-file-pdf"></i>';
                iconDiv.style.backgroundColor = '#d32f2f'; // Красный для PDF
                break;
            default:
                iconDiv.innerHTML = '<i class="fas fa-file"></i>';
                iconDiv.style.backgroundColor = '#607d8b'; // Серый для других типов файлов
        }

        return iconDiv;
    }

    // Функция для создания подписи с именем файла
    function createFileNameLabel(file, container) {

        const label = document.createElement('p');
        label.classList.add('caption');


        // undefined = авто скрин
        if(file.name === undefined){
            label.textContent = "auto_screenshot"
        }
        else{
            label.textContent = file.name; // Assuming file.name is the caption
        }
        container.appendChild(label);
    }

    // Функция для открытия изображения в модальном окне
    function openBugPreviewImageModal(imageSrc) {
        bugFullImageElement.src = imageSrc;
        bugPreviewImageModal.show();
        document.querySelector('.bugModal-image-full').style.display = 'flex'; // Центрируем модальное окно
    }

    // Функция для скачивания файла при нажатии на иконку
    function downloadFile(file) {
        const link = document.createElement('a');
        link.href = URL.createObjectURL(file);
        link.download = file.name;
        link.click();
        URL.revokeObjectURL(link.href);
    }





    // Отправка файлов на сервер через AJAX
    // uploadButton.addEventListener('click', () => {
    //     if (filesToUpload.length === 0) {
    //         alert('Нет файлов для загрузки');
    //         return;
    //     }
    //
    //     const formData = new FormData();
    //     filesToUpload.forEach((file, index) => {
    //         formData.append(`file${index}`, file);
    //     });
    //
    //     // Отправка AJAX-запроса на сервер
    //     fetch('/upload', {
    //         method: 'POST',
    //         body: formData,
    //     })
    //         .then(response => response.json())
    //         .then(data => {
    //             alert('Файлы успешно загружены!');
    //         })
    //         .catch(error => {
    //             console.error('Ошибка при загрузке файлов:', error);
    //         });
    // });
