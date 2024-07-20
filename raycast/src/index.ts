import { showHUD, Clipboard, LaunchProps } from "@raycast/api";
import fetch from "node-fetch";
import { Endpoint, SecretKey } from "./secret";

interface CaptureProps {
  content: string;
}

export default async function main(props: LaunchProps<{ arguments: CaptureProps }>) {
  const requestOptions = {
    method: "POST",
    headers: {
      Authorization: `Bearer ${SecretKey}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      message: props.arguments.content,
    }),
  };

  try {
    const res = await fetch(Endpoint, requestOptions);
    const data = await res.json();

    console.log(data);
  } catch (e) {
    console.log(e);
  }

  // console.log(res)

  // const now = new Date();
  // await Clipboard.copy(now.toLocaleDateString());
  // await showHUD("Copied date to clipboard");
}
