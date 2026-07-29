import SyntaxHighlighter from "react-syntax-highlighter"
import { gruvboxDark } from "react-syntax-highlighter/dist/esm/styles/hljs"

export function Json({ data }: { data: any }) {
  return (
    <SyntaxHighlighter language="json" style={gruvboxDark}>
      {JSON.stringify(data, null, 2)}
    </SyntaxHighlighter>
  )
}
