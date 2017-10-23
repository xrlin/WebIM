import React from 'react';
import {Button, Icon, Input, message, Modal} from 'antd';
import ReactDOMServer from 'react-dom/server';
import ReactDom from 'react-dom';
import styles from './EditArea.less';
import {htmlToMessage} from "../../utils/message";
import {Scrollbars} from "react-custom-scrollbars";
import {uploadImage} from "../../utils/request";
import {dataURLtoBlob} from "../../utils/common";
import {connect} from 'dva';
import {getMusicIdFromLink} from "../../utils/url";

function Face({faceID, faceText, onClick}) {
  let clickHandler = onClick || (() => {
  });
  return (
    <img className={styles['face']} key={`${faceText}_${faceID}`} data-faceID={faceID} data-faceText={faceText}
         src={`/${faceID}.gif`} onClick={clickHandler}/>
  )
}

class InputMenu extends React.Component {
  constructor(props) {
    super(props);
    this.state = {faceSelectorVisible: false};
  }

  clickHandler = (e) => {
    let {faceid, facetext} = e.target.dataset;
    let faceString = ReactDOMServer.renderToString(<Face faceID={faceid} faceText={facetext}/>);
    let editArea = document.getElementById('message-input');
    editArea.focus();
    let selection = window.getSelection(),
      range = selection.getRangeAt(0);
    range.deleteContents();
    let documentFragment = range.createContextualFragment(faceString);
    range.insertNode(documentFragment);
    range.collapse(false);
    selection.removeAllRanges();
    selection.addRange(range);
    // Trigger input event to update inputMessage
    let event = new Event('input', {
      'bubbles': true,
      'cancelable': true,
      'simulated': true
    });
    editArea.dispatchEvent(event);
  };

  showFaceSelector = () => {
    this.setState({faceSelectorVisible: true});
  };

  hideFaceSelector = () => {
    this.setState({faceSelectorVisible: false});
  };

  render() {
    let emojs = [];
    for (let i = 1; i < 92; i++) {
      emojs.push(<Face key={i} faceID={i} faceText={'test'} onClick={this.clickHandler}/>)
    }
    return (
      <div>
        <div className={styles['emoj-menu']}>
          <Icon type="smile-o" onClick={this.showFaceSelector}/>
          <div className={`${styles['emoj-selector']} ${this.state.faceSelectorVisible ? styles['visible'] : ''}`}>
            {emojs}
          </div>
        </div>
        <FileInputMenu dispatch={this.props.dispatch}/>
        <MusicMenu dispatch={this.props.dispatch}/>
        <div id="face-selector__mask"
             className={`${styles['mask']} ${this.state.faceSelectorVisible ? styles['visible'] : ''}`}
             onClick={this.hideFaceSelector}/>
      </div>
    )
  }
}

class FileInputMenu extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      file: null
    }
  }

  handleFile = (e) => {
    let reader = new FileReader();
    let file = e.target.files[0];

    if (!file) return;

    reader.onload = async function (img) {
      ReactDom.findDOMNode(this.refs.in).value = '';
      let fileBlob = dataURLtoBlob(img.target.result);
      let {data} = await uploadImage(fileBlob, file.fileName);
      this.props.dispatch({
        type: 'users/sendMessage',
        payload: {content: data['hash'], msgType: 4}
      });
    }.bind(this);
    reader.readAsDataURL(file);
  };

  render() {
    return (
      <div className={styles["file-menu__wrapper"]}>
        <Icon type="folder"/>
        <input ref="in" type="file" accept="image/*" onChange={this.handleFile}/>
      </div>
    )
  }
}

