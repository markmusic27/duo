import { showHUD, Clipboard, LaunchProps } from "@raycast/api";
import fetch from "node-fetch";
import { SecretKey } from "./secret";

interface CaptureProps {
  content: string;
}

export default async function main(props: LaunchProps<{ arguments: CaptureProps }>) {
  //   const requestOptions = {
  //     method: 'POST',
  //     headers: {
  //         'Authorization': `Bearer ${token}`,
  //         'Content-Type': 'application/json'
  //     },
  //     body: JSON.stringify({
  //         title: 'foo',
  //         body: 'bar',
  //         userId: 1
  //     })
  // };
  console.log(SecretKey);

  // let res = await fetch("https://jsonplaceholder.typicode.com/posts")

  // console.log(res)

  // const now = new Date();
  // await Clipboard.copy(now.toLocaleDateString());
  // await showHUD("Copied date to clipboard");
}
