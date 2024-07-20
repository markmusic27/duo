import { showHUD, Clipboard, LaunchProps } from "@raycast/api";

interface CaptureProps {
  content: string
}

export default async function main(props: LaunchProps<{arguments: CaptureProps}>) {
  console.log(props.arguments.content)

  const now = new Date();
  await Clipboard.copy(now.toLocaleDateString());
  await showHUD("Copied date to clipboard");
}
