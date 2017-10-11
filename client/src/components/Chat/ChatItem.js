/**
 * Created by xr_li on 2017/8/6.
 */

import React from 'react';
import style from './ChatItem.less';
import {connect} from 'dva';
import {ContextMenu, ContextMenuTrigger, MenuItem} from "react-contextmenu";

class ChatItem extends React.Component {
  constructor(props) {
    super(props)
  }

  setCurrentRoom = () => {
    this.props.dispatch({
      type: 'users/setCurrentRoom',
      payload: this.props.room['id']
    })
  };

  leave = () => {
    this.props.dispatch({
      type: 'users/leaveRoom',
      payload: this.props.room.id
    })
  };

  render() {
    return (
      <div>
        <ContextMenuTrigger id={`contextmenu__${this.props.room.id}`}>
          <div
            className={`${style['chat-item']} ${this.props.currentRoom.id === this.props.room.id ? style['active'] : ''}`}
            onClick={this.setCurrentRoom}>
            <div className={style['avatar']}>
              <img src="https://xrlin.github.io/assets/img/crown-logo.png" className={style['img']}/>
              <div
                className={`${style['message-notification']} ${this.props.newMessages.length > 0 ? style['visible'] : ''}`}>
                {this.props.newMessages.length}
              </div>
            </div>
            <div className={style['info']}>
              <span className={style['nickname']}>{this.props.room.name}</span>
            </div>
          </div>
        </ContextMenuTrigger>
        <ContextMenu id={`contextmenu__${this.props.room.id}`}>
          <MenuItem onClick={this.leave}>
            退出群聊
          </MenuItem>
          <MenuItem onClick={this.handleClick}>
            ContextMenu Item 2
          </MenuItem>
          <MenuItem divider/>
          <MenuItem onClick={this.handleClick}>
            ContextMenu Item 3
          </MenuItem>
        </ContextMenu>
      </div>
    )
  }
}

function mapStateToProps({users}, ownProps) {
  let currentRoom = users.currentRoom;
  let newMessages = users.newMessages[ownProps.room.id] || [];
  return {currentRoom, newMessages, ...ownProps}
}


export default connect(mapStateToProps)(ChatItem);
