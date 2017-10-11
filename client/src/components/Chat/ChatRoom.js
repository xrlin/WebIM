import React from "react";
import styles from "./ChatRoom.less"
import ChatSidebar from "./ChatSidebar"
import ChatBox from "./ChatBox";
import {connect} from "dva";

class ChatRoom extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <div className={styles['chat-room']}>
        <ChatSidebar/>
        <ChatBox/>
      </div>
    );
  }
}

function mapStateToProps(state, {props}) {
  return {
    messages: state.users.messages,
    ...props
  };
}

export default connect(mapStateToProps)(ChatRoom);
