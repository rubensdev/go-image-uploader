import alpine from 'alpinejs';
import { v4 as uuidv4 } from 'uuid';

alpine.data('UploadForm', () => ({
  items: [],

  init() {
    const fileInput = this.$el.querySelector('.file_input');
    const uploadForm = this.$el.querySelector('#form');
    const browseBtn = this.$el.querySelector('.btn');

    uploadForm.onchange = ({ target }) => {
      Array.from(target.files).forEach((file) => {
        // if (!file) {
        //   // TODO: Return an error or something.
        //   return;
        // }
        this.items.unshift({
          id: uuidv4(),
          file,
        });
      });
      target.value = '';
    };

    browseBtn.addEventListener('click', () => {
      fileInput.click();
    });

    this.$el.addEventListener('dismiss', this.onDismiss.bind(this));
  },

  onDrop(ev) {
    if (!ev.dataTransfer.items) return;

    // Use DataTransferItemList interface to access the file(s)
    [...ev.dataTransfer.items].forEach((item, i) => {
      if (item.kind !== 'file') return;

      const file = item.getAsFile();
      this.items.unshift({
        id: uuidv4(),
        file,
      });
    });
  },

  onDismiss(event) {
    const id = event.detail;
    const index = this.items.findIndex((item) => item.id === id);

    if (index !== -1) {
      this.items.splice(index, 1);
    }
  },
}));
