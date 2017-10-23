import {
  ackReadNotifications,
  addFriend,
  applyFriendship,
  createRoom,
  getRecentRooms,
  getUserInfo,
  initWebSocket,
  register,
  requestToken,
  retrieveFriends,
  sendMessage,
  TokenMissingError,
  updateAvatar
} from "../services/users";
import {routerRedux} from "dva/router";
import {leaveRoom} from "../services/rooms";
import {generateUUID} from "../utils/common";
import {ackMessages, retrieveUnreaOfflinedMessages} from "../services/messages";

const SingleMessage = 1, RoomMessage = 2;

export default {
  namespace: 'users',
  state: {
    roomIDs: [],
    rooms: new Map(),
    currentRoom: {},
    messages: {},
    newMessages: {},
    notifications: [],
    roomInputMessages: {},
    friends: [],
    info: {}
  },
  reducers: {
    updateInfo(state, {payload: info}) {
      return Object.assign({}, state, {...state, info})
    },
    updateRecentRooms(state, {payload: rooms}) {
      let currentRooms = state.rooms;
      for (let room of rooms) {
        currentRooms.set(room['id'], room)
      }
      let roomIDs = state.roomIDs;
      roomIDs = Array.from(currentRooms.keys()).sort((a, b) => (b - a));
      return Object.assign({}, state, {...state, rooms: currentRooms, roomIDs})
    },
    addMessage(state, {payload: message}) {
      const FriendshipMessage = 7;
      // Add notification
      if (message.msg_type === FriendshipMessage) {
        let {notifications} = state;
        notifications.push(message);
        return Object.assign({}, state, {...state, notifications: notifications});
      }
      let {messages, newMessages} = state;
      let room_id = message.room_id;
      if (messages[room_id]) {
        messages[room_id] = messages[room_id].concat(message);
      } else {
        messages[room_id] = [message]
      }
      if (room_id !== state.currentRoom.id) {
        if (newMessages[room_id]) {
          newMessages[room_id] = newMessages[room_id].concat(message)
        } else {
          newMessages[room_id] = [message]
        }
      }
      let roomIDs = state.roomIDs;
      bubbleRoom(roomIDs, room_id);
      return Object.assign({}, state, {...state, roomIDs, messages, newMessages});
    },
    stickRoom(state, {payload: roomID}) {
      let roomIDs = state.roomIDs;
      bubbleRoom(roomIDs, roomID);
      let currentRoom = state.rooms.get(roomID) || {};
      let newMessages = clearUnreadMessage(state.newMessages, id);
      return Object.assign({}, state, {...state, roomIDs, currentRoom, newMessages});
    },
    updateFriends(state, {payload: friends}) {
      return Object.assign({}, state, {...state, friends});
    },
    setCurrentRoom(state, {payload: id}) {
      let currentRoom = state.rooms.get(id) || {};
      let newMessages = clearUnreadMessage(state.newMessages, id);
      return Object.assign({}, state, {...state, currentRoom, newMessages})
    },
    removeRoom(state, {payload: roomID}) {
      let rooms = state.rooms;
      let roomIDs = state.roomIDs;
      let currentRoom = state.currentRoom;
      if (state.currentRoom.id === roomID) {
        currentRoom = {}
      }
      rooms.delete(roomID);
      let idx = roomIDs.indexOf(roomID);
      if (idx !== -1) {
        roomIDs.splice(idx, 1);
      }
      return Object.assign({}, state, {...state, roomIDs, rooms, currentRoom})
    },
    updateRoomInputMessages(state, {payload: {message, roomID}}) {
      let roomInputMessages = state.roomInputMessages;
      roomInputMessages[roomID] = message;
      return Object.assign({}, state, {...state, roomInputMessages});
    }
  },
  effects: {
    * login({payload: {username, password}}, {call, put}) {
      let {data} = yield call(requestToken, username, password);
      sessionStorage.setItem("token", data.token);
      yield put(routerRedux.push({pathname: '/'}));
    },
    * register({payload: {username, password}}, {call, put}) {
      let {data} = yield call(register, username, password);
      yield put(routerRedux.push({pathname: '/login'}));
    },
    * sendMessage({payload: {content, msgType}}, {call, put, select}) {
      let [current_user, current_room] = yield select(({users}) => [users.info, users.currentRoom]);
      let message = {
        content: content,
        msg_type: msgType,
        uuid: generateUUID(),
        user_id: current_user.id,
        room_id: current_room.id,
        from_user: current_user.id
      };
      yield call(sendMessage, message);
    },
    * getUserInfo({}, {call, put}) {
      let {data} = yield call(getUserInfo);
      yield put({type: 'updateInfo', payload: data})
    },
    * getRecentRooms({}, {put, call}) {
      let {data} = yield call(getRecentRooms);
      yield put({type: 'updateRecentRooms', payload: data.rooms});
    },
    * addRemoteMessage({payload: message}, {put, call, select}) {
      let state = yield select(({users}) => users);
      let room_id = message.room_id;
      if (!state.roomIDs.includes(room_id)) {
        // If room is not loaded, load the room info before adding message.
        let {data} = yield call(getRecentRooms);
        let rooms = data.rooms.filter((room) => room.id === room_id);
        yield put({type: 'updateRecentRooms', payload: rooms})
      }
      yield put({type: 'addMessage', payload: message})
    },
    * addFriend({payload: {friend_id}}, {put, call}) {
      let {data} = yield call(addFriend, friend_id);
      yield put({type: 'retrieveFriends'});
      yield put({type: 'updateRecentRooms', payload: [data.room]})
    },
    * retrieveFriends({}, {put, call}) {
      let {data} = yield call(retrieveFriends);
      yield put({type: 'updateFriends', payload: data.friends})
    },
    * createRoom({payload: userIds}, {put, call}) {
      let {data} = yield call(createRoom, userIds);
      yield put({type: 'updateRecentRooms', payload: [data['room']]});
      yield put({type: 'stickRoom', payload: data['room'].id})
    },
    * leaveRoom({payload: roomID}, {put, call}) {
      yield call(leaveRoom, roomID);
      yield put({type: 'removeRoom', payload: roomID});
      yield call(retrieveFriends);
    },
    * retrieveOfflineMessages({}, {call, put}) {
      let {data} = yield call(retrieveUnreaOfflinedMessages);
      let {messages} = data;
      for (let msg of messages) {
        yield put({type: 'addMessage', payload: msg})
      }
    },
    * ackMessages({payload: messageIds}, {call}) {
      yield call(ackMessages, messageIds);
    },
    * updateAvatar({payload: avatarHash}, {call, put}) {
      let {data} = yield call(updateAvatar, avatarHash);
      yield put({type: 'updateInfo', payload: data['user']})
    },
    * applyFriendship({payload: userID}, {call, put}) {
      yield call(applyFriendship, userID);
      yield put(retrieveOfflineMessages)
    },
    * ackReadNotifications({payload: uuidArray}, {call}) {
      yield call(ackReadNotifications, uuidArray)
    }
  },
  subscriptions: {
    setup({dispatch, history}) {
      history.listen(({pathname}) => {
        if (pathname === '/') {
          // init ws
          let ws;
          try {
            ws = initWebSocket();
            ws.addEventListener('message', function (event) {
              dispatch({
                type: 'addRemoteMessage',
                payload: JSON.parse(event.data)
              });
            });
            ws.addEventListener('close', function () {
              dispatch(routerRedux.push({pathname: '/login'}));
            });
          } catch (e) {
            if (e instanceof TokenMissingError) {
              dispatch(routerRedux.push({pathname: '/login'}));
              return
            }
          }
          // init user data
          dispatch({
            type: 'getUserInfo',
          });
          dispatch({
            type: 'getRecentRooms',
          });
          dispatch({
            type: 'retrieveFriends'
          });
          dispatch({
            type: 'retrieveOfflineMessages'
          });
        }
      });
    }
  },
};

/**
 * Set roomID as the first item in list.
 * @param roomIDs
 * @param roomID
 */
function bubbleRoom(roomIDs, roomID) {
  let idx = roomIDs.indexOf(roomID);
  if (idx !== -1) {
    roomIDs.splice(idx, 1);
    roomIDs.unshift(roomID);
  }
}

function clearUnreadMessage(unreadMessages, roomID) {
  unreadMessages[roomID] = [];
  return unreadMessages
}

