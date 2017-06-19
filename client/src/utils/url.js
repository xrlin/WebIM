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

export {generateURL};
