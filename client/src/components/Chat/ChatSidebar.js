import React from 'react';
import style from './ChatSidebar.less';
import {Dropdown, Icon, Input, Menu} from 'antd';
import ChatItem from './ChatItem';
import {Scrollbars} from 'react-custom-scrollbars';
import {connect} from 'dva';
import SearchModal from "../Users/SearchModal";
import ChatModal from "./ChatModal";
import RoomSearch from './RoomSearch';
import UserAvatar from "../Users/UserAvatar";

const Search = Input.Search;

class ChatSidebar extends React.Component {
  constructor(props) {
    super(props);
  }

  displaySearchUsersModal = ({key}) => {
    switch (key) {
      case "0":
        this.props.dispatch({
          type: 'search_modal/showModal'
        });
        break;
      case "1":
        this.props.dispatch({
          type: 'chat_modal/showModal'
        });
        break;
    }

  };

  render() {
    let chatItems = [];
    for (let room of this.props.rooms) {
      chatItems.push(
        <ChatItem key={room.id} room={room}/>
      )
    }
    const menu = (
      <Menu onClick={this.displaySearchUsersModal}>
        <Menu.Item key="0">
          <span>添加朋友</span>
        </Menu.Item>
        <Menu.Item key="1">
          <span>发起聊天</span>
        </Menu.Item>
        <Menu.Divider/>
        <Menu.Item key="3">3d menu item</Menu.Item>
      </Menu>
    );
    return (
      <aside className={style['sidebar']}>
        <div className={style['header']}>
          <div className={style['avatar']}>
            <UserAvatar user={this.props.current_user} width="40px" height="40px"/>
          </div>
          <div className={style['info']}>
            <h3>
              <span className={style['nickname']}>{this.props.current_user && this.props.current_user.name}</span>
              <div className={style['opt-menu']}>
                <Dropdown overlay={menu} trigger={['click']}>
                  <a className="ant-dropdown-link" href="#">
                    <Icon type="bars"/>
                  </a>
                </Dropdown>
              </div>
            </h3>
          </div>
        </div>
        <div className={style['search-bar']}>
          <RoomSearch friendRooms={this.props.friendRooms} multiRooms={this.props.multiRooms}
                      dispatch={this.props.dispatch}/>
        </div>
        <div className={style["chat-list"]}>
          <Scrollbars
            autoHideTimeout={1} autoHide={true} hideTracksWhenNotNeeded={true}
            renderThumbVertical={props => <div {...props} className={style['thumb-vertical']}/>}>
            {chatItems}
          </Scrollbars>
        </div>
        <SearchModal/>
        <ChatModal/>
      </aside>
    )
  }
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

function mapStateToProps({users}) {
  let rooms = [];
  for (let id of users.roomIDs) {
    rooms.push(users.rooms.get(id));
  }
  return {
    current_user: users.info,
    rooms: rooms,
    friendRooms: selectRooms(rooms, FriendRoom),
    multiRooms: selectRooms(rooms, MultiRoom)
  }
}

export default connect(mapStateToProps)(ChatSidebar);
