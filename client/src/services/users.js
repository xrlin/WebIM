/**
 * Created by xr_li on 2017/6/17.
 */
import request from '../utils/request';
import {generateURL} from '../utils/url';

/**
 * Get token with username and password
 * @param {String} username
 * @param {String} password
 * @return {Promise}
 */
export async function requestToken(username, password) {
  let options = {
    method: 'POST',
    body: JSON.stringify({user_name: username, password: password})
  };
  let url = generateURL('/api/user/token');
  return request(url, options);
}

export async function register(username, password) {
  let options = {
    method: 'POST',
    body: JSON.stringify({user_name: username, password: password})
  };
  let url = generateURL('/api/users');
  return request(url, options)
}
