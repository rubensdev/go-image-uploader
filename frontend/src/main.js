import './style.css';
import.meta.glob('./svg/*.svg', {
  eager: true,
  query: '?no-inline',
});
import alpine from 'alpinejs';
import './components/image_upload.component';
import './components/upload_form.component';

window.Alpine = alpine;
window.Alpine.start();
