import React from "react";
import {connect} from "dva";
import ApplicationLayout from "../components/Layouts/ApplicationLayout";
import RegisterForm from "../components/Users/RegisterForm";

function Register({dispatch}) {
  return (
    <ApplicationLayout>
      <RegisterForm dispatch={dispatch}/>
    </ApplicationLayout>
  );
}

function mapStateToProps() {
  return {};
}

export default connect(mapStateToProps)(Register);
