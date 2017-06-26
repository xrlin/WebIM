import React from 'react';
import {connect} from 'dva';
import LoginForm from '../components/Users/LoginForm'
import ApplicationLayout from '../components/Layouts/ApplicationLayout'

function Login({dispatch}) {
  return (
    <ApplicationLayout>
      <LoginForm dispatch={dispatch}/>
    </ApplicationLayout>
  );
}

function mapStateToProps() {
  return {};
}

export default connect(mapStateToProps)(Login);
