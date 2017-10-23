import fetch from "dva/fetch";

function parseJSON(response) {
  return response.json();
}

function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }
  const error = new Error(response.statusText);
  error.response = response;
  throw error;
}

/**
 * Requests a URL, returning a promise.
 * This function contains default options and authorization headers(if have token),
 * but can easily override them or append them.
 *
 * @param  {string} url       The URL we want to request
 * @param  {object} [options] The options we want to pass to "fetch"
 * @return {object}           An object containing either "data" or "err"
 */
async function request(url, options) {
  let defaultHeaders = new Headers({"Content-Type": "application/json"});
  if (sessionStorage.getItem('token')) {
    defaultHeaders.append('Authorization', `Bearer ${sessionStorage.getItem('token')}`);
  }
  let defaultOptions = {
    mode: 'cors',
    cache: 'default',
    headers: defaultHeaders
  };
  if (options.headers) {
    for (let pair of options.headers.entries()) {
      defaultHeaders.set(pair[0], pair[1]);
    }
    options.headers = defaultHeaders;
  }
  options = Object.assign(defaultOptions, options);
  if (options.headers && options.headers.get('content-type') === 'false') {
    options.headers.delete('content-type')
  }
  const response = await fetch(url, options);

  checkStatus(response);
  const data = await response.json();
  refreshLocalToken(response);

  return {data}
}

function refreshLocalToken(response) {
  let token = response.headers.get('Token');
  if (token) {
    sessionStorage.setItem('token', token);
  }
}

export default request;

export async function getUploadToken() {
  let {data} = await request('http://localhost:8080/api/qiniu/uptoken', {method: 'POST'});
  return data['uptoken']
}

export async function uploadImage(fileObj, fileName) {
  let QiniuUploadUrl = 'https://upload-z2.qiniup.com';
  let upToken = await getUploadToken();
  let formData = new FormData();
  formData.append('token', upToken);
  formData.append('file', fileObj, fileName);
  let header = new Headers({"Content-Type": false});
  let {data} = await request(QiniuUploadUrl, {method: 'POST', body: formData, headers: header});
  return {data}
}
