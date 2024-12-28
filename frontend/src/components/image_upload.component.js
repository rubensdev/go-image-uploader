import alpine from 'alpinejs';

const MAX_WIDTH_SIZE = 48;

alpine.data('ImageUpload', (_data) => ({
  bytesLoaded: 0,
  id: _data.id,
  name: _data.name,
  size: _data.size,
  thumbnail: '',

  init() {
    this.previewImage();
  },

  previewImage() {
    const fileReader = new FileReader();

    fileReader.onload = (event) => {
      const blob = new Blob([event.target.result]);
      const url = URL.createObjectURL(blob);
      const img = new Image();

      img.onload = () => {
        this.thumbnail = this.generateThumbnail(img, 128, 128);
        URL.revokeObjectURL(url);
        this.uploadFile();
      };

      img.src = url;
    };

    fileReader.readAsArrayBuffer(_data.imgRaw);
  },

  generateThumbnail(srcImg, width, height) {
    if (srcImg.width <= MAX_WIDTH_SIZE) {
      return srcImg.src;
    }

    const srcWidth = srcImg.width,
      srcHeight = srcImg.height;

    let canvas = document.createElement('canvas'),
      ctx = canvas.getContext('2d');

    canvas.width = width;
    canvas.height = height;

    if (srcWidth === srcHeight) {
      ctx.drawImage(srcImg, 0, 0, width, height);
      return canvas.toDataURL();
    }

    const minVal = Math.min(srcWidth, srcHeight);

    if (srcWidth > srcHeight) {
      ctx.drawImage(srcImg, (srcWidth - minVal) / 2, 0, minVal, minVal, 0, 0, width, height);
    } else {
      ctx.drawImage(srcImg, 0, (srcHeight - minVal) / 2, minVal, minVal, 0, 0, width, height);
    }

    return canvas.toDataURL();
  },

  uploadFile() {
    let xhr = new XMLHttpRequest();

    const onUploadProgress = ({ loaded, total }) => {
      this.bytesLoaded = Math.floor((loaded / total) * 100);
    };

    xhr.onload = (event) => {
      const status = xhr.responseText.length ? xhr.responseText : 'Uploaded';
      const statusCode = xhr.status;

      xhr.upload.removeEventListener('progress', onUploadProgress);
      this.onUploadCompleted(status, statusCode, this.thumbnail);
    };

    xhr.open('POST', '/upload');
    xhr.upload.addEventListener('progress', onUploadProgress);

    let data = new FormData();
    data.append('image', _data.imgRaw);
    xhr.send(data);
  },

  onUploadCompleted(status, statusCode, thumbnail) {
    this.$dispatch('uploaded', {
      id: this.id,
      status,
      statusCode,
      thumbnail,
    });
  },
}));
