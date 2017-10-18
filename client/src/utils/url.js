/**
 * Created by xr_li on 2017/6/17.
 */

/**
 * Receive relative path and return the url with domain
 * @param {String} path   path's format is like /api/user/token
 * @return {String}
 */
function generateURL(path) {
  return 'http://localhost:8080' + path
}

function getAvatarUrl(avatarHash) {
  if (!avatarHash) return 'https://xrlin.github.io/assets/img/crown-logo.png';
  return `http://oxupzzce5.bkt.clouddn.com/${avatarHash}`
}

function getImageUrls(imageHash) {
  let thumbnailUrl = `http://oxupzzce5.bkt.clouddn.com/${imageHash}?imageView2/2/w/100/h/150`;
  let imageUrl = `http://oxupzzce5.bkt.clouddn.com/${imageHash}`;
  return {thumbnailUrl, imageUrl}
}

function getMusicIdFromLink(link) {
  let pattern = /https?:\/\/music\.163\.com\/#\/song\?id=(\d+)/i;
  if (!pattern.test(link)) return false;
  return RegExp.$1;
}

export {generateURL, getAvatarUrl, getImageUrls, getMusicIdFromLink};

