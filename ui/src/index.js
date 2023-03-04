import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import { BrowserRouter } from "react-router-dom";
import axios from "axios";
import { config } from "./config";

function a() {
  axios
    .get(config.getPrefix() + `/api/config`)
    .then(
      (resolve) => {
        console.log(resolve.data);
        let configs = resolve.data;
        const root = ReactDOM.createRoot(document.getElementById("root"));
        root.render(
          <React.StrictMode>
            <BrowserRouter>
              <App configs={configs} />
            </BrowserRouter>
          </React.StrictMode>
        );
      },
      (reject) => {
        console.log(reject);
      }
    )
    .catch((e) => {
      console.log(e);
    });
}
a();
