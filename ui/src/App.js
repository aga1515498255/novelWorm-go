import "./app.css";
import Router from "./components/router";
import { Route, Routes, Navigate } from "react-router-dom";
import Config from "./components/config";
import Request from "./components/request";
import React, { Component } from "react";
import configContext from "./components/context/configContext";
import Task from "./components/task";
// import axois from "axios";
// import { config } from "./config.js";

export default class App extends Component {
  state = {
    path: [
      { path: "/request", text: "爬取小说" },
      { path: "/config", text: "设置" },
      { path: "/task", text: "任务" },
    ],
  };

  render() {
    let configs = this.props.configs;

    return (
      <div className="App">
        <Router paths={this.state.path}></Router>
        <div className="container">
          <configContext.Provider value={{ configs: configs }}>
            <Routes>
              <Route path="/config" element={<Config />}></Route>
              <Route path="/request" element={<Request />}></Route>

              <Route
                path="/"
                element={<Navigate to="/request" replace={true} />}
              />
              <Route path="/task" element={<Task />} />
            </Routes>
          </configContext.Provider>
        </div>
      </div>
    );
  }
}
