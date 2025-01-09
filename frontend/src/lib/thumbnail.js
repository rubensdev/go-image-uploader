/**
 *
 * @param {HTMLImageElement} srcImg - image to convert to thumbnail
 * @param {Number} width - thumbnail width in pixels
 * @param {Number} height - thumbnail height in pixels
 * @param {Number} [minWidth=128] - if the image has lower size than the thumbnail
 * @returns {string} the result image source (url or data URL)
 */
export const generateThumbnail = (srcImg, width, height, minWidth = 128) => {
  if (srcImg.width <= minWidth) {
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
};
