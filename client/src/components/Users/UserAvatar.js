import styles from './UserAvatar.less';
import {connect} from 'dva';
import {Button, Modal, Slider} from 'antd';
import ReactDom from 'react-dom';
import AvatarEditor from 'react-avatar-editor'
import {uploadImage} from "../../utils/request";
import {dataURLtoBlob} from "../../utils/common";
import {getAvatarUrl} from "../../utils/url";
import enhanceWithClickOutside from "react-click-outside"

class UserAvatar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      modalVisible: false,
      profileX: 0,
      profileY: 0,
      addFriendBtnDisabled: !props.isFriend,
      addFriendBtnLoading: false,
      img: null,
      croppedImg: getAvatarUrl(props.user.avatar),
      cropperOpen: false
    };
    this.showModal = this.showModal.bind(this);
  }

  handleClickOutside = () => {
    this.setState({modalVisible: false});
  };

  componentWillUpdate(nextProps, nextState) {
    let croppedImg = getAvatarUrl(nextProps.user.avatar);
    if (nextState.croppedImg !== croppedImg) this.setState({croppedImg});
  }

  showModal({clientX, clientY}) {
    let profileX, profileY;
    profileX = clientX + 20;
    profileY = clientY + 20;
    if (profileX + 180 > window.innerWidth) {
      profileX = clientX - 180 - 20;
    }
    this.setState({modalVisible: true, profileX: profileX, profileY: profileY})
  };

  addFriend = () => {
    this.setState({addFriendBtnLoading: true});
    this.props.dispatch({
      type: 'users/applyFriendship',
      payload: this.props.user.id
    });
    this.setState({addFriendBtnLoading: false, addFriendBtnDisabled: true});
  };

  render() {
    let styleAttrs = {};
    let {width, height} = this.props;
    if (width) {
      styleAttrs['width'] = width;
    }
    if (height) {
      styleAttrs['height'] = height;
    }
    return (
      <div className={`${styles['avatar']} ${this.props.className}`}>
        <img src={this.state.croppedImg} style={styleAttrs}
             className={this.props.imgClassName} onClick={this.showModal}/>
        <div className={`${styles['profile_mini']} ${this.state.modalVisible ? styles['visible'] : ''}`}
             style={{top: this.state.profileY, left: this.state.profileX}}>
          <div className={styles['profile_mini__header']}>
            <img src={this.state.croppedImg}/>
            <Cropper dispatch={this.props.dispatch}/>
          </div>
          <div className={styles['profile_mini__body']}>
            <div className={styles['nickname_area']}>
              <h4>{this.props.user.name}</h4>
              <div style={{display: `${this.props.isFriend ? "none" : 'block'}`}}>
                <Button type="primary" icon="plus" size="small" onClick={this.addFriend}
                        loading={this.state.addFriendBtnLoading}
                        disabled={this.state.addFriendBtnDisabled}>添加好友</Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

class FileUpload extends React.Component {
  handleFile = (e) => {
    let reader = new FileReader();
    let file = e.target.files[0];

    if (!file) return;

    reader.onload = function (img) {
      ReactDom.findDOMNode(this.refs.in).value = '';
      this.props.handleFileChange(img.target.result);
    }.bind(this);
    reader.readAsDataURL(file);
  };

  render() {
    return (
      <input ref="in" type="file" accept="image/*" onChange={this.handleFile}/>
    );
  }
}


class Cropper extends React.Component {
  state = {
    visible: false,
    image: null,
    croppedImage: null,
    confirmLoading: false,
    scale: 1
  };

  handleFileChange = (dataURI) => {
    this.setState({
      image: dataURI,
      visible: true
    });
  };

  handleCancel = () => {
    this.setState({visible: false, scale: 1})
  };

  submit = async () => {
    const {dispatch} = this.props;
    let dataURI = this.refs.avatarEditor.getImageScaledToCanvas().toDataURL();
    let fileBlob = dataURLtoBlob(dataURI);
    this.setState({confirmLoading: true});
    let {data} = await uploadImage(fileBlob);
    dispatch({
      type: 'users/updateAvatar',
      payload: data['hash']
    });
    this.setState({visible: false, confirmLoading: false, scale: 1})
  };

  scaleChange = (scale) => {
    this.setState({scale});
  };

  render() {
    const {visible, confirmLoading, image, scale} = this.state;
    return (
      <div>
        <FileUpload handleFileChange={this.handleFileChange}/>
        <Modal style={{top: 20}} title="头像编辑" visible={visible} confirmLoading={confirmLoading}
               onCancel={this.handleCancel} onOk={this.submit}>
          <div className={styles['cropper']}>
            <AvatarEditor
              ref="avatarEditor"
              image={image}
              width={250}
              height={250}
              scale={scale}
              border={0}
              rotate={0}
            />
            <Slider min={1} max={2} defaultValue={1} step={0.01} value={scale} onChange={this.scaleChange}/>
          </div>
        </Modal>
      </div>
    )
  }
}

function checkIsFriend(friends, user, current_user) {
  if (current_user.id === user.id) {
    return true
  }
  for (let friend of friends) {
    if (user.id === friend.id) return true;
  }
  return false
}

function mapStateToProps({users}, ownProps) {
  let friends = users.friends;
  let isFriend = checkIsFriend(friends, ownProps.user, users.info);
  return {isFriend, friends, ...ownProps}
}

export default connect(mapStateToProps)(enhanceWithClickOutside(UserAvatar));
