// webui/src/utils/resizeImage.js
// Downscale + compress lato client. Ritorna { file, dataURL }.
// - file: File pronto per upload (multipart).
// - dataURL: base64 (utile quando invii immagini nel messaggio).
export async function resizeImageIfNeeded(file, {
  maxWidth = 1600,
  maxHeight = 1600,
  maxBytes = 800 * 1024,     // 800 KB
  outputMime = 'image/jpeg', // meglio per la compressione
  quality = 0.9
} = {}) {
  if (!file || !file.type || !file.type.startsWith('image/')) {
    return { file, dataURL: null };
  }

  const createImage = (src) => new Promise((res, rej) => {
    const img = new Image();
    img.onload = () => res(img);
    img.onerror = rej;
    img.src = src;
  });

  const objectURL = URL.createObjectURL(file);
  const img = await createImage(objectURL);
  URL.revokeObjectURL(objectURL);

  const w = img.naturalWidth;
  const h = img.naturalHeight;

  const scale = Math.min(1, maxWidth / w, maxHeight / h);
  const targetW = Math.round(w * scale);
  const targetH = Math.round(h * scale);

  const canvas = document.createElement('canvas');
  canvas.width = targetW;
  canvas.height = targetH;
  const ctx = canvas.getContext('2d');
  ctx.drawImage(img, 0, 0, targetW, targetH);

  // PNG trasparente? resta PNG; altrimenti JPEG
  if (file.type === 'image/png' && outputMime === 'image/jpeg') {
    try {
      const a = ctx.getImageData(0, 0, 1, 1).data[3];
      if (a < 255) outputMime = 'image/png';
    } catch (_) {}
  }

  let q = quality;
  let dataURL = canvas.toDataURL(outputMime, q);

  const bytesFromDataURL = (d) => {
    const b64 = d.split(',')[1] || '';
    return Math.floor(b64.length * 3 / 4);
  };

  while (bytesFromDataURL(dataURL) > maxBytes && q > 0.5) {
    q = Math.max(0.5, q - 0.1);
    dataURL = canvas.toDataURL(outputMime, q);
  }

  const bin = atob(dataURL.split(',')[1]);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) arr[i] = bin.charCodeAt(i);
  const blob = new Blob([arr], { type: outputMime });

  const ext = outputMime === 'image/png' ? 'png' : 'jpg';
  const newName = (file.name?.replace(/\.[^.]+$/, '') || 'image') + '.' + ext;
  const outFile = new File([blob], newName, { type: outputMime, lastModified: Date.now() });

  return { file: outFile, dataURL };
}
