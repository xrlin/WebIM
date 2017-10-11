import React from 'react';
import {connect} from 'dva';
import ChatRoom from '../components/Chat/ChatRoom';
import ApplicationLayout from '../components/Layouts/ApplicationLayout'

function IndexPage() {
  return (
    <ApplicationLayout>
      <ChatRoom/>
    </ApplicationLayout>
  );
}

IndexPage.propTypes = {};

export default connect()(IndexPage);
