import dva from 'dva';
import {routerRedux} from "dva/router";
import './index.css';

// 1. Initialize
const app = dva({
  onError(e, dispatch) {
    if (e.response) {
      e.response.json().then((data) => {
        if (e.response.status === 401) {
          dispatch(routerRedux.push({pathname: '/login'}));
          return
        }
        dispatch({
          type: 'errors/add',
          payload: data.errors
        });
      });
    } else {
      dispatch({
        type: 'errors/add',
        payload: e
      });
    }

  },
});

app.model(require("./models/users"));

app.model(require("./models/chat_modal"));

app.model(require("./models/search_modal"));

app.model(require("./models/errors"));

// 2. Plugins
// app.use({});

// 3. Model
// app.model(require('./models/example'));

// 4. Router
app.router(require('./router'));

// 5. Start
app.start('#root');
