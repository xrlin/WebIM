import React from 'react';
import styles from './ChatModal.less';
import {Button, Checkbox, Modal, Tabs} from 'antd';
import {connect} from 'dva';
import {Scrollbars} from "react-custom-scrollbars";

const TabPane = Tabs.TabPane;

function ChatModal({dispatch, visible, friendRooms, multiRooms, selectedUsers}) {
  const closeModal = () => {
    dispatch({
      type: 'chat_modal/closeModal'
    })
  };

  return (
    <Modal
      visible={visible}
      title="发起聊天"
      footer={null}
      closable={true}
      onCancel={closeModal}
    >
      <div className={styles['wrapper']}>
        <Tabs defaultActiveKey="1" size="small">
          <TabPane tab="好友" key="1"><FriendRoomTab rooms={friendRooms} dispatch={dispatch}
                                                   selectedUsers={selectedUsers}/></TabPane>
          <TabPane tab="群聊" key="2"><RoomTab rooms={multiRooms} dispatch={dispatch}/></TabPane>
        </Tabs>
      </div>
    </Modal>
  );
}

function userIndex(users, user) {
  for (let i = 0; i < users.length; i++) {
    let current = users[i];
    if (current.roomID === user.roomID) {
      return i
    }
  }
  return -1
}

class FriendRoomTab extends React.Component {
  constructor(props) {
    super(props);
  }

  checkedHandler = (event) => {
    let target = event.target;
    let item = target.value;
    let selectedItems = this.props.selectedUsers;
    if (target.checked) {
      selectedItems.push(item);
    } else {
      let idx = userIndex(selectedItems, item);
      if (idx === -1) return;
      selectedItems.splice(idx, 1);
    }
    this.props.dispatch({
      type: 'chat_modal/updateSelectedUsers',
      payload: selectedItems
    });
  };

  btnClickHandler = () => {
    if (this.props.selectedUsers.length === 1) {
      this.props.dispatch({
          type: 'users/stickRoom',
          payload: this.props.selectedUsers[0].roomID
        }
      );
    } else {
      let userIds = this.props.selectedUsers.map(function (info) {
        return info.userID
      });
      this.props.dispatch({
          type: 'users/createRoom',
          payload: userIds
        }
      );
    }
    this.props.dispatch({
      type: 'chat_modal/closeModal'
    });
  };

  generateList = (rooms, handler, selectedUsers) => {
    let items = [];
    rooms.forEach(function (room) {
      let user = room.users[0];
      if (!user) {
        return
      }
      items.push(
        <li key={room.id} className={styles['list-item']}>
          <div className={styles['operation']}>
            <Checkbox value={{roomID: room.id, userID: user.id}} onChange={handler}
                      checked={userIndex(selectedUsers, {roomID: room.id, userID: user.id}) !== -1}/>
          </div>
          <div className={styles['info']}>
            <div className={styles['avatar']}>
              <img src="https://xrlin.github.io/assets/img/crown-logo.png"/>
            </div>
            <div className={styles['nickname']}>
              <h4>{user.name}</h4>
            </div>
          </div>
        </li>
      )
    });
    return items
  };

  render() {
    return (
      <div className={styles['rooms-wrapper']}>
        <ul className={styles['list-wrapper']}>
          <Scrollbars
            autoHideTimeout={1} autoHide={true} hideTracksWhenNotNeeded={true}
            renderThumbVertical={props => <div {...props} className={styles['thumb-vertical']}/>}>
            {this.generateList(this.props.rooms, this.checkedHandler, this.props.selectedUsers)}
          </Scrollbars>
        </ul>
        <div className={styles['footer']}>
          <Button type="primary" onClick={this.btnClickHandler}
                  disabled={this.props.selectedUsers.length <= 0}>开始聊天</Button>
        </div>
      </div>
    )
  }
}

function RoomTab({dispatch, rooms}) {
  const clickHandler = (roomId) => {
    dispatch({
      type: 'users/stickRoom',
      payload: roomId
    });
    dispatch({
      type: 'chat_modal/closeModal'
    });
  };
  const generateList = () => {
    let items = [];
    rooms.forEach(function (room) {
      items.push(
        <li key={room.id} className={styles['list-item']} onClick={() => clickHandler(room.id)}>
          <div className={styles['info']}>
            <div className={styles['avatar']}>
              <img src="https://xrlin.github.io/assets/img/crown-logo.png"/>
            </div>
            <div className={styles['nickname']}>
              <h4>{room.name}</h4>
            </div>
          </div>
        </li>
      )
    });
    return items
  };
  return (
    <div className={styles['rooms-wrapper']}>
      <ul className={styles['list-wrapper']}>
        <Scrollbars
          autoHideTimeout={1} autoHide={true} hideTracksWhenNotNeeded={true}
          renderThumbVertical={props => <div {...props} className={styles['thumb-vertical']}/>}>
          {generateList()}
        </Scrollbars>
      </ul>
    </div>
  )
}

/**
 * Select the rooms by type
 * @param {Map} rooms
 * @param {int} roomType
 * @return {Map}
 */
function selectRooms(rooms, roomType) {
  let results = new Map();
  rooms.forEach(function (v, k) {
    if (v.room_type === roomType) {
      results.set(k, v);
    }
  });
  return results;
}

const MultiRoom = 0, FriendRoom = 1;

function mapStateToProps({chat_modal, users}) {
  return {
    visible: chat_modal.visible,
    selectedUsers: chat_modal.selectedUsers,
    friendRooms: selectRooms(users.rooms, FriendRoom),
    multiRooms: selectRooms(users.rooms, MultiRoom)
  }
}

export default connect(mapStateToProps)(ChatModal);
