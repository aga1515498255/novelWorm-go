import * as React from "react";
import PropTypes from "prop-types";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
// import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import configContext from "../context/configContext";
import Typography from "@mui/material/Typography";

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

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.number.isRequired,
  value: PropTypes.number.isRequired,
};

function a11yProps(index) {
  return {
    id: `vertical-tab-${index}`,
    "aria-controls": `vertical-tabpanel-${index}`,
  };
}

export default function VerticalTabs() {
  const [num, setNum] = React.useState(0);
  const handleChange = (event, newValue) => {
    setNum(newValue);
  };

  return (
    <configContext.Consumer>
      {(value) => {
        return (
          <Box
            sx={{
              flexGrow: 1,
              bgcolor: "background.paper",
              display: "flex",
              height: "100%",
            }}
          >
            <Tabs
              orientation="vertical"
              variant="scrollable"
              value={num}
              onChange={handleChange}
              aria-label="Vertical tabs example"
              sx={{ borderRight: 1, borderColor: "divider" }}
            >
              <Tab label="如何配置" key={-1} {...a11yProps(0)} />
              {value.configs.map((c, index) => {
                return (
                  <Tab label={c.name} key={index} {...a11yProps(index + 1)} />
                );
              })}
            </Tabs>
            {value.configs.map((c, index) => {
              console.log("num:", num, "index:", index);
              return (
                <TabPanel value={num} key={index} index={index + 1}>
                  <Stack direction="row" spacing={2}>
                    <span>名称：</span>
                    <span>{c.name}</span>
                  </Stack>
                  <Stack direction="row" spacing={2}>
                    <span>地址：</span>
                    <span>{c.websetURl}</span>
                  </Stack>
                  <Stack direction="row" spacing={2}>
                    <span>章节选择器：</span>
                    {c.chapterSelector.map((s, i) => {
                      return <span key={i}>{s}</span>;
                    })}
                  </Stack>
                  <Stack direction="row" spacing={2}>
                    <span>小说内容选择器：</span>
                    {c.contentSelector.map((s, i) => {
                      return <span key={i}>{s}</span>;
                    })}
                  </Stack>
                </TabPanel>
              );
            })}
          </Box>
        );
      }}
    </configContext.Consumer>
  );
}