class MusicMenu extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      visible: false,
      musicLink: ''
    }
  }

  showModal = () => {
    this.setState({visible: true});
  };

  handleCancel = () => {
    this.setState({visible: false, musicLink: ''});
  };

  handleOk = () => {
    let link = this.state.musicLink;
    let result = getMusicIdFromLink(link);
    if (!result) {
      message.error('不能解析音乐信息，请核对输入！');
      return
    }
    this.props.dispatch({
      type: 'users/sendMessage',
      payload: {content: result, msgType: 6}
    });
    this.setState({visible: false, musicLink: ''})
  };

  updateMusicLink = (event) => {
    this.setState({musicLink: event.target.value})
  };

  render() {
    return (
      <div>
        <Button type="primary" icon="share-alt" size="small" onClick={this.showModal}/>
        <Modal title="Title"
               visible={this.state.visible}
               onOk={this.handleOk}
               onCancel={this.handleCancel}
        >
          <Input value={this.state.musicLink} onChange={this.updateMusicLink} placeholder="请输入音乐链接，暂时只支持网易云音乐"/>
        </Modal>
      </div>
    )
  }
}
class EditArea extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      lastEditRange: null
    };
    this.emitChange = this.emitChange.bind(this);
  }

  render() {
    let {html, id, ...props} = this.props;
    return (
      <div className={styles['edit-wrapper']}>
        <InputMenu dispatch={this.props.dispatch}/>
        <div className={styles['pre-wrapper']}>
          <Scrollbars
            autoHideTimeout={1} autoHide={true} hideTracksWhenNotNeeded={true}
            renderThumbVertical={props => <div {...props} className={styles['thumb-vertical']}/>}
          >
            <pre {...props} id={id || 'message-input'} ref={(e) => this.htmlEl = e} onInput={this.emitChange}
                 onKeyPress={this.keyDownHandler}
                 onBlur={this.props.onBlur || this.emitChange}
                 onFocus={this.props.onFocus || this.emitChange}
                 contentEditable={!this.props.disabled}
                 dangerouslySetInnerHTML={{__html: html}}
            >
              {this.props.children}
            </pre>
          </Scrollbars>
        </div>
        <div className={styles['edit__footer']}>
          <div className={styles['pull-right']}>
            <span>按下Ctrl+Enter换行</span>
            <Button onClick={this.sendButtonHandler}>发送</Button>
          </div>
        </div>
      </div>
    );
  }

  componentDidMount() {
    document.addEventListener('keydown', this.keyDownHandler);
  }

  componentWillUnmount() {
    document.removeEventListener('keydown', this.keyDownHandler);
  }

  keyDownHandler = (e) => {
    if (e.keyCode === 13 && e.ctrlKey) {
      e.preventDefault();
      let selection = window.getSelection(),
        range = selection.getRangeAt(0),
        br = document.createElement("br");
      range.deleteContents();
      range.insertNode(br);
      if (this.htmlEl.innerHTML.search(/.*?<br><br>$/) === -1) {
        range.insertNode(br.cloneNode(true));
      }
      range.collapse(false);

      selection.removeAllRanges();
      selection.addRange(range);
      return
    }
    if (e.keyCode === 13 && this.props.onEnter) {
      this.props.onEnter(htmlToMessage(this.props.html));
      e.preventDefault();
    }
  };

  sendButtonHandler = () => {
    this.props.onEnter(htmlToMessage(this.props.html));
  };

  initRange = () => {
    let range = document.createRange();
    if (!this.htmlEl) return;
    range.setStart(this.htmlEl, this.htmlEl.childElementCount);
    range.setEnd(this.htmlEl, this.htmlEl.childElementCount);
    return range
  };

  getSelectionRange = () => {
    if (!this.htmlEl) return;
    let selection = window.getSelection();
    console.log('Get range:', selection.getRangeAt(0));
    this.setState({lastEditRange: selection.getRangeAt(0)})
  };

  updateSelectionRangeWithLast = () => {
    if (!this.htmlEl) return;
    let range = this.state.lastEditRange;
    if (!range) return;
    this.htmlEl.focus();
    let selection = window.getSelection();
    selection.removeAllRanges();
    console.log('Set range:', range);
    selection.addRange(range);
  };

  shouldComponentUpdate(nextProps) {
    // We need not rerender if the change of props simply reflects the user's
    // edits. Rerendering in this case would make the cursor/caret jump.
    return (
      // Rerender if there is no element yet... (somehow?)
      !this.htmlEl
      // ...or if html really changed... (programmatically, not by user edit)
      || ( nextProps.html !== this.htmlEl.innerHTML
        && nextProps.html !== this.props.html )
      // ...or if editing is enabled or disabled.
      || this.props.disabled !== nextProps.disabled
      // ...or if className changed
      || this.props.className !== nextProps.className
    );
  }

  componentDidUpdate() {
    let selection = window.getSelection();
    selection.removeAllRanges();
    let r = document.createRange();
    let p = this.htmlEl;
    let lastNode = p.lastChild;
    if (lastNode) {
      r.setStartAfter(p.lastChild);
      r.setEndAfter(p.lastChild);
    }
    selection.addRange(r);
    if (this.htmlEl && this.props.html !== this.htmlEl.innerHTML) {
      // Perhaps React (whose VDOM gets outdated because we often prevent
      // rerendering) did not update the DOM. So we update it manually now.
      this.htmlEl.innerHTML = this.props.html;
      // this.updateSelectionRangeWithLast();
    }
  }

  emitChange(evt) {
    if (!this.htmlEl) return;
    let html = this.htmlEl.innerHTML;
    if (this.props.onChange && html !== this.lastHtml) {
      evt.target = {value: html};
      this.props.onChange(evt);
    }
    this.lastHtml = html;
  }

}

function mapStateToProps({}, ownProps) {
  return {...ownProps}
}

export default connect(mapStateToProps)(EditArea);
