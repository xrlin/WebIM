export default {
  namespace: 'errors',
  state: {
    errors: []
  },
  reducers: {
    add(state, {payload: error_message}) {
      return Object.assign({}, state, {errors: state.errors.concat(error_message)});
    },
    clear(state) {
      return Object.assign({}, state, {errors: []});
    }
  },
  effects: {},
  subscriptions: {},
};
