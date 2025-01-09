import alpine from 'alpinejs';
import { generateThumbnail } from '../lib/thumbnail';

/** @type {import('../../types').UploadStatus} */
const Status = Object.freeze({
  cancelled: 'cancelled',
  error: 'error',
  file_size_exceeded: 'file_size_exceeded',
  uploading: 'uploading',
  success: 'success',
  unsupported: 'unsupported',
});

/**
 *
 * @param {Object} data
 * @param {string} data.id
 * @param {File} data.file
 * @returns
 */
function UploadItem({ id, file }) {
  const globalConfig = window.globalConfig;
  const statusMessages = i18n.upload_status;
  let xhr = new XMLHttpRequest();

  return {
    progress: 0,
    id,
    thumbnail: '',
    status: Status.idle,
    statusMsg: statusMessages.uploading,

    init() {
      if (!this.isValidMimetype()) {
        this.setStatus(Status.unsupported);
        this.statusMsg = statusMessages.unsupported;
        this.progress = 100;
        return;
      }

      xhr = new XMLHttpRequest();
      this.previewImage();
    },

    destroy() {
      xhr = null;
    },

    isValidMimetype() {
      return globalConfig.allowedMimetypes.includes(file.type);
    },

    previewImage() {
      const fileReader = new FileReader();

      fileReader.onload = (event) => {
        const blob = new Blob([event.target.result]);
        const url = URL.createObjectURL(blob);
        const img = new Image();

        img.onload = () => {
          this.thumbnail = generateThumbnail(img, 128, 128);
          URL.revokeObjectURL(url);

          if (!this.validateMaxFileSize()) {
            return;
          }

          this.uploadFile();
        };

        img.onerror = (e) => {
          this.setStatus(Status.error);
          this.statusMsg = i18n.cannot_load_img;
          this.progress = 100;
        };

        img.src = url;
      };

      fileReader.readAsArrayBuffer(file);
    },

    uploadFile() {
      this.progress = 0;
      this.setStatus(Status.uploading);
      this.statusMsg = statusMessages.uploading;

      const onUploadProgress = ({ loaded, total }) => {
        this.progress = Math.floor((loaded / total) * 100);
      };

      xhr.onload = (event) => {
        xhr.upload.removeEventListener('progress', onUploadProgress);

        if (xhr.status === 200) {
          this.setStatus(Status.success);
          this.statusMsg = statusMessages.success;
        } else {
          this.setStatus(Status.error);
          this.statusMsg = xhr.responseText ?? statusMessages.error;
        }
      };

      xhr.open('POST', globalConfig.uploadEndpoint);
      xhr.upload.addEventListener('progress', onUploadProgress);

      let data = new FormData();
      data.append('image', file);
      xhr.send(data);
    },

    cancelUpload() {
      this.setStatus(Status.cancelled);
      this.statusMsg = statusMessages.cancelled;
      this.progress = 0;
      xhr.abort();
    },

    dismiss() {
      this.$dispatch('dismiss', this.id);
    },

    retry() {
      if (!this.validateMaxFileSize()) {
        return;
      }
      this.uploadFile();
    },

    validateMaxFileSize() {
      const maxFileSizeInBytes = globalConfig.maxFileSize * 1024 * 1024;

      if (maxFileSizeInBytes >= file.size) {
        return true;
      }

      this.setStatus(Status.file_size_exceeded);
      this.statusMsg = statusMessages.file_size_exceeded;
      this.progress = 100;
      return false;
    },

    setStatus(status) {
      if (typeof Status[status] === 'undefined') {
        throw new Error(`status ${status} is not a valid status`);
      }
      this.status = status;
    },

    get filename() {
      if (file.name.length < 16) {
        return file.name.toLowerCase();
      }

      const splitName = file.name.toLowerCase().split('.');
      return splitName[0].substring(0, 12) + '... .' + splitName[1];
    },

    get isUploading() {
      return this.status === Status.uploading;
    },

    get isSuccess() {
      return this.status === Status.success;
    },

    get isPaused() {
      return this.status === Status.paused;
    },

    get hasError() {
      return this.status === Status.error;
    },

    get size() {
      const fileTotal = Math.floor(file.size / 1000);

      if (fileTotal < 1024) {
        return fileTotal + 'KB';
      }
      return (fileTotal / 1000).toFixed(2) + 'MB';
    },
  };
}

alpine.data('UploadItem', UploadItem);
