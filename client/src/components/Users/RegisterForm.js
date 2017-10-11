import React from "react";
import {Button, Form, Icon, Input} from "antd";
import {Link} from 'dva/router';
import styles from './RegisterForm.less';

const FormItem = Form.Item;

class RegisterForm extends React.Component {
  constructor(props) {
    super(props);
  }

  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        let {username, password} = values;
        this.props.dispatch({
          type: 'users/register',
          payload: {username, password}
        })
      }
    })
  };

  render() {
    const {getFieldDecorator} = this.props.form;
    return (
      <Form onSubmit={this.handleSubmit} className={styles['register-form']}>
        <FormItem>
          {getFieldDecorator('username', {
            rules: [{required: true, message: 'Please input your username!'}],
          })(
            <Input prefix={<Icon type="user" style={{fontSize: 13}}/>} placeholder="Username"/>
          )}
        </FormItem>
        <FormItem>
          {getFieldDecorator('password', {
            rules: [{required: true, message: 'Please input your Password!'}],
          })(
            <Input prefix={<Icon type="lock" style={{fontSize: 13}}/>} type="password" placeholder="Password"/>
          )}
        </FormItem>
        <FormItem>
          <Button type="primary" htmlType="submit" className={styles['register-form-button']}>
            Register
          </Button>
          Or <Link to="/login">login now!</Link>
        </FormItem>
      </Form>
    )
  }
}
export default Form.create()(RegisterForm);
