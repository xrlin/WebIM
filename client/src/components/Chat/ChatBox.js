/**
 * Created by xr_li on 2017/8/7.
 */

import React from 'react';
import style from './ChatBox.less';
import {Scrollbars} from 'react-custom-scrollbars';
import {Icon, Input} from 'antd';
import {connect} from 'dva';
import ChatMessage from "./ChatMessage";
import UserAvatar from "../Users/UserAvatar";
import EditArea from './EditArea';

const {TextArea} = Input;

function RoomUserGallery({room, visible}) {
  let users = [];
  if (room && room.users) {
    users = room.users;
  }

  return (
    <div className={`${style['members']} ${visible ? style['visible'] : ''}`}>
      <div className={style['wrapper']}>
        <Scrollbars
          autoHeight
          autoHeightMax={200}
          autoHideTimeout={1} autoHide={true} hideTracksWhenNotNeeded={true}
          renderThumbVertical={props => <div {...props} className={style['thumb-vertical']}/>}>
          {users.map((user, i) => {
            return (
              <div key={i} className={style['member-info']}>
                <UserAvatar className={style['avatar']} user={user}/>
                <p className={style['nickname']}>{user.name}</p>
              </div>
            )
          })}
        </Scrollbars>
      </div>
    </div>
  )
}

class MessageBox extends React.Component {
  constructor(props) {
    super(props);
  }

  shouldComponentUpdate(nextProps) {
    return nextProps.messages !== this.props.messages
  }

  componentDidUpdate() {
    this.refs.messageBox.scrollToBottom();
  }

  render() {
    let messageItems = [];
    for (let i = 0; i < this.props.messages.length; i++) {
      let msg = this.props.messages[i];
      messageItems.push(
        <ChatMessage key={i} message={msg}/>
      );
    }
    return (
      <div className={style["messages-box"]}>
        <Scrollbars
          ref="messageBox"
          autoHideTimeout={1} autoHide={true} hideTracksWhenNotNeeded={true}
          renderThumbVertical={props => <div {...props} className={style['thumb-vertical']}/>}>
          {messageItems}
        </Scrollbars>
      </div>
    )
  }
}

class ChatBox extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      visible: false
    };
  }

  componentWillReceiveProps() {
    this.setState({visible: false});
  }

  send = (message) => {
    this.props.dispatch({
      type: 'users/sendMessage',
      payload: {content: message, msgType: 2}
    });
    this.props.dispatch({
      type: 'users/updateRoomInputMessages',
      payload: {message: '', roomID: this.props.room.id}
    });
    let mask = document.getElementById('face-selector__mask');
    mask && mask.click();
  };

  updateInputMessage = (event) => {
    this.props.dispatch({
      type: 'users/updateRoomInputMessages',
      payload: {message: event.target.value, roomID: this.props.room.id}
    });
  };

  toggleGallery = () => {
    this.setState({visible: !this.state.visible})
  };

  render() {
    return (
      <div className={style["chat-box"]}>
        <div className={style['title-wrapper']}>
          <div className={style['info']}>
            {this.props.room.name ? (<div className={style["name-wrapper"]} onClick={this.toggleGallery}><span
              className={style['title']}>{this.props.room.name}</span><Icon type={this.state.visible ? 'down' : 'up'}/>
            </div>) : ''}
          </div>
          <RoomUserGallery room={this.props.room} visible={this.state.visible}/>
          <div className={`${style['mask']} ${this.state.visible ? style['visible'] : ''}`}
               onClick={this.toggleGallery}/>
        </div>
        <MessageBox messages={this.props.messages}/>
        <div className={style['chat-input-area']}>
          <EditArea html={this.props.inputMessage} onChange={this.updateInputMessage} onEnter={this.send}/>
        </div>
      </div>
    )
  }
}

function mapStateToProps(state, {props}) {
  return {
    inputMessage: state.users.roomInputMessages[state.users.currentRoom.id] || '',
    current_user: state.users.info,
    room: state.users.currentRoom,
    messages: (state.users.messages && state.users.currentRoom && state.users.messages[state.users.currentRoom.id]) || [],
    ...props
  };
}

export default connect(mapStateToProps)(ChatBox);
