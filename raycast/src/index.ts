import { showHUD, Clipboard, LaunchProps } from "@raycast/api";
import fetch from "node-fetch";
import { Endpoint, SecretKey } from "./secret";

interface CaptureProps {
  content: string;
}

interface ResponseData {
  error: string | undefined;
  message: string | undefined;
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

    if (res.ok) {
      await showHUD("✅ Logged");
      return;
    }

    const data = await res.json();
    const responseData = data as ResponseData;

    await showHUD(`🚨 Error: ${responseData.error}`);
  } catch (e) {
    await showHUD("🚨 Unable to send response");
  }
}
