import React from 'react';
import { Router } from 'react-router-dom';
import {Route, Switch, Redirect} from "react-router-dom"
import './App.css';
import {Layout} from 'antd';
import Sidebar from './Components/Sidebar/Sidebar'
import Header from './Components/Header/Header'
import HomePage from "./Pages/HomePage";
import SettingsPage from "./Pages/SettingsPage";
import {createBrowserHistory} from 'history'
import UserListPage from "./Pages/UserListPage";
import TaskListPage from "./Pages/TaskListPage";
import TeamSettingsPage from "./Pages/TeamSettingsPage";
import WeekendSettingsPage from "./Pages/WeekendSettingsPage";
import UserAddPage from "./Pages/UserAddPage";
import UserEditPage from "./Pages/UserEditPage";
import TaskAddPage from "./Pages/TaskAddPage";
import TaskEditPage from "./Pages/TaskEditPage";
import TeamAddPage from "./Pages/TeamAddPage";
import TaskHistoryPage from "./Pages/TaskHistoryPage";

const { Content, Footer } = Layout;
// создаём кастомную историю
const history = createBrowserHistory()

class App extends React.Component {
    state = {
        collapsed: false,
    };

    toggle = () => {
        this.setState({
            collapsed: !this.state.collapsed,
        });
    };

    render() {
        return (
            <Router history={history}>
            <Layout className={"global-layout"} style={{ marginLeft: 200 }}>
                <Sidebar />
                <Layout className="site-layout">
                    <Header />
                    <Content
                        className=""
                        style={{
                            margin: '24px 16px',
                            minHeight: 280,
                            overflow: 'initial'
                        }}
                    >
                        <div className="site-layout-background">
                        <Switch>
                            <Route history={history} path='/team/add' component={TeamAddPage} />
                            <Route history={history} path='/team/:teamId/settings' component={TeamSettingsPage} />
                            <Route history={history} path='/team/weekend' component={WeekendSettingsPage} />
                            <Route history={history} path='/settings' component={SettingsPage} />

                            <Route history={history} path='/team/:teamId/users/add' component={UserAddPage} />
                            <Route history={history} path='/team/:teamId/users/edit/:userId' component={UserEditPage} />
                            <Route history={history} path='/team/:teamId/users' component={UserListPage} />

                            <Route history={history} path='/team/:teamId/tasks/add' component={TaskAddPage} />
                            <Route history={history} path='/team/:teamId/tasks/edit/:taskId' component={TaskEditPage} />
                            <Route history={history} path='/team/:teamId/tasks/:taskId/history' component={TaskHistoryPage} />
                            <Route history={history} path='/team/:teamId/tasks' component={TaskListPage} />
                            <Route history={history} exact path='/' component={HomePage} />
                            <Redirect to='/' />
                        </Switch>
                        </div>
                    </Content>
                    <Footer style={{ textAlign: 'center' }}>PolkovnikBot ©2020</Footer>
                </Layout>
            </Layout>
            </Router>
        );
    }
}

export default App;
