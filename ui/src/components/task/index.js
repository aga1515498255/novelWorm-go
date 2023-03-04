import { Link } from "react-router-dom";

import * as React from "react";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Button from "@mui/material/Button";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Typography from "@mui/material/Typography";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import { nanoid } from "nanoid";
// import Divider from "@mui/material/Divider";
// import InboxIcon from "@mui/icons-material/Inbox";s
// import DraftsIcon from "@mui/icons-material/Drafts";

function a11yProps(index) {
  return {
    id: `vertical-tab-${index}`,
    "aria-controls": `vertical-tabpanel-${index}`,
  };
}

function TabPanel(props) {
  const { value, index, children } = props;
  // console.log("in render tabePanel:", value, index);

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`vertical-tabpanel-${index}`}
      aria-labelledby={`vertical-tab-${index}`}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

export default function Task() {
  return (
    <Box>
      <Tabs
        orientation="vertical"
        variant="scrollable"
        aria-label="task catagorise"
        sx={{ borderRight: 1, borderColor: "divider" }}
      >
        <Tab label="已完成" key={0} {...a11yProps(0)} />

        <Tab label="未完成" key={1} {...a11yProps(1)} />
      </Tabs>
      <TabPanel></TabPanel>
    </Box>
  );
}
