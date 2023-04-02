import { Link } from "react-router-dom";
import axios from "axios";
import { config } from "../../config";

import * as React from "react";
import Box from "@mui/material/Box";
import { Tabs, TabPane, Card, Table } from "@douyinfe/semi-ui";

import { nanoid } from "nanoid";
// import Divider from "@mui/material/Divider";
// import InboxIcon from "@mui/icons-material/Inbox";s
// import DraftsIcon from "@mui/icons-material/Drafts";

const statusText = {
  0: "已暂停",
  1: "爬取中",
  2: "已完成",
};

function a11yProps(index) {
  return {
    id: `vertical-tab-${index}`,
    "aria-controls": `vertical-tabpanel-${index}`,
  };
}

export default function Task() {
  const [finished, setFinished] = React.useState([]);
  const [unfinished, setUnFinished] = React.useState([]);

  const fetchTask = () => {
    axios.get(config.getPrefix() + `/api/tasks`).then((res) => {
      let tasks = [];

      tasks = res.data.map((v) => JSON.parse(v));

      const finished = tasks.filter((v) => v.status === 2);
      setFinished(finished);

      const unfinished = tasks.filter((v) => v.status !== 2);

      setUnFinished(unfinished);
    });
  };

  const getTasks = () => {
    const get = setInterval(() => {
      console.log("get tasks");

      fetchTask();
    }, 1000);

    return get;
  };

  React.useEffect(() => {
    const get = getTasks();

    return () => {
      clearInterval(get);
    };
  }, []);

  return (
    <div style={{ width: "100%", height: "100%" }} s>
      <Card style={{ width: "100%", height: "100%" }}>
        <Tabs type="line">
          <TabPane tab="未完成" itemKey="1">
            <Table dataSource={unfinished}>
              <Table.Column dataIndex="name" key="name" title="任务名"></Table.Column>
              <Table.Column dataIndex="status" key="status" title="状态" render={(v) => statusText[v]}></Table.Column>
              <Table.Column
                dataIndex="currentIndex"
                key="currentIndex"
                title="当前进度"
                render={(v) => "第" + (v + 1) + "章"}
              ></Table.Column>
            </Table>
          </TabPane>
          <TabPane tab="已完成" itemKey="2">
            <Table dataSource={finished}>
              <Table.Column dataIndex="name" key="name" title="任务名"></Table.Column>
              <Table.Column dataIndex="status" key="status" title="状态" render={(v) => statusText[v]}></Table.Column>
              <Table.Column
                dataIndex="currentIndex"
                key="currentIndex"
                title="当前进度"
                render={(v) => "第" + (v + 1) + "章"}
              ></Table.Column>
            </Table>
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
}
