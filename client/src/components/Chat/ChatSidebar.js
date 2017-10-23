import React from 'react';
import style from './ChatSidebar.less';
import {Button, Dropdown, Icon, Input, Menu, Tabs} from 'antd';
import ChatItem from './ChatItem';
import {Scrollbars} from 'react-custom-scrollbars';
import {connect} from 'dva';
import SearchModal from "../Users/SearchModal";
import ChatModal from "./ChatModal";
import RoomSearch from './RoomSearch';
import UserAvatar from "../Users/UserAvatar";
import {checkFriendshipApplication} from "../../services/users";

const Search = Input.Search;
const TabPane = Tabs.TabPane;

class Notification extends React.Component {
  checkHandler = async (uuid, action) => {
    let {data} = await checkFriendshipApplication(uuid, action);
    const {dispatch} = this.props;
    dispatch({
      type: 'users/getRecentRooms'
    })
  };

  render() {
    const {notification} = this.props;
    return (
      <div>
        {notification.source_user.name}申请添加你为好友
        <Button type="primary" onClick={() => this.checkHandler(notification.uuid, 'pass')}>通过</Button>
        <Button type="primary" onClick={() => this.checkHandler(notification.uuid, 'reject')}>不通过</Button>
      </div>
    )
  }
}

class ChatSidebar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hasUnreadNotifications: !!props.hasUnreadNotifications
    }
  }

  componentWillUpdate(nextProps) {
    if (this.state.hasUnreadNotifications !== nextProps.hasUnreadNotifications) {
      this.setState({hasUnreadNotifications: !!nextProps.hasUnreadNotifications})
    }
  }

  ackReadNotifications = (activeKey) => {
    if (activeKey !== "2") return;
    const {dispatch, unreadNotifications} = this.props;
    let uuidArray = [];
    unreadNotifications.forEach(notify => {
      uuidArray.push(notify['uuid'])
    });
    if (uuidArray.length === 0) return;
    dispatch({
      type: 'users/ackReadNotifications',
      payload: uuidArray
    });
    // this.setState({hasUnreadNotifications: false});
  };

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
    let {hasUnreadNotifications} = this.state;
    let notificationTabClass = hasUnreadNotifications ? style['tab--unread'] : '';
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
    const notificationItems = [];
    const {notifications} = this.props;
    notifications.forEach((notification, idx) => {
      notificationItems.push(
        <Notification key={idx} notification={notification} dispatch={this.props.dispatch}/>
      )
    });
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
        <Tabs defaultActiveKey="1" onTabClick={this.ackReadNotifications}>
          <TabPane tab={<span><Icon type="apple"/>Tab 1</span>} key="1">
            <div className={style["chat-list"]}>
              <Scrollbars
                autoHideTimeout={1} autoHide={true} hideTracksWhenNotNeeded={true}
                renderThumbVertical={props => <div {...props} className={style['thumb-vertical']}/>}>
                {chatItems}
              </Scrollbars>
            </div>
          </TabPane>
          <TabPane tab={<span className={notificationTabClass}><Icon type="android"/>通知</span>} key="2">
            {notificationItems}
          </TabPane>
        </Tabs>
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
  let unreadNotifications = [];
  users.notifications.forEach(notify => {
    if (!notify['read']) unreadNotifications.push(notify)
  });
  let hasUnreadNotifications = unreadNotifications.length > 0;
  return {
    current_user: users.info,
    rooms: rooms,
    friendRooms: selectRooms(rooms, FriendRoom),
    multiRooms: selectRooms(rooms, MultiRoom),
    notifications: users.notifications,
    hasUnreadNotifications,
    unreadNotifications
  }
}

export default connect(mapStateToProps)(ChatSidebar);

