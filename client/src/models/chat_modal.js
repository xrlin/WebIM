export default {
  namespace: 'chat_modal',
  state: {
    visible: false,
    selectedUsers: [],
  },
  reducers: {
    showModal(state) {
      return Object.assign({}, state, {...state, visible: true})
    },
    closeModal(state) {
      return Object.assign({}, state, {...state, visible: false, selectedUsers: []})
    },
    updateSelectedUsers(state, {payload: userInfo}) {
      return Object.assign({}, state, {...state, selectedUsers: userInfo})
    }
  },
  effects: {},
  subscriptions: {},
};
