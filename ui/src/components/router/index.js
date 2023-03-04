import { Link } from "react-router-dom";

import * as React from "react";
import Box from "@mui/material/Box";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import { nanoid } from "nanoid";
// import Divider from "@mui/material/Divider";
// import InboxIcon from "@mui/icons-material/Inbox";s
// import DraftsIcon from "@mui/icons-material/Drafts";

export default function BasicList(props) {
  let paths = props.paths;

  return (
    <Box sx={{ width: "100%", maxWidth: 200, bgcolor: "background.paper" }}>
      <nav aria-label="main mailbox folders">
        <List>
          {paths.map((i) => {
            return (
              <Link
                style={{ textDecoration: "none" }}
                to={i.path}
                key={nanoid()}
              >
                <ListItem disablePadding>
                  <ListItemButton>
                    <ListItemIcon></ListItemIcon>

                    <ListItemText primary={i.text} />
                  </ListItemButton>
                </ListItem>
              </Link>
            );
          })}
        </List>
      </nav>
    </Box>
  );
}
