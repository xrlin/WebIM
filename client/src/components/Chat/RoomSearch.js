import {Select} from 'antd';
import styles from './RoomSearch.less';

class RoomSearch extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      value: '',
      friendRooms: props.friendRooms,
      multiRooms: props.multiRooms
    };
    this.search = this.search.bind(this);
    this.selectHandler = this.selectHandler.bind(this);
  }

  children = (rooms) => {
    const Option = Select.Option;
    let options = [];
    for (let room of rooms.values()) {
      options.push(
        <Option key={room.id}>
          <div className={styles['avatar']}>
            <img src="https://xrlin.github.io/assets/img/crown-logo.png"/>
          </div>
          <div className={styles['nickname']}>
            <span>{room.name}</span>
          </div>
        </Option>
      );
    }
    return options;
  };

  search(value) {
    let friendRooms = new Map();
    for (let [key, room] of this.props.friendRooms) {
      if (room.name.includes(value)) {
        friendRooms.set(key, room);
      }
    }

    let multiRooms = new Map();
    for (let [key, room] of this.props.multiRooms) {
      if (room.name.includes(value)) {
        multiRooms.set(key, room);
      }
    }
    this.setState({
      value: value,
      friendRooms, multiRooms
    });
  }

  selectHandler(value) {
    this.props.dispatch({
      type: 'users/stickRoom',
      payload: parseInt(value)
    });
    this.setState({value: ''})
  }

  render() {
    const OptGroup = Select.OptGroup;
    return (
      <div className={styles['search']}>
        <Select
          value={this.state.value}
          style={{width: '100%'}} mode="combobox"
          onSearch={this.search}
          dropdownClassName={styles['dropdown']}
          notFoundContent=""
          showArrow={false}
          onSelect={this.selectHandler}
          filterOption={false}
          placeholder="搜索..."
        >
          <OptGroup label="好友">
            {this.children(this.state.friendRooms)}
          </OptGroup>
          <OptGroup label="群聊">
            {this.children(this.state.multiRooms)}
          </OptGroup>
        </Select>
      </div>
    )
  }
}

export default RoomSearch;
