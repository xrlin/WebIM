import React from "react";
import {Route, Router} from "dva/router";
import IndexPage from "./routes/IndexPage";

import Login from "./routes/Login.js";

import Register from "./routes/Register.js";

function RouterConfig({history}) {
  return (
    <Router history={history}>
      <Route path="/" component={IndexPage} onEnter={checkLogin}/>
      <Route path="/login" component={Login}/>
      <Route path="/register" component={Register}/>
    </Router>
  );
}

/**
 * Check is user is login, redirect to login page if not.
 * @param {Object} nextState The object contains the path and other information.
 * @param {Function} transition The function that can used to replace ang redirect to another path.
 * @return undefined
 */
const checkLogin = (nextState, transition) => {
  if (sessionStorage.getItem('token') === null) {
    transition('/login');
  }
};

export default RouterConfig;
