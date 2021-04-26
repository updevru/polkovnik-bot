import React from 'react';
import {Layout, Menu } from "antd";
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    TeamOutlined,
    SettingOutlined,
    UserAddOutlined
} from "@ant-design/icons";
import './Sidebar.scss'
import {Link} from "react-router-dom";

const {Sider} = Layout;

class Sidebar extends React.Component {
    state = {
        collapsed: false,
        menuSelected: "2"
    };


    handleClick = e => {
        if (e.key === "toggle") {
            this.setState({
                collapsed: !this.state.collapsed,
            });
        } else {
            this.setState({ menuSelected: e.key });
        }
    };

    render() {
        const { menuSelected } = this.state;
        return (
            <Sider trigger={null} collapsible collapsed={this.state.collapsed}>
                <div className="logo">PolkovnikBot</div>
                <Menu theme="dark" mode="inline" onClick={this.handleClick} selectedKeys={[menuSelected]}>
                    <Menu.Item key="1" icon={<UserAddOutlined />}>
                        <Link to={"/team/add"}>Добавить команду</Link>
                    </Menu.Item>
                    <Menu.Item key="2" icon={<TeamOutlined />}>
                        <Link to={"/"}>Команды</Link>
                    </Menu.Item>
                    <Menu.Item key="3" icon={<SettingOutlined />}>
                        <Link to={"/settings"}>Настройки</Link>
                    </Menu.Item>
                    <Menu.Item key="toggle" icon={React.createElement(this.state.collapsed ? MenuUnfoldOutlined : MenuFoldOutlined)}>
                        Свернуть
                    </Menu.Item>
                </Menu>
            </Sider>
        )
    }
}

export default Sidebar;