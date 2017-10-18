import styles from './ChatMessage.less';
import {connect} from 'dva';
import UserAvatar from "../Users/UserAvatar";
import {messagetToHtml} from "../../utils/message";

function ChatMessage({message, current_user}) {
  return (
    <div className={`${styles.message} ${styles.clearfix} ${current_user.id === message.from_user ? styles.me : ''}`}>
      <UserAvatar user={current_user} className={styles['avatar']}/>
      <div className={styles.content}>
        <div className={styles.bubble}>
          <div className={styles.bubble_content}>
            <div className={styles.plain}>
              <pre dangerouslySetInnerHTML={{__html: messagetToHtml(message)}}/>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function mapStateToProps(state) {
  return {
    current_user: state.users.info
  }
}

export default connect(mapStateToProps)(ChatMessage);
