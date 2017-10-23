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
    body: JSON.stringify({user_name: username, password: password}),
    headers: new Headers({'Authorization': ''})
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

export async function getUserInfo() {
  let options = {
    method: 'POST'
  };
  let url = generateURL('/api/user/info');
  return request(url, options)
}

export function TokenMissingError() {
  this.message = 'Token is missing';
  this.name = 'TokenMissingError';
}

let _ws = null;

export function initWebSocket() {
  if (_ws !== null && _ws.readyState === 1) {
    return _ws;
  }
  let token = getLocalToken();
  _ws = new WebSocket(`ws://localhost:8080/ws/chat?token=${token}`);
  return _ws
}

function getLocalToken() {
  let token = sessionStorage.getItem("token");
  if (!token) {
    throw new TokenMissingError();
  }
  return token;
}

function getWebSocket() {
  return initWebSocket();
}

export function sendMessage(message) {
  let ws = getWebSocket();
  console.log(ws);
  ws.send(JSON.stringify({...message}));
}

export function getRecentRooms() {
  let options = {
    method: 'GET'
  };
  let url = generateURL('/api/user/rooms');
  return request(url, options)
}

export function searchUsers(name) {
  let options = {
    method: 'GET'
  };
  let url = generateURL('/api/users/search');
  return request(`${url}?name=${name}`, options)
}

export function addFriend(friend_id) {
  let options = {
    method: 'POST',
    body: JSON.stringify({friend_id: friend_id})
  };
  let url = generateURL('/api/friends');
  return request(url, options)
}

export function updateAvatar(avatarHash) {
  let options = {
    method: 'PUT',
    body: JSON.stringify({avatar: avatarHash})
  };
  let url = generateURL('/api/user/avatar');
  return request(url, options)
}

export function retrieveFriends() {
  let options = {
    method: 'GET'
  };
  let url = generateURL('/api/friends');
  return request(url, options)
}

export function createRoom(userIds) {
  let options = {
    method: 'POST',
    body: JSON.stringify({user_ids: userIds})
  };
  let url = generateURL('/api/user/rooms');
  return request(url, options)
}

export function applyFriendship(userID) {
  let options = {
    method: 'POST',
    body: JSON.stringify({user_id: userID})
  };
  let url = generateURL('/api/friendship/apply');
  return request(url, options)
}

/**
 *
 * @param uuid the uuid of message(the application)
 * @param action 'pass' or 'reject'
 */
export function checkFriendshipApplication(uuid, action) {
  let options = {
    method: 'POST',
    body: JSON.stringify({uuid, action})
  };
  let url = generateURL('/api/friendship/check');
  return request(url, options)
}

/**
 *
 * @param uuidArray An array of notifications'(also call messages) uuid
 */
export function ackReadNotifications(uuidArray) {
  let options = {
    method: 'POST',
    body: JSON.stringify({uuid_array: uuidArray})
  };
  let url = generateURL('/api/notifications/read');
  return request(url, options)
}
