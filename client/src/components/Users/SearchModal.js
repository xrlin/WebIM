import {Button, Input, Modal} from 'antd';
import {connect} from 'dva';
import styles from './SearchModal.less';

class UserItem extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      disabled: this.buttonDisabled()
    }
  }

  clickHandler = () => {
    this.setState({loading: true});
    this.props.dispatch({
      type: 'users/applyFriendship',
      payload: this.props.user.id
    });
    this.setState({loading: false, disabled: true})
  };

  buttonHide = () => {
    return this.props.current_user.id === this.props.user.id
  };

  buttonDisabled = () => {
    if (this.buttonHide()) {
      return true
    }
    for (let user of this.props.friends) {
      if (user.id === this.props.user.id) {
        return true
      }
    }
    return false
  };

  render() {
    return (
      <li className={styles['user-list-item']}>
        <div className={styles['avatar']}>
          <img src="https://xrlin.github.io/assets/img/crown-logo.png"/>
        </div>
        <div className={styles['user-info']}>
          <div className={styles['nickname']}><h4>{this.props.user.name}</h4></div>
          <div className={styles['user-operation']}>{
            this.buttonHide() ||
            <Button type="primary" icon="plus" size="small" onClick={this.clickHandler} loading={this.state.loading}
                    disabled={this.state.disabled}>添加好友</Button>
          }</div>
        </div>
      </li>
    )
  }
}

class SearchModal extends React.Component {
  constructor(props) {
    super(props);
  }

  closeHandler = () => {
    this.props.dispatch({
      type: 'search_modal/closeModal'
    })
  };

  searchUsersByName = (e) => {
    let name = e.target.value;
    if (name === "") {
      return
    }
    this.props.dispatch({
      type: 'search_modal/searchUsers',
      payload: {name}
    })
  };

  render() {
    let userItems = [];
    console.log(this.props);
    for (let user of this.props.users) {
      userItems.push(<UserItem user={user} key={user.id} dispatch={this.props.dispatch}
                               current_user={this.props.current_user} friends={this.props.friends}/>)
    }
    return (
      <Modal
        visible={this.props.visible}
        title="搜索/添加用户"
        footer={null}
        closable={true}
        onCancel={this.closeHandler}
      >
        <Input placeholder="请输入用户名" onInput={this.searchUsersByName}/>
        <div className={styles['user-list']}>
          <ul>
            {userItems}
          </ul>
        </div>
      </Modal>
    )
  }
}

function mapStateToProps({users, search_modal}) {
  return {visible: search_modal.visible, users: search_modal.users, friends: users.friends, current_user: users.info}
}

export default connect(mapStateToProps)(SearchModal);
