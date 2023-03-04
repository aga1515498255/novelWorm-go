import React, { Component } from "react";
import TextField from "@mui/material/TextField";
import Style from "./request.module.css";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import { config } from "../../config.js";

import axois from "axios";
// import { resolveBreakpointValues } from "@mui/system/breakpoints";

export default class Header extends Component {
  state = { previewText: "这是章节预览区" };

  preview = () => {
    let url = document.getElementById("outlined-basic").value;
    console.log(config.getPrefix());
    axois.get(config.getPrefix() + `/api/preview?url=${url}`).then(
      (resolve) => {
        this.setState({ previewText: resolve.data });
      },
      (reject) => {
        console.log(reject);
      }
    );
  };

  submit = () => {
    let url = document.getElementById("outlined-basic").value;
    console.log(config.getPrefix());
    axois.get(config.getPrefix() + `/api/getNovel?url=${url}`).then(
      (resolve) => {
        this.setState({ previewText: resolve.data });
      },
      (reject) => {
        console.log(reject);
      }
    );
  };

  render() {
    return (
      <div className={Style.requestSection}>
        <div>
          <Stack direction="row" spacing={2}>
            <TextField
              ref={(c) => {
                this.urlInput = c;
              }}
              id="outlined-basic"
              label="小说地址"
              variant="outlined"
            />
            <Button onClick={this.preview}>预览</Button>
            <Button onClick={this.submit}>提交</Button>
          </Stack>
        </div>
        <div className={Style.inexPreview}>{this.state.previewText}</div>
      </div>
    );
  }
}
