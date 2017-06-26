import {register, requestToken} from "../services/users";
import {routerRedux} from "dva/router";

export default {
  namespace: 'users',
  state: {},
  reducers: {},
  effects: {
    *login({payload: {username, password}}, {call, put}) {
      let {data} = yield call(requestToken, username, password);
      sessionStorage.setItem("token", data.token);
      yield put(routerRedux.push({pathname: '/'}));
    },
    *register({payload: {username, password}}, {call, put}) {
      let {data} = yield call(register, username, password);
      yield put(routerRedux.push({pathname: '/login'}));
    }
  },
  subscriptions: {},
};
