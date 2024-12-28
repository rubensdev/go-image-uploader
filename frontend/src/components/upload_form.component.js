import './upload_form.component.css';
import alpine from 'alpinejs';
import { v4 as uuidv4 } from 'uuid';

alpine.data('UploadForm', () => ({
  files: [],
  uploadHistory: [],

  init() {
    const fileInput = this.$el.querySelector('.fu_form__file_input');
    const uploadForm = this.$el.querySelector('#form');

    uploadForm.onchange = ({ target }) => {
      const file = target.files[0];

      if (!file) {
        // TODO: Return an error or something.
        return;
      }

      const fileTotal = Math.floor(file.size / 1000);
      let fileName = file.name;
      let fileSize;

      if (fileName.length >= 12) {
        let splitName = fileName.split('.');
        fileName = splitName[0].substring(0, 8) + '... .' + splitName[1];
      }

      if (fileTotal < 1024) {
        fileSize = fileTotal + 'KB';
      } else {
        fileSize = (fileTotal / 1000).toFixed(2) + ' MB';
      }

      const fileData = this.newFileData(fileName, fileSize, file);
      this.files.unshift(fileData);
      target.value = '';
    };

    uploadForm.addEventListener('click', () => {
      fileInput.click();
    });

    this.$el.addEventListener('uploaded', this.onUploadCompleted.bind(this));
  },

  addUploadHistory(data, status, statusCode) {
    this.uploadHistory.unshift({ ...data, status, statusCode });
  },

  newFileData(name, size, imgRaw) {
    return {
      id: uuidv4(),
      name,
      size,
      imgRaw,
    };
  },

  onUploadCompleted(event) {
    const { id, status, statusCode, thumbnail } = event.detail;

    const index = this.files.findIndex((file) => file.id === id);
    if (index !== -1) {
      const file = this.files[index];

      this.files.splice(index, 1);
      this.uploadHistory.unshift({
        id,
        name: file.name,
        size: file.size,
        status,
        statusCode,
        thumbnail,
      });
    }
  },
}));
