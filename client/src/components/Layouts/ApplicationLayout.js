import React from "react";
import {connect} from 'dva';
import {Layout, message} from "antd";
import styles from './ApplicationLayout.css'

const {Content, Footer} = Layout;

class ApplicationLayout extends React.Component {
  constructor(props) {
    super(props)
  }

  componentDidUpdate() {
    let errors = this.props.errors;
    if (errors.length === 0) return;
    errors.forEach(error => {
      message.error(error)
    });
    this.props.dispatch({
      type: 'errors/clear'
    })
  }

  render() {
    return (
      <Layout className={`${styles['layout']} ${styles['layout-background']}`}>
        <Content style={{padding: '0 50px'}}>
          { this.props.children }
        </Content>
        <Footer style={{textAlign: 'center'}}>
          Ant Design ©2016 Created by Ant UED
        </Footer>
      </Layout>
    );
  }
}

function mapStateToProps({errors}) {
  return {
    errors: errors.errors
  };
}

export default connect(mapStateToProps)(ApplicationLayout);
