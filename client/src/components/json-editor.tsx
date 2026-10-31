import { useCallback, useEffect, useMemo, useRef } from "react"
import CodeMirror, { EditorView, type Extension } from "@uiw/react-codemirror"
import { json, jsonParseLinter } from "@codemirror/lang-json"
import { linter, lintGutter } from "@codemirror/lint"
// import { useTheme } from "next-themes"
import { cn } from "@/lib/utils"

function formatJsonString(raw: string, indent = 2): string {
  if (!raw || !raw.trim()) return raw
  try {
    return JSON.stringify(JSON.parse(raw), null, indent)
  } catch {
    return raw
  }
}

interface JsonEditorProps {
  value: string
  onChange?: (value: string) => void
  className?: string
  readOnly?: boolean
  minHeight?: string
}

export function JsonEditor({
  value,
  onChange,
  className,
  readOnly = false,
  minHeight = "200px",
}: JsonEditorProps) {
  // const { resolvedTheme } = useTheme()
  const isFirstMount = useRef(true)

  const extensions: Extension[] = useMemo(
    () => [
      json(),
      linter(jsonParseLinter()),
      lintGutter(),
      EditorView.lineWrapping,
    ],
    []
  )

  const handleChange = useCallback(
    (val: string) => {
      onChange?.(val)
    },
    [onChange]
  )

  useEffect(() => {
    if (isFirstMount.current && value) {
      isFirstMount.current = false
      const formatted = formatJsonString(value)
      if (formatted !== value) {
        onChange?.(formatted)
      }
    }
  }, [value, onChange])

  return (
    <div
      className={cn(
        "relative rounded-md border border-input bg-transparent text-sm shadow-xs transition-colors",
        "focus-within:border-ring focus-within:ring-[3px] focus-within:ring-ring/50",
        "overflow-hidden font-mono",
        className
      )}
    >
      <CodeMirror
        value={value}
        height="100%"
        minHeight={minHeight}
        // theme={resolvedTheme === "dark" ? "dark" : "light"}
        theme="dark"
        extensions={extensions}
        onChange={handleChange}
        readOnly={readOnly}
        basicSetup={{
          lineNumbers: true,
          foldGutter: true,
          highlightActiveLineGutter: true,
          highlightActiveLine: true,
          bracketMatching: true,
          closeBrackets: true,
          autocompletion: true,
        }}
      />
    </div>
  )
}
