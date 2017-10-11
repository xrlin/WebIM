import {searchUsers} from "../services/users";

export default {
  namespace: 'search_modal',
  state: {visible: false, users: []},
  reducers: {
    showModal(state) {
      return Object.assign({}, state, {...state, visible: true})
    },
    closeModal(state) {
      return Object.assign({}, state, {...state, visible: false})
    },
    updateUsers(state, {payload: users}) {
      return Object.assign({}, state, {...state, users});
    }
  },
  effects: {
    * searchUsers({payload: {name}}, {call, put}) {
      let {data} = yield call(searchUsers, name);
      yield put({type: 'updateUsers', payload: data.users})
    }
  },
  subscriptions: {},
};
