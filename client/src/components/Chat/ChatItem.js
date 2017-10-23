/**
 * Created by xr_li on 2017/8/6.
 */

import React from 'react';
import style from './ChatItem.less';
import {connect} from 'dva';
import {ContextMenu, ContextMenuTrigger, MenuItem} from "react-contextmenu";
import {getAvatarUrl} from "../../utils/url";

class RoomAvatar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      avatars: null
    }
  }

  async componentDidMount() {
    let {room} = this.props;
    let imgList = [];
    // 非群聊
    room.users.forEach(user => imgList.push(user.avatar));
    this.combineAvatars(imgList);
  }

  async combineAvatars(imgList) {
    const {width, height, padding} = this.props;
    let len = imgList.length;
    let {cellWidth, cellHeight, cellPadding, totalCells} = this.calculateCells(width, height, padding, len);
    let avatars = [];
    let count = 0;
    for (let idxY = cellPadding; idxY < height;) {
      for (let idxX = cellPadding; idxX < width;) {
        let cellStyle = {
          width: cellWidth,
          height: cellHeight,
          position: 'absolute',
          left: idxX + 'px',
          top: idxY + 'px'
        };
        const img = <img key={count} src={getAvatarUrl(imgList[count])} style={cellStyle}/>;
        avatars.push(img);
        idxX = idxX + cellPadding + cellWidth;
        count += 1;
        if (count === totalCells) break;
      }
      idxY = idxY + cellPadding + cellHeight;
    }
    this.setState({avatars})
  }

  calculateCells = (width, height, padding, total) => {
    let cellWidth, cellHeight, cellPadding, rows, totalCells;
    cellPadding = padding;
    if (total === 1) {
      cellPadding = 0
    }
    if (total / 2 <= 1) {
      rows = 1;
    } else if (total / 3 <= 1) {
      rows = 2;
    } else {
      rows = 3;
    }
    cellHeight = (height - (rows + 1) * cellPadding) / rows;
    cellWidth = (width - (rows + 1) * cellPadding) / rows;
    totalCells = total - rows * rows >= 0 ? rows * rows : total;
    return {cellWidth, cellHeight, rows, cellPadding, totalCells}
  };

  render() {
    let {width, height, background} = this.props;
    let wrapperStyle = {
      width: width + 'px',
      height: height + 'px',
      position: 'relative',
      background: (background || '#eeeeee')
    };
    return (
      <div style={wrapperStyle}>
        {this.state.avatars}
      </div>
    )
  }
}

class ChatItem extends React.Component {
  constructor(props) {
    super(props)
  }

  componentDidUpdate(prevProps) {
    let active = this.props.currentRoom.id === this.props.room.id;
    let {dispatch} = this.props;
    let prevNewMessages = prevProps.newMessages;
    if (active && prevNewMessages.length > 0) {
      let messageIds = [];
      for (let msg of prevNewMessages) {
        if (!msg.id) continue;
        messageIds.push(msg.id);
      }
      dispatch({
        type: 'users/ackMessages',
        payload: messageIds
      })
    }
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
              <RoomAvatar className={style['img']} room={this.props.room} width={40} height={40} padding={2}/>
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
