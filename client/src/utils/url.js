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

function getAvatarUrl(avatarHash, width, height) {
  let avatarUrl = '';
  if (!avatarHash) {
    avatarUrl = 'https://xrlin.github.io/assets/img/crown-logo.png';
  } else {
    avatarUrl = `http://oxupzzce5.bkt.clouddn.com/${avatarHash}`
  }
  if (width && height) {
    avatarUrl = `${avatarUrl}?imageView2/2/w/${width}/h/${height}`
  }
  return avatarUrl
}

function getImageUrls(imageHash, width = 100, height = 150) {
  let thumbnailUrl = `http://oxupzzce5.bkt.clouddn.com/${imageHash}?imageView2/2/w/${width}/h/${height}`;
  let imageUrl = `http://oxupzzce5.bkt.clouddn.com/${imageHash}`;
  return {thumbnailUrl, imageUrl}
}

function getMusicIdFromLink(link) {
  let pattern = /https?:\/\/music\.163\.com\/#\/song\?id=(\d+)/i;
  if (!pattern.test(link)) return false;
  return RegExp.$1;
}

export {generateURL, getAvatarUrl, getImageUrls, getMusicIdFromLink};

