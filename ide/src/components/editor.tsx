import { Editor } from "@monaco-editor/react";

interface EditorProps {
  code: string;
  language?: string;
  onChange?: (value: string) => void;
}

export function EditorWindow({ code, language }: EditorProps) {
  return (
    <div className="h-full w-full">
      <Editor
        height="100%"
        width="100%"
        value={code}
        theme="vs-dark"
        language={language}
      />
    </div>
  );
}
