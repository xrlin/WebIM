import styles from './UserAvatar.less';
import {connect} from 'dva';
import {Button} from 'antd';
import ReactDom from 'react-dom';
import AvatarCropper from "react-avatar-cropper";
import {uploadImage} from "../../utils/request";
import {dataURLtoBlob} from "../../utils/common";
import {getAvatarUrl} from "../../utils/url";

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
    this.hideModal = this.hideModal.bind(this);
  }

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
      type: 'users/addFriend',
      payload: {friend_id: this.props.user.id}
    });
    this.setState({addFriendBtnLoading: false, addFriendBtnDisabled: true});
  };

  hideModal() {
    this.setState({modalVisible: false})
  };

  handleRequestHide = () => {
    this.setState({cropperOpen: false});
  };

  handleCrop = async (dataURI, fileName) => {
    const {dispatch} = this.props;
    let fileBlob = dataURLtoBlob(dataURI);
    let {data} = await uploadImage(fileBlob, fileName);
    dispatch({
      type: 'users/updateAvatar',
      payload: data['hash']
    });
    this.setState({cropperOpen: false});
  };

  handleFileChange = (dataURI) => {
    this.setState({
      img: dataURI,
      croppedImg: this.state.croppedImg,
      cropperOpen: true
    });
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
            <FileUpload handleFileChange={this.handleFileChange}/>
            {this.state.cropperOpen && <AvatarCropper
              onRequestHide={this.handleRequestHide}
              onCrop={this.handleCrop}
              cropperOpen={this.state.cropperOpen}
              image={this.state.img}
              width={400}
              height={400}
            />}
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
        <div className={`${styles['mask']} ${this.state.modalVisible ? styles['visible'] : ''}`}
             onClick={this.hideModal}/>
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
      this.props.handleFileChange(img.target.result, img.target.fileName);
    }.bind(this);
    reader.readAsDataURL(file);
  };

  render() {
    return (
      <input ref="in" type="file" accept="image/*" onChange={this.handleFile}/>
    );
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

export default connect(mapStateToProps)(UserAvatar);
